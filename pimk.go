package main

import (
    "keymap"
    "gadget"
    "config"
    "fmt"
    "os"
    "log"
    "github.com/gvalkov/golang-evdev"
    "strings"
    "time"
    "path/filepath"
    "errors"
)

type Keystate struct {
    Code uint16
    State bool
    Keybind bool
    Layerbind bool
}

/*
var definitions

keybind - combination of keys that when pressed type other keys or activate a function
keybinds - a group of keybinds

bind_input_key - an individual key used to activate keybind
bind_input_keys - key(s) required to activate keybind

bind_output_key - an individual key activated when keybind is pressed
bind_output_keys - key(s) or function activated when keybind is pressed

layerbind - keybind used to change layers
*/

/*
var keybinds = [][]config.Keybind {
    {
        // a, b -> 1, 2
        {[]uint16{30, 48}, []uint16{2, 3}},
        // ctrl, h -> 3, 4
        {[]uint16{1029, 35}, []uint16{6, 7}},
        // u -> k
        {[]uint16{22}, []uint16{37}},
        // rctrl, lshift, f -> comma
        {[]uint16{97, 42, 33}, []uint16{51}},
    },
    {
        {[]uint16{30, 48}, []uint16{18, 33}},
        {[]uint16{30, 48}, []uint16{18, 33}},
    },
}
*/

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

    var hooked_keyboards []string
    for {
        // I don't think this will ever error for no keyboards being plugged in
        // errors would probably be related to permissions issues
        keyboard_paths, err := get_keyboard_paths()
        if err != nil {
            fmt.Println(err)
        } else {
            next:
            for _, keyboard_path := range keyboard_paths {
                // don't hook keyboards that are already hooked
                for _, hooked_keyboard := range hooked_keyboards {
                    if keyboard_path == hooked_keyboard {
                        continue next
                    }
                }

                keyboard_device, err := evdev.Open(keyboard_path)
                if err != nil {
                    fmt.Println(err)
                } else {
                    // load default.conf config if keyboard_device specific one does not exist
                    var keyboard_config string
                    config_dir := "/etc/pimk"
                    default_keyboard_config := filepath.Join(config_dir, "default.conf")
                    custom_keyboard_config := filepath.Join(config_dir, filepath.Base(keyboard_path) + ".conf")

                    _, err = os.Stat(custom_keyboard_config)
                    if errors.Is(err, os.ErrNotExist) {
                        fmt.Println("Using default config can't find: ", custom_keyboard_config)
                        keyboard_config = default_keyboard_config
                    } else {
                        keyboard_config = custom_keyboard_config
                    }

                    // track keyboards that are currently connected and hooked in hooked_keyboards slice
                    hooked_keyboards = append(hooked_keyboards, keyboard_path)
                    // hook the keyboard
                    go hook_keyboard(keyboard_device, keyboard_config, gadget_device, keyboard_path, &hooked_keyboards)
                }
            }
        }
        time.Sleep(1 * time.Second)
    }
}

