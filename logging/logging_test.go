package logging

import (
	"context"
	"errors" //nolint:depguard // used to test wrapping
	"reflect"
	"testing"

	"github.com/gsmcwhirter/go-util/v9/request"
)

func Test_logger_Log(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name      string
		l         *logger
		args      args
		wantLines [][]interface{}
		wantErr   bool
	}{
		{
			name:      "test pass through",
			l:         &logger{base: &dummyLogger{}},
			args:      args{[]interface{}{"foo", "bar"}},
			wantLines: [][]interface{}{{"foo", "bar"}},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Log(tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("logger.Log() error = %v, wantErr %v", err, tt.wantErr)
			}

			lines := tt.l.base.(*dummyLogger).lines
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}

func Test_logger_Printf(t *testing.T) {
	type args struct {
		f    string
		args []interface{}
	}
	tests := []struct {
		name      string
		l         *logger
		args      args
		wantLines [][]interface{}
	}{
		{
			name:      "test basic print",
			l:         &logger{base: &dummyLogger{}},
			args:      args{"%s %s", []interface{}{"foo", "bar"}},
			wantLines: [][]interface{}{{"message", "foo bar"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Printf(tt.args.f, tt.args.args...)

			lines := tt.l.base.(*dummyLogger).lines
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}

func Test_logger_Message(t *testing.T) {
	type args struct {
		msg  string
		args []interface{}
	}
	tests := []struct {
		name      string
		l         *logger
		args      args
		wantLines [][]interface{}
	}{
		{
			name:      "test basic message",
			l:         &logger{base: &dummyLogger{}},
			args:      args{"m", []interface{}{"foo", "bar"}},
			wantLines: [][]interface{}{{"message", "m", "foo", "bar"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Message(tt.args.msg, tt.args.args...)

			lines := tt.l.base.(*dummyLogger).lines
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}

func Test_logger_Err(t *testing.T) {
	testErr := errors.New("test")
	type args struct {
		msg  string
		err  error
		args []interface{}
	}
	tests := []struct {
		name      string
		l         *logger
		args      args
		wantLines [][]interface{}
	}{
		{
			name:      "test basic error",
			l:         &logger{base: &dummyLogger{}},
			args:      args{"m", testErr, []interface{}{"foo", "bar"}},
			wantLines: [][]interface{}{{"message", "m", "error", testErr, "foo", "bar"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Err(tt.args.msg, tt.args.err, tt.args.args...)

			lines := tt.l.base.(*dummyLogger).lines
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}

func TestNewFrom(t *testing.T) {
	dummy := &dummyLogger{}
	type args struct {
		l BaseLogger
	}
	tests := []struct {
		name string
		args args
		want Logger
	}{
		{
			name: "test plain base",
			args: args{
				l: dummy,
			},
			want: &logger{base: dummy},
		},
		{
			name: "test Logger base",
			args: args{
				l: &logger{base: dummy},
			},
			want: &logger{base: dummy},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummy.reset()

			if got := NewFrom(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFrom() = %v, want %v", got, tt.want)
			}

			if tt.want.(*logger).base != dummy {
				t.Errorf("NewFrom() didn't preserve exact base")
			}
		})
	}
}

func TestBaseFrom(t *testing.T) {
	dummy := &dummyLogger{}

	type args struct {
		l Logger
	}
	tests := []struct {
		name string
		args args
		want BaseLogger
	}{
		{
			name: "test Logger base",
			args: args{
				l: &logger{base: dummy},
			},
			want: dummy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummy.reset()

			if got := BaseFrom(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWith(t *testing.T) {
	dummy := &dummyLogger{}
	type args struct {
		l          Logger
		keyvals1   []interface{}
		keyvals2   []interface{}
		logKeyvals []interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantLines [][]interface{}
		wantErr   bool
	}{
		{
			name: "with tags test",
			args: args{
				l:          &logger{base: dummy},
				keyvals1:   []interface{}{"foo", "bar", "caller", DefaultCaller},
				keyvals2:   []interface{}{"test", "baz"},
				logKeyvals: []interface{}{"message", "test"},
			},
			wantLines: [][]interface{}{
				// NOTE: When adding code, you'll probably have to change the line numbers here
				{"foo", "bar", "caller", "logging_test.go:243", "message", "test"},
				{"foo", "bar", "caller", "logging_test.go:249", "test", "baz", "message", "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummy.reset()

			l := With(tt.args.l, tt.args.keyvals1...)

			if err := l.Log(tt.args.logKeyvals...); (err != nil) != tt.wantErr {
				t.Errorf("logger.Log() 1 error = %v, wantErr %v", err, tt.wantErr)
			}

			l = With(l, tt.args.keyvals2...)

			if err := l.Log(tt.args.logKeyvals...); (err != nil) != tt.wantErr {
				t.Errorf("logger.Log() 2 error = %v, wantErr %v", err, tt.wantErr)
			}

			lines := dummy.lines
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}

func TestWithContext(t *testing.T) {
	dummy := &dummyLogger{}
	rid := request.GenerateRequestID()

	type args struct {
		ctx     context.Context
		logger  Logger
		keyvals []interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantLines [][]interface{}
		wantErr   bool
	}{
		{
			name: "without request_id",
			args: args{
				logger:  With(&logger{base: dummy}, "caller", DefaultCaller),
				ctx:     context.Background(),
				keyvals: []interface{}{"message", "test"},
			},
			wantLines: [][]interface{}{
				// NOTE: When adding code, you'll probably have to change the line numbers here
				{"caller", "logging_test.go:309", "request_id", "unknown", "message", "test"},
			},
			wantErr: false,
		},
		{
			name: "with request_id",
			args: args{
				logger:  With(&logger{base: dummy}, "caller", DefaultCaller),
				ctx:     request.NewRequestContextWithRequestID(context.Background(), rid),
				keyvals: []interface{}{"message", "test"},
			},
			wantLines: [][]interface{}{
				// NOTE: When adding code, you'll probably have to change the line numbers here
				{"caller", "logging_test.go:309", "request_id", rid, "message", "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummy.reset()

			l := WithContext(tt.args.ctx, tt.args.logger)

			if err := l.Log(tt.args.keyvals...); (err != nil) != tt.wantErr {
				t.Errorf("logger.Log() 1 error = %v, wantErr %v", err, tt.wantErr)
			}

			lines := dummy.lines
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}
