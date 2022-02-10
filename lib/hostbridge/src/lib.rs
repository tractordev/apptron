use std::ffi::CStr;
use std::mem::{ManuallyDrop, forget, size_of};
use std::option::Option;
use std::cell::RefCell;

use wry::{
	application::{
		event::{Event, WindowEvent},
		event_loop::{ControlFlow, EventLoop},
		window::{WindowBuilder, Fullscreen},
	},
	webview::{WebViewBuilder},
};

struct Window {
	id: i32,
	webview: wry::webview::WebView,
}

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
}

#[repr(C)]
pub struct CEvent {
	pub event_type: CInt,
	pub window_id:  CInt,
	pub position:   CPosition,
	pub size:       CSize,
}

// NOTE(nick): even though this stuct is not FFI compatible, we use it as an opaque handle on the C/Go side
// so the layout of the data shouldn't matter
type CEventLoop = EventLoop<()>;

#[repr(C)]
pub struct CWindow_Options {
	pub transparent: CBool,
	pub decorations: CBool,
	pub html: CString,
}

thread_local! {
	static GLOBAL_WINDOWS: RefCell<Vec<Window>> = RefCell::new(Vec::new());
}

fn string_from_cstr(cstr: CString) -> String {
	let buffer = unsafe { CStr::from_ptr(cstr).to_bytes() };
	String::from_utf8(buffer.to_vec()).unwrap()
}

fn find_window_by_id(windows: &Vec<Window>, window_id: i32) -> Option<&Window> {
	return windows.iter().find(|&it| it.id == window_id);
}

macro_rules! find_local_window {
	($window_id: expr, $func: expr) => {{
		let mut result = false;

		GLOBAL_WINDOWS.with(|windows| {
			let array = windows.borrow();
			let it = find_window_by_id(&array, $window_id);

			if let Some(it) = it {
				$func(it);
				result = true;
			}
		});

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
  let mut r2 = ManuallyDrop::new(result);

  unsafe {
	ManuallyDrop::take(&mut r2)
  }
}

#[no_mangle]
#[allow(improper_ctypes_definitions)]
pub extern "C" fn window_create(event_loop: CEventLoop, options: CWindow_Options) -> i32 {
	let maybe_window = WindowBuilder::new()
	.with_title("")
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

	let mut result: i32 = -1;

	GLOBAL_WINDOWS.with(|windows| {
		let mut array = windows.borrow_mut();

		result = array.len() as i32;
		let the_window = Window{ id: result, webview };
		array.push(the_window);
	});

	return result;
}

#[no_mangle]
pub extern "C" fn window_destroy(window_id: CInt) -> CBool {
	let mut result = false;

	GLOBAL_WINDOWS.with(|windows| {
		let mut array = windows.borrow_mut();

		let found = array.iter().position(|it| it.id == window_id);
		if let Some(index) = found {
			array.remove(index);
			result = true;
		}
	});

	result
}

#[no_mangle]
pub extern "C" fn window_set_title(window_id: CInt, title: CString) -> CBool {
	let title = string_from_cstr(title);
	find_local_window!(window_id, |it: &Window| it.webview.window().set_title(&title))
}

#[no_mangle]
pub extern "C" fn window_set_visible(window_id: CInt, is_visible: CBool) -> CBool {
	find_local_window!(window_id, |it: &Window| it.webview.window().set_visible(is_visible))
}

#[no_mangle]
pub extern "C" fn window_set_fullscreen(window_id: CInt, is_fullscreen: CBool) -> CBool {
	let fullscreen = if is_fullscreen { Some(Fullscreen::Borderless(None)) } else { None };
	find_local_window!(window_id, |it: &Window| it.webview.window().set_fullscreen(fullscreen))
}

#[no_mangle]
pub extern "C" fn window_get_outer_position(window_id: CInt) -> CPosition {
	let mut result = CPosition{ x: 0.0, y: 0.0 };

	find_local_window!(window_id, |it: &Window| {
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

	find_local_window!(window_id, |it: &Window| {
		let size = it.webview.window().outer_size();
		result.width = size.width as f64;
		result.height = size.height as f64;
	});

	result
}

#[no_mangle]
pub extern "C" fn window_get_inner_position(window_id: CInt) -> CPosition {
	let mut result = CPosition{ x: 0.0, y: 0.0 };

	find_local_window!(window_id, |it: &Window| {
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

	find_local_window!(window_id, |it: &Window| {
		let size = it.webview.window().inner_size();
		result.width = size.width as f64;
		result.height = size.height as f64;
	});

	result
}

#[no_mangle]
pub extern "C" fn window_get_dpi_scale(window_id: CInt) -> CDouble {
	let mut result = 1.0;
	find_local_window!(window_id, |it: &Window| result = it.webview.window().scale_factor());
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
		};

		match event {
			Event::WindowEvent { event, window_id, .. } => {
				// @Incomplete: when a window is being destroyed we still want to get its user_window_id
				// Right now, it will be -1
				let user_window_id = GLOBAL_WINDOWS.with(|windows| {
					let array = windows.borrow();
					let it = array.iter().find(|&it| it.webview.window().id() == window_id);

					if let Some(it) = it {
						it.id
					} else {
						-1
					}
				});

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

						GLOBAL_WINDOWS.with(|windows| {
							let array = windows.borrow();
							let it = array.iter().find(|&it| it.webview.window().id() == window_id);

							if let Some(it) = it {
								let size = it.webview.inner_size();
								result.size = CSize{width: size.width as f64, height: size.height as f64};

								let _ = it.webview.resize();
							}
						});
					},
					_ => {}
				};
			},
			_ => (),
		}

		unsafe {
			user_callback(result);
		}
	});
}
