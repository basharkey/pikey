package main

import (
    "testing"
    "github.com/basharkey/pimk/config"
)

func benchmark_detect_keybinds(pressed_keys []Keystate, pressed_keybinds []config.Keybind, layer_keybinds []config.Keybind, b *testing.B) {
    for i := 0; i < b.N; i++ {
        detect_keybinds(&pressed_keys, &pressed_keybinds, layer_keybinds)
    }
}

func benchmark_detect_keybinds2(pressed_keys []Keystate, pressed_keybinds []config.Keybind, layer_keybinds []config.Keybind, b *testing.B) {
    for i := 0; i < b.N; i++ {
        pressed_keybinds = detect_keybinds2(&pressed_keys, pressed_keybinds, layer_keybinds)
    }
}

func benchmark_detect_keybinds3(pressed_keys []Keystate, pressed_keybinds []config.Keybind, layer_keybinds []config.Keybind, b *testing.B) {
    for i := 0; i < b.N; i++ {
        pressed_keys, pressed_keybinds = detect_keybinds3(pressed_keys, pressed_keybinds, layer_keybinds)
    }
}

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
var pressed_keys = []Keystate {
    {30, true},
    {48, true},
    {97, true},
}
var pressed_keybinds []config.Keybind

func Benchmark_pointer2(b *testing.B) {benchmark_detect_keybinds(pressed_keys, pressed_keybinds, keybinds[0], b)}
func Benchmark_pointer1(b *testing.B) {benchmark_detect_keybinds2(pressed_keys, pressed_keybinds, keybinds[0], b)}
func Benchmark_no_pointer(b *testing.B) {benchmark_detect_keybinds3(pressed_keys, pressed_keybinds, keybinds[0], b)}
