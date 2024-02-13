package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	limiter "github.com/danluki/uddug-test"
	"github.com/gin-gonic/gin"
)

var (
	port = flag.String("port", "8080", "port")
)

func LimitOnRps() gin.HandlerFunc {
	//5 Request per day based on client ip
	flowRate := float64(5) / float64(24*60*60)
	l := limiter.NewRateLimiter(flowRate, 5, time.Hour*24)

	return func(ctx *gin.Context) {
		ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !l.GetVisitor(ip).Allow() {
			ctx.AbortWithStatus(http.StatusTooManyRequests)

			return
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/ping", LimitOnRps(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Run(fmt.Sprintf(":%s", *port))
}
