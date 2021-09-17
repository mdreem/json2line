package processor

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"text/template"
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
		{
			name: "multiple lines of JSON",
			args: args{
				r: strings.NewReader("{\"key\": \"line1\"}\n{\"key\": \"line2\"}"),
			},
			wantW:   "line1\nline2\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := ProcessInput(tt.args.r, w, nil, nil)
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

func TestProcessInputWithTemplate(t *testing.T) {
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
			wantW:   "->value<--><no value><-[<no value>-<no value>]\n",
			wantErr: false,
		},
		{
			name: "simple JSON with two elements",
			args: args{
				r: strings.NewReader("{\"key\": \"value\", \"anotherKey\": \"anotherValue\"}"),
			},
			wantW:   "->value<-->anotherValue<-[<no value>-<no value>]\n",
			wantErr: false,
		},
		{
			name: "nested JSON",
			args: args{
				r: strings.NewReader("{\"key\": \"value\", \"newLevel\" : {\"nestedKey1\": \"nestedValue1\", \"nestedKey2\": \"nestedValue2\" }}"),
			},
			wantW:   "->value<--><no value><-[nestedValue1-nestedValue2]\n",
			wantErr: false,
		},
		{
			name: "multiple lines of JSON",
			args: args{
				r: strings.NewReader("{\"key\": \"line1\"}\n{\"key\": \"line2\"}"),
			},
			wantW:   "->line1<--><no value><-[<no value>-<no value>]\n->line2<--><no value><-[<no value>-<no value>]\n",
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				r: strings.NewReader("NOT_A_JSON"),
			},
			wantW:   "NOT_A_JSON\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parse, err := template.New("test_template").Parse("->{{ .key }}<-->{{ .anotherKey }}<-[{{ .newLevel.nestedKey1 }}-{{ .newLevel.nestedKey2 }}]")
			if err != nil {
				t.Errorf("could not parse template: %v", err)
				return
			}

			w := &bytes.Buffer{}
			err = ProcessInput(tt.args.r, w, parse, nil)
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

func TestProcessInputWithTemplateHandlingSpecialCharacters(t *testing.T) {
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
				r: strings.NewReader("{\"@key\": \"value\"}"),
			},
			wantW:   "->value<-[<no value>-<no value>]\n",
			wantErr: false,
		},
		{
			name: "nested JSON",
			args: args{
				r: strings.NewReader("{\"@key\": \"value\", \"@newLevel\" : {\"@nestedKey1\": \"nestedValue1\", \"@nestedKey2\": \"nestedValue2\" }}"),
			},
			wantW:   "->value<-[nestedValue1-nestedValue2]\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parse, err := template.New("test_template").Parse("->{{ .at_key }}<-[{{ .at_newLevel.at_nestedKey1 }}-{{ .at_newLevel.at_nestedKey2 }}]")
			if err != nil {
				t.Errorf("could not parse template: %v", err)
				return
			}

			w := &bytes.Buffer{}
			err = ProcessInput(tt.args.r, w, parse, map[string]string{"@": "at_"})
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
