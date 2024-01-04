package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_parseFields(t *testing.T) {
	type args struct {
		fieldsStr string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{name: "nameMe", args: args{fieldsStr: "1,2,3"}, want: []int{1, 2, 3}, wantErr: false},
		{name: "nameMe", args: args{fieldsStr: "1,2,3,"}, want: nil, wantErr: true},
		{name: "nameMe", args: args{fieldsStr: ""}, want: nil, wantErr: true},
		{name: "nameMe", args: args{fieldsStr: "5"}, want: []int{5}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFields(tt.args.fieldsStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_app(t *testing.T) {
	type args struct {
		conf   *config
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "nameMe",
			args: args{
				conf:   &config{fieldsStr: nil, delimiter: nil, separated: nil},
				reader: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app(tt.args.conf, tt.args.reader)
		})
	}
}

func Test_app1(t *testing.T) {
	type args struct {
		conf   *config
		reader io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{}, // TODO
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			app(tt.args.conf, tt.args.reader, writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("app() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
