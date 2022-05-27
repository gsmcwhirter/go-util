package level

import (
	"reflect"
	"testing"

	"github.com/gsmcwhirter/go-util/v10/logging"
)

func TestDebug(t *testing.T) {
	dummy := &dummyLogger{}

	type args struct {
		logger  logging.Logger
		keyvals []interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantLines [][]interface{}
		wantErr   bool
	}{
		{
			name: "with level",
			args: args{
				logger:  logging.With(logging.NewFrom(dummy), "caller", logging.DefaultCaller),
				keyvals: []interface{}{"message", "test"},
			},
			wantLines: [][]interface{}{
				// NOTE: When adding code, you'll probably have to change the line numbers here
				{"level", "debug", "caller", "level_test.go:42", "message", "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummy.reset()

			l := Debug(tt.args.logger)

			if err := l.Log(tt.args.keyvals...); (err != nil) != tt.wantErr {
				t.Errorf("logger.Log() 1 error = %v, wantErr %v", err, tt.wantErr)
			}

			lines := dummy.Lines()
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	dummy := &dummyLogger{}

	type args struct {
		logger  logging.Logger
		keyvals []interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantLines [][]interface{}
		wantErr   bool
	}{
		{
			name: "with level",
			args: args{
				logger:  logging.With(logging.NewFrom(dummy), "caller", logging.DefaultCaller),
				keyvals: []interface{}{"message", "test"},
			},
			wantLines: [][]interface{}{
				// NOTE: When adding code, you'll probably have to change the line numbers here
				{"level", "info", "caller", "level_test.go:86", "message", "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummy.reset()

			l := Info(tt.args.logger)

			if err := l.Log(tt.args.keyvals...); (err != nil) != tt.wantErr {
				t.Errorf("logger.Log() 1 error = %v, wantErr %v", err, tt.wantErr)
			}

			lines := dummy.Lines()
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}

func TestError(t *testing.T) {
	dummy := &dummyLogger{}

	type args struct {
		logger  logging.Logger
		keyvals []interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantLines [][]interface{}
		wantErr   bool
	}{
		{
			name: "with level",
			args: args{
				logger:  logging.With(logging.NewFrom(dummy), "caller", logging.DefaultCaller),
				keyvals: []interface{}{"message", "test"},
			},
			wantLines: [][]interface{}{
				// NOTE: When adding code, you'll probably have to change the line numbers here
				{"level", "error", "caller", "level_test.go:130", "message", "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dummy.reset()

			l := Error(tt.args.logger)

			if err := l.Log(tt.args.keyvals...); (err != nil) != tt.wantErr {
				t.Errorf("logger.Log() 1 error = %v, wantErr %v", err, tt.wantErr)
			}

			lines := dummy.Lines()
			if !reflect.DeepEqual(lines, tt.wantLines) {
				t.Errorf("logger.Log() output = %v, want %v", lines, tt.wantLines)
			}
		})
	}
}
