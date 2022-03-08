mod memory;
use memory::*;

use std::collections::HashMap;
use once_cell::unsync::Lazy;
use std::str::FromStr;
use std::mem::{ManuallyDrop, forget, size_of, transmute};

use wry::{
	application::{
		accelerator::{Accelerator},
		clipboard::Clipboard,
		dpi::{LogicalSize, LogicalPosition, PhysicalPosition},
		event::{Event, WindowEvent},
		event_loop::{ControlFlow, EventLoop},
		global_shortcut::ShortcutManager,
	  menu::{ContextMenu, MenuBar, MenuItemAttributes},
	  system_tray::SystemTrayBuilder,
		window::{WindowBuilder, Fullscreen},
	},
	webview::{WebViewBuilder},
};

use raw_window_handle::{HasRawWindowHandle, RawWindowHandle};

#[cfg(target_os = "windows")]
mod win32 {
  type HWND = *const libc::c_void;

  #[link(name = "user32")]
  extern "C" {
    pub fn GetActiveWindow() ->  HWND;
  }
}

#[cfg(target_os = "macos")]
use objc::{
	msg_send,
	sel,
	sel_impl,
  runtime::{Object, BOOL, YES},
};

//
// C Types
//

type CInt    = libc::c_int;
type CString = *const libc::c_char;
type CBool   = bool;
type CDouble = f64;

#[repr(C)]
pub struct CPosition {
pub x: f64,
pub y: f64,
}

#[repr(C)]
pub struct CSize {
	pub width: f64,
	pub height: f64,
}

#[repr(C)]
#[allow(non_camel_case_types)]
enum CEventType {
	None      = 0,
	Close     = 1,
	Destroyed = 2,
	Focused   = 3,
	Blurred   = 4,
	Resized   = 5,
	Moved     = 6,
	MenuItem  = 7,
	Shortcut  = 8,
}

#[repr(C)]
pub struct CEvent {
	pub event_type: CInt,
	pub window_id:  CInt,
	pub position:   CPosition,
	pub size:       CSize,
	pub menu_id:    CInt,
	pub shortcut:   CString,
}

// NOTE(nick): even though this stuct is not FFI compatible, we use it as an opaque handle on the C/Go side
// so the layout of the data shouldn't matter
type CEventLoop = EventLoop<()>;
type CMenu = MenuBar;
type CContextMenu = ContextMenu;

#[repr(C)]
pub struct CWindow_Options {
	pub always_on_top: CBool,
	pub frameless:     CBool,
	pub fullscreen:    CBool,
	pub size:          CSize,
	pub min_size:      CSize,
	pub max_size:      CSize,
	pub maximized:     CBool,
	pub position:      CPosition,
	pub resizable:     CBool,
	pub title:         CString,
	pub transparent:   CBool,
	pub visible:       CBool,
	pub center:        CBool,
	pub icon:          CIcon,
	pub url:           CString,
	pub html:          CString,
	pub script:        CString,
}

#[repr(C)]
pub struct CMenu_Item {
	pub id: CInt,
	pub title: CString,
	pub enabled: CBool,
	pub selected: CBool,
	pub accelerator: CString,
}

#[repr(C)]
pub struct CIcon {
	pub data: *mut u8,
	pub size: CInt,
}

// NOTE(nick): generic array return value (can be cast into any other typed array)
#[repr(C)]
pub struct CArray {
	pub data: *mut libc::c_void,
	pub count: CInt,
}

#[repr(C)]
pub struct CDisplay {
	pub name: CString,
	pub size: CSize,
	pub position: CPosition,
	pub scale_factor: CDouble,
}

//
// Globals
//

struct Window {
	id: i32,
	webview: wry::webview::WebView,
}

struct Shortcut {
	id: u16,
	accelerator: String,
	shortcut: wry::application::global_shortcut::GlobalShortcut,
	menu_id: i32,
}

static mut WINDOWS: Vec<Window> = Vec::new();
static mut NEXT_WINDOW_ID: i32 = 1;
static mut SHORTCUT_MANAGER: Option<ShortcutManager> = None;
static mut SHORTCUTS: Lazy<HashMap<u16, Shortcut>> = Lazy::new(|| {
  HashMap::new()
});

