use std::ffi::CStr;

#[no_mangle]
pub extern "C" fn hello(name: *const libc::c_char) {
    let buf_name = unsafe { CStr::from_ptr(name).to_bytes() };
    let str_name = String::from_utf8(buf_name.to_vec()).unwrap();
    println!("Hello {}!", str_name);
}

/*
#[no_mangle]
pub extern "C" fn window_create(width: i32, height: i32, title: *const libc::c_char) -> i32 {
  use wry::{
      application::{
          event::{Event, StartCause, WindowEvent},
          event_loop::{ControlFlow, EventLoop},
          window::WindowBuilder,
          platform::run_return::EventLoopExtRunReturn,
      },
      webview::WebViewBuilder,
  };

  let title_buf = unsafe { CStr::from_ptr(title).to_bytes() };
  let title_str = String::from_utf8(title_buf.to_vec()).unwrap();

  let mut event_loop = EventLoop::new();
  let maybe_window = WindowBuilder::new()
      .with_title(title_str)
      .with_decorations(true)
      //.with_transparent(false)
      //.with_inner_size(width, height)
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

  return 1
}
*/

#[no_mangle]
pub extern "C" fn run(user_callback: unsafe extern "C" fn(i32)) -> i32 {
  use wry::{
      application::{
          event::{Event, WindowEvent},
          event_loop::{ControlFlow, EventLoop},
          window::WindowBuilder,
      },
      webview::WebViewBuilder,
  };

  let event_loop = EventLoop::new();
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
      *control_flow = ControlFlow::Wait; // Or Poll

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

  return 1
}

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