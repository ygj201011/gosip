package log

import (
	"regexp"
	"strconv"

	"github.com/ghettovoice/gosip/runtime"
	"github.com/ghettovoice/logrus"
)

const UndefStack = "???"
const (
	stackNumFrames = 20
	// constant value of stack offset from logrus.Logger.* fn call to current hook Fire call
	hookStackDelta = 5
)

// CallInfoHook is an hook for logrus logger that adds file, line, function info.
type CallInfoHook struct {
}

// NewCallInfoHook creates new `CallInfoHook`.
func NewCallInfoHook() *CallInfoHook {
	return &CallInfoHook{}
}

// Fire is an callback that will be called by logrus for each log entry.
func (hook *CallInfoHook) Fire(entry *logrus.Entry) error {
	file, line, fn := GetStackInfo()

	entry.SetField("file", file)
	entry.SetField("line", line)
	entry.SetField("func", fn)

	return nil
}

// Levels returns `CallInfoHook` working levels.
func (hook *CallInfoHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// GetStackInfo returns file, line and function name from runtime call stack.
func GetStackInfo() (string, string, string) {
	// Get information about the stack.
	// Try and find the first stack frame outside the logging package.
	// Only search up a few frames, it should never be very far.
	file := UndefStack
	line := UndefStack
	fn := UndefStack

	for depth := hookStackDelta; depth < stackNumFrames+hookStackDelta; depth++ {
		if frame, ok := runtime.GetFrameOffset(depth); ok {
			fnName := frame.Func.Name()
			if isLog, _ := regexp.MatchString(`(log\w*\..*)`, fnName); isLog {
				continue
			}

			file = frame.File
			line = strconv.Itoa(frame.Line)
			fn = fnName
			break
		}
		// If we get here, we failed to retrieve the stack information.
		// Just give up.
		break
	}

	return file, line, fn
}