func hook_keyboard(keyboard_device *evdev.InputDevice, keyboard_config string, gadget_device *os.File, keyboard_path string, hooked_keyboards *[]string) error {
    var layerbinds [][]config.Layerbind
    var keybinds [][]config.Keybind
    var rebinds []config.Rebind

    rebinds, layerbinds, keybinds = config.Parse(keyboard_config)

    fmt.Println(keyboard_device)
    fmt.Println(layerbinds)
    fmt.Println(keybinds)
    for key, value := range rebinds[0].Modifiers {
        fmt.Println(key, value)
    }

    var layer int = 0
    var pressed_keys []Keystate
    var pressed_layerbinds []config.Layerbind
    var pressed_keybinds []config.Keybind
    var index_bind_input_keys []uint16

    // main keyboard_device event loop
    for {
        // check if events can be read from keyboard_device (if keyboard is still connected)
        key_events, err := keyboard_device.Read()
        if err != nil {
            // remove keyboard from hooked keyboards if it has been disconnected
            for i, hooked_keyboard := range *hooked_keyboards {
                if keyboard_path == hooked_keyboard {
                    (*hooked_keyboards)[i] = (*hooked_keyboards)[len(*hooked_keyboards)-1]
                    *hooked_keyboards = (*hooked_keyboards)[:len(*hooked_keyboards)-1]
                }
            }
            return err
        }

        for _, event := range key_events {
            key_type := event.Type
            key_code := event.Code
            key_state := event.Value

            if key_type != 0 && key_type != 4 {
                if key_state == 1 {
                    // TODO dont add keys if 6 are currently pressed

                    /* Only add newly pressed key to pressed_keys if it doesn't already contain it
                    Prevents duplicate key entries in pressed_keys when a keybind output key is the same as a pressed key */
                    if contains_key(key_code, &pressed_keys) == false {
                        pressed_keys = append(pressed_keys, Keystate{key_code, true, false, false})
                    }

                    pressed_layerbinds, layer = detect_layerbinds(&pressed_keys, pressed_layerbinds, layerbinds[layer], layer)
                    pressed_keybinds = detect_keybinds(&pressed_keys, pressed_keybinds, keybinds[layer])

                    type_bytes(gadget_device, keys_to_bytes(&pressed_keys, rebinds[layer]))
                } else if key_state == 2 {
                    pressed_layerbinds, layer = detect_layerbinds(&pressed_keys, pressed_layerbinds, layerbinds[layer], layer)
                    pressed_keybinds = detect_keybinds(&pressed_keys, pressed_keybinds, keybinds[layer])

                    type_bytes(gadget_device, keys_to_bytes(&pressed_keys, rebinds[layer]))
                } else {
                    // remove released key from pressed_keys
                    for i, key := range pressed_keys {
                        if key.Code == key_code {
                            pressed_keys[i] = pressed_keys[len(pressed_keys)-1]
                            pressed_keys = pressed_keys[:len(pressed_keys)-1]
                        }
                    }

                    pressed_layerbinds, layer = remove_layerbinds(&pressed_keys, pressed_layerbinds, layerbinds[layer], layer)
                    pressed_keybinds, index_bind_input_keys = remove_keybinds(&pressed_keys, pressed_keybinds)
                    // if no keys pressed clear buffer
                    if len(pressed_keys) == 0 {
                        type_bytes(gadget_device, make([]byte, 8))
                    // else update with currently pressed keys
                    } else {
                        type_bytes(gadget_device, keys_to_bytes(&pressed_keys, rebinds[layer]))
                        for _, bind_input_key := range index_bind_input_keys {
                            for i, key := range pressed_keys {
                                if bind_input_key == key.Code {
                                    pressed_keys[i].State = true
                                }
                            }
                        }
                    }
                }
                fmt.Println("keys", pressed_keys)
                fmt.Println("keybinds", pressed_keybinds)
                fmt.Println("layerbinds", pressed_layerbinds)
                fmt.Println("-------------------------------")
            }
        }
    }
}

func contains_key(key_code uint16, pressed_keys *[]Keystate) bool {
    for _, key := range *pressed_keys {
        if key.Code == key_code {
            return true
        }
    }
    return false
}

func get_keyboard_paths() ([]string, error)  {
    dev_path := "/dev/input/by-id/"
    dir, err := os.Open(dev_path)
    if err != nil {
        return nil, err
    }
    devices, err := dir.Readdir(0)
    if err != nil {
        return nil, err
    }

    var keyboards []string
    for i := range devices {
        device := devices[i].Name()
        if strings.Contains(device, "event-kbd") && !strings.Contains(device, "if01") {
            keyboards = append(keyboards, dev_path + device)
        }
    }
    return keyboards, nil
}

func keycode_equals_bindkey(keycode uint16, bind_input_key uint16) bool {
    if multicode_mod, ok := keymap.Multicode_modifiers[bind_input_key]; ok {
        // check that keycode is equal to either the left or right Modifier keycodes
        if keycode == multicode_mod.Leftkey || keycode == multicode_mod.Rightkey {
            return true
        }
    } else if keycode == bind_input_key {
        return true
    }
    return false
}

func detect_bind(pressed_keys []Keystate, bind_input_keys []uint16, is_layerbind bool) bool {
    num_pressed_bind_input_keys := 0
    for _, key := range pressed_keys {
        for _, bind_input_key := range bind_input_keys {
            if (key.State || is_layerbind) && !key.Keybind && keycode_equals_bindkey(key.Code, bind_input_key) {
                num_pressed_bind_input_keys += 1
            }
        }
    }
    // if all bind_input_keys are pressed return true
    if num_pressed_bind_input_keys == len(bind_input_keys) {
        return true
    }
    return false
}

