use std::ffi::CStr;
use std::mem::ManuallyDrop;
use std::mem::forget;
use std::option::Option;
use std::cell::RefCell;

use wry::{
    application::{
        event::{Event, WindowEvent},
        event_loop::{ControlFlow, EventLoop},
        window::{WindowBuilder},
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
pub struct CVector2 {
  pub x: f64,
  pub y: f64,
}

// NOTE(nick): even though this stuct is not FFI compatible, we use it as an opaque handle on the C/Go side
// so the layout of the data shouldn't matter
type CEventLoop = EventLoop<()>;

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
      let found = find_window_by_id(&array, $window_id);

      if let Some(it) = found {
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
pub extern "C" fn create_window(event_loop: CEventLoop) -> i32 {
  let maybe_window = WindowBuilder::new()
      .with_title("")
      .with_decorations(true)
      .build(&event_loop);

  forget(event_loop);

  if !maybe_window.is_ok() {
    return 0
  }

  let window = maybe_window.unwrap();

  let maybe_webview_builder = WebViewBuilder::new(window);

  if !maybe_webview_builder.is_ok() { return 0; }

  let maybe_webview = maybe_webview_builder.unwrap()
      .with_url(
          r#"data:text/html,
          <!doctype html>
          <html>
            <body style="font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Ubuntu, roboto, noto, arial, sans-serif; background-color:rgba(87,87,87,0.75);"></body>
            <script>
              window.onload = function() {
                document.body.innerHTML = `<div style="padding: 30px">Transparency Test<br><br>${navigator.userAgent}</div>`;
              };
            </script>
          </html>"#,
      );

  if !maybe_webview.is_ok() { return 0; }

  let result = maybe_webview.unwrap().build();

  if !result.is_ok() { return 0; }

  let webview = result.unwrap();

  let mut result: i32 = -1;

  GLOBAL_WINDOWS.with(|windows| {
    let mut array = windows.borrow_mut();

    result = array.len() as i32;
    let the_window = Window{ id: result, webview };
    array.push(the_window);

    //println!("[rust] push new window id {:?}, new windows length: {:?}", result, array.len());
  });

  return result;
}

#[no_mangle]
pub extern "C" fn destroy_window(window_id: CInt) -> CBool {
  false
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
pub extern "C" fn window_get_outer_position(window_id: CInt) -> CVector2 {
  let mut result = CVector2{ x: 0.0, y: 0.0 };

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
pub extern "C" fn window_get_outer_size(window_id: CInt) -> CVector2 {
  let mut result = CVector2{ x: 0.0, y: 0.0 };

  find_local_window!(window_id, |it: &Window| {
    let size = it.webview.window().outer_size();
    result.x = size.width as f64;
    result.y = size.height as f64;
  });

  result
}

#[no_mangle]
pub extern "C" fn window_get_inner_position(window_id: CInt) -> CVector2 {
  let mut result = CVector2{ x: 0.0, y: 0.0 };

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
pub extern "C" fn window_get_inner_size(window_id: CInt) -> CVector2 {
  let mut result = CVector2{ x: 0.0, y: 0.0 };

  find_local_window!(window_id, |it: &Window| {
    let size = it.webview.window().inner_size();
    result.x = size.width as f64;
    result.y = size.height as f64;
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
pub extern "C" fn run(event_loop: CEventLoop, user_callback: unsafe extern "C" fn(i32)) {
  event_loop.run(move |event, _, control_flow| {
    *control_flow = ControlFlow::Poll;

    //println!("{:?}", event);

    let event_type = match event {
      Event::WindowEvent { event: WindowEvent::Resized{ .. }, .. } => 1,
      Event::WindowEvent { event: WindowEvent::Moved{ .. }, .. } => 2,
      Event::WindowEvent { event: WindowEvent::Focused{ .. }, .. } => 3,
      Event::WindowEvent { event: WindowEvent::MouseInput{ .. }, .. } => 4,
      Event::WindowEvent { event: WindowEvent::KeyboardInput{ .. }, .. } => 5,
      _ => 0,
    };

    unsafe {
      user_callback(event_type);
    }
  });
}
