package errors

import (
	"errors" //nolint:depguard // this is the package that wraps the stdlib
	"reflect"
	"testing"
)

func Test_errStruct_Msg(t *testing.T) {
	t.Parallel()

	type fields struct {
		msg  string
		data []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				msg:  "test",
				data: nil,
			},
			want: "test",
		},
		{
			name: "with data",
			fields: fields{
				msg:  "test",
				data: []interface{}{"foo", "bar", "baz", 1},
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &errStruct{
				msg:  tt.fields.msg,
				data: tt.fields.data,
			}
			if got := e.Msg(); got != tt.want {
				t.Errorf("errStruct.Msg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errStruct_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		msg   string
		data  []interface{}
		cause error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				msg:  "test",
				data: nil,
			},
			want: "test",
		},
		{
			name: "with data",
			fields: fields{
				msg:  "test",
				data: []interface{}{"foo", "bar", "baz", 1},
			},
			want: "test foo=bar baz=1",
		},
		{
			name: "double wrapped",
			fields: fields{
				msg:   "another level",
				data:  []interface{}{"data2", "yep"},
				cause: Wrap(New("cause"), "first level", "data1", "woo!"),
			},
			want: "another level: first level: cause data1=woo! data2=yep",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &errStruct{
				msg:   tt.fields.msg,
				data:  tt.fields.data,
				cause: tt.fields.cause,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("errStruct.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errStruct_Unwrap(t *testing.T) {
	t.Parallel()

	type fields struct {
		msg  string
		data []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "basic",
			fields: fields{
				msg:  "test",
				data: nil,
			},
			want: nil,
		},
		{
			name: "with data",
			fields: fields{
				msg:  "test",
				data: []interface{}{"foo", "bar", "baz", 1},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &errStruct{
				msg:  tt.fields.msg,
				data: tt.fields.data,
			}
			if err := e.Unwrap(); err != tt.want {
				t.Errorf("errStruct.Unwrap() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func Test_errStruct_Data(t *testing.T) {
	t.Parallel()

	type fields struct {
		msg   string
		data  []interface{}
		cause error
	}
	tests := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		{
			name: "basic",
			fields: fields{
				msg:  "test",
				data: nil,
			},
			want: nil,
		},
		{
			name: "with data",
			fields: fields{
				msg:  "test",
				data: []interface{}{"foo", "bar", "baz", 1},
			},
			want: []interface{}{"foo", "bar", "baz", 1},
		},
		{
			name: "double wrapped",
			fields: fields{
				msg:   "another level",
				data:  []interface{}{"data2", "yep"},
				cause: Wrap(New("cause"), "first level", "data1", "woo!"),
			},
			want: []interface{}{"data1", "woo!", "data2", "yep"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &errStruct{
				msg:   tt.fields.msg,
				data:  tt.fields.data,
				cause: tt.fields.cause,
			}
			if got := e.Data(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("errStruct.Data() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errStruct_addDetails(t *testing.T) {
	t.Parallel()

	type fields struct {
		msg  string
		data []interface{}
	}
	type args struct {
		data []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
	}{
		{
			name: "start nil",
			fields: fields{
				msg:  "test",
				data: nil,
			},
			args: args{
				data: []interface{}{"foo", "bar"},
			},
			want: []interface{}{"foo", "bar"},
		},
		{
			name: "with data",
			fields: fields{
				msg:  "test",
				data: []interface{}{"foo", "bar", "baz", 1},
			},
			args: args{
				data: []interface{}{"foo", "bar"},
			},
			want: []interface{}{"foo", "bar", "baz", 1, "foo", "bar"},
		},
		{
			name: "no parity fix",
			fields: fields{
				msg:  "test",
				data: nil,
			},
			args: args{
				data: []interface{}{"foo"},
			},
			want: []interface{}{"foo"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &errStruct{
				msg:  tt.fields.msg,
				data: tt.fields.data,
			}
			e.addDetails(tt.args.data)
			if got := e.Data(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("errStruct.Data() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapped_Error(t *testing.T) {
	t.Parallel()

	testErr := errors.New("cause")
	testData := WithDetails(New("cause"), "quux", "foobar")

	type fields struct {
		errStruct errStruct
		cause     error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "cause basic",
			fields: fields{
				errStruct: errStruct{
					msg:  "test",
					data: nil,
				},
				cause: testErr,
			},
			want: "test: cause",
		},
		{
			name: "cause with data",
			fields: fields{
				errStruct: errStruct{
					msg:  "test",
					data: []interface{}{"foo", "bar", "baz", 1},
				},
				cause: testErr,
			},
			want: "test: cause foo=bar baz=1",
		},
		{
			name: "nest data basic",
			fields: fields{
				errStruct: errStruct{
					msg:  "test",
					data: nil,
				},
				cause: testData,
			},
			want: "test: cause quux=foobar",
		},
		{
			name: "nest with data",
			fields: fields{
				errStruct: errStruct{
					msg:  "test",
					data: []interface{}{"foo", "bar", "baz", 1},
				},
				cause: testData,
			},
			want: "test: cause quux=foobar foo=bar baz=1",
		},
		{
			name: "no new message",
			fields: fields{
				errStruct: errStruct{
					msg:  "",
					data: []interface{}{"foo", "bar", "baz", 1},
				},
				cause: testData,
			},
			want: "cause quux=foobar foo=bar baz=1",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := &errStruct{
				msg:   tt.fields.errStruct.msg,
				data:  tt.fields.errStruct.data,
				cause: tt.fields.cause,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("errStruct.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wrapped_Unwrap(t *testing.T) {
	t.Parallel()

	c := errors.New("test")

	type fields struct {
		errStruct errStruct
		cause     error
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "basic",
			fields: fields{
				errStruct: errStruct{
					msg:  "test",
					data: nil,
				},
				cause: c,
			},
			want: c,
		},
		{
			name: "basic 2",
			fields: fields{
				errStruct: errStruct{
					msg:  "test",
					data: nil,
				},
				cause: nil,
			},
			want: nil,
		},
		{
			name: "with data",
			fields: fields{
				errStruct: errStruct{
					msg:  "test",
					data: []interface{}{"foo", "bar", "baz", 1},
				},
				cause: c,
			},
			want: c,
		},
		{
			name: "with data 2",
			fields: fields{
				errStruct: errStruct{
					msg:  "test",
					data: []interface{}{"foo", "bar", "baz", 1},
				},
				cause: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := &errStruct{
				msg:   tt.fields.errStruct.msg,
				data:  tt.fields.errStruct.data,
				cause: tt.fields.cause,
			}
			if err := e.Unwrap(); err != tt.want {
				t.Errorf("errStruct.Unwrap() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func Test_formatData(t *testing.T) {
	t.Parallel()

	type args struct {
		data []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil",
			args: args{
				data: nil,
			},
			want: "",
		},
		{
			name: "good parity",
			args: args{
				data: []interface{}{"foo", "bar"},
			},
			want: "foo=bar",
		},
		{
			name: "good parity 2",
			args: args{
				data: []interface{}{"foo", "bar", "baz", 1},
			},
			want: "foo=bar baz=1",
		},
		{
			name: "bad parity",
			args: args{
				data: []interface{}{"foo", "bar", "baz"},
			},
			want: "foo=bar",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := formatData(tt.args.data); got != tt.want {
				t.Errorf("formatData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantMsg string
	}{
		{
			name: "happy path",
			args: args{
				msg: "test msg",
			},
			wantMsg: "test msg",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := New(tt.args.msg)
			if err.Error() != tt.wantMsg {
				t.Errorf("New() error = %v, wantMsg %s", err, tt.wantMsg)
			}
		})
	}
}

func TestNewf(t *testing.T) {
	t.Parallel()

	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantMsg string
	}{
		{
			name: "no format",
			args: args{
				msg: "test msg",
			},
			wantMsg: "test msg",
		},
		{
			name: "formatted",
			args: args{
				msg:  "test %d %d %d",
				args: []interface{}{1, 2, 3},
			},
			wantMsg: "test 1 2 3",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Newf(tt.args.msg, tt.args.args...)
			if err.Error() != tt.wantMsg {
				t.Errorf("Newf() error = %v, wantMsg %s", err, tt.wantMsg)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	t.Parallel()

	testErr := errors.New("cause")
	testData := WithDetails(New("cause"), "quux", "foobar")

	type args struct {
		err  error
		msg  string
		data []interface{}
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "nil",
			args: args{
				err:  nil,
				msg:  "test",
				data: []interface{}{"foo", "bar"},
			},
			want: nil,
		},
		{
			name: "basic error",
			args: args{
				err:  testErr,
				msg:  "test",
				data: []interface{}{"foo", "bar"},
			},
			want: &errStruct{
				msg:   "test",
				data:  []interface{}{"foo", "bar"},
				cause: testErr,
			},
		},
		{
			name: "data error",
			args: args{
				err:  testData,
				msg:  "test",
				data: []interface{}{"foo", "bar"},
			},
			want: &errStruct{
				msg:   "test",
				data:  []interface{}{"foo", "bar"},
				cause: testData,
			},
		},
		{
			name: "bad parity",
			args: args{
				err:  testErr,
				msg:  "test",
				data: []interface{}{"foo", "bar", "baz"},
			},
			want: &errStruct{
				msg:   "test",
				data:  []interface{}{"foo", "bar", "baz", ""},
				cause: testErr,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := Wrap(tt.args.err, tt.args.msg, tt.args.data...); !reflect.DeepEqual(err, tt.want) {
				t.Errorf("Wrap() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestWithDetails(t *testing.T) {
	t.Parallel()

	testErr := errors.New("cause")
	testData := WithDetails(New("cause"), "quux", "foobar")

	type args struct {
		err  error
		data []interface{}
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "nil",
			args: args{
				err:  nil,
				data: []interface{}{"foo", "bar"},
			},
			want: nil,
		},
		{
			name: "basic error",
			args: args{
				err:  testErr,
				data: []interface{}{"foo", "bar"},
			},
			want: &errStruct{
				msg:   "",
				data:  []interface{}{"foo", "bar"},
				cause: testErr,
			},
		},
		{
			name: "data error",
			args: args{
				err:  testData,
				data: []interface{}{"foo", "bar"},
			},
			want: &errStruct{"cause", []interface{}{"quux", "foobar", "foo", "bar"}, nil},
		},
		{
			name: "bad parity",
			args: args{
				err:  testErr,
				data: []interface{}{"foo", "bar", "baz"},
			},
			want: &errStruct{
				msg:   "",
				data:  []interface{}{"foo", "bar", "baz", ""},
				cause: testErr,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := WithDetails(tt.args.err, tt.args.data...); !reflect.DeepEqual(err, tt.want) {
				t.Errorf("WithDetails() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestWithDetailsMsg(t *testing.T) {
	t.Parallel()

	testErr := errors.New("cause")
	testData := WithDetails(New("cause"), "quux", "foobar")

	type args struct {
		err  error
		data []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic error",
			args: args{
				err:  testErr,
				data: []interface{}{"foo", "bar"},
			},
			want: "cause",
		},
		{
			name: "data error",
			args: args{
				err:  testData,
				data: []interface{}{"foo", "bar"},
			},
			want: "cause",
		},
		{
			name: "bad parity",
			args: args{
				err:  testErr,
				data: []interface{}{"foo", "bar", "baz"},
			},
			want: "cause",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := WithDetails(tt.args.err, tt.args.data...)
			if err == nil {
				t.Errorf("WithDetails() returned nil")
				return
			}

			e, ok := err.(Error)
			if !ok {
				t.Errorf("Wrap() returned a non-Error")
				return
			}

			if msg := e.Msg(); !reflect.DeepEqual(msg, tt.want) {
				t.Errorf("Msg() error = %v, wantErr %v", e.Msg(), tt.want)
			}
		})
	}
}

func TestWrappedMsg(t *testing.T) {
	t.Parallel()

	testErr := errors.New("cause")
	testData := WithDetails(New("cause"), "quux", "foobar")
	testDouble := Wrap(errors.New("cause"), "first level", "data1", "woo!")

	type args struct {
		err  error
		msg  string
		data []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic error",
			args: args{
				err:  testErr,
				msg:  "test",
				data: []interface{}{"foo", "bar"},
			},
			want: "test: cause",
		},
		{
			name: "data error",
			args: args{
				err:  testData,
				msg:  "test",
				data: []interface{}{"foo", "bar"},
			},
			want: "test: cause",
		},
		{
			name: "bad parity",
			args: args{
				err:  testErr,
				msg:  "test",
				data: []interface{}{"foo", "bar", "baz"},
			},
			want: "test: cause",
		},
		{
			name: "double wrapped",
			args: args{
				err:  testDouble,
				msg:  "another level",
				data: []interface{}{"data2", "yep"},
			},
			want: "another level: first level: cause",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Wrap(tt.args.err, tt.args.msg, tt.args.data...)
			if err == nil {
				t.Errorf("Wrap() returned nil")
				return
			}

			e, ok := err.(Error)
			if !ok {
				t.Errorf("Wrap() returned a non-Error")
				return
			}

			if msg := e.Msg(); !reflect.DeepEqual(msg, tt.want) {
				t.Errorf("Msg() = %v, want = %v", msg, tt.want)
			}
		})
	}
}