const STORAGE_SIZE: usize = kilobytes!(256);
static mut STORAGE_BUFFER: [u8; STORAGE_SIZE] = [0; STORAGE_SIZE];
static mut TEMPORARY_STORAGE: Arena = unsafe { Arena::new(&STORAGE_BUFFER as *const u8 as *mut u8, STORAGE_SIZE) };

//
// Helpers
//

fn string_from_cstr(cstr: CString) -> String {
	use std::ffi::CStr;
	let buffer = unsafe { CStr::from_ptr(cstr).to_bytes() };
	String::from_utf8(buffer.to_vec()).unwrap()
}

fn str_from_cstr(cstr: CString) -> &'static str {
	use std::ffi::CStr;
	let result: &CStr = unsafe { CStr::from_ptr(cstr) };
	result.to_str().unwrap()
}

fn cstr_from_string(it: String) -> CString {
	let ptr = unsafe {
		let result = TEMPORARY_STORAGE.write_str(it.as_str());
		TEMPORARY_STORAGE.write_str("\0");
		result
	};

	ptr as *const libc::c_char
}

fn carray_from_string_array(vec: Vec<String>) -> CArray {
	let count = vec.len();
	let array = unsafe { TEMPORARY_STORAGE.push_aligned(size_of::<usize>() * count, 8) };

	let result = CArray { data: array as *mut libc::c_void, count: count as i32 };

	for (index, it) in vec.iter().enumerate() {
		unsafe {
			let ptr = TEMPORARY_STORAGE.write_str(it.as_str());
			TEMPORARY_STORAGE.write_str("\0");

			Arena::copy(transmute(&ptr), (result.data as usize + (size_of::<usize>() * index)) as *mut u8, size_of::<usize>());
		}
	}

	result
}

fn carray_from_vec<T>(vec: Vec<T>) -> CArray {
	// @Robustness: what should this alignment be?
	let count = vec.len();
	let array = unsafe { TEMPORARY_STORAGE.write_aligned(vec.as_ptr() as *mut u8, size_of::<T>() * count as usize, 16) };
	let result = CArray { data: array as *mut libc::c_void, count: count as i32 };
	result
}

macro_rules! find_item {
	($array: expr, $id: expr, $func: expr) => {{
		let mut result = false;

		unsafe {
			let it = $array.iter().find(|&it| it.id == $id);

			if let Some(it) = it {
				$func(it);
				result = true;
			}
		}

		result
	}};
} 

fn register_shortcut(accel_str: &'static str, menu_id: i32) -> (bool, u16, Option<Accelerator>) {
	unsafe {
		if SHORTCUT_MANAGER.is_none() {
			println!("SHORTCUT_MANAGER not initialized!");
			return (false, 0, None);
		}
	}

	let accelerator = Accelerator::from_str(accel_str);

	if accelerator.is_ok() {
		let accelerator = accelerator.unwrap();
		let id = accelerator.clone().id().0;
		let result = unsafe { SHORTCUT_MANAGER.as_mut().unwrap().register(accelerator.clone()) };

		if result.is_ok() {
			let item = Shortcut{ id: id as u16, menu_id, shortcut: result.unwrap(), accelerator: accel_str.to_string() };
			unsafe { SHORTCUTS.insert(id, item); };

			return (true, id, Some(accelerator));
		}
	}

	(false, 0, None)
}

