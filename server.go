package main

import (
	"bytes"
	"cy17-starter-go-api/web"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pbnjay/memory"
)

var startTime = time.Now();


func main() {
	e := echo.New()
	web.RegisterHandlers(e)
	RegisterRatesLimiter(e)
	e.GET("/api/hello",getHello)
	e.GET("/actuator/health", getHealth)
	e.GET("/actuator/status", getStatus)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
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
	buffer.WriteString("\n")
	buffer.WriteString("Rates\n")
	buffer.WriteString("~~~~~\n")
	buffer.WriteString("Total active: " + fmt.Sprint(len(Rates)) + "\n")
	for key, value := range Rates {
		buffer.WriteString(fmt.Sprintf("%-15s -> %d\n", key, value))
    }
	buffer.WriteString("\n")
	return c.String(http.StatusOK, buffer.String());
}
