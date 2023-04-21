package main

import "github.com/gin-gonic/gin"

func main() {
    // Create a new Gin router
    r := gin.Default()

    // Define a route
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, world!",
        })
    })

	r.GET("/getCourses", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Here is your IELTS course",
        })
    })

    // Start the server
    r.Run(":8080")
}
