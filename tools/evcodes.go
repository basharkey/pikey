package main

import (
    "fmt"
    "encoding/csv"
    "os"
    "log"
    "os/exec"
    "github.com/gvalkov/golang-evdev"
)

func main() {
	hid_csv := os.Args[1]
	host := os.Args[2]
	keyboard_path := os.Args[3]
	report_byte := os.Args[4]

    keyboard_device, err := evdev.Open(keyboard_path)
	check_err(err)
	keyboard_device.Grab()

	f, err := os.Open(hid_csv)
	check_err(err)

	csv_reader := csv.NewReader(f)
	hid_table, err := csv_reader.ReadAll()
	check_err(err)

	for _, hid_entry := range hid_table {
		usage_id := hid_entry[0]
		usage_name := hid_entry[1]

		cmd := exec.Command("ssh", host, "sudo", "type-byte", report_byte, usage_id, "0")
		err := cmd.Run()
		check_err(err)

	loop:
		for {
			key_events, err := keyboard_device.Read()
			check_err(err)
			for _, event := range key_events {
				key_type := event.Type
				key_code := event.Code
				key_state := event.Value

				if key_type != 0 && key_type != 4 && key_type != 17 {
					if key_state == 1 {
						fmt.Printf("%d: {%s, \"%s\"},\n", key_code, usage_id, usage_name)
						break loop
					}
				}
			}
		}

		cmd = exec.Command("ssh", host, "sudo", "type-byte", report_byte, "0", "0")
		err = cmd.Run()
		check_err(err)
	}
}

func check_err(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
