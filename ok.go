package main

import (
	"bytes"
	"gopkg.in/gin-gonic/gin.v1"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var twoHundredMegofA = bytes.Repeat([]byte("a"), 1024*1024*200)

func CommonHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", "OKGO")
	}
}

func Respond(c *gin.Context) {
	statusValue := c.Query("status")
	status, err := strconv.Atoi(statusValue)
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse given status.\nYou provided: %s\n", statusValue)
	}
	if status > 299 && status < 400 {
		c.Redirect(status, "/nowhere")
		return
	}
	c.Status(status)
}

func Chunked(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain; charset=UTF-8", twoHundredMegofA)
}

func Strict(c *gin.Context) {
	c.Header("Content-Length", "209715200")
	c.Data(http.StatusOK, "text/plain; charset=UTF-8", twoHundredMegofA)
}

func Success(c *gin.Context) {
	io.Copy(ioutil.Discard, c.Request.Body)
}

func Redirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/nowhere")
}

func Upload(c *gin.Context) {
	expectedSize := c.Query("size")
	expectedBytes, _ := strconv.ParseInt(expectedSize, 10, 64)
	totBytes, _ := io.Copy(ioutil.Discard, c.Request.Body)
	if expectedBytes == totBytes {
		c.String(http.StatusOK, "Hey Buddy,\nEverything is good, expected %d bytes and I got %d bytes.\n", expectedBytes, totBytes)
	} else {
		c.String(http.StatusBadRequest, "Expected %d bytes, got %d bytes instead.\n", expectedBytes, totBytes)
	}
}

func Delay(c *gin.Context) {
	duration := c.Query("duration")
	d, err := time.ParseDuration(duration)
	if err != nil {
		c.String(http.StatusBadRequest, "Could not parse given duration. Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.\nYou provided: %s\n", duration)
		return
	}
	time.Sleep(d)
}

func main() {
	logFile, _ := os.Create("/tmp/okgo.log")
	router := gin.New()
	router.Use(gin.LoggerWithWriter(logFile))
	router.Use(CommonHeaders())
	router.Use(gin.Recovery())

	router.GET("/success/:testNum", Success)
	router.POST("/success/:testNum", Success)

	router.GET("/respond/:testNum", Respond)
	router.POST("/respond/:testNum", Respond)
	router.PUT("/respond/:testNum", Respond)
	router.HEAD("/respond/:testNum", Respond)
	router.DELETE("/respond/:testNum", Respond)
	router.PATCH("/respond/:testNum", Respond)
	router.OPTIONS("/respond/:testNum", Respond)

	router.GET("/chunkedString/:testNum", Chunked)
	router.GET("/strictString/:testNum", Strict)

	router.GET("/redirect/:testNum", Redirect)
	router.POST("/redirect/:testNum", Redirect)

	router.POST("/upload/:testNum", Upload)

	router.GET("/delay/:testNum", Delay)

	router.Run(":50001")
}
