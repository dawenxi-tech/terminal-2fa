package main

import (
	"testing"
)

func Fuzz_sliceMoveElement(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte, index int, offset int) {
		newData := make([]byte, 0, len(data))
		newData = append(newData, data...)
		newData = sliceMoveElement(newData, index, offset)
		if 0 <= index && index < len(data) && 0 < index+offset && index+offset < len(data) {
			if data[index] != newData[index+offset] {
				t.Errorf("err: %v, %d, %d", data, index, offset)
			}
		}
	})
}
