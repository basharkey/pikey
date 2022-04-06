package config

import (
    "keymap"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "strconv"
)

type Layerbind struct {
    Input_keys []uint16
    To_layer int
    From_layer int
    Suppress bool
    State int
    Type string
}

type Keybind struct {
    Input_keys []uint16
    Output_keys []uint16
}

type Rebind struct {
    Multicode_modifiers map[uint16]keymap.Multi_mod_key
    Modifiers map[uint16]keymap.Key
    Keys map[uint16]keymap.Key
}

type Layer struct {
    Rebinds [][]string
    Layerbinds [][][]string
    Keybinds [][][]string
}

func Parse(config_file string) ([]Rebind, [][]Layerbind, [][]Keybind) {
    config, err := ioutil.ReadFile(config_file)
    if err != nil {
        fmt.Println(err)
    }

    var layers []Layer
    json.Unmarshal([]byte(string(config)), &layers)

    var rebinds []Rebind
    var layerbinds [][]Layerbind
    var keybinds [][]Keybind
    for i, layer := range layers {
        layer_rebinds := new(Rebind)
        layer_rebinds.Multicode_modifiers = make(map[uint16]keymap.Multi_mod_key)
        layer_rebinds.Modifiers = make(map[uint16]keymap.Key)
        layer_rebinds.Keys = make(map[uint16]keymap.Key)
        for key, value := range keymap.Multicode_modifiers {
            layer_rebinds.Multicode_modifiers[key] = value
        }
        for key, value := range keymap.Modifiers {
            layer_rebinds.Modifiers[key] = value
        }
        for key, value := range keymap.Keys {
            layer_rebinds.Keys[key] = value
        }

        for _, rebind := range layer.Rebinds {
            rebind_output_keycode := keyname_to_keycode(rebind[1])

            var rebind_output_keystruct keymap.Key
            if _, ok := keymap.Multicode_modifiers[rebind_output_keycode]; ok {
                rebind_output_keystruct.Scancode = keymap.Multicode_modifiers[rebind_output_keycode].Scancode
                rebind_output_keystruct.Keyname = keymap.Multicode_modifiers[rebind_output_keycode].Keyname
                fmt.Println(rebind[1], rebind_output_keycode)
            } else if _, ok := keymap.Modifiers[rebind_output_keycode]; ok {
                rebind_output_keystruct = keymap.Modifiers[rebind_output_keycode]
                fmt.Println(rebind[1], rebind_output_keycode)
            } else if _, ok := keymap.Keys[rebind_output_keycode]; ok {
                rebind_output_keystruct = keymap.Keys[rebind_output_keycode]
                fmt.Println(rebind[1], rebind_output_keycode)
            }

            rebind_input_keycode := keyname_to_keycode(rebind[0])

            if multicode_mod, ok := keymap.Multicode_modifiers[rebind_input_keycode]; ok {
                if _, ok := keymap.Keys[rebind_output_keycode]; ok {
                    // if you are rebinding a multicode modifier to a key
                    // remove the key from modifiers and add both modifiers to keys
                    delete(layer_rebinds.Modifiers, multicode_mod.Leftkey)
                    delete(layer_rebinds.Modifiers, multicode_mod.Rightkey)
                    layer_rebinds.Keys[multicode_mod.Leftkey] = rebind_output_keystruct
                    layer_rebinds.Keys[multicode_mod.Rightkey] = rebind_output_keystruct
                } else if  multicode_mod2, ok := keymap.Multicode_modifiers[rebind_output_keycode]; ok {
                    // if you are rebinding a  multicode modifier to a multicode modifier
                    // add the modifier as the left mulitcode modifier to modifiers
                    layer_rebinds.Modifiers[multicode_mod.Leftkey] = keymap.Key{
                        Scancode: keymap.Modifiers[multicode_mod2.Leftkey].Scancode,
                        Keyname: keymap.Modifiers[multicode_mod2.Leftkey].Keyname}
                    layer_rebinds.Modifiers[multicode_mod.Rightkey] = keymap.Key{
                        Scancode: keymap.Modifiers[multicode_mod2.Leftkey].Scancode,
                        Keyname: keymap.Modifiers[multicode_mod2.Leftkey].Keyname}
                } else {
                    layer_rebinds.Modifiers[multicode_mod.Leftkey] = rebind_output_keystruct
                    layer_rebinds.Modifiers[multicode_mod.Rightkey] = rebind_output_keystruct
                }

            } else if _, ok := keymap.Modifiers[rebind_input_keycode]; ok {
                if _, ok := keymap.Keys[rebind_output_keycode]; ok {
                    // if you are rebinding a modifier to a key
                    // remove the key from modifiers and add it to keys
                    delete(layer_rebinds.Modifiers, rebind_input_keycode)
                    layer_rebinds.Keys[rebind_input_keycode] = rebind_output_keystruct
                } else if  multicode_mod, ok := keymap.Multicode_modifiers[rebind_output_keycode]; ok {
                    // if you are rebinding a modifier to a multicode modifier
                    // add the modifier as the left mulitcode modifier to modifiers
                    layer_rebinds.Modifiers[rebind_input_keycode] = keymap.Key{
                        Scancode: keymap.Modifiers[multicode_mod.Leftkey].Scancode,
                        Keyname: keymap.Modifiers[multicode_mod.Leftkey].Keyname}
                } else {
                    layer_rebinds.Modifiers[rebind_input_keycode] = rebind_output_keystruct
                }

            } else if _, ok := keymap.Keys[rebind_input_keycode]; ok {
                if _, ok := keymap.Modifiers[rebind_output_keycode]; ok {
                    // if you are rebinding a key to a modifier
                    // remove the key from keys and add it to modifiers
                    delete(layer_rebinds.Keys, rebind_input_keycode)
                    layer_rebinds.Modifiers[rebind_input_keycode] = rebind_output_keystruct
                } else if  multicode_mod, ok := keymap.Multicode_modifiers[rebind_output_keycode]; ok {
                    // if you are rebinding a key to a multicode modifier
                    // remove the key from keys and add it to modifiers as the left multicode modifier
                    delete(layer_rebinds.Keys, rebind_input_keycode)
                    layer_rebinds.Modifiers[rebind_input_keycode] = keymap.Key{
                        Scancode: keymap.Modifiers[multicode_mod.Leftkey].Scancode,
                        Keyname: keymap.Modifiers[multicode_mod.Leftkey].Keyname}
                } else {
                    layer_rebinds.Keys[rebind_input_keycode] = rebind_output_keystruct
                }
            }
        }

        var layer_layerbinds []Layerbind
        for _, layerbind := range layer.Layerbinds {
            var layerbind_input_keys []uint16
            for _, layerbind_input_key := range layerbind[0] {
                keycode := keyname_to_keycode(layerbind_input_key)
                layerbind_input_keys = append(layerbind_input_keys, keycode)
            }

            layerbind_to_layer, _ := strconv.Atoi(layerbind[1][0])
            layerbind_from_layer := i
            layerbind_suppress, _ := strconv.ParseBool(layerbind[1][2])
            layerbind_type := layerbind[1][1]
            layer_layerbinds = append(layer_layerbinds, Layerbind{
                Input_keys: layerbind_input_keys,
                To_layer: layerbind_to_layer,
                From_layer: layerbind_from_layer,
                Suppress: layerbind_suppress,
                State: 0,
                Type: layerbind_type})
        }

        var layer_keybinds []Keybind
        for _, keybind := range layer.Keybinds {
            var keybind_input_keys []uint16
            for _, keybind_input_key := range keybind[0] {
                keycode := keyname_to_keycode(keybind_input_key)
                keybind_input_keys = append(keybind_input_keys, keycode)
            }

            var bind_output_keys []uint16
            for _, bind_output_key := range keybind[1] {
                keycode := keyname_to_keycode(bind_output_key)
                if multi_mod, ok := keymap.Multicode_modifiers[keycode]; ok {
                    keycode = multi_mod.Leftkey
                }

                bind_output_keys = append(bind_output_keys, keycode)
            }
            layer_keybinds = append(layer_keybinds, Keybind{
                Input_keys: keybind_input_keys,
                Output_keys: bind_output_keys})
        }
        rebinds = append(rebinds, *layer_rebinds)
        layerbinds = append(layerbinds, layer_layerbinds)
        keybinds = append(keybinds, layer_keybinds)
    }
    return rebinds, layerbinds, keybinds
}

func keyname_to_keycode(keyname string) uint16 {
    for keycode, key := range keymap.Multicode_modifiers {
        if key.Keyname == keyname {
            return keycode
        }
    }

    for keycode, key := range keymap.Modifiers {
        if key.Keyname == keyname {
            return keycode
        }
    }

    for keycode, key := range keymap.Keys {
        if key.Keyname == keyname {
            return keycode
        }
    }
    return 0
}
