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

	pub fn push(&mut self, size: usize) -> *mut u8 {
		assert!(size + self.offset < self.size);

		let prev_offset = self.offset;
		self.offset += size;
		self.at(prev_offset)
	}

	pub fn push_aligned(&mut self, size: usize, alignment: usize) -> *mut u8 {
		assert!(alignment >= 1);
		// NOTE(nick): pow2
		assert!((alignment & !(alignment - 1)) == alignment);

		let base_address: usize = self.data as usize + self.offset;
		let mut align_offset: usize = alignment - (base_address & (alignment - 1));
		align_offset &= alignment - 1;

		let size = size + align_offset;

		assert!(self.offset + size < self.size);

		let result = self.offset + align_offset;
		self.offset += size;
		self.at(result)
	}

	pub fn reset(&mut self) {
		self.offset = 0;
	}

	pub fn write_aligned(&mut self, ptr: *mut u8, size: usize, alignment: usize) -> *mut u8 {
		let result = self.push_aligned(size, alignment);
		Arena::copy(ptr, result, size);
		result
	}

	pub fn write(&mut self, ptr: *mut u8, size: usize) -> *mut u8 {
		let result = self.push(size);
		Arena::copy(ptr, result, size);
		result
	}

	pub fn write_str(&mut self, str: &str) -> *mut u8 {
		self.write(str.as_ptr() as *mut u8, str.len())
	}

	pub fn copy(from: *mut u8, to: *mut u8, size: usize) {
		unsafe {
			libc::memcpy(to as *mut libc::c_void, from as *mut libc::c_void, size as usize);
		}
	}

	pub fn set_alignment(&mut self, alignment: usize) {
		self.push_aligned(0, alignment);
	}

	pub fn write_ptr(&mut self) -> *mut u8 {
		unsafe { self.data.offset(self.offset as isize) }
	}

	pub fn at(&mut self, offset: usize) -> *mut u8 {
		assert!(offset < self.size);

		unsafe { self.data.offset(offset as isize) }
	}
}