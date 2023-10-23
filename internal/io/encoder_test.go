package io_test

import (
	"diceDB/internal/io"
	"reflect"
	"testing"
)

func Test_respEncoder_Encode(t *testing.T) {
	type args struct {
		value    interface{}
		isSimple bool
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test encode for plain string with is simple true",
			args: args{isSimple: true, value: "PING"},
			want: []byte("+PING\r\n"),
		},
		{
			name: "test encode for plain string with is simple false",
			args: args{isSimple: false, value: "PING"},
			want: []byte("$4\r\nPING\r\n"),
		},
		{
			name: "test encode for integer",
			args: args{isSimple: true, value: 1},
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := io.NewRESPEncoder()
			if got := encoder.Encode(tt.args.value, tt.args.isSimple); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
