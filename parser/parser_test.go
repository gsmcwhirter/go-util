package parser

import (
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	t.Parallel()

	type args struct {
		opts Options
	}
	tests := []struct {
		name string
		args args
		want Parser
	}{
		{
			name: "heppy path",
			args: args{
				opts: Options{
					CmdIndicator:  "!",
					KnownCommands: nil,
					CaseSensitive: false,
				},
			},
			want: &parser{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{},
				caseSensitive: false,
			},
		},
		{
			name: "preload",
			args: args{
				opts: Options{
					CmdIndicator:  "!",
					KnownCommands: []string{"foo", "bar"},
					CaseSensitive: true,
				},
			},
			want: &parser{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true, "bar": true},
				caseSensitive: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewParser(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_IsCaseSensitive(t *testing.T) {
	t.Parallel()

	type fields struct {
		CmdIndicator  string
		knownCommands map[string]bool
		caseSensitive bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "yes",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: nil,
				caseSensitive: true,
			},
			want: true,
		},
		{
			name: "no",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: nil,
				caseSensitive: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &parser{
				CmdIndicator:  tt.fields.CmdIndicator,
				knownCommands: tt.fields.knownCommands,
				caseSensitive: tt.fields.caseSensitive,
			}
			if got := p.IsCaseSensitive(); got != tt.want {
				t.Errorf("parser.IsCaseSensitive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_KnownCommand(t *testing.T) {
	t.Parallel()

	type fields struct {
		CmdIndicator  string
		knownCommands map[string]bool
		caseSensitive bool
	}
	type args struct {
		cmd string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "yes - insensitive match",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: false,
			},
			args: args{"foo"},
			want: true,
		},
		{
			name: "yes - insensitive nomatch",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: false,
			},
			args: args{"FOO"},
			want: true,
		},
		{
			name: "yes - sensitive",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: true,
			},
			args: args{"foo"},
			want: true,
		},
		{
			name: "no - sensitive nomatch",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: true,
			},
			args: args{"FOO"},
			want: false,
		},
		{
			name: "no - insensitive nomatch",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: true,
			},
			args: args{"fool"},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &parser{
				CmdIndicator:  tt.fields.CmdIndicator,
				knownCommands: tt.fields.knownCommands,
				caseSensitive: tt.fields.caseSensitive,
			}
			if got := p.KnownCommand(tt.args.cmd); got != tt.want {
				t.Errorf("parser.KnownCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_LearnCommand(t *testing.T) {
	t.Parallel()

	type fields struct {
		CmdIndicator  string
		knownCommands map[string]bool
		caseSensitive bool
	}
	type args struct {
		cmd string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantKnown map[string]bool
	}{
		{
			name: "new - insensitive",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: false,
			},
			args:      args{"BAR"},
			wantKnown: map[string]bool{"foo": true, "bar": true},
		},
		{
			name: "new - sensitive",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: true,
			},
			args:      args{"BAR"},
			wantKnown: map[string]bool{"foo": true, "BAR": true},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &parser{
				CmdIndicator:  tt.fields.CmdIndicator,
				knownCommands: tt.fields.knownCommands,
				caseSensitive: tt.fields.caseSensitive,
			}
			p.LearnCommand(tt.args.cmd)

			if !reflect.DeepEqual(p.knownCommands, tt.wantKnown) {
				t.Errorf("LearnCommand() post was %v, want %v", p.knownCommands, tt.wantKnown)
			}
		})
	}
}

func Test_parser_LeadChar(t *testing.T) {
	t.Parallel()

	type fields struct {
		CmdIndicator  string
		knownCommands map[string]bool
		caseSensitive bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{},
				caseSensitive: false,
			},
			want: "!",
		},
		{
			name: "not always !",
			fields: fields{
				CmdIndicator:  "%",
				knownCommands: map[string]bool{},
				caseSensitive: false,
			},
			want: "%",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &parser{
				CmdIndicator:  tt.fields.CmdIndicator,
				knownCommands: tt.fields.knownCommands,
				caseSensitive: tt.fields.caseSensitive,
			}
			if got := p.LeadChar(); got != tt.want {
				t.Errorf("parser.LeadChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parser_ParseCommand(t *testing.T) {
	t.Parallel()

	type fields struct {
		CmdIndicator  string
		knownCommands map[string]bool
		caseSensitive bool
	}
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantCmd string
		wantErr bool
	}{
		{
			name: "is known command",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: false,
			},
			args:    args{"!foo"},
			wantCmd: "foo",
			wantErr: false,
		},
		{
			name: "is known command insensitive",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: false,
			},
			args:    args{"!FOO"},
			wantCmd: "foo",
			wantErr: false,
		},
		{
			name: "is not known command",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: false,
			},
			args:    args{"!bar"},
			wantCmd: "bar",
			wantErr: true,
		},
		{
			name: "is not known command - sensitive",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: true,
			},
			args:    args{"!FOO"},
			wantCmd: "FOO",
			wantErr: true,
		},
		{
			name: "is known command - sensitive",
			fields: fields{
				CmdIndicator:  "!",
				knownCommands: map[string]bool{"FOO": true},
				caseSensitive: true,
			},
			args:    args{"!FOO"},
			wantCmd: "FOO",
			wantErr: false,
		},
		{
			name: "is not command",
			fields: fields{
				CmdIndicator:  "?",
				knownCommands: map[string]bool{"foo": true},
				caseSensitive: false,
			},
			args:    args{"!foo"},
			wantCmd: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &parser{
				CmdIndicator:  tt.fields.CmdIndicator,
				knownCommands: tt.fields.knownCommands,
				caseSensitive: tt.fields.caseSensitive,
			}
			gotCmd, err := p.ParseCommand(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.ParseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("parser.ParseCommand() = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}

func TestMaybeCount(t *testing.T) {
	t.Parallel()

	type args struct {
		line string
	}
	tests := []struct {
		name  string
		args  args
		wantL string
		wantC string
	}{
		{
			name:  "normal",
			args:  args{"test 123"},
			wantL: "test ",
			wantC: "123",
		},
		{
			name:  "with +",
			args:  args{"test +123"},
			wantL: "test ",
			wantC: "123",
		},
		{
			name:  "with x",
			args:  args{"test x123"},
			wantL: "test ",
			wantC: "123",
		},
		{
			name:  "with -",
			args:  args{"test -123"},
			wantL: "test ",
			wantC: "-123",
		},
		{
			name:  "fake",
			args:  args{"test123"},
			wantL: "test123",
			wantC: "",
		},
		{
			name:  "with + fake",
			args:  args{"test+123"},
			wantL: "test+123",
			wantC: "",
		},
		{
			name:  "with x fake",
			args:  args{"testx123"},
			wantL: "testx123",
			wantC: "",
		},
		{
			name:  "with - fake",
			args:  args{"test-123"},
			wantL: "test-123",
			wantC: "",
		},
		{
			name:  "just numerals",
			args:  args{"123"},
			wantL: "123",
			wantC: "",
		},
		{
			name:  "no numerals",
			args:  args{"test"},
			wantL: "test",
			wantC: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotL, gotC := MaybeCount(tt.args.line)
			if gotL != tt.wantL {
				t.Errorf("MaybeCount() gotL = %v, want %v", gotL, tt.wantL)
			}
			if gotC != tt.wantC {
				t.Errorf("MaybeCount() gotC = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
