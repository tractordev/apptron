use std::ffi::CStr;
use std::mem::size_of;

use wry::{
    application::{
        event::{Event, WindowEvent},
        event_loop::{ControlFlow, EventLoop, EventLoopProxy},
        window::WindowBuilder,
    },
    webview::WebViewBuilder,
};

unsafe fn any_as_u8_slice<T: Sized>(p: &T) -> &[u8] {
    ::std::slice::from_raw_parts(
        (p as *const T) as *const u8,
        ::std::mem::size_of::<T>(),
    )
}

#[no_mangle]
pub extern "C" fn create_event_loop() -> EventLoop<()> {
  //let raw_mt_one_ptr = Box::into_raw(Box::new(mt_one)) as *const c_void;
  let result = Box::new(EventLoop::new());

  let bytes: &[u8] = unsafe { any_as_u8_slice(&result) };
  println!("[rust] create_event_loop {:?}", bytes);

  std::mem::forget(result);

  result
}

#[no_mangle]
//pub extern "C" fn create_window(data: *const libc::c_void) -> i32 {
pub extern "C" fn create_window(event_loop: EventLoop<()>) -> i32 {

  println!("[rust] event_loop {:?}", event_loop);
  
  let bytes: &[u8] = unsafe { any_as_u8_slice(&event_loop) };
  println!("[rust] bytes {:?}", bytes);

  return 42;
}

/*
pub extern "C" fn create_window(event_loop: EventLoop<()>) -> i32 {
  println!("{:?}", event_loop);
  
  return 42;
}
*/

fn string_from_cstr(cstr: *const libc::c_char) -> String {
  let buffer = unsafe { CStr::from_ptr(cstr).to_bytes() };
  String::from_utf8(buffer.to_vec()).unwrap()
}

#[no_mangle]
pub extern "C" fn run(event_loop: EventLoop<()>, user_callback: unsafe extern "C" fn(i32)) -> i32 {
  println!("{}", size_of::<EventLoop<()>>());

  //GLOBAL_WINDOWS.with(|windows| windows.push(42));

  /*
  let event_loop21 = EventLoop::new();

  let event_loop = EventLoop::new();
  */

  println!("[rust] window_create");

  let maybe_window = WindowBuilder::new()
      .with_title("Progrium Test")
      .with_decorations(true)
      //.with_transparent(false)
      .build(&event_loop);

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

  // @Incomplete: return window id
  return 1
}

/*
#[no_mangle]
pub extern "C" fn gomain() {
    oldmain().ok();
}

fn oldmain() -> wry::Result<()> {
  use wry::{
      application::{
          event::{Event, StartCause, WindowEvent},
          event_loop::{ControlFlow, EventLoop},
          window::WindowBuilder,
      },
      webview::WebViewBuilder,
  };

  let event_loop = EventLoop::new();
  let window = WindowBuilder::new()
      .with_title("Progrium Test")
      .with_decorations(false)
      .with_transparent(true)
      .build(&event_loop)?;

  let _webview = WebViewBuilder::new(window)?
      .with_transparent(true)
      //.with_url("https://progrium.com")?
      .with_url(
          r#"data:text/html,
          <!doctype html>
          <html>
            <body style="font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Ubuntu, roboto, noto, arial, sans-serif; background-color:rgba(87,87,87,0.5);"></body>
            <script>
              window.onload = function() {
                document.body.innerHTML = `<div style="padding: 30px">Transparency Test<br><br>${navigator.userAgent}</div>`;
              };
            </script>
          </html>"#,
      )?
      .build()?;

  event_loop.run(move |event, _, control_flow| {
      *control_flow = ControlFlow::Wait;

      match event {
          Event::NewEvents(StartCause::Init) => println!("Started"),
          Event::WindowEvent {
              event: WindowEvent::CloseRequested,
              ..
          } => *control_flow = ControlFlow::Exit,
          _ => (),
      }
  });
}
*/