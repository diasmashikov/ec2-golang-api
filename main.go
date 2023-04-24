package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type course struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}


func main() {
    // Create a new Gin router
    r := gin.Default()

    // Create a new database connection
    db, err := sql.Open("postgres", "postgres://ec2admin:ec2mylove@localhost:5432/courses?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Test the database connection
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    // Define a route
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, world!",
        })
    })

	r.GET("/getCourses", func(c *gin.Context) {
        // Execute a SQL query and return the result as JSON
        rows, err := db.Query("SELECT * FROM course_details")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        courses := []course{}
        for rows.Next() {
            var course course
            err := rows.Scan(&course.ID, &course.Name, &course.Description)
            if err != nil {
                log.Fatal(err)
            }
            courses = append(courses, course)
        }

        c.JSON(http.StatusOK, courses)
    })

    // Start the server
    r.Run(":8080")
}
