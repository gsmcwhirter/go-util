package level

type stringer interface {
	String() string
}

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

func (l *dummyLogger) Lines() [][]interface{} {
	lines := make([][]interface{}, 0, len(l.lines))
	for _, line := range l.lines {
		newln := make([]interface{}, 0, len(line))
		for _, item := range line {
			if i, ok := item.(stringer); ok {
				newln = append(newln, i.String())
			} else {
				newln = append(newln, item)
			}
		}

		lines = append(lines, newln)
	}

	return lines
}
