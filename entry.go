package pkglog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

// Entry is a single log entry.
type Entry struct {
	Line     int
	File     string
	pkgPath  string
	Level    Level
	Message  string
	Function string
	logger   *Logger
	Time     time.Time
}

// newEntry is an internal helper for creating an entry
func (l *Logger) newEntry(calldepth int) *Entry {
	e := &Entry{
		Level:  l.level,
		logger: l,
		Time:   time.Now(),
	}
	pc, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		return e
	}
	f := runtime.FuncForPC(pc)
	name := f.Name()
	info := strings.Split(name, ".")
	pkgPath, _ := info[0], info[1]
	if pkgPath != "main" {
		pkgPath, _ = path.Split(file)
	}
	e.Function = name
	e.File = file
	e.Line = line
	e.pkgPath = pkgPath
	e.setLogLevel()
	return e
}

// setLogLevel is an internal helper to set loglevel based on caller's package.
func (e *Entry) setLogLevel() {
	// this is absolutely dreadful

	logger := e.logger
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if ll, ok := logger.cache[e.pkgPath]; ok {
		// we have seen this before so no need for the horrors below
		e.Level = ll
		return
	}

	level := logger.level

	for p, ll := range logger.packages {
		// I apologize in advance
		//strip trailing /
		strings.TrimSuffix(e.pkgPath, "/")

		// this is broken and for quick demo purposes only
		// ie "main" will match "/foo/hello/world/fakemain
		// issue is finding the "Real" package path is hard, especially with vendoring, etc
		// perhaps a trie with longest match?
		if strings.HasSuffix(e.pkgPath, p) {
			// we found it. first match wins
			logger.cache[e.pkgPath] = ll
			level = ll
			break
		}
	}

	e.Level = level
}

// Output emits a single entry using the outputter.
func (e *Entry) Output(level Level, format string, v ...interface{}) {
	// any reason for this to be exported?
	if v != nil {
		e.Message = fmt.Sprintf(format, v)
	} else {
		e.Message = fmt.Sprintf(format)
	}
	// level at which we are logging. should this be a seperate field??
	e.Level = level
	e.logger.output.Output(e)
}
