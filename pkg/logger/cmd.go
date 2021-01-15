package logger

import (
	"fmt"
	"strings"
)

func scrubSecrets(message string, secrets []string) string {
	scrub := "*****"
	for _, secret := range secrets {
		message = strings.ReplaceAll(message, secret, scrub)
	}
	return message
}

// AddSecret add secret to be scrubbed from logging
func (l CMD) AddSecret(secret string) Logger {
	l.secrets = append(l.secrets, secret)
	return l
}

// ClearSecrets clear secrets
func (l CMD) ClearSecrets() Logger {
	l.secrets = []string{}
	return l
}

// Writef to stdout
func (l CMD) Writef(format string, a ...interface{}) {
	if l.Enabled() {
		message := fmt.Sprintf(format, a...)
		message = scrubSecrets(message, l.secrets)
		fmt.Printf("%s", message)
	}
}

// Writeln to stdout
func (l CMD) Writeln(message string) {
	if l.Enabled() {
		message = scrubSecrets(message, l.secrets)
		fmt.Println(message)
	}
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
