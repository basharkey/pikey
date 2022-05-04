package main

import (
    "fmt"
    "gadget"
    "os"
    "log"
    "strconv"
    "encoding/binary"
)

func main () {
    //gadget.Destroy()
    // initialize usb gadget
    gadget.Initialize()

    // open usb gadget device for write only
    gadget_device, err := os.OpenFile(
        "/dev/hidg0",
        os.O_WRONLY,
        0000,
    )
    defer gadget_device.Close()
    check_err(err)

    report_byte, _ := strconv.Atoi(os.Args[1])
    key_byte1, _ := strconv.ParseUint(os.Args[2], 0, 16)
    key_byte2, _ := strconv.Atoi(os.Args[3])

    if report_byte == 1 {
        type_bytes(gadget_device, []byte{byte(report_byte), 0, 0, byte(key_byte1), 0, 0, 0, 0, 0})
    } else if report_byte >= 2 {
        if key_byte1 > 255 {
            bs := make([]byte, report_byte)
            binary.LittleEndian.PutUint16(bs, uint16(key_byte1))
            type_bytes(gadget_device, []byte{byte(report_byte), byte(bs[0]), byte(bs[1])})
        } else {
            type_bytes(gadget_device, []byte{byte(report_byte), byte(key_byte1), byte(key_byte2)})
        }
    }
}

func type_bytes(gadget_device *os.File, key_bytes []byte) {
    fmt.Println("typing:", key_bytes)
    //key_bytes = make([]byte, 8)
    _, err := gadget_device.Write(key_bytes)
	check_err(err)
}

func check_err(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
