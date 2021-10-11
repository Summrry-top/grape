package main

import (
	"grape"
	"log"
)

type Data struct {
	Code int    `form:"code"`
	Msg  string `form:"msg"`
}

func test(c *grape.Context) {
	log.Println("test")
	var d Data
	err := c.ShouldBind(&d)
	if err != nil {
		log.Println(err)
		return
	}
	c.Json(d)
}

func main() {
	//r := grape.New()
	//v1:=r.Group("/v1")
	//v1.Use(grape.Logger(),grape.Recovery())
	//v1.Use(grape.cors())
	//v1.GET("/test", func(c *grape.Context) {
	//	log.Println("test")
	//	c.String("v1OKm")
	//})

	r := grape.Default()
	{
		r.GET("/test", test)
		r.POST("/test", test)
		r.GET("/panic", func(c *grape.Context) {
			panic("panic")
		})
	}

	v1 := r.Group("/v1")
	{
		v1.GET("/test", test)
		v1.POST("/test", test)
	}

	v2 := v1.Group("/v2")
	v2.GET("/test", test)

	r.Run(":8080")

}
