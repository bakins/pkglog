package foo

import (
	"log"

	"github.com/bakins/pkglog"
)

func Log() {
	pkglog.Printf("hello from foo")
}

func StandardLog() {
	log.Printf("hello from foo standardlog")
}
