package grape

const (
	AllowCredentials = "true"
	AllowMethods     = "POST, GET"
	AllowHeaders     = "Origin, X-Requested-With, Content-Type, Accept, Authorization"
	ExposeHeaders    = "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type"
)

func Cors() HandlerFunc {
	return func(c *Context) {
		c.Header("Access-Control-Allow-Credentials", AllowCredentials)
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Methods", AllowMethods)
		c.Header("Access-Control-Allow-Headers", AllowHeaders)
		c.Header("Access-Control-Expose-Headers", ExposeHeaders)
		c.Next()
	}
}
