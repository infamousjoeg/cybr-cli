package logger

// Logger interface should be used to log to different locations
type Logger interface {
	Writef(string, ...interface{})
	Writeln(...interface{})
	Enabled() bool
	LogHeader() bool
	LogBody() bool
}

// CMD log to the stdout
type CMD struct {
	LoggerEnabled    bool
	LogHeaderEnabled bool
	LogBodyEnabled   bool
}

// File log to a file
type File struct {
	Path string
}
