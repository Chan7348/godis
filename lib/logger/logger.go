package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Settings struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Ext        string `yaml:"ext"`
	TimeFormat string `yaml:"time-format`
}

type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const (
	flags              = log.LstdFlags
	defaultCallerDepth = 2
	bufferSize         = 1e5
)

type logEntry struct {
	msg   string
	level logLevel
}

var (
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

type Logger struct {
	logFile   *os.File
	logger    *log.Logger
	entryChan chan *logEntry
	entryPool *sync.Pool
}

var DefaultLogger = NewStdoutLogger()

func NewStdoutLogger() *Logger {
	logger := &Logger{
		logFile:   nil,
		logger:    log.New(os.Stdout, "", flags),
		entryChan: make(chan *logEntry, bufferSize),
		entryPool: &sync.Pool{
			New: func() interface{} {
				return &logEntry{}
			},
		},
	}

	go func() {
		for e := range logger.entryChan {
			_ = logger.logger.Output(0, e.msg)
			logger.entryPool.Put(e)
		}
	}()

	return logger
}

func NewFileLogger(settings *Settings) (*Logger, error) {
	fileName := fmt.Sprintf(
		"%s-%s.%s",
		settings.Name,
		time.Now().Format(settings.TimeFormat),
		settings.Ext,
	)

	logFile, err := mustOpen(fileName, settings.Path)
	if err != nil {
		return nil, fmt.Errorf("Logging.join err: %s", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := &Logger{
		logFile:   logFile,
		logger:    log.New(multiWriter, "", flags),
		entryChan: make(chan *logEntry, bufferSize),
		entryPool: &sync.Pool{
			New: func() interface{} {
				return &logEntry{}
			},
		},
	}

	go func() {
		for e := range logger.entryChan {
			logFilename := fmt.Sprintf(
				"%s-%s.%s",
				settings.Name,
				time.Now().Format(settings.TimeFormat),
				settings.Ext,
			)
			if path.Join(settings.Path, logFilename) != logger.logFile.Name() {
				logFile, err := mustOpen(logFilename, settings.Path)
				if err != nil {
					panic("open log " + logFilename + " failed: " + err.Error())
				}

				logger.logFile = logFile
				logger.logger = log.New(io.MultiWriter(os.Stdout, logFile), "", flags)
			}
			_ = logger.logger.Output(0, e.msg) // msg includes call stack, no need for call depth
			logger.entryPool.Put(e)
		}
	}()
	return logger, nil
}

func Setup(settings *Settings) {
	logger, err := NewFileLogger(settings)
	if err != nil {
		panic(err)
	}
	DefaultLogger = logger
}

func (logger *Logger) Output(level logLevel, callerDepth int, msg string) {
	var formattedMsg string

	_, file, line, ok := runtime.Caller(callerDepth)
	if ok {
		formattedMsg = fmt.Sprintf("[%s][%s:%d] %s", levelFlags[level], filepath.Base(file), line, msg)
	} else {
		formattedMsg = fmt.Sprintf("[%s] %s", levelFlags[level], msg)
	}

	entry := logger.entryPool.Get().(*logEntry)
	entry.msg = formattedMsg
	entry.level = level
	logger.entryChan <- entry
}

func Debug(format string, v ...interface{}) {
	DefaultLogger.Output(DEBUG, defaultCallerDepth, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	DefaultLogger.Output(INFO, defaultCallerDepth, fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Output(INFO, defaultCallerDepth, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	DefaultLogger.Output(WARNING, defaultCallerDepth, fmt.Sprintln(v...))
}

func Error(v ...interface{}) {
	DefaultLogger.Output(ERROR, defaultCallerDepth, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Output(ERROR, defaultCallerDepth, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	DefaultLogger.Output(FATAL, defaultCallerDepth, fmt.Sprintln(v...))
}
