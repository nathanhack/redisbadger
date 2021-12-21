package command

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"

	"github.com/nathanhack/redisbadger/command/commandname"
)

func TestParseCommand(t *testing.T) {
	type args struct {
		tokens [][]byte
	}
	tests := []struct {
		args    args
		want    *Command
		wantErr bool
	}{
		0: {args: args{bytes.Split([]byte("Set Yes there"), []byte{' '})}, want: &Command{Cmd: commandname.Set, Args: []string{"Yes", "there"}, Options: []Option{}}, wantErr: false},
		1: {args: args{bytes.Split([]byte("Set Yes there GET"), []byte{' '})}, want: &Command{Cmd: commandname.Set, Args: []string{"Yes", "there"}, Options: []Option{{Flag: true, Name: "GET", Value: ""}}}, wantErr: false},
		2: {args: args{bytes.Split([]byte("Set Yes there bob"), []byte{' '})}, want: nil, wantErr: true},
		3: {args: args{bytes.Split([]byte("Set Yes there EX"), []byte{' '})}, want: nil, wantErr: true},
		4: {args: args{bytes.Split([]byte("Set Yes"), []byte{' '})}, want: nil, wantErr: true},
		5: {args: args{bytes.Split([]byte("Set Yes there EX 1"), []byte{' '})}, want: &Command{Cmd: commandname.Set, Args: []string{"Yes", "there"}, Options: []Option{{Flag: false, Name: "EX", Value: "1"}}}, wantErr: false},
		6: {args: args{bytes.Split([]byte("Set Yes there EX 1 GET"), []byte{' '})}, want: &Command{Cmd: commandname.Set, Args: []string{"Yes", "there"}, Options: []Option{{Flag: false, Name: "EX", Value: "1"}, {Flag: true, Name: "GET"}}}, wantErr: false},
		7: {args: args{bytes.Split([]byte("Ping"), []byte{' '})}, want: &Command{Cmd: commandname.Ping, Args: []string{}, Options: []Option{}}, wantErr: false},
		8: {args: args{bytes.Split([]byte("Ping"), []byte{' '})}, want: &Command{Cmd: commandname.Ping, Args: []string{}, Options: []Option{}}, wantErr: false},
	}
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			got, err := ParseCommand(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommand_String(t *testing.T) {
	type fields struct {
		Cmd     commandname.CommandName
		Args    []string
		Options []Option
	}
	tests := []struct {
		fields fields
		want   string
	}{
		{fields: fields{Cmd: commandname.Set, Args: []string{"Hello there"}, Options: []Option{{Flag: true, Name: "GET"}, {Flag: false, Name: "EX", Value: "1"}}}, want: "SET Hello there GET EX 1"},
		{fields: fields{Cmd: commandname.Set, Args: []string{"Hello there"}, Options: []Option{{Flag: false, Name: "EX", Value: "1"}, {Flag: true, Name: "GET"}}}, want: "SET Hello there EX 1 GET"},
	}
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			c := Command{
				Cmd:     tt.fields.Cmd,
				Args:    tt.fields.Args,
				Options: tt.fields.Options,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Command.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
