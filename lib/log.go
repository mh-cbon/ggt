package lib

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPLogger handle logs.
type HTTPLogger interface {
	Handle(http.ResponseWriter, *http.Request, error, ...string)
}

// VoidLog doesn ot log
type VoidLog struct{}

// Handle does not handle log.
func (l *VoidLog) Handle(w http.ResponseWriter, r *http.Request, err error, subjects ...string) {}

// WriteLog to a sink writer
type WriteLog struct {
	Sink io.Writer
}

// Handle does not handle log.
func (l *WriteLog) Handle(w http.ResponseWriter, r *http.Request, err error, subjects ...string) {
	fmt.Fprintf(l.Sink, "%v %v %v\n", time.Now(), subjects, err)
}
