package laniakea

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

type LoggerWriter interface {
	Close() error
	Write(p []byte) (n int, err error)
	Print(level LogLevel, prefix string, traceback []*MethodTraceback, messages ...any) error
}
type LoggerTextWriter struct {
	LoggerWriter
	writer         io.Writer
	printTraceback bool
	printTime      bool
}

func (w *LoggerTextWriter) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	if err != nil {
		return n, err
	}
	err = bufio.NewWriter(w.writer).Flush()
	return n, err
}
func (w *LoggerTextWriter) Print(level LogLevel, prefix string, tb []*MethodTraceback, messages ...any) error {
	s := buildString(level, prefix, true, w.printTraceback, w.printTime, messages...)
	_, err := w.Write([]byte(s))
	return err
}
func (w *LoggerTextWriter) Close() error {
	return w.writer.(io.Closer).Close()
}

type LoggerJsonWriter struct {
	LoggerWriter
	writer io.Writer
	pretty bool
}
type LoggerJsonMessage struct {
	Time      time.Time          `json:"time"`
	Level     string             `json:"level"`
	Prefix    string             `json:"prefix"`
	Message   string             `json:"message"`
	Traceback []*MethodTraceback `json:"traceback"`
}

func (w *LoggerJsonWriter) Write(data []byte) (int, error) {
	n, err := w.writer.Write(data)
	if err != nil {
		return n, err
	}
	err = bufio.NewWriter(w.writer).Flush()
	return n, err
}
func (w *LoggerJsonWriter) Print(level LogLevel, prefix string, traceback []*MethodTraceback, messages ...any) error {
	msg := Map(messages, func(el any) string {
		return fmt.Sprintf("%v", el)
	})
	m := LoggerJsonMessage{
		Time:      time.Now(),
		Level:     level.GetName(),
		Prefix:    prefix,
		Message:   strings.Join(msg, " "),
		Traceback: traceback,
	}
	var data []byte
	var err error
	if w.pretty {
		data, err = json.MarshalIndent(m, "", "  ")
	} else {
		data, err = json.Marshal(m)
	}
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, []byte("\n")...))
	return err
}
func (w *LoggerJsonWriter) Close() error {
	return w.writer.(io.Closer).Close()
}

type Logger struct {
	prefix         string
	level          LogLevel
	printTraceback bool
	printTime      bool
	jsonPretty     bool
	writers        []LoggerWriter
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
	Package   string `json:"package"`
	Method    string `json:"method"`
	fullPath  string
	signature string
	Filename  string `json:"filename"`
	Line      int    `json:"line"`
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

func (l *Logger) CreateTextFileWriter(path string) (*LoggerTextWriter, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	writer := &LoggerTextWriter{
		writer: file, printTraceback: l.printTraceback, printTime: l.printTime,
	}
	return writer, nil
}
func (l *Logger) CreateTextStdoutWriter() *LoggerTextWriter {
	writer := &LoggerTextWriter{
		writer: os.Stdout, printTraceback: l.printTraceback, printTime: l.printTime,
	}
	return writer
}
func (l *Logger) CreateJsonStdoutWriter() *LoggerJsonWriter {
	writer := &LoggerJsonWriter{
		writer: os.Stdout, pretty: l.jsonPretty,
	}
	return writer
}
func (l *Logger) CreateJsonFileWriter(path string) (*LoggerJsonWriter, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	writer := &LoggerJsonWriter{
		writer: file, pretty: l.jsonPretty,
	}
	return writer, nil
}
func (l *Logger) CreateTextWriter(w io.Writer) *LoggerTextWriter {
	writer := &LoggerTextWriter{
		writer: w, printTraceback: l.printTraceback, printTime: l.printTime,
	}
	return writer
}
func (l *Logger) CreateJsonWriter(w io.Writer) *LoggerJsonWriter {
	writer := &LoggerJsonWriter{
		writer: w, pretty: l.jsonPretty,
	}
	return writer
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
func (l *Logger) JsonPretty(b bool) *Logger {
	l.jsonPretty = b
	return l
}
func (l *Logger) AddWriters(writers ...LoggerWriter) *Logger {
	l.writers = append(l.writers, writers...)
	return l
}
func (l *Logger) AddWriter(writer LoggerWriter) *Logger {
	l.writers = append(l.writers, writer)
	return l
}

func (l *Logger) Infof(format string, args ...any) {
	l.print(INFO, fmt.Sprintf(format, args...))
}
func (l *Logger) Info(m ...any) {
	l.println(INFO, m...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.print(WARN, fmt.Sprintf(format, args...))
}
func (l *Logger) Warn(m ...any) {
	l.println(WARN, m...)
}

func (l *Logger) Error(m ...any) {
	l.println(ERROR, m...)
}

func (l *Logger) Fatal(m ...any) {
	l.println(FATAL, m...)
	os.Exit(1)
}

func (l *Logger) Debug(m ...any) {
	l.println(DEBUG, m...)
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%02d.%02d.%02d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
}
func formatTraceback(mt *MethodTraceback) string {
	return fmt.Sprintf("%s:%s:%d", mt.Filename, mt.Method, mt.Line)
}
func FormatFullTraceback(tracebacks []*MethodTraceback) string {
	formatted := make([]string, 0)
	for _, tb := range tracebacks {
		formatted = append(formatted, formatTraceback(tb))
	}
	return strings.Join(formatted, "->")
}

func getTraceback() *MethodTraceback {
	caller, _, _, _ := runtime.Caller(4)
	details := runtime.FuncForPC(caller)
	signature := details.Name()
	path, line := details.FileLine(caller)
	splitPath := strings.Split(path, "/")

	splitSignature := strings.Split(signature, ".")
	pkg, method := splitSignature[0], splitSignature[len(splitSignature)-1]

	tb := &MethodTraceback{
		Filename:  splitPath[len(splitPath)-1],
		fullPath:  path,
		Line:      line,
		signature: signature,
		Package:   pkg,
		Method:    method,
	}

	return tb
}
func getFullTraceback(skip int) []*MethodTraceback {
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
			Filename:  splitPath[len(splitPath)-1],
			fullPath:  path,
			Line:      line,
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

func buildString(level LogLevel, prefix string, newline, printTime, printTraceback bool, m ...any) string {
	args := []string{
		fmt.Sprintf("[%s]", prefix),
		fmt.Sprintf("[%s]", strings.ToUpper(level.t)),
	}

	if printTraceback {
		args = append(args, fmt.Sprintf("[%s]", formatTraceback(getTraceback())))
	}

	if printTime {
		args = append(args, fmt.Sprintf("[%s]", formatTime(time.Now())))
	}

	msg := Map(m, func(el any) string {
		return fmt.Sprintf("%v", el)
	})
	s := fmt.Sprintf("%s %s", strings.Join(args, " "), strings.Join(msg, " "))
	if newline {
		s += "\n"
	}
	return s
}

func (l *Logger) print(level LogLevel, m ...any) {
	if l.level.n < level.n {
		return
	}

	tb := getFullTraceback(0)
	for _, writer := range l.writers {
		err := writer.Print(level, l.prefix, tb, m...)
		if err != nil {
			l.Error(err)
		}
	}
}

// Docker requires "\n" at end to write to log.
// print not work for docker, otherwise it will work and write into stdout
func (l *Logger) println(level LogLevel, m ...any) {
	if l.level.n < level.n {
		return
	}

	tb := getFullTraceback(0)
	for _, writer := range l.writers {
		err := writer.Print(level, l.prefix, tb, m...)
		if err != nil {
			l.Error(err)
		}
	}
}
