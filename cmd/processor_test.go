package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestProcessInput(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "simple JSON",
			args: args{
				r: strings.NewReader("{\"key\": \"value\"}"),
			},
			wantW:   "value\n",
			wantErr: false,
		},
		{
			name: "simple JSON with two elements",
			args: args{
				r: strings.NewReader("{\"key\": \"value\", \"anotherKey\": \"anotherValue\"}"),
			},
			wantW:   "anotherValue value\n",
			wantErr: false,
		},
		{
			name: "nested JSON",
			args: args{
				r: strings.NewReader("{\"key\": \"value\", \"newLevel\" : {\"nestedKey1\": \"nestedValue1\", \"nestedKey2\": \"nestedValue2\" }}"),
			},
			wantW:   "value nestedValue1 nestedValue2\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := ProcessInput(tt.args.r, w)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ProcessInput() gotW = '%v', want '%v'", gotW, tt.wantW)
			}
		})
	}
}
