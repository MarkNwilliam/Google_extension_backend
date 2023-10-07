package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello, Gin!")
    })

    r.Run() // By default, it serves on :8080 unless a PORT environment variable is defined.
}
