package logger

// Logger interface should be used to log to different locations
type Logger interface {
	Writef(string, ...interface{})
	Writeln(string)
	Enabled() bool
	LogHeader() bool
	LogBody() bool
	AddSecret(string) Logger
	ClearSecrets() Logger
}

// CMD log to the stdout
type CMD struct {
	LoggerEnabled    bool
	LogHeaderEnabled bool
	LogBodyEnabled   bool
	secrets          []string
}

// File log to a file
type File struct {
	Path string
}
