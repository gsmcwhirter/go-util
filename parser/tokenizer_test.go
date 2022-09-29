package parser

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	t.Parallel()

	type args struct {
		msg    string
		delim  rune
		escape rune
		quot   []rune
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "no escapes",
			args: args{
				msg:    "test 1 2 3",
				delim:  ' ',
				escape: '\\',
				quot:   []rune{'"'},
			},
			want:    []string{"test", "1", "2", "3"},
			wantErr: false,
		},
		{
			name: "no splits",
			args: args{
				msg:    "test 1 2 3",
				delim:  '%',
				escape: '\\',
				quot:   []rune{'"'},
			},
			want:    []string{"test 1 2 3"},
			wantErr: false,
		},
		{
			name: "quotes",
			args: args{
				msg:    `test "1 2 3"`,
				delim:  ' ',
				escape: '\\',
				quot:   []rune{'"'},
			},
			want:    []string{"test", "1 2 3"},
			wantErr: false,
		},
		{
			name: "escaped quote",
			args: args{
				msg:    `test \"1 2 3\"`,
				delim:  ' ',
				escape: '\\',
				quot:   []rune{'"'},
			},
			want:    []string{`test`, `"1`, `2`, `3"`},
			wantErr: false,
		},
		{
			name: "escaped in quotes",
			args: args{
				msg:    `test "1 \"2 3"`,
				delim:  ' ',
				escape: '\\',
				quot:   []rune{'"'},
			},
			want:    []string{`test`, `1 "2 3`},
			wantErr: false,
		},
		{
			name: "unmatched quotes",
			args: args{
				msg:    `test "1 2 3`,
				delim:  ' ',
				escape: '\\',
				quot:   []rune{'"'},
			},
			want:    []string{`test`, `1 2 3`},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Tokenize(tt.args.msg, tt.args.delim, tt.args.escape, tt.args.quot)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
