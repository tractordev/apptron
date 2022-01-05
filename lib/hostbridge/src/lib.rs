use std::ffi::CStr;

#[no_mangle]
pub extern "C" fn hello(name: *const libc::c_char) {
    let buf_name = unsafe { CStr::from_ptr(name).to_bytes() };
    let str_name = String::from_utf8(buf_name.to_vec()).unwrap();
    println!("Hello {}!", str_name);
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