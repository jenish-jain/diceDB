package io_test

import (
	"diceDB/internal/io"
	"errors"
	"reflect"
	"testing"
)

func Test_Decode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name:    "test decode with empty input",
			args:    args{data: []byte("")},
			want:    nil,
			wantErr: errors.New("no data"),
		},
		{
			name:    "test simple string decode",
			args:    args{data: []byte("+OK\r\n")},
			want:    "OK",
			wantErr: nil,
		},
		{
			name:    "test err decode",
			args:    args{data: []byte("-Error message\r\n")},
			want:    "Error message",
			wantErr: nil,
		},
		{
			name:    "test integer decode 1",
			args:    args{data: []byte(":99\r\n")},
			want:    int64(99),
			wantErr: nil,
		},
		{
			name:    "test integer decode 2",
			args:    args{data: []byte(":0\r\n")},
			want:    int64(0),
			wantErr: nil,
		},
		{
			name:    "test bulk string decode 1",
			args:    args{data: []byte("$5\r\nhello\r\n")},
			want:    "hello",
			wantErr: nil,
		},
		{
			name:    "test bulk string decode 2",
			args:    args{data: []byte("$0\r\\n\r\n")},
			want:    "",
			wantErr: nil,
		},
		{
			name:    "test bulk array decode 1",
			args:    args{data: []byte("*0\r\n")},
			want:    []interface{}{},
			wantErr: nil,
		},
		{
			name:    "test bulk array decode 2",
			args:    args{data: []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n")},
			want:    []interface{}{"hello", "world"},
			wantErr: nil,
		},
		{
			name:    "test bulk array decode 3",
			args:    args{data: []byte("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n")},
			want:    []interface{}{[]interface{}{int64(1), int64(2), int64(3)}, []interface{}{"Hello", "World"}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := io.NewRESPDecoder()
			got, err := decoder.Decode(tt.args.data)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