// compare toggle binds only
func check_opposite_toggle_layerbinds(layerbind1 config.Layerbind, layerbind2 config.Layerbind) bool {
    if layerbind1.Type != layerbind2.Type {
        return false
    }
    if layerbind1.To_layer != layerbind2.From_layer {
        return false
    }
    if layerbind1.From_layer != layerbind2.To_layer {
        return false
    }
    if len(layerbind1.Input_keys) != len(layerbind2.Input_keys) {
        return false
    }
    for i := range layerbind1.Input_keys {
        if layerbind1.Input_keys[i] != layerbind2.Input_keys[i] {
            return false
        }
    }
    return true
}

func detect_layerbinds(pressed_keys *[]Keystate, pressed_layerbinds []config.Layerbind, layer_layerbinds []config.Layerbind, layer int) ([]config.Layerbind, int) {
    // detect newly pressed layerbinds
    for _, layerbind := range layer_layerbinds {
        if detect_bind(*pressed_keys, layerbind.Input_keys, true) {
            for _, bind_input_key := range layerbind.Input_keys {
                for i, key := range *pressed_keys {
                    if keycode_equals_bindkey(key.Code, bind_input_key) {
                        // tag layerbind input keys
                        (*pressed_keys)[i].Layerbind = true
                        // if layerbind suppression enabled set the input keys pressed state to false
                        if layerbind.Suppress {
                            (*pressed_keys)[i].State = false
                        }
                    }
                }
            }
            if layerbind.Type == "tap" {
                layer = layerbind.To_layer
            } else if layerbind.Type == "momentary"{
                layer = layerbind.To_layer
                pressed_layerbinds = append(pressed_layerbinds, layerbind)
            } else if layerbind.Type == "toggle"{
                for i, pressed_layerbind := range pressed_layerbinds {
                    if check_opposite_toggle_layerbinds(layerbind, pressed_layerbind) {
                        if pressed_layerbind.State == 2 {
                            (pressed_layerbinds)[i].State = 3
                        }
                        return pressed_layerbinds, layer
                    }
                }
                layer = layerbind.To_layer
                layerbind.State = 1
                pressed_layerbinds = append(pressed_layerbinds, layerbind)
            }
        }
    }
    return pressed_layerbinds, layer
}

func remove_layerbinds(pressed_keys *[]Keystate, pressed_layerbinds []config.Layerbind, layer_layerbinds []config.Layerbind, layer int) ([]config.Layerbind, int) {
    for {
        check_again := false
        for i, layerbind := range pressed_layerbinds {
            if layerbind.Type == "momentary" && !detect_bind(*pressed_keys, layerbind.Input_keys, true) && layer == layerbind.To_layer {
                pressed_layerbinds[i] = pressed_layerbinds[len(pressed_layerbinds)-1]
                pressed_layerbinds = pressed_layerbinds[:len(pressed_layerbinds)-1]
                layer = layerbind.From_layer
                check_again = true
            }
            if layerbind.Type == "toggle" && !detect_bind(*pressed_keys, layerbind.Input_keys, true) {
                if layerbind.State == 1 {
                    pressed_layerbinds[i].State = 2
                } else if layerbind.State == 3 {
                    pressed_layerbinds[i] = pressed_layerbinds[len(pressed_layerbinds)-1]
                    pressed_layerbinds = pressed_layerbinds[:len(pressed_layerbinds)-1]
                    layer = layerbind.From_layer
                    check_again = true
                }
            }
        }
        if !check_again {
            break
        }
    }
    return pressed_layerbinds, layer
}

func detect_keybinds(pressed_keys *[]Keystate, pressed_keybinds []config.Keybind, layer_keybinds []config.Keybind) []config.Keybind {
    for _, keybind := range layer_keybinds {
        if detect_bind(*pressed_keys, keybind.Input_keys, false) {
            /*
            if all bind_input_keys for a keybind have been pressed:
                add bind_output_keys to pressed keys
                set bind_input_keys state to false
            */
            pressed_keybinds = append(pressed_keybinds, keybind)
            for _, bind_input_key := range keybind.Input_keys {
                for i, key := range *pressed_keys{
                    if keycode_equals_bindkey(key.Code, bind_input_key) {
                        (*pressed_keys)[i].State = false
                    }
                }
            }

            for _, bind_output_key := range keybind.Output_keys {
                found := false
                for i, key := range *pressed_keys {
                    // don't true layerbind input keys
                    if !(*pressed_keys)[i].Layerbind && key.Code == bind_output_key {
                        (*pressed_keys)[i].State = true
                        (*pressed_keys)[i].Keybind = true
                        found = true
                    }
                }
                if !found {
                    *pressed_keys = append(*pressed_keys, Keystate{bind_output_key, true, true, false})
                }
            }
        }
    }
    //fmt.Println(pressed_keybinds)
    return pressed_keybinds
}

