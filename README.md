# PiMK

PiMK allows you to turn a Raspberry Pi 4 Model B into a keyboard remapping device. By plugging your keyboard(s) into the Pi, and plugging the Pi into your computer it can act as a keyboard device and translate keystrokes allowing you to remap keys.


If you plan to plug in multiple keyboards or your keyboard(s) require higher amounts of power I reccommend using either a deciated Pi power supply such as [Waveshare's UPS](https://www.waveshare.com/wiki/UPS_HAT_(B)) or plugging your Pi into a [powered USB hub](https://plugable.com/products/usbc-hub7bc).

## Completed Features

- [x] Key remapping
- [x] Custom hotkeys
- [x] Layer support
- [x] Multiple keyboards with separate config support
- [ ] Macros

## Security

While running the Pi has the ability to see all keystrokes entered. This can become a security concern when entering sensitive information such as passwords.


As a result I do NOT recommend having your Pi connected to your network as this could provide an entry point for attackers. If you must ensure you are hardening the OS and installing the latest security updates.

## Installation

```
git clone https://github.com/basharkey/pimk.git
sudo apt install -y ansible
cd pimk/
```

## Configuration

### Start/Stop Service

```
sudo systemctl start pimk
sudo systemctl stop pimk
```

## Config Files

### Default

```
/etc/pimk/default.conf
```

### Custom

```
/etc/pimk/usb-04d9_USB_Keyboard-event-kbd.conf
```


For more information on creating config files see the [Wiki](https://github.com/basharkey/pimk/wiki)
