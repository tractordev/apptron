#[macro_export]
macro_rules! kilobytes { ($x: expr) => { $x * 1024 } }

pub struct Arena {
  pub data: *mut u8,
  pub offset: usize,
  pub size: usize,
}

impl Arena {
  pub const fn new(data: *mut u8, size: usize) -> Arena {
    Arena { data, size, offset: 0 }
  }

  pub fn push(&mut self, count: usize) -> *mut u8 {
    assert!(count + self.offset < self.size);

    let prev_offset = self.offset;
    self.offset += count;
    self.at(prev_offset)
  }

  pub fn reset(&mut self) {
    self.offset = 0;
  }

  pub fn write_raw(&mut self, ptr: *mut u8, count: usize) -> *mut u8 {
    let result = self.push(count);
    Arena::copy(ptr, result, count);
    result
  }

  pub fn write(&mut self, str: &str) -> *mut u8 {
    self.write_raw(str.as_ptr() as *mut u8, str.len())
  }

  pub fn copy(from: *mut u8, to: *mut u8, size: usize) {
    unsafe {
      libc::memcpy(to as *mut libc::c_void, from as *mut libc::c_void, size as usize);
    }
  }

  pub fn at(&mut self, offset: usize) -> *mut u8 {
    assert!(offset < self.size);

    unsafe { self.data.offset(offset as isize) }
  }
}