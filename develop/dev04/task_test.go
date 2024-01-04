package main

import (
	"reflect"
	"testing"
)

func TestAnagramLookup(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  *map[string][]string
	}{
		{name: "InitTestCase1",
			input: []string{"пятак", "пятка", "тяпка"},
			want:  &map[string][]string{"пятак": {"пятка", "тяпка"}}},
		{name: "InitTestCase2",
			input: []string{"листок", "слиток", "столик"},
			want:  &map[string][]string{"листок": {"слиток", "столик"}}},
		{name: "InitTestCase3",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			want: &map[string][]string{
				"пятак":  {"пятка", "тяпка"},
				"листок": {"слиток", "столик"},
			},
		},
		{name: "NoAnagrams",
			input: []string{"a", "b", "c", "d", "e", "f"},
			want:  &map[string][]string{},
		},
		{name: "SingleAnagramSet1",
			input: []string{"abc", "bca", "cab", "acb", "bac", "cba", "nba", "bff", "zzz"},
			want:  &map[string][]string{"abc": {"acb", "bac", "bca", "cab", "cba"}},
		},
		{name: "SingleAnagramSet2",
			input: []string{"bff", "zzz", "abc", "bca", "cab", "acb", "bac", "cba", "nba"},
			want:  &map[string][]string{"abc": {"acb", "bac", "bca", "cab", "cba"}},
		},
		{name: "TripleAnagramSet",
			input: []string{"зз", "цц", "аб", "бв", "вб", "ба", "хф", "фа", "ут", "дд", "дд", "аф"},
			want: &map[string][]string{
				"аб": {"ба"},
				"бв": {"вб"},
				"дд": {"дд"},
				"фа": {"аф"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnagramLookup(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnagramLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}
