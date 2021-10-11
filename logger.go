package grape

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("%s in %v\n", c.Request.RequestURI, time.Since(t))
	}
}
