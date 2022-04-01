package main

import (
    "github.com/basharkey/pikey/gadget"
    "os"
    "log"
)

func main() {
    // initialize usb gadget
    gadget.Initialize()

    // open usb gadget device for write only
    gadget_device, err := os.OpenFile(
        "/dev/hidg0",
        os.O_WRONLY,
        0000,
    )
    check_err(err)
    defer gadget_device.Close()
    type_bytes(gadget_device, make([]byte, 8))
}
func type_bytes(gadget_device *os.File, key_bytes []byte) {
    //fmt.Println("typing:", key_bytes)
    //key_bytes = make([]byte, 8)
    _, err := gadget_device.Write(key_bytes)
    check_err(err)
}

func check_err(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