//
// API
//

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn create_event_loop() -> CEventLoop {
	//
	// NOTE(nick): If this changes, go and update hostbridge.h EventLoop size
	// @Robustness: make this a static assertion
	//
	assert_eq!(size_of::<CEventLoop>(), 40);

	let result = EventLoop::new();


	// @Cleanup: move this to a app_init() method or something?
	// Calling this multiple times will be result in the wrong ShortcutManager being used
	unsafe {
		if SHORTCUT_MANAGER.is_none() {
			SHORTCUT_MANAGER = Some(ShortcutManager::new(&result));
		}
	}
	
	//
	// NOTE(nick): prevent the EventLoop's destructor from being called here
	// Other places that take the event_loop as an argument will also need to call `std::mem::forget`
	//
	let mut result = ManuallyDrop::new(result);

	unsafe {
		ManuallyDrop::take(&mut result)
	}
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn window_create(event_loop: CEventLoop, options: CWindow_Options, menu: CMenu) -> i32 {
	let title = str_from_cstr(options.title);
	let fullscreen = if options.fullscreen { Some(Fullscreen::Borderless(None)) } else { None };

	let mut window_builder = WindowBuilder::new()
		.with_menu(menu)
		.with_always_on_top(options.always_on_top)
		.with_decorations(!options.frameless)
		.with_fullscreen(fullscreen)
		.with_maximized(options.maximized)
		.with_resizable(options.visible)
		.with_title(title)
		.with_transparent(options.transparent)
		.with_visible(options.visible);

	if options.size.width > 0.0 || options.size.height > 0.0 {
		window_builder = window_builder.with_inner_size(LogicalSize::new(options.size.width, options.size.height));
	}

	if options.min_size.width > 0.0 || options.min_size.height > 0.0 {
		window_builder = window_builder.with_min_inner_size(LogicalSize::new(options.min_size.width, options.min_size.height));
	}

	if options.max_size.width > 0.0 || options.max_size.height > 0.0 {
		window_builder = window_builder.with_max_inner_size(LogicalSize::new(options.max_size.width, options.max_size.height));
	}

	if !options.center {
		window_builder = window_builder.with_position(LogicalPosition::new(options.position.x, options.position.y));
	}

	// @Incomplete: add support for icons
	/*
	if options.icon.size > 0 {
		let icon = options.icon;
		let icon_buf = unsafe { Vec::<u8>::from_raw_parts(icon.data, icon.size as usize, icon.size as usize) };

		// @Incomplete: size is in pixels but we only have bytes here
		let icon = Icon::from_rgba(icon_buf, 0, 0);
		if let Ok(icon) = icon {
			window_builder = window_builder.with_window_icon(Some(icon));
		}
	}
	*/

	let maybe_window = window_builder.build(&event_loop);
	forget(event_loop);

	if !maybe_window.is_ok() {
		return -2;
	}

	let window = maybe_window.unwrap();

	if options.center {
		let monitor = window.current_monitor();

		if let Some(monitor) = monitor {
			let size = window.outer_size();
			let monitor_size = monitor.size();
			let center = PhysicalPosition::new((monitor_size.width - size.width) / 2, (monitor_size.height - size.height) / 2);
			window.set_outer_position(center);
		}
	}

	let webview_builder = WebViewBuilder::new(window);

	if !webview_builder.is_ok() {
		return -3;
	}

	let webview_builder = webview_builder.unwrap();

	let html = string_from_cstr(options.html);
	let url = str_from_cstr(options.url);
	let script = str_from_cstr(options.script);

	let mut webview_builder = webview_builder
		.with_transparent(options.transparent);

	if script.len() > 0 {
		webview_builder = webview_builder.with_initialization_script(script);
	}

	if url.len() > 0 {
		let success = webview_builder.with_url(url);

		if success.is_ok() {
			webview_builder = success.unwrap();
		} else {
			return -4;
		}
	} else {
		// From the wry docs: This will be ignored if url is already provided.
		let success = webview_builder.with_html(html);

		if success.is_ok() {
			webview_builder = success.unwrap();
		} else {
			return -4;
		}
	}

	let result = webview_builder.build();

	if !result.is_ok() {
		return -5;
	}

	let webview = result.unwrap();

	let result: i32;

	unsafe {
		result = NEXT_WINDOW_ID;
		NEXT_WINDOW_ID += 1;
		let the_window = Window{ id: result, webview };

		WINDOWS.push(the_window);
	}

	return result;
}

#[no_mangle]
pub extern "C" fn window_destroy(window_id: CInt) -> CBool {
	let mut result = false;

	unsafe {
		let found = WINDOWS.iter().position(|it| it.id == window_id);
		if let Some(index) = found {
			WINDOWS.remove(index);
			result = true;
		}
	}

	result
}

#[no_mangle]
pub extern "C" fn window_set_title(window_id: CInt, title: CString) -> CBool {
	let title = string_from_cstr(title);
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_title(&title))
}

