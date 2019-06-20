package logging

type dummyLogger struct {
	lines [][]interface{}
}

func (l *dummyLogger) Log(keyvals ...interface{}) error {
	l.lines = append(l.lines, keyvals)
	return nil
}

func (l *dummyLogger) reset() {
	l.lines = nil
}
