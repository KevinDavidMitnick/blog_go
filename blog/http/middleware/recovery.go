package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"

	"github.com/toolkits/web/errors"
	"github.com/ulricqin/blog/http/render"
)

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	Logger     *log.Logger
	PrintStack bool
	StackAll   bool
	StackSize  int
}

// NewRecovery returns a new instance of Recovery
func NewRecovery(out io.Writer) *Recovery {
	return &Recovery{
		Logger:     log.New(out, "", 0),
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {

			if er, ok := err.(errors.Error); ok {

				if isAjax(r) {
					render.Message(w, er.Msg)
					return
				}

				if er.Code == http.StatusUnauthorized || er.Code == http.StatusForbidden {
					http.Redirect(w, r, "/auth/login?callback="+r.URL.String(), 302)
					return
				}

				render.Put(r, "Error", er.Msg)
				render.HTML(r, w, "common/error")
				return
			}

			// Negroni part
			w.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "PANIC: %s\n%s"
			rec.Logger.Printf(f, err, stack)

			if rec.PrintStack {
				fmt.Fprintf(w, f, err, stack)
			}
		}
	}()

	next(w, r)
}

func isAjax(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}
