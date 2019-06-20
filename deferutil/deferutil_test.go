package deferutil

import "testing"

var ran int

func run() error {
	ran++
	return nil
}

func reset() {
	ran = 0
}

func TestCheckDefer(t *testing.T) {
	type args struct {
		fs []func() error
	}
	tests := []struct {
		name   string
		args   args
		wantCt int
	}{
		{
			name: "once",
			args: args{
				fs: []func() error{run},
			},
			wantCt: 1,
		},
		{
			name: "twice",
			args: args{
				fs: []func() error{run, run},
			},
			wantCt: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reset()

			if ran != 0 {
				t.Errorf("wrong starting run count %d", ran)
			}

			CheckDefer(tt.args.fs...)

			if ran != tt.wantCt {
				t.Errorf("wrong ending run count; want %v got %v", tt.wantCt, ran)
			}
		})
	}
}
