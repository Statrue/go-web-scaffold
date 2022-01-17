package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-web-scaffold/routes/middle"
	"go.uber.org/zap"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()

	switch viper.GetString("app.mode") {
	case "prod", "release":
		gin.SetMode("release")
	case "test":
		gin.SetMode("test")
	case "dev", "debug":
		fallthrough
	default:
		gin.SetMode("debug")
	}

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
