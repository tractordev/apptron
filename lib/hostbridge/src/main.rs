use hostbridge::*;
use std::mem::{size_of};

fn main() {
  let result = hostbridge::screen_get_available_displays();

  for i in 0..result.count {
    let it: *mut CDisplay = unsafe { std::mem::transmute(result.data as usize + size_of::<CDisplay>() * i as usize) };

    println!("Display {}", i);
    unsafe {
      println!("  name {}", std::ffi::CStr::from_ptr((*it).name).to_str().unwrap());
      println!("  size {}, {}", (*it).size.width, (*it).size.height);
      println!("  position {}, {}", (*it).position.x, (*it).position.y);
      println!("  scale_factor {}", (*it).scale_factor);
    }
  }

  println!("GOODBYE.");
}