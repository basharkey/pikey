package config

import (
    "keymap"
    "fmt"
    "encoding/json"
    "io/ioutil"
    //"strconv"
    "errors"
    "log"
)

type Rebind struct {
    Multi_mods map[uint16]keymap.Multi_mod
    Mods map[uint16]keymap.Key
    Keys map[uint16]keymap.Key
    Consumer_keys map[uint16]keymap.Key
}

type Keybind struct {
    InputKeys []string
    InputKeyCodes []uint16
    OutputKeys []string
    OutputKeyCodes []uint16
}

type Layerbind struct {
    InputKeys []string
    InputKeyCodes []uint16
    Layer int
    FromLayer int
    Type string
    OutputInputKeys bool
    State int
}

type MacroCmd struct {
    Type string
    Params []string
}

type Macro struct {
    InputKeys []string
    InputKeyCodes []uint16
    Cmds []MacroCmd
}

type Layer struct {
    Rebinds [][]string
    Layerbinds []Layerbind
    Keybinds []Keybind
    Macros []Macro
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
        layer_rebinds.Multi_mods = make(map[uint16]keymap.Multi_mod)
        layer_rebinds.Mods = make(map[uint16]keymap.Key)
        layer_rebinds.Keys = make(map[uint16]keymap.Key)
        layer_rebinds.Consumer_keys = make(map[uint16]keymap.Key)

        for key, value := range keymap.Multi_mods {
            layer_rebinds.Multi_mods[key] = value
        }
        for key, value := range keymap.Mods {
            layer_rebinds.Mods[key] = value
        }
        for key, value := range keymap.Keys {
            layer_rebinds.Keys[key] = value
        }
        for key, value := range keymap.Consumer_keys {
            layer_rebinds.Consumer_keys[key] = value
        }

        for _, rebind := range layer.Rebinds {
            // if ALL is used rebind all mods, keys, and consumer_keys
            if rebind[0] == "ALL" {
                for _, mod := range keymap.Mods {
                    rebind_input_keycode := keyNameToKeyCode(mod.Name)
                    rebind_output_keycode := keyNameToKeyCode(rebind[1])

                    rebind_key(layer_rebinds, rebind_input_keycode, rebind_output_keycode)
                }
                for _, key := range keymap.Keys {
                    rebind_input_keycode := keyNameToKeyCode(key.Name)
                    rebind_output_keycode := keyNameToKeyCode(rebind[1])

                    rebind_key(layer_rebinds, rebind_input_keycode, rebind_output_keycode)
                }
                for _, consumer_key := range keymap.Consumer_keys {
                    rebind_input_keycode := keyNameToKeyCode(consumer_key.Name)
                    rebind_output_keycode := keyNameToKeyCode(rebind[1])

                    rebind_key(layer_rebinds, rebind_input_keycode, rebind_output_keycode)
                }
            } else {
                rebind_input_keycode := keyNameToKeyCode(rebind[0])
                rebind_output_keycode := keyNameToKeyCode(rebind[1])

                rebind_key(layer_rebinds, rebind_input_keycode, rebind_output_keycode)
            }
        }

        /*
        var layer_layerbinds []Layerbind
        for _, layerbind := range layer.Layerbinds {
            var layerbind_input_keys []uint16
            for _, layerbind_input_key := range layerbind[0] {
                keycode := keyNameToKeyCode(layerbind_input_key)
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
        */

        /*
        var layer_keybinds []Keybind
        for _, keybind := range layer.Keybinds {
            var keybind_input_keys []uint16
            for _, keybind_input_key := range keybind[0] {
                keycode := keyNameToKeyCode(keybind_input_key)
                keybind_input_keys = append(keybind_input_keys, keycode)
            }

            var bind_output_keys []uint16
            for _, bind_output_key := range keybind[1] {
                keycode := keyNameToKeyCode(bind_output_key)
                if multi_mod, ok := keymap.Multi_mods[keycode]; ok {
                    keycode = multi_mod.Left_code
                }

                bind_output_keys = append(bind_output_keys, keycode)
            }
            layer_keybinds = append(layer_keybinds, Keybind{
                Input_keys: keybind_input_keys,
                Output_keys: bind_output_keys})
        }
        */

        for j, keybind := range layer.Keybinds {
            for _, key := range keybind.InputKeys {
                keybind.InputKeyCodes = append(keybind.InputKeyCodes, keyNameToKeyCode(key))
            }
            layers[i].Keybinds[j].InputKeyCodes = keybind.InputKeyCodes

            for _, key := range keybind.OutputKeys {
                keyCode := keyNameToKeyCode(key)
                if multi_mod, ok := keymap.Multi_mods[keyCode]; ok {
                    keyCode = multi_mod.Left_code
                }
                keybind.OutputKeyCodes = append(keybind.OutputKeyCodes, keyCode)
            }
            layers[i].Keybinds[j].OutputKeyCodes = keybind.OutputKeyCodes
        }

        for j, layerbind := range layer.Layerbinds {
            for _, key := range layerbind.InputKeys {
                layerbind.InputKeyCodes = append(layerbind.InputKeyCodes, keyNameToKeyCode(key))
            }
            layers[i].Layerbinds[j].InputKeyCodes = layerbind.InputKeyCodes

            layers[i].Layerbinds[j].FromLayer = i
        }

        for j, macro := range layer.Macros {
            for _, key := range macro.InputKeys {
                macro.InputKeyCodes = append(macro.InputKeyCodes, keyNameToKeyCode(key))
            }
            layers[i].Macros[j].InputKeyCodes = macro.InputKeyCodes
        }

        rebinds = append(rebinds, *layer_rebinds)
        //layerbinds = append(layerbinds, layer_layerbinds)
        //keybinds = append(keybinds, layer_keybinds)
        fmt.Println(layer)
    }
    return rebinds, layerbinds, keybinds
}

