package pkglog

import (
	"bytes"
	"io"
	"os"
)

// Outputter interface implements a custom outputter.
type Outputter interface {
	Output(*Entry)
}

type DefaultOutputter struct {
	DisableTimestamps bool
	Writer            io.Writer
}

func (d *DefaultOutputter) Output(e *Entry) {
	w := d.Writer
	if w == nil {
		w = os.Stdout
	}

	b := &bytes.Buffer{}
	if !d.DisableTimestamps {
		b.WriteString(e.Time.Format("2006/01/02 15:04:05 "))
	}
	b.WriteString(e.Level.String())
	b.WriteString(": ")
	b.WriteString(e.Message)
	b.WriteByte('\n')
	w.Write(b.Bytes())
}