#[no_mangle]
pub extern "C" fn window_set_visible(window_id: CInt, is_visible: CBool) -> CBool {
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_visible(is_visible))
}

#[no_mangle]
pub extern "C" fn window_set_focused(window_id: CInt) -> CBool {
	find_item!(WINDOWS, window_id, |it: &Window| {
		it.webview.focus();
		it.webview.window().set_focus();
	})
}

#[no_mangle]
pub extern "C" fn window_set_fullscreen(window_id: CInt, is_fullscreen: CBool) -> CBool {
	let fullscreen = if is_fullscreen { Some(Fullscreen::Borderless(None)) } else { None };
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_fullscreen(fullscreen))
}

#[no_mangle]
pub extern "C" fn window_set_maximized(window_id: CInt, is_maximized: CBool) -> CBool {
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_maximized(is_maximized))
}

#[no_mangle]
pub extern "C" fn window_set_minimized(window_id: CInt, is_minimized: CBool) -> CBool {
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_minimized(is_minimized))
}

#[no_mangle]
pub extern "C" fn window_set_size(window_id: CInt, size: CSize) -> CBool {
	let size = LogicalSize::new(size.width, size.height);
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_inner_size(size))
}

#[no_mangle]
pub extern "C" fn window_set_min_size(window_id: CInt, size: CSize) -> CBool {
	let size = LogicalSize::new(size.width, size.height);
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_min_inner_size(Some(size)))
}

#[no_mangle]
pub extern "C" fn window_set_max_size(window_id: CInt, size: CSize) -> CBool {
	let size = LogicalSize::new(size.width, size.height);
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_max_inner_size(Some(size)))
}

#[no_mangle]
pub extern "C" fn window_set_resizable(window_id: CInt, is_resizable: CBool) -> CBool {
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_resizable(is_resizable))
}

#[no_mangle]
pub extern "C" fn window_set_always_on_top(window_id: CInt, is_on_top: CBool) -> CBool {
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_always_on_top(is_on_top))
}

#[no_mangle]
pub extern "C" fn window_set_position(window_id: CInt, position: CPosition) -> CBool {
	let position = LogicalPosition::new(position.x, position.y);
	find_item!(WINDOWS, window_id, |it: &Window| it.webview.window().set_outer_position(position))
}

#[no_mangle]
pub extern "C" fn window_get_outer_position(window_id: CInt) -> CPosition {
	let mut result = CPosition{ x: 0.0, y: 0.0 };

	find_item!(WINDOWS, window_id, |it: &Window| {
		let position = it.webview.window().outer_position();
		if position.is_ok() {
			let position = position.unwrap();
			result.x = position.x as f64;
			result.y = position.y as f64;
		}
	});

	result
}

#[no_mangle]
pub extern "C" fn window_get_outer_size(window_id: CInt) -> CSize {
	let mut result = CSize{ width: 0.0, height: 0.0 };

	find_item!(WINDOWS, window_id, |it: &Window| {
		let size = it.webview.window().outer_size();
		result.width = size.width as f64;
		result.height = size.height as f64;
	});

	result
}

#[no_mangle]
pub extern "C" fn window_get_inner_position(window_id: CInt) -> CPosition {
	let mut result = CPosition{ x: 0.0, y: 0.0 };

	find_item!(WINDOWS, window_id, |it: &Window| {
		let position = it.webview.window().inner_position();
		if position.is_ok() {
			let position = position.unwrap();
			result.x = position.x as f64;
			result.y = position.y as f64;
		}
	});

	result
}

#[no_mangle]
pub extern "C" fn window_get_inner_size(window_id: CInt) -> CSize {
	let mut result = CSize{ width: 0.0, height: 0.0 };

	find_item!(WINDOWS, window_id, |it: &Window| {
		let size = it.webview.window().inner_size();
		result.width = size.width as f64;
		result.height = size.height as f64;
	});

	result
}

#[no_mangle]
pub extern "C" fn window_get_dpi_scale(window_id: CInt) -> CDouble {
	let mut result = 1.0;
	find_item!(WINDOWS, window_id, |it: &Window| result = it.webview.window().scale_factor());
	result
}

