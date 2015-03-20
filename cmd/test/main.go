package main

import (
	stdlog "log"

	log "github.com/bakins/pkglog"
	"github.com/bakins/pkglog/cmd/test/foo"
)

func main() {

	w := log.StandardLogger().Writer()
	stdlog.SetOutput(w)

	log.Printf("hello world")
	log.SetLogLevel(log.DebugLevel)
	log.Printf("hello world")

	log.SetLogLevel(log.WarnLevel)
	foo.Log()

	log.SetLogLevel(log.DebugLevel)

	foo.Log()

	log.SetPackageLogLevel("github.com/bakins/pkglog/cmd/test/foo", log.PanicLevel)

	foo.Log()

	log.SetPackageLogLevel("main", log.PanicLevel)
	log.Printf("hello world")

	log.SetLogLevel(log.DebugLevel)

	foo.StandardLog()

}
