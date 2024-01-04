package main

import "testing"

func TestPrintTime(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test1"},
		{"test2"},
		{"test3"},
		{"test4"},
		{"test5"},
		{"test6"},
		{"test7"},
		{"test8"},
		{"test9"},
		{"test10"},
		{"test11"},
		{"test12"},
		{"test13"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintTime()
		})
	}
}
