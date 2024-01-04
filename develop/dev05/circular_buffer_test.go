package main

import (
	"reflect"
	"testing"
)

func TestCircularBuffer_Push_Retrieve(t *testing.T) {
	tests := []struct {
		name        string
		size        int
		values2push []any
		want        []any
	}{
		{name: "SemiFilled_1", size: 4, values2push: []interface{}{1, 2}, want: []interface{}{1, 2}},
		{name: "SemiFilled_2", size: 4, values2push: []interface{}{"a", "b", "c"}, want: []interface{}{"a", "b", "c"}},
		{name: "Filled_1", size: 3, values2push: []interface{}{"a", "b", "c"}, want: []interface{}{"a", "b", "c"}},
		{name: "Filled_2", size: 1, values2push: []interface{}{"a"}, want: []interface{}{"a"}},
		{name: "NoPush", size: 5, values2push: []interface{}{}, want: []interface{}{}},
		{name: "ZeroSized_NoPush", size: 0, values2push: []interface{}{}, want: []interface{}{}},
		{name: "ZeroSized_Push", size: 0, values2push: []interface{}{"a", "b", "c"}, want: []interface{}{}},
		{name: "OverBuffing_1", size: 3, values2push: []interface{}{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"}, want: []interface{}{"i", "o", "p"}},
		{name: "OverBuffing_2", size: 2, values2push: []interface{}{"q", "w", "e"}, want: []interface{}{"w", "e"}},
		{name: "SemiFilled_3", size: 5, values2push: []interface{}{1, 2, 3, 4}, want: []interface{}{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCircularBuffer(tt.size)
			for _, v := range tt.values2push {
				c.Push(v)
			}
			if got := c.Retrieve(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Retrieve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircularBuffer_Push_Pop_Retrieve(t *testing.T) {
	tests := []struct {
		name         string
		size         int
		values2push  []any
		wantPop      []any
		wantRetrieve []any
	}{
		{
			name:         "nameMe",
			size:         4,
			values2push:  []interface{}{1, 2, 3, 4},
			wantPop:      []interface{}{4, 3},
			wantRetrieve: []interface{}{1, 2},
		},

		{
			name:         "nameMe",
			size:         4,
			values2push:  []interface{}{1, 2, 3, 4, 5, 6, 7, 8},
			wantPop:      []interface{}{8, 7, 6, 5},
			wantRetrieve: []interface{}{},
		},

		{
			name:         "nameMe",
			size:         4,
			values2push:  []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantPop:      []interface{}{10, 9, 8},
			wantRetrieve: []interface{}{7},
		},

		{
			name:         "nameMe",
			size:         4,
			values2push:  []interface{}{1, 2, 3, 4},
			wantPop:      []interface{}{4, 3},
			wantRetrieve: []interface{}{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCircularBuffer(tt.size)
			for _, v := range tt.values2push {
				c.Push(v)
			}
			poped := make([]any, 0, len(tt.wantPop))
			for _ = range tt.wantPop {
				poped = append(poped, c.Pop())
			}
			if got := c.Retrieve(); !reflect.DeepEqual(got, tt.wantRetrieve) {
				t.Errorf("\nRetrieve() = %v, wantRetrieve %v\nPop() = %v, wantPop %v", got, tt.wantRetrieve, poped, tt.wantPop)
			}
			if !reflect.DeepEqual(poped, tt.wantPop) {
				t.Errorf("\nPop() = %v, wantPop %v", poped, tt.wantPop)
			}
		})
	}
}
