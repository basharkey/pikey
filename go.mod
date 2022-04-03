module github.com/basharkey/pimk

go 1.15

require (
	config v1.0.0 // indirect
	gadget v1.0.0 // indirect
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	keymap v1.0.0 // indirect
)

replace gadget v1.0.0 => ./gadget

replace keymap v1.0.0 => ./keymap

replace config v1.0.0 => ./config
