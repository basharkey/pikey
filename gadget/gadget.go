package gadget

import (
    "os"
    //"log"
    "fmt"
    "path/filepath"
)

func Initialize() {
    var base_path string = "/sys/kernel/config/usb_gadget/pikey"

    var files = [][]string {
        {"idVendor", "0x1d6b"}, // Linux Foundation
        {"idProduct", "0x0104"}, // Multifunction Composite Gadget
        {"bcdDevice", "0x0100"}, // v1.0.0
        {"bcdUSB", "0x0200"}, // USB2

        {"strings/0x409/serialnumber", "fedcba9876543210"},
        {"strings/0x409/manufacturer", "Pikey"},
        {"strings/0x409/product", "Pikey USB Keyboard Device"},

        {"configs/c.1/strings/0x409/configuration", "Config 1: ECM network"},
        {"configs/c.1/MaxPower", "250"},

        {"functions/hid.usb0/protocol", "1"},
        {"functions/hid.usb0/subclass", "1"},
        {"functions/hid.usb0/report_length", "8"},
        {"functions/hid.usb0/report_desc", "\x05\x01\x09\x06\xa1\x01\x05\x07\x19\xe0\x29\xe7\x15\x00\x25\x01\x75\x01\x95\x08\x81\x02\x95\x01\x75\x08\x81\x03\x95\x05\x75\x01\x05\x08\x19\x01\x29\x05\x91\x02\x95\x01\x75\x03\x91\x03\x95\x06\x75\x08\x15\x00\x25\x65\x05\x07\x19\x00\x29\x65\x81\x00\xc0"},
    }

    for _, file := range files {
        write_to_file(filepath.Join(base_path, file[0]), file[1])
    }

    link := filepath.Join(base_path, "functions/hid.usb0")
    target := filepath.Join(base_path, "configs/c.1/hid.usb0")
    os.Symlink(link, target)

    dir, err := os.Open("/sys/class/udc/")
    check_err(err)
    file, err := dir.Readdir(0)
    check_err(err)
    write_to_file(filepath.Join(base_path, "UDC"), file[0].Name())
}

func check_err(err error) {
    if err != nil {
        //log.Fatal(err)
        fmt.Println(err)
    }
}

func write_to_file(file_path string, file_content string) {
    file_dir := filepath.Dir(file_path)
    _, err := os.Stat(file_dir)
    if os.IsNotExist(err) {
        err := os.MkdirAll(file_dir, 0755)
        check_err(err)
    }

    file, err := os.Create(file_path)
    check_err(err)
    _, err = file.WriteString(file_content)
    check_err(err)
}
