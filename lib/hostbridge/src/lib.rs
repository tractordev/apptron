mod memory;
use memory::*;

use std::ffi::CStr;
use std::str::FromStr;
use std::mem::{ManuallyDrop, forget, size_of};

use wry::{
	application::{
		accelerator::{Accelerator},
		event::{Event, WindowEvent},
		event_loop::{ControlFlow, EventLoop},
		//global_shortcut::ShortcutManager,
	  menu::{ContextMenu, MenuBar, MenuItemAttributes},
	  system_tray::SystemTrayBuilder,
		window::{WindowBuilder, Fullscreen},
	},
	webview::{WebViewBuilder},
};

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
	Resized   = 4,
	Moved     = 5,
	MenuItem  = 6,
}

#[repr(C)]
pub struct CEvent {
	pub event_type: CInt,
	pub window_id:  CInt,
	pub position:   CPosition,
	pub size:       CSize,
	pub menu_id:    CInt,
}

// NOTE(nick): even though this stuct is not FFI compatible, we use it as an opaque handle on the C/Go side
// so the layout of the data shouldn't matter
type CEventLoop = EventLoop<()>;
type CMenu = MenuBar;
type CContextMenu = ContextMenu;

#[repr(C)]
pub struct CWindow_Options {
	pub transparent: CBool,
	pub decorations: CBool,
	pub html: CString,
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

#[repr(C)]
pub struct CStringArray {
	pub data: *mut *mut u8,
	pub count: CInt,
}

struct Window {
	id: i32,
	webview: wry::webview::WebView,
}

static mut GLOBAL_WINDOWS: Vec<Window> = Vec::new();

const STORAGE_SIZE: usize = kilobytes!(256);
static mut STORAGE_BUFFER: [u8; STORAGE_SIZE] = [0; STORAGE_SIZE];
static mut TEMPORARY_STORAGE: Arena = unsafe { Arena::new(&STORAGE_BUFFER as *const u8 as *mut u8, STORAGE_SIZE) };

fn string_from_cstr(cstr: CString) -> String {
	let buffer = unsafe { CStr::from_ptr(cstr).to_bytes() };
	String::from_utf8(buffer.to_vec()).unwrap()
}

fn str_from_cstr(cstr: CString) -> &'static str {
  let result: &CStr = unsafe { CStr::from_ptr(cstr) };
  result.to_str().unwrap()
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

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn create_event_loop() -> CEventLoop {
	//
	// NOTE(nick): If this changes, go and update hostbridge.h EventLoop size
	// @Robustness: make this a static assertion
	//
	assert_eq!(size_of::<CEventLoop>(), 40);

	let result = EventLoop::new();
	
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
	let maybe_window = WindowBuilder::new()
		.with_title("")
		.with_menu(menu)
		.with_decorations(options.decorations)
		.with_transparent(options.transparent)
		.build(&event_loop);


	forget(event_loop);

	if !maybe_window.is_ok() {
		return -2;
	}

	let window = maybe_window.unwrap();

	let maybe_webview_builder = WebViewBuilder::new(window);

	if !maybe_webview_builder.is_ok() {
		return -3;
	}

	let html = string_from_cstr(options.html);

	let maybe_webview = maybe_webview_builder.unwrap()
		.with_transparent(options.transparent)
		.with_html(html);

	if !maybe_webview.is_ok() {
		return -4;
	}

	let result = maybe_webview.unwrap().build();

	if !result.is_ok() {
		return -5;
	}

	let webview = result.unwrap();

	let result: i32;

	unsafe {
		result = GLOBAL_WINDOWS.len() as i32;
		let the_window = Window{ id: result, webview };

		GLOBAL_WINDOWS.push(the_window);
	}

	return result;
}

#[no_mangle]
pub extern "C" fn window_destroy(window_id: CInt) -> CBool {
	let mut result = false;

	unsafe {
		let found = GLOBAL_WINDOWS.iter().position(|it| it.id == window_id);
		if let Some(index) = found {
			GLOBAL_WINDOWS.remove(index);
			result = true;
		}
	}

	result
}

#[no_mangle]
pub extern "C" fn window_set_title(window_id: CInt, title: CString) -> CBool {
	let title = string_from_cstr(title);
	find_item!(GLOBAL_WINDOWS, window_id, |it: &Window| it.webview.window().set_title(&title))
}

#[no_mangle]
pub extern "C" fn window_set_visible(window_id: CInt, is_visible: CBool) -> CBool {
	find_item!(GLOBAL_WINDOWS, window_id, |it: &Window| it.webview.window().set_visible(is_visible))
}

#[no_mangle]
pub extern "C" fn window_set_fullscreen(window_id: CInt, is_fullscreen: CBool) -> CBool {
	let fullscreen = if is_fullscreen { Some(Fullscreen::Borderless(None)) } else { None };
	find_item!(GLOBAL_WINDOWS, window_id, |it: &Window| it.webview.window().set_fullscreen(fullscreen))
}

#[no_mangle]
pub extern "C" fn window_get_outer_position(window_id: CInt) -> CPosition {
	let mut result = CPosition{ x: 0.0, y: 0.0 };

	find_item!(GLOBAL_WINDOWS, window_id, |it: &Window| {
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

	find_item!(GLOBAL_WINDOWS, window_id, |it: &Window| {
		let size = it.webview.window().outer_size();
		result.width = size.width as f64;
		result.height = size.height as f64;
	});

	result
}

#[no_mangle]
pub extern "C" fn window_get_inner_position(window_id: CInt) -> CPosition {
	let mut result = CPosition{ x: 0.0, y: 0.0 };

	find_item!(GLOBAL_WINDOWS, window_id, |it: &Window| {
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

	find_item!(GLOBAL_WINDOWS, window_id, |it: &Window| {
		let size = it.webview.window().inner_size();
		result.width = size.width as f64;
		result.height = size.height as f64;
	});

	result
}

#[no_mangle]
pub extern "C" fn window_get_dpi_scale(window_id: CInt) -> CDouble {
	let mut result = 1.0;
	find_item!(GLOBAL_WINDOWS, window_id, |it: &Window| result = it.webview.window().scale_factor());
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

	let accelerator = string_from_cstr(item.accelerator);
	let mut success = true;

	if accelerator.len() > 0 {
		let accelerator: &str = &accelerator[..];
		let parsed = Accelerator::from_str(accelerator);

		if let Ok(it) = parsed {
			result = result.with_accelerators(&it);
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

	let accelerator = string_from_cstr(item.accelerator);
	let mut success = true;

	if accelerator.len() > 0 {
		let accelerator: &str = &accelerator[..];
		let parsed = Accelerator::from_str(accelerator);

		if let Ok(it) = parsed {
			result = result.with_accelerators(&it);

			/*
			let mut shortcut_manager = ShortcutManager::new(&event_loop);
			shortcut_manager.register(it.clone()).unwrap();
			// @MemoryLeak
			forget(shortcut_manager);
			*/
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
	let icon = unsafe { Vec::<u8>::from_raw_parts(icon.data, icon.size as usize, icon.size as usize) };

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
pub extern fn shell_show_file_picker(title: CString, directory: CString, filename: CString, mode: CString, filters: CString) -> CStringArray {
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

	let count = paths.len();
	let array = unsafe { TEMPORARY_STORAGE.push(size_of::<usize>() * count) };

	let result: CStringArray = CStringArray { data: array as *mut *mut u8, count: count as i32 };

	for (i, path) in paths.iter().enumerate() {
		let path = path.clone().into_os_string().into_string();

		if path.is_ok() {
			let path = path.unwrap();

			unsafe {
				let ptr = TEMPORARY_STORAGE.write(path.as_str());
				TEMPORARY_STORAGE.write("\0");

				Arena::copy(std::mem::transmute(&ptr), (result.data as usize + (size_of::<usize>() * i)) as *mut u8, size_of::<usize>());
		  }
		}
	}

	result
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
		};

		match event {
			Event::WindowEvent { event, window_id, .. } => {
				// @Incomplete: when a window is being destroyed we still want to get its user_window_id
				// Right now, it will be -1
				let user_window_id = unsafe {
					let it = GLOBAL_WINDOWS.iter().find(|&it| it.webview.window().id() == window_id);

					if let Some(it) = it {
						it.id
					} else {
						-1
					}
				};

				let event_type = match event {
					WindowEvent::CloseRequested{ .. } => CEventType::Close as i32,
					WindowEvent::Destroyed{ .. }      => CEventType::Destroyed as i32,
					WindowEvent::Focused{ .. }        => CEventType::Focused as i32,
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
							let it = GLOBAL_WINDOWS.iter().find(|&it| it.webview.window().id() == window_id);

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
						let it = GLOBAL_WINDOWS.iter().find(|&it| it.webview.window().id() == window_id);

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
				// @Incomplete: need a way to match these hotkey_id's back to their original MenuItems (or possibly keyboard shortcuts)
				// that triggered them
				println!("GlobalShortcutEvent {:?}", hotkey_id);
			},
			_ => (),
		}

		unsafe {
			user_callback(result);
		}
	});
}
