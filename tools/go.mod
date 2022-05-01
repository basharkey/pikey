module evcodes

go 1.15

require (
	config v1.0.0
	gadget v1.0.0
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	keymap v1.0.0
)

replace gadget v1.0.0 => ../gadget

replace keymap v1.0.0 => ../keymap

replace config v1.0.0 => ../config
