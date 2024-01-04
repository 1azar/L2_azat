package main

import (
	"strings"
	"testing"
)

func TestUnpack(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "validInput_Example1", args: args{"a4bc2d5e"}, want: "aaaabccddddde", wantErr: false},
		{name: "validInput_Example2", args: args{"abcd"}, want: "abcd", wantErr: false},
		{name: "invalidInput_Example3", args: args{"45"}, want: "", wantErr: true},
		{name: "validInput_Example4", args: args{""}, want: "", wantErr: false},
		{name: "invalidInput_Spaces", args: args{"69 - Zodiac sign - Pisces"}, want: "", wantErr: true},
		{name: "validInput_Cyrillic", args: args{"a0c3"}, want: "ccc", wantErr: false},
		{name: "validInput_Cyrillic", args: args{"a30"}, want: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", wantErr: false},
		{name: "validInput_Cyrillic", args: args{"a03"}, want: "aaa", wantErr: false},
		{name: "validInput_Cyrillic", args: args{"а1б2в3г4д5е6ж7"}, want: "аббвввггггдддддеееееежжжжжжж", wantErr: false},
		{name: "validInput_Cyrillic", args: args{"аб2в3г4д5е6ж7"}, want: "аббвввггггдддддеееееежжжжжжж", wantErr: false},
		{name: "invalidInput_Cyrillic", args: args{"1аб2в3г4д5е6ж7"}, want: "", wantErr: true},
		{name: "validInput_Chinese", args: args{"货币5"}, want: "货币币币币币", wantErr: false},
		{name: "validInput_Chinese", args: args{"货0"}, want: "", wantErr: false},
		{name: "invalidInput_Chinese", args: args{"100货"}, want: "", wantErr: true},
		{name: "validInput_Marks", args: args{":)10"}, want: ":))))))))))", wantErr: false},
		{name: "validInput_Marks", args: args{"❤3"}, want: "❤❤❤", wantErr: false},
		{name: "validInput_Marks", args: args{"(^/3^)"}, want: "(^///^)", wantErr: false},
		{name: "validInput_Marks", args: args{"<3"}, want: "<<<", wantErr: false},
		{name: "invalidInput_Marks", args: args{"8D"}, want: "", wantErr: true},
		{name: "hugeInput_Cyrillic", args: args{"Ж9999999"}, want: strings.Repeat("Ж", 9999999), wantErr: false},
		{name: "validInput_EscapeChar", args: args{`qwe\4\5`}, want: "qwe45", wantErr: false},
		{name: "validInput_EscapeChar", args: args{`qwe\45`}, want: "qwe44444", wantErr: false},
		{name: "validInput_EscapeChar", args: args{`qwe\\5`}, want: `qwe\\\\\`, wantErr: false},
		{name: "validInput_EscapeChar", args: args{`\5\6\7\8\10`}, want: "5678", wantErr: false},
		{name: "validInput_EscapeChar", args: args{`\00`}, want: "", wantErr: false},
		{name: "validInput_EscapeChar", args: args{`0\90`}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Unpack() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractNumber(t *testing.T) {
	tests := []struct {
		name            string
		args            []rune
		wantNumber      int
		wantDigitsCount int
		wantErr         bool
	}{
		{name: "validInput", args: []rune("34aaa"), wantNumber: 34, wantDigitsCount: 2, wantErr: false},
		{name: "validInput", args: []rune("4ff"), wantNumber: 4, wantDigitsCount: 1, wantErr: false},
		{name: "validInput", args: []rune("0fa"), wantNumber: 0, wantDigitsCount: 1, wantErr: false},
		{name: "invalidInput", args: []rune("as09"), wantNumber: 0, wantDigitsCount: 0, wantErr: true},
		{name: "validInput", args: []rune("09sa"), wantNumber: 9, wantDigitsCount: 2, wantErr: false},
		{name: "validInput", args: []rune("0009dff"), wantNumber: 9, wantDigitsCount: 4, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNumber, gotDigitsCount, err := extractNumber(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNumber != tt.wantNumber {
				t.Errorf("extractNumber() gotNumber = %v, want %v", gotNumber, tt.wantNumber)
			}
			if gotDigitsCount != tt.wantDigitsCount {
				t.Errorf("extractNumber() gotDigitsCount = %v, want %v", gotDigitsCount, tt.wantDigitsCount)
			}
		})
	}
}
