package main

import (
	"bytes"
	"cy17-starter-go-api/web"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pbnjay/memory"
)

var startTime = time.Now();
var statics = [...]string{".html", ".map", ".js", ".css", ".json", ".jpg", ".png", ".gif", ".ttf", ".svg", ".ico"}
var debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

func main() {
	e := echo.New()
	web.RegisterHandlers(e)
	RegisterRatesLimiter(e)
	e.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
	e.Use(CacheHandler)
	e.GET("/api/hello",getHello)
	e.GET("/actuator/health", getHealth)
	e.GET("/actuator/status", getStatus)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	if debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.ERROR)
	}
	e.Logger.Fatal(e.Start(":" + httpPort))
}

func CacheHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		if path != "/" && path != "/index.html" {
			for _,v := range statics {
				if strings.HasSuffix(path, v) {
					c.Response().Header().Add("last-modified",startTime.Format(time.RFC822))
				}
			}
		}
		
		return next(c)
	}
}

func getHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func getHealth(c echo.Context) error {
  	return c.JSON(http.StatusOK, struct{ Status string }{Status: "UP"})
}

func getStatus(c echo.Context) error {
	apiKey := os.Getenv("SECURITY_API_KEY")
	if (apiKey == "") {
		return c.String(http.StatusServiceUnavailable, "")
	}
	requestApiKey := c.Request().Header.Get("X-API-KEY")
	if (requestApiKey == "") {
		return c.String(http.StatusUnauthorized, "")
	}
	if (apiKey != requestApiKey) {
		return c.String(http.StatusForbidden, "")
	}
	var buffer bytes.Buffer
	buffer.WriteString("Go\n")
	buffer.WriteString("~~~~~\n")
	buffer.WriteString("Version: " + runtime.Version() + " " +runtime.GOOS + " " + runtime.GOARCH + "\n")
	buffer.WriteString("Free memory: " + fmt.Sprint(memory.FreeMemory()) + "\n")
	buffer.WriteString("Total memory: " + fmt.Sprint(memory.FreeMemory()) + "\n")
	buffer.WriteString("Available processors: " + fmt.Sprint(runtime.NumCPU()) + "\n")
	buffer.WriteString("\n")
	buffer.WriteString("App\n")
	buffer.WriteString("~~~~~\n")
	buffer.WriteString("Server time: " + time.Now().Format(time.RFC3339Nano) + "\n")
	buffer.WriteString("Uptime: " + time.Since(startTime).String() + "\n")
	buffer.WriteString("Max simultaneous request: " + fmt.Sprint(runtime.GOMAXPROCS(-1)) + "\n")
	buffer.WriteString("\n")
	buffer.WriteString("Rates\n")
	buffer.WriteString("~~~~~\n")
	Lock.RLock()
	defer Lock.RUnlock()
	buffer.WriteString("Total active: " + fmt.Sprint(len(Rates)) + "\n")
	for key, value := range Rates {
		buffer.WriteString(fmt.Sprintf("%-15s -> %d\n", key, value))
    }
	buffer.WriteString("\n")
	return c.String(http.StatusOK, buffer.String());
}
