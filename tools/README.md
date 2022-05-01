## Extract HID usage tables from specification PDF

https://usb.org/sites/default/files/hut1_3_0.pdf


keyboard page

./extract-hid-tables.py ~/Downloads/hut1_3_0.pdf 89-95 keyboard-hut.csv

consumer page

./extract-hid-tables.py ~/Downloads/hut1_3_0.pdf 124-135 consumer-hut.csv


Manually edited huts

cleaned-keyboard-hut.csv

cleaned-consumer-hut.csv

