package main

import "testing"

//TestCalcCrcCcitt test checksum calculation
func TestCalcCrcCcitt(t *testing.T) {
	expected := [2]byte{0xc0, 0x6f}
	result := calcCrcCcitt([]byte{0xb6, 0x49, 0x22, 0xb7, 0x07, 0x01, 0x81}, 0x0000)
	if result != expected {
		t.Errorf("Invalid checksum: <% x> expected: <% x>", result, expected)
	}
}
