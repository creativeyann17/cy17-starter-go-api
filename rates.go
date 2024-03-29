package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Rates = make(map[string]int)
var Lock = sync.RWMutex{}

var config = middleware.RateLimiterConfig{
    Skipper: func(c echo.Context) bool {
		return !strings.HasPrefix(c.Request().URL.Path, "/api")
	},
    Store: middleware.NewRateLimiterMemoryStoreWithConfig(
        middleware.RateLimiterMemoryStoreConfig{Rate: 0.16, Burst: 10, ExpiresIn: 60 * time.Second},
    ),
    IdentifierExtractor: func(ctx echo.Context) (string, error) {
        id := ctx.RealIP()
		updateRate(id)
        return id, nil
    },
    ErrorHandler: func(context echo.Context, err error) error {
        return context.String(http.StatusForbidden, "")
    },
    DenyHandler: func(context echo.Context, identifier string,err error) error {
        return context.String(http.StatusTooManyRequests, "")
    },
}

func RegisterRatesLimiter(e *echo.Echo) {
	e.Use(middleware.RateLimiterWithConfig(config))
}

func updateRate(id string) {
	Lock.Lock()
	defer Lock.Unlock()
	if v, ok := Rates[id]; ok {
		rate := v-1
		if rate <0 {
			Rates[id] = 0
		} else {
			Rates[id] = rate
		}
	} else {
		Rates[id] = 10
	}
}
