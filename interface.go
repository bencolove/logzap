package logzap

import "fmt"

func (l *Logger) Printf(template string, values ...any) {
	msg := template
	if len(values) > 0 {
		msg = fmt.Sprintf(template, values...)
	}
	l.l.Info(msg)
}
