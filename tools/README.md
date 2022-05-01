## Extract HID usage tables from specification PDF

https://usb.org/sites/default/files/hut1_3_0.pdf


keyboard page
```
./extract-hid-tables.py ~/Downloads/hut1_3_0.pdf 89-95 tables/keyboard-hut.csv
```
consumer page
```
./extract-hid-tables.py ~/Downloads/hut1_3_0.pdf 124-135 tables/consumer-hut.csv
```

Manually edited huts

cleaned-keyboard-hut.csv

cleaned-consumer-hut.csv

## Generate final keymap

evcodes, SSH's into the Pi and types a key then checks if the host detects the key press. This allows you to determine which Usage IDs actually register on Linux as well as correlate decimal Usage IDs to decimal evcodes

### Requirements

For Pi `type-byte` binary needs to be in path (/usr/local/bin/)

Ability for host to SSH to Pi without a password (e.g. using SSH keys)


where 192.168.1.0 is Pi's IP

keyboard
```
./evcodes tables/registered-keyboard-hut.csv 192.168.1.0 /dev/input/by-id/usb-Pimk_USB_Keyboard_Device_fedcba9876543210-event-kbd 1 > keyboard-keymap.txt
```

consumer
```
./evcodes tables/registered-consumer-hut.csv 192.168.1.0 /dev/input/by-id/usb-Pimk_USB_Keyboard_Device_fedcba9876543210-event-kbd 2 > consumer-keymap.txt
```

evcode of 240 means there is no mapping and should be removed from keymap file