func keycode_to_keytype(code uint16) (string, error) {
    if _, ok := keymap.Multi_mods[code]; ok {
        return "multicode_mod", nil
    } else if _, ok := keymap.Mods[code]; ok {
        return "mod", nil
    } else if _, ok := keymap.Keys[code]; ok {
        return "key", nil
    } else if _, ok := keymap.Consumer_keys[code]; ok {
        return "consumer_key", nil
    }
    return "", errors.New("keycode not in keymap")
}

func rebind_key(layer_rebinds *Rebind, input_keycode uint16, output_keycode uint16) {
    input_keytype, err := keycode_to_keytype(input_keycode)
    if err != nil {
        log.Fatal(err)
    }

    output_keytype, err := keycode_to_keytype(output_keycode)
    if err != nil {
        log.Fatal(err)
    }

    if input_keytype == "multicode_mod" {
        input_keystruct := keymap.Multi_mods[input_keycode]

        if output_keytype == "multicode_mod" {
            output_keystruct := keymap.Multi_mods[output_keycode]
            layer_rebinds.Mods[input_keystruct.Left_code] = keymap.Key{
                Code: keymap.Mods[output_keystruct.Left_code].Code,
                Name: keymap.Mods[output_keystruct.Left_code].Name}

            layer_rebinds.Mods[input_keystruct.Right_code] = keymap.Key{
                Code: keymap.Mods[output_keystruct.Left_code].Code,
                Name: keymap.Mods[output_keystruct.Left_code].Name}

        } else if output_keytype == "mod" {
            output_keystruct := keymap.Mods[output_keycode]
            layer_rebinds.Mods[input_keystruct.Left_code] = output_keystruct
            layer_rebinds.Mods[input_keystruct.Right_code] = output_keystruct

        } else if output_keytype == "key" {
            output_keystruct := keymap.Keys[output_keycode]
            delete(layer_rebinds.Mods, input_keystruct.Left_code)
            delete(layer_rebinds.Mods, input_keystruct.Right_code)
            layer_rebinds.Keys[input_keystruct.Left_code] = output_keystruct
            layer_rebinds.Keys[input_keystruct.Right_code] = output_keystruct

        } else if output_keytype == "consumer_key" {
            output_keystruct := keymap.Consumer_keys[output_keycode]
            delete(layer_rebinds.Mods, input_keystruct.Left_code)
            delete(layer_rebinds.Mods, input_keystruct.Right_code)
            layer_rebinds.Consumer_keys[input_keystruct.Left_code] = output_keystruct
            layer_rebinds.Consumer_keys[input_keystruct.Right_code] = output_keystruct
        }

    } else if input_keytype == "mod" {
        if output_keytype == "multicode_mod" {
            output_keystruct := keymap.Multi_mods[output_keycode]
            layer_rebinds.Mods[input_keycode] = keymap.Key{
                Code: keymap.Mods[output_keystruct.Left_code].Code,
                Name: keymap.Mods[output_keystruct.Left_code].Name}

        } else if output_keytype == "mod" {
            output_keystruct := keymap.Mods[output_keycode]
            layer_rebinds.Mods[input_keycode] = output_keystruct

        } else if output_keytype == "key" {
            output_keystruct := keymap.Keys[output_keycode]
            delete(layer_rebinds.Mods, input_keycode)
            layer_rebinds.Keys[input_keycode] = output_keystruct

        } else if output_keytype == "consumer_key" {
            output_keystruct := keymap.Consumer_keys[output_keycode]
            delete(layer_rebinds.Mods, input_keycode)
            layer_rebinds.Consumer_keys[input_keycode] = output_keystruct
        }

    } else if input_keytype == "key" {
        if output_keytype == "multicode_mod" {
            output_keystruct := keymap.Multi_mods[output_keycode]
            delete(layer_rebinds.Keys, input_keycode)
            layer_rebinds.Mods[input_keycode] = keymap.Key{
                Code: keymap.Mods[output_keystruct.Left_code].Code,
                Name: keymap.Mods[output_keystruct.Left_code].Name}

        } else if output_keytype == "mod" {
            output_keystruct := keymap.Mods[output_keycode]
            delete(layer_rebinds.Keys, input_keycode)
            layer_rebinds.Mods[input_keycode] = output_keystruct

        } else if output_keytype == "key" {
            output_keystruct := keymap.Keys[output_keycode]
            layer_rebinds.Keys[input_keycode] = output_keystruct

        } else if output_keytype == "consumer_key" {
            output_keystruct := keymap.Consumer_keys[output_keycode]
            delete(layer_rebinds.Keys, input_keycode)
            layer_rebinds.Consumer_keys[input_keycode] = output_keystruct
        }
    }
}

func keyNameToKeyCode(name string) uint16 {
    for code, key := range keymap.Multi_mods {
        if key.Name == name {
            return code
        }
    }

    for code, key := range keymap.Mods {
        if key.Name == name {
            return code
        }
    }

    for code, key := range keymap.Keys {
        if key.Name == name {
            return code
        }
    }

    for code, key := range keymap.Consumer_keys {
        if key.Name == name {
            return code
        }
    }
    return 0
}
