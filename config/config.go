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
        layer_rebinds := Rebind{
            Multicode_modifiers: keymap.Multicode_modifiers,
            Modifiers: keymap.Modifiers,
            Keys: keymap.Keys}
        for _, rebind := range layer.Rebinds {
            rebind_key1_keycode := keyname_to_keycode(rebind[0])
            rebind_key1_scancode := keymap.Keys[rebind_key1_keycode]
            rebind_key2_keycode := keyname_to_keycode(rebind[1])
            rebind_key2_scancode := keymap.Keys[rebind_key2_keycode]

            layer_rebinds.Keys[rebind_key1_keycode] = rebind_key2_scancode
            layer_rebinds.Keys[rebind_key2_keycode] = rebind_key1_scancode
        }

        var layer_layerbinds []Layerbind
        for _, layerbind := range layer.Layerbinds {
            var layerbind_input_keys []uint16
            for _, layerbind_input_key := range layerbind[0] {
                layerbind_input_keys = append(layerbind_input_keys, keyname_to_keycode(layerbind_input_key))
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
                keybind_input_keys = append(keybind_input_keys, keyname_to_keycode(keybind_input_key))
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
        rebinds = append(rebinds, layer_rebinds)
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
