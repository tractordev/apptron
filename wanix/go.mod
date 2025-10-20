module apptron.dev/wanix

go 1.25.0

replace golang.org/x/sys => github.com/progrium/sys-wasm v0.0.0-20240620081741-5ccc4fc17421

replace github.com/hugelgupf/p9 => github.com/progrium/p9 v0.0.0-20251014191851-920015933007

// replace tractor.dev/wanix => ../../wanix

require tractor.dev/wanix v0.0.0-20251020193138-e75c1328ccc9

require (
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/hugelgupf/p9 v0.3.1-0.20240118043522-6f4f11e5296e // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	tractor.dev/toolkit-go v0.0.0-20250103001615-9a6753936c19 // indirect
)
