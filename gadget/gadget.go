package gadget

import (
    "os"
    //"log"
    "fmt"
    "path/filepath"
    "os/exec"
)

var base_dir = "/sys/kernel/config/usb_gadget/pimk"

var usb_string = "0x409"
var usb_config = "c.1"
var usb_device = "hid.usb0"

var strings_dir = filepath.Join("strings", usb_string)
var configs_dir = filepath.Join("configs", usb_config)
var functions_dir = filepath.Join("functions", usb_device)

func Initialize() {

    var files = [][]string {
        {"idVendor", "0xbeaf"}, // Linux Foundation
        {"idProduct", "0x0104"}, // Multifunction Composite Gadget
        {"bcdDevice", "0x0100"}, // v1.0.0
        {"bcdUSB", "0x0200"}, // USB2

        {filepath.Join(strings_dir, "serialnumber"), "fedcba9876543210"},
        {filepath.Join(strings_dir, "manufacturer"), "Pimk"},
        {filepath.Join(strings_dir, "product"), "USB Keyboard Device"},

        {filepath.Join(configs_dir, strings_dir, "configuration"), "Config 1: ECM network"},
        {filepath.Join(configs_dir, "MaxPower"), "250"},

        {filepath.Join(functions_dir, "protocol"), "1"},
        {filepath.Join(functions_dir, "subclass"), "1"},
        {filepath.Join(functions_dir, "report_length"), "8"},
        // sudo usbhid-dump -d beaf | tail -n +2 | xxd -r -p | hidrd-convert -o spec
        {filepath.Join(functions_dir, "report_desc"),
            "\x05\x01" + // Usage Page (Desktop)
            "\x09\x06" + // Usage (Keyboard)
            "\xa1\x01" + // Collection (Application)
            "\x85\x01" + // Report ID (1)
            // modifier report byte
            "\x75\x01" + // Report Size (1)
            "\x95\x08" + // Report Count (8)
            "\x05\x07" + // Usage Page (Keyboard)
            "\x19\xe0" + // Usage Minimum (KB Leftcontrol)
            "\x29\xe7" + // Usage Maximum (KB Right GUI)
            "\x15\x00" + // Logical Minimum (0)
            "\x25\x01" + // Logical Maximum (1)
            "\x81\x02" + // Input (Variable)
            // reserved byte
            "\x95\x01" + // Report Count (1)
            "\x75\x08" + // Report Size (8)
            "\x81\x03" + // Input (Constant, Variable)
            // LED report 5 bits
            "\x95\x05" + // Report Count (5)
            "\x75\x01" + // Report Size (1)
            "\x05\x08" + // Usage Page (LED)
            "\x19\x01" + // Usage Minimum (01)
            "\x29\x05" + // Usage Maximum (05)
            "\x91\x02" + // Output (Variable)
            // LED report padding 3 bits
            "\x95\x01" + // Report Count (1)
            "\x75\x03" + // Report Size (3)
            "\x91\x03" + // Output (Constant, Variable)
            // key report byte
            "\x95\x06" + // Report Count (6)
            "\x75\x08" + // Report Size (8)
            "\x15\x00" + // Logical Minimum (0)
            "\x25\x01" + // Logical Maximum (01)
            "\x05\x07" + // Usage Page (Keyboard)
            "\x19\x00" + // Usage Minimum (None)
            "\x29\xf7" + // Usage Maximum (f7h)
            "\x81\x00" + // Input
            "\xc0" + // End Collection

            "\x05\x0c" + // Usage Page (Consumer)
            "\x09\x01" + // Usage (Consumer Control)
            "\xa1\x01" + // Collection (Application)
            "\x85\x02" + // Report ID (2)
            "\x16\x01\x00" + // Logical Minimum (01) 
            "\x26\x14\x05" + // Logical Maximum (514)
            "\x1a\x01\x00" + // Usage Minimum (01)
            "\x2a\x14\x05" + // Usage Maximum (514)
            "\x95\x01" + // Report Count (1)           num of fields
            "\x75\x10" + // Report Size (16)             num of bits per field
            "\x81\x00" + // Input
            "\xc0"}, // End Collection

            /* probably not useful, does not provide many detectable functions
            "\x05\x01" + // Usage Page (Desktop)
            "\x09\x80" + // Usage (Sys Control)
            "\xa1\x01" + // Collection (Application)
            "\x85\x03" + // Report ID (3)
            "\x16\x01\x00" + // Logical Minimum (01)
            "\x26\xe2\x00" + // Logical Maximum (55)
            "\x1a\x01\x00" + // Usage Minimum (Sys Power Down)
            "\x2a\xe2\x00" + // Usage Maximum (Sys Dspl LCD Autoscale)
            "\x75\x10" + // Report Size (16)
            "\x95\x01" + // Report Count (1)
            "\x81\x00" + // Input
            "\xC0"}, // End Collection
            */
    }

    for _, file := range files {
        write_to_file(filepath.Join(base_dir, file[0]), file[1])
    }

    link := filepath.Join(base_dir, functions_dir)
    target := filepath.Join(base_dir, filepath.Join(configs_dir, usb_device))
    os.Symlink(link, target)

    dir, err := os.Open("/sys/class/udc/")
    check_err(err)
    file, err := dir.Readdir(0)
    check_err(err)
    write_to_file(filepath.Join(base_dir, "UDC"), file[0].Name())
}

func Destroy() {
    var base_dir string = "/sys/kernel/config/usb_gadget/pimk"
    // clear UDC file data, don't think there is a way to do this with pure go
    //https://askubuntu.com/questions/823380/cannot-delete-residual-system-files-period-even-after-changing-permissions-as-r
    cmd := exec.Command("/usr/bin/env", "bash", "-c", "echo '' > " + filepath.Join(base_dir, "UDC"))
    cmd.Run()

    var gadget_files =  []string {
        // remove strings from configs
        filepath.Join(base_dir, configs_dir, strings_dir),
        // remove functions from configs
        filepath.Join(base_dir, configs_dir, usb_device),
        // remove configs
        filepath.Join(base_dir, configs_dir),
        // remove functions
        filepath.Join(base_dir, functions_dir),
        // remove strings
        filepath.Join(base_dir, strings_dir),
        // remove gadget
        base_dir,
    }

    for _, file := range gadget_files {
        os.Remove(file)
    }
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
        err := os.MkdirAll(file_dir, 0644)
        check_err(err)
    }

    file, err := os.Create(file_path)
    check_err(err)
    _, err = file.WriteString(file_content)
    check_err(err)
}
