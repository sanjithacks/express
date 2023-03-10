package express

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}


func Recovery(errorHandler func (err interface{}, c *Context)) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				if errorHandler != nil {
					errorHandler(err, c)
				} else {
					if httpError, ok := err.(HttpError); ok {
						c.Fail(httpError.Status, httpError.Error())
					} else {
						c.Fail(http.StatusInternalServerError, "Internal Server Error")
					}
				}
			}
		}()

		c.Next()
	}
}
