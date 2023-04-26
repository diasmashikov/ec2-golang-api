package main

import (
	"database/sql"
	"log"
	"net/http"
	"syscall"
	"time"

	db "ec2-go-api/db"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type course struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

func main() {
      // Create a new database connection
    db, err := db.NewDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    server := NewServer(db)
    server.Run(":8080")
}

// NewServer creates a new instance of the Go server.
func NewServer(db *sql.DB) *gin.Engine {
    r := gin.Default()

    registry := prometheus.NewRegistry()
    cpuUsageGauge := newCPUUsageGauge()
    counter := newCounter()

    registry.MustRegister(cpuUsageGauge, counter)

    r.GET("/", func(c *gin.Context) {
        handleHello(c, counter)
    })
    r.GET("/getCourses", func(c *gin.Context) {
        handleGetCourses(c, db)
    })

    r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(registry, promhttp.HandlerOpts{})))

    return r
}

// newCPUUsageGauge creates a new gauge vector to track CPU usage.
func newCPUUsageGauge() *prometheus.GaugeVec {
    cpuUsageGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Name: "cpu_usage",
        Help: "Current CPU usage in percentage",
    }, []string{"mode"})

    go func() {
        for {
            var rusage syscall.Rusage
            if err := syscall.Getrusage(syscall.RUSAGE_SELF, &rusage); err == nil {
                cpuUsageGauge.With(prometheus.Labels{"mode": "user"}).Set(float64(rusage.Utime.Nano()) / 1e9)
                cpuUsageGauge.With(prometheus.Labels{"mode": "system"}).Set(float64(rusage.Stime.Nano()) / 1e9)
            }
            time.Sleep(time.Second)
        }
    }()

    return cpuUsageGauge
}

// newCounter creates a new Prometheus counter metric.
func newCounter() prometheus.Counter {
    counter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "my_counter",
        Help: "This is my counter",
    })

    return counter
}

// handleHello increments the Prometheus counter metric and returns a JSON response.
func handleHello(c *gin.Context, counter prometheus.Counter) {
    counter.Inc()
    c.JSON(http.StatusOK, gin.H{
        "message": "Hello, world!",
    })
}

// handleGetCourses retrieves the courses from the database and returns them as a JSON response.
func handleGetCourses(c *gin.Context, db *sql.DB) {

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
}