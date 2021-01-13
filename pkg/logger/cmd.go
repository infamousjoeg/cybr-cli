package logger

import "fmt"

// Writef to stdout
func (l CMD) Writef(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// Writeln to stdout
func (l CMD) Writeln(a ...interface{}) {
	fmt.Println(a...)
}

// Enabled to log to stdout
func (l CMD) Enabled() bool {
	return l.LoggerEnabled
}

// LogHeader to log to stdout
func (l CMD) LogHeader() bool {
	return l.LogHeaderEnabled
}

// LogBody to log to stdout
func (l CMD) LogBody() bool {
	return l.LogBodyEnabled
}