#[no_mangle]
pub extern "C" fn window_is_visible(window_id: CInt) -> CBool {
	let mut result = false;
	find_item!(WINDOWS, window_id, |it: &Window| result = it.webview.window().is_visible());
	result
}

#[no_mangle]
pub extern "C" fn window_is_focused(window_id: CInt) -> CBool {
	let mut result = false;

	find_item!(WINDOWS, window_id, |it: &Window| {

		let handle = it.webview.window().raw_window_handle();

		#[cfg(target_os = "windows")]
		{
			if let RawWindowHandle::Win32(handle) = handle {
				result = win32::GetActiveWindow() == handle.hwnd;
			}
		}

		#[cfg(target_os = "macos")]
		{
			if let RawWindowHandle::AppKit(handle) = handle {
				let is_key_window: BOOL = msg_send!(handle.ns_window as *const Object, isKeyWindow);

				result = is_key_window == YES;
			}
		}

	});

	result
}


#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn menu_create() -> CMenu {
	//
	// NOTE(nick): If this changes, go and update hostbridge.h Menu size
	// @Robustness: make this a static assertion
	//
	#[cfg(target_os = "macos")]
	assert_eq!(size_of::<CMenu>(), 16);
	#[cfg(target_os = "windows")]
	assert_eq!(size_of::<CMenu>(), 64);

	let result = MenuBar::new();
	
	let mut result = ManuallyDrop::new(result);

	unsafe {
		ManuallyDrop::take(&mut result)
	}
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn menu_add_item(mut menu: CMenu, item: CMenu_Item) -> CBool {
	// @Copypaste
	let title = str_from_cstr(item.title);

	// @Incomplete: use `role` to handle native items
	//menu.add_native_item(MenuItem::About("Todos".to_string()));

	let mut result =
		MenuItemAttributes::new(title)
			.with_id(wry::application::menu::MenuId(item.id as u16))
			.with_enabled(item.enabled)
			.with_selected(item.selected);

	let mut success = true;

	let accel_str = str_from_cstr(item.accelerator);
	if accel_str.len() > 0 {
		let (ok, _, accelerator) = register_shortcut(accel_str, item.id);

		if ok {
			let accelerator = accelerator.unwrap();
			result = result.with_accelerators(&accelerator);
		} else {
			success = false;
		}
	}

	menu.add_item(result);

	forget(menu);

	success
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn menu_add_submenu(mut menu: CMenu, title: CString, enabled: CBool, submenu: CMenu) -> CBool {
	let title = str_from_cstr(title);
	menu.add_submenu(title, enabled, submenu);

	forget(menu);

	true
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn context_menu_create() -> CContextMenu {
	//
	// NOTE(nick): If this changes, go and update hostbridge.h Menu size
	// @Robustness: make this a static assertion
	//
	#[cfg(target_os = "macos")]
	assert_eq!(size_of::<CContextMenu>(), 16);
	#[cfg(target_os = "windows")]
	assert_eq!(size_of::<CContextMenu>(), 64);

	let result = CContextMenu::new();
	
	let mut result = ManuallyDrop::new(result);

	unsafe {
		ManuallyDrop::take(&mut result)
	}
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn context_menu_add_item(mut menu: CContextMenu, item: CMenu_Item) -> CBool {
	let title = str_from_cstr(item.title);

	// @Incomplete: use `role` to handle native items
	//menu.add_native_item(MenuItem::About("Todos".to_string()));

	let mut result =
		MenuItemAttributes::new(title)
			.with_id(wry::application::menu::MenuId(item.id as u16))
			.with_enabled(item.enabled)
			.with_selected(item.selected);

	let mut success = true;

	let accel_str = str_from_cstr(item.accelerator);
	if accel_str.len() > 0 {
		let (ok, _, accelerator) = register_shortcut(accel_str, item.id);

		if ok {
			let accelerator = accelerator.unwrap();
			result = result.with_accelerators(&accelerator);
		} else {
			success = false;
		}
	}

	menu.add_item(result);

	forget(menu);

	success
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn context_menu_add_submenu(mut menu: CContextMenu, title: CString, enabled: CBool, submenu: CContextMenu) -> CBool {
	let title = str_from_cstr(title);
	menu.add_submenu(title, enabled, submenu);

	forget(menu);

	true
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn tray_set_system_tray(event_loop: CEventLoop, icon: CIcon, tray_menu: CContextMenu) -> CBool {
	let mut icon = unsafe { Vec::<u8>::from_raw_parts(icon.data, icon.size as usize, icon.size as usize) };

	#[cfg(target_os = "windows")]
	{
		use std::io::Cursor;

		let mut icon_dir = ico::IconDir::new(ico::ResourceType::Icon);
		let image = ico::IconImage::read_png(&mut Cursor::new(&mut icon));
		
		if image.is_ok() {
			let image = image.unwrap();

			let entry = ico::IconDirEntry::encode(&image);
			if entry.is_ok() {
				let entry = entry.unwrap();
				icon_dir.add_entry(entry);

				let mut bytes: Vec<u8> = Vec::new();
				let result = icon_dir.write(&mut Cursor::new(&mut bytes));

				if result.is_ok() {
					forget(icon); // Prevent dealloc
					icon = bytes; // @MemoryLeak
				}
			}
		}
	}

	let system_tray = SystemTrayBuilder::new(icon.clone(), Some(tray_menu)).build(&event_loop).unwrap();

	// @Incomplete: in the future we will want to store the system_tray somewhere (probably in Go), so that
	// you can call system_tray.set_icon to change the icon dynamically

	forget(event_loop);
	forget(icon); // NOTE(nick): prevent rust from trying to dealloc something it doesn't own
	forget(system_tray); // @MemoryLeak

	true
}

#[no_mangle]
pub extern fn shell_show_notification(title: CString, subtitle: CString, body: CString) -> CBool {
	let title = str_from_cstr(title);
	let subtitle = str_from_cstr(subtitle);
	let body = str_from_cstr(body);

	let result = notify_rust::Notification::new()
		.summary(title)
		.subtitle(subtitle)
		.body(body)
		.show();

	result.is_ok()
}

#[no_mangle]
pub extern fn shell_show_dialog(title: CString, body: CString, level: CString, buttons: CString) -> CBool {
	let title = str_from_cstr(title);
	let body = str_from_cstr(body);
	let level = str_from_cstr(level);
	let buttons = str_from_cstr(buttons);

	let level = match level {
		"info" => rfd::MessageLevel::Info,
		"warning" => rfd::MessageLevel::Warning,
		"error" => rfd::MessageLevel::Error,
		_ => rfd::MessageLevel::Info,
	};

	let buttons = match buttons {
		"ok" => rfd::MessageButtons::Ok,
		"okcancel" => rfd::MessageButtons::OkCancel,
		"yesno" => rfd::MessageButtons::YesNo,
		_ => rfd::MessageButtons::Ok,
	};

	let result = rfd::MessageDialog::new()
		.set_title(title)
		.set_description(body)
		.set_buttons(buttons)
		.set_level(level)
		.show();

	result
}

#[no_mangle]
pub extern fn reset_temporary_storage() {
	unsafe {
		TEMPORARY_STORAGE.reset();
	}
}

#[no_mangle]
pub extern fn shell_show_file_picker(title: CString, directory: CString, filename: CString, mode: CString, filters: CString) -> CArray {
	let title = str_from_cstr(title);
	let directory = str_from_cstr(directory);
	let filename = str_from_cstr(filename);
	let mode = str_from_cstr(mode);

	let mut picker = rfd::FileDialog::new();

	if title.len() > 0 {
		picker = picker.set_title(title);
	}

	if directory.len() > 0 {
		picker = picker.set_directory(&directory);
	}

	if filename.len() > 0 {
		picker = picker.set_file_name(filename);
	}

	let filters = string_from_cstr(filters);
	let filters = filters.split("|");
	
	for filter in filters {
		if filter.len() > 0 {
			let mut label = "";
			let mut extensions = filter;

			let index = filter.find(':');
			if index.is_some() {
				let index = index.unwrap();
				label = &filter[0..index];
				extensions = &filter[index+1..];
			}

			let extensions: Vec<&str> = extensions.split(",").collect();

			picker = picker.add_filter(label, &extensions);
		}
	}

	let mut paths: Vec<std::path::PathBuf> = Vec::new();

	match mode {
		"pickfolder" => {
			let res = picker.pick_folder();
			if res.is_some() { paths.push(res.unwrap()); }
		},
		"savefile" => {
			let res = picker.save_file();
			if res.is_some() { paths.push(res.unwrap()); }
		},
		"pickfiles" => {
			let res = picker.pick_files();
			if res.is_some() { paths.append(&mut res.unwrap()); }
		},
		_ => { // "pickfile"
			let res = picker.pick_file();
			if res.is_some() { paths.push(res.unwrap()); }
		},
	};

	let array: Vec<String> = paths.into_iter().map(|it| it.clone().into_os_string().into_string().unwrap()).collect();
	carray_from_string_array(array)
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn screen_get_available_displays(event_loop: CEventLoop) -> CArray {
	let mut monitors: Vec<wry::application::monitor::MonitorHandle> = Vec::new();

	let first_user_window = unsafe { WINDOWS.first() };

	if let Some(first_user_window) = first_user_window {
		let window = first_user_window.webview.window();
		monitors = window.available_monitors().collect::<Vec<_>>();
	} else {
		// @Incomplete: does this cause any visual artifcats on any operating systems?
		let window_builder = WindowBuilder::new()
			.with_visible(false)
			.with_decorations(false)
			.with_transparent(true)
			.build(&event_loop);

		if window_builder.is_ok() {
			let window = window_builder.unwrap();
			monitors = window.available_monitors().collect::<Vec<_>>();
		}
	}

	forget(event_loop);

	let array: Vec<CDisplay> = monitors.into_iter().map(|it| {
		let name = it.name().unwrap();
		let size = it.size();
		let position = it.position();
		let scale_factor = it.scale_factor();

		CDisplay{
			name: cstr_from_string(name),
			size: CSize{width: size.width as f64, height: size.height as f64},
			position: CPosition{x: position.x as f64, y: position.y as f64},
			scale_factor,
		}
	}).collect();

	carray_from_vec::<CDisplay>(array)
}

#[no_mangle]
pub extern "C" fn shell_read_clipboard() -> CString {
	let cliboard = Clipboard::new();
	let content = cliboard.read_text();

	if !content.is_some() {
		return std::ptr::null();
	}

	let content = content.unwrap();
	cstr_from_string(content)
}

#[no_mangle]
pub extern "C" fn shell_write_clipboard(text: CString) -> CBool {
	let mut clipboard = Clipboard::new();

	let text = str_from_cstr(text);
	clipboard.write_text(&text);

	// @Speed: don't most OSs tell you if this succeeds? At least windows does
	let written = clipboard.read_text();
	if written.is_some() {
		let written = written.unwrap();
		return written == text;
	}

	false
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn shell_register_shortcut(accelerator: CString) -> CBool {
	let accelerator = str_from_cstr(accelerator);
	let (success, _, _) = register_shortcut(accelerator, 0);
	success
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn shell_is_shortcut_registered(accelerator: CString) -> CBool {
	unsafe {
		if SHORTCUT_MANAGER.is_none() {
			return false;
		}
	}

	let accelerator = str_from_cstr(accelerator);
	let shortcut = Accelerator::from_str(accelerator);

	if shortcut.is_ok() {
		let shortcut = shortcut.unwrap();
		let result = unsafe { SHORTCUT_MANAGER.as_mut().unwrap().is_registered(&shortcut) };
		return result;
	}

	false
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn shell_unregister_shortcut(accelerator: CString) -> CBool {
	unsafe {
		if SHORTCUT_MANAGER.is_none() {
			return false;
		}
	}

	let accelerator = str_from_cstr(accelerator);
	let accelerator = Accelerator::from_str(accelerator);

	if accelerator.is_ok() {
		let accelerator = accelerator.unwrap();

		unsafe {
			let id = accelerator.id().0;
			let it = SHORTCUTS.get_mut(&id);

			if it.is_some() {
				let it = it.unwrap();

				let result = SHORTCUT_MANAGER.as_mut().unwrap().unregister(it.shortcut.clone());
				SHORTCUTS.remove(&id);

				return result.is_ok();
			}
		}
	}

	false
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn shell_unregister_all_shortcuts() -> CBool {
	unsafe {
		if SHORTCUT_MANAGER.is_some() {
			return false;
		}
	}

	unsafe {
		SHORTCUTS.clear();
	}

	let result = unsafe { SHORTCUT_MANAGER.as_mut().unwrap().unregister_all() };
	return result.is_ok();
}


#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn run(event_loop: CEventLoop, user_callback: unsafe extern "C" fn(CEvent)) {
	event_loop.run(move |event, _, control_flow| {
		*control_flow = ControlFlow::Poll;

		let mut result = CEvent{
			event_type: CEventType::None as i32,
			window_id: -1,
			position: CPosition{x: 0.0, y: 0.0},
			size: CSize{width: 0.0, height: 0.0},
			menu_id: 0,
			shortcut: std::ptr::null(),
		};

		match event {
			Event::WindowEvent { event, window_id, .. } => {
				// @Incomplete: when a window is being destroyed we still want to get its user_window_id
				// Right now, it will be -1
				let user_window_id = unsafe {
					let it = WINDOWS.iter().find(|&it| it.webview.window().id() == window_id);

					if let Some(it) = it {
						it.id
					} else {
						-1
					}
				};

				let event_type = match event {
					WindowEvent::CloseRequested{ .. } => CEventType::Close as i32,
					WindowEvent::Destroyed{ .. }      => CEventType::Destroyed as i32,
					WindowEvent::Focused(mut focus)   => {
						// NOTE(nick): the "focus" argument _should_ be according to the docs:
						// > The parameter is true if the window has gained focus, and false if it has lost focus.
						// But in reality the _opposite_ seems to be true at the moment on win32
						if cfg!(windows) {
							focus = !focus; // @Hack
						}

						if focus {
							CEventType::Focused as i32
						} else {
							CEventType::Blurred as i32
						}
					},
					WindowEvent::Resized{ .. }        => CEventType::Resized as i32,
					WindowEvent::Moved{ .. }          => CEventType::Moved as i32,
					_ => CEventType::None as i32,
				};

				result.window_id  = user_window_id;
				result.event_type = event_type;

				match event {
					WindowEvent::Moved(position)   => {
						result.position = CPosition{x: position.x as f64, y: position.y as f64}
					},
					WindowEvent::Resized(_) => {
						// NOTE(nick): Resized event doesn't currently return the correct window size
						// result.size = CSize{width: size.width as f64, height: size.height as f64}

						unsafe {
							let it = WINDOWS.iter().find(|&it| it.webview.window().id() == window_id);

							if let Some(it) = it {
								let size = it.webview.inner_size();
								result.size = CSize{width: size.width as f64, height: size.height as f64};

								let _ = it.webview.resize();
							}
						}
					},
					_ => {}
				};
			},
			Event::MenuEvent { window_id, menu_id, .. } => {
				result.event_type = CEventType::MenuItem as i32;
				result.menu_id = menu_id.0 as i32;

				if let Some(window_id) = window_id {
					let user_window_id = unsafe {
						let it = WINDOWS.iter().find(|&it| it.webview.window().id() == window_id);

						if let Some(it) = it {
							it.id
						} else {
							-1
						}
					};

					result.window_id = user_window_id;
				}
			},
			Event::GlobalShortcutEvent(hotkey_id) => {
				let id = hotkey_id.0;

				unsafe {
					let it = SHORTCUTS.get(&id);

					if let Some(it) = it {
						if it.menu_id == 0 {
							result.event_type = CEventType::Shortcut as i32;
						} else {
							result.event_type = CEventType::MenuItem as i32;
							result.menu_id = it.menu_id as i32;
						}

						result.shortcut = cstr_from_string(it.accelerator.clone());
					}
				}
			},
			_ => (),
		}

		unsafe {
			user_callback(result);
		}
	});
}
