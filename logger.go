package laniakea

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

type LoggerWriter func(level LogLevel, prefix, traceback string, message []any)

type Logger struct {
	prefix         string
	level          LogLevel
	printTraceback bool
	printTime      bool
	writers        []LoggerWriter

	f *os.File
}

type LogLevel struct {
	n uint8
	t string
	c color.Attribute
}

func (l *LogLevel) GetName() string {
	return l.t
}

type MethodTraceback struct {
	Package   string
	Method    string
	fullPath  string
	signature string
	filename  string
	line      int
}

var (
	INFO  = LogLevel{n: 0, t: "info", c: color.FgWhite}
	WARN  = LogLevel{n: 1, t: "warn", c: color.FgHiYellow}
	ERROR = LogLevel{n: 2, t: "error", c: color.FgHiRed}
	FATAL = LogLevel{n: 3, t: "fatal", c: color.FgRed}
	DEBUG = LogLevel{n: 4, t: "debug", c: color.FgGreen}
)

func CreateLogger() *Logger {
	return &Logger{
		prefix:         "LOG",
		level:          FATAL,
		printTraceback: false,
		printTime:      true,
	}
}

func (l *Logger) OpenFile(name string) *Logger {
	err := os.MkdirAll(filepath.Dir(name), os.ModePerm)
	if err != nil {
		l.Fatal(err)
	}
	l.f, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.Fatal(err)
	}
	return l
}

func (l *Logger) Prefix(prefix string) *Logger {
	l.prefix = prefix
	return l
}
func (l *Logger) Level(level LogLevel) *Logger {
	l.level = level
	return l
}
func (l *Logger) PrintTraceback(b bool) *Logger {
	l.printTraceback = b
	return l
}
func (l *Logger) PrintTime(b bool) *Logger {
	l.printTime = b
	return l
}
func (l *Logger) AddWriters(writers []LoggerWriter) *Logger {
	l.writers = append(l.writers, writers...)
	return l
}

func (l *Logger) Info(m ...any) {
	l.print(INFO, m)
}

func (l *Logger) Warn(m ...any) {
	l.print(WARN, m)
}

func (l *Logger) Error(m ...any) {
	l.print(ERROR, m)
}

func (l *Logger) Fatal(m ...any) {
	l.print(FATAL, m)
	os.Exit(1)
}

func (l *Logger) Debug(m ...any) {
	l.print(DEBUG, m)
}

func (l *Logger) formatTime(t time.Time) string {
	return fmt.Sprintf("%02d.%02d.%02d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
}

func (l *Logger) getTraceback() *MethodTraceback {
	caller, _, _, _ := runtime.Caller(4)
	details := runtime.FuncForPC(caller)
	signature := details.Name()
	path, line := details.FileLine(caller)
	splitPath := strings.Split(path, "/")

	splitSignature := strings.Split(signature, ".")
	pkg, method := splitSignature[0], splitSignature[len(splitSignature)-1]

	tb := &MethodTraceback{
		filename:  splitPath[len(splitPath)-1],
		fullPath:  path,
		line:      line,
		signature: signature,
		Package:   pkg,
		Method:    method,
	}

	return tb
}
func (l *Logger) formatTraceback(mt *MethodTraceback) string {
	return fmt.Sprintf("%s:%s:%d", mt.filename, mt.Method, mt.line)
}

func (l *Logger) getFullTraceback(skip int) []*MethodTraceback {
	pc := make([]uintptr, 15)
	runtime.Callers(skip, pc)
	list := make([]*MethodTraceback, 0)
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		details := runtime.FuncForPC(frame.PC)
		signature := details.Name()
		path, line := details.FileLine(frame.PC)
		splitPath := strings.Split(path, "/")

		splitSignature := strings.Split(signature, ".")
		pkg, method := splitSignature[0], splitSignature[len(splitSignature)-1]

		tb := &MethodTraceback{
			filename:  splitPath[len(splitPath)-1],
			fullPath:  path,
			line:      line,
			signature: signature,
			Package:   pkg,
			Method:    method,
		}
		list = append(list, tb)
	}
	sort.Slice(list, func(i, j int) bool {
		return j < i
	})
	return list
}
func (l *Logger) formatFullTraceback(tracebacks []*MethodTraceback) string {
	formatted := make([]string, 0)
	for _, tb := range tracebacks {
		formatted = append(formatted, l.formatTraceback(tb))
	}
	return strings.Join(formatted, "->")
}

func (l *Logger) buildString(level LogLevel, m []any) string {
	args := []string{
		fmt.Sprintf("[%s]", l.prefix),
		fmt.Sprintf("[%s]", strings.ToUpper(level.t)),
	}

	if l.printTraceback {
		args = append(args, fmt.Sprintf("[%s]", l.formatTraceback(l.getTraceback())))
	}

	if l.printTime {
		args = append(args, fmt.Sprintf("[%s]", l.formatTime(time.Now())))
	}

	msg := Map(m, func(el any) string {
		return fmt.Sprintf("%v", el)
	})

	return fmt.Sprintf("%s %v", strings.Join(args, " "), strings.Join(msg, " "))
}

func (l *Logger) print(level LogLevel, m []any) {
	if l.level.n < level.n {
		return
	}
	_, err := color.New(level.c).Println(l.buildString(level, m))
	if err != nil {
		l.Fatal(err)
		return
	}

	for _, writer := range l.writers {
		writer(level, l.prefix, l.formatFullTraceback(l.getFullTraceback(4)), m)
	}

	if l.f != nil {
		writeToFiles := os.Getenv("WRITE_TO_FILE")
		if writeToFiles != "false" {
			if _, err := l.f.Write([]byte(l.buildString(level, m) + "\n")); err != nil {
				l.Fatal(err)
			}
		}
	}
}
