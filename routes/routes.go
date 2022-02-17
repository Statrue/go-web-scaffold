package routes

import (
	"github.com/gin-gonic/gin"
	"go-web-scaffold/routes/middle"
	"go.uber.org/zap"
	"net/http"
)

func Setup(mode string) *gin.Engine {
	r := gin.New()

	setEngineMode(mode)

	r.Use(
		middle.Ginzap(zap.L()),
		middle.RecoveryWithZap(zap.L(), true),
	)

	// register routes here
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})

	return r
}

func setEngineMode(mode string) {
	switch mode {
	case "prod", "release":
		gin.SetMode("release")
	case "test":
		gin.SetMode("test")
	case "dev", "debug":
		fallthrough
	default:
		gin.SetMode("debug")
	}
}
