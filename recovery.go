package grape

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				//message := fmt.Sprintf("%s", err)
				//log.Printf("%s\n\n", trace(message))
				log.Println("发生恐慌")
				c.String("Internal Server Error")
				c.Abort()
				return
			}
		}()
		log.Printf("recovery")
		c.Next()
	}
}

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