func remove_keybinds(pressed_keys *[]Keystate, pressed_keybinds []config.Keybind) ([]config.Keybind, []uint16) {
    var index_bind_input_keys []uint16

    for i, keybind := range pressed_keybinds {
        num_bind_input_keys := len(keybind.Input_keys)
        num_pressed_bind_input_keys := 0
        for _, bind_input_key := range keybind.Input_keys {
            for _, key := range *pressed_keys {
                if keycode_equals_bindkey(key.Code, bind_input_key) {
                    num_pressed_bind_input_keys += 1
                }
            }
        }

        if num_bind_input_keys != num_pressed_bind_input_keys {
            pressed_keybinds[i] = pressed_keybinds[len(pressed_keybinds)-1]
            pressed_keybinds = pressed_keybinds[:len(pressed_keybinds)-1]

            for _, bind_output_key := range keybind.Output_keys {
                for i, key := range *pressed_keys {
                    if bind_output_key == key.Code {
                        (*pressed_keys)[i] = (*pressed_keys)[len(*pressed_keys)-1]
                        *pressed_keys = (*pressed_keys)[:len(*pressed_keys)-1]
                        break
                    }
                }
            }

            for _, bind_input_key := range keybind.Input_keys {
                for _, key := range *pressed_keys {
                    if keycode_equals_bindkey(key.Code, bind_input_key) {
                        index_bind_input_keys = append(index_bind_input_keys, key.Code)
                    }
                }
            }
        }
    }
    return pressed_keybinds, index_bind_input_keys
}

func type_bytes(gadget_device *os.File, key_bytes []byte) {
    fmt.Println("typing:", key_bytes)
    //key_bytes = make([]byte, 8)
    _, err := gadget_device.Write(key_bytes)
    check_err(err)
}

// efficient way of prepending to slice
// https://stackoverflow.com/questions/53737435/how-to-prepend-int-to-slice
func prepend_byte(x []byte, y byte) []byte {
    x = append(x, 0)
    copy(x[1:], x)
    x[0] = y
    return x
}

func keys_to_bytes(pressed_keys *[]Keystate, rebinds config.Rebind) []byte {
    /*
    key_bytes slice
    [1, 0, 42, 35, 78, 0, 0, 0]
    byte 1 = modifier byte (bitwise OR of each modifier bit)
    byte 2 = reserved byte
    bytes 3-8 = key bytes
    */

    var key_bytes []byte
    var pressed_mods []byte

    for _, key := range *pressed_keys {
        // remove keys with state of false
        if key.State {
            // if key is apart of a keybind use default keymap
            if key.Keybind {
                // check if key is key and not a modifier
                if _, ok := keymap.Keys[key.Code]; ok {
                    key_bytes = append(key_bytes, keymap.Keys[key.Code].Scancode)
                } else {
                    // add modifiers in pressed keys to pressed_mods
                    for mod, _ := range keymap.Modifiers {
                        if key.Code == mod {
                            pressed_mods = append(pressed_mods, keymap.Modifiers[mod].Scancode)
                        }
                    }
                }
            // if key is not apart of a keybind use rebinds keymap
            } else {
                // check if key is key and not a modifier
                if _, ok := rebinds.Keys[key.Code]; ok {
                    key_bytes = append(key_bytes, rebinds.Keys[key.Code].Scancode)
                } else {
                    // add modifiers in pressed keys to pressed_mods
                    for mod, _ := range rebinds.Modifiers {
                        if key.Code == mod {
                            pressed_mods = append(pressed_mods, rebinds.Modifiers[mod].Scancode)
                        }
                    }
                }
            }
        }
    }

    // prepend reserved byte
    key_bytes = prepend_byte(key_bytes, 0)

    // generate modifier byte
    if len(pressed_mods) > 0 {
        // bitwise OR each modifier bit to create modifier byte
        var mod_byte byte
        for i, _ := range pressed_mods {
            mod_byte = mod_byte | pressed_mods[i]
        }
        // prepend modifier byte
        key_bytes = prepend_byte(key_bytes, mod_byte)
    } else {
        // prepend empty modifier byte
        key_bytes = prepend_byte(key_bytes, 0)
    }

    // pad remaining space with null bytes
    key_bytes = append(key_bytes, make([]byte, 8-len(key_bytes))...)
    return key_bytes
}

func check_err(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
