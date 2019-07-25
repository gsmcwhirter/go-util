package logging

import "strings"

type writer struct {
	l Logger
}

func (w writer) Write(d []byte) (int, error) {
	n := len(d)

	dstr := strings.TrimRight(string(d), "\n\r")
	if err := w.l.Log("message", dstr); err != nil {
		return 0, err
	}

	return n, nil
}
