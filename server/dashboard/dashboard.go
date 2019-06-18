package dashboard

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/goomadao/serverstatus/assets/statik"
	"github.com/goomadao/serverstatus/server/status"
	"github.com/goomadao/serverstatus/util/logger"
	"github.com/rakyll/statik/fs"

	"go.uber.org/zap"
)

type postData struct {
	From, UUID, Action string
}

var (
	Port int
)

func handlePost(c *gin.Context) bool {
	result, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Logger.Error("Parse data from post body failed",
			zap.Error(err))
		return false
	}
	defer c.Request.Body.Close()
	var res = new(postData)
	err = json.Unmarshal(result, res)
	if err != nil {
		logger.Logger.Error("Turn json data from post body to struct fail",
			zap.Error(err))
		return false
	}
	if status.MoveServer(res.From, res.UUID, res.Action) {
		return true
	} else {
		return false
	}
}

func Dashboard() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(ginzap.Ginzap(logger.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.Logger, true))
	r.GET("/api/servers", func(c *gin.Context) {
		if c.ClientIP() == "127.0.0.1" {
			servers := new(status.Servers)
			status.GetServers(servers)
			c.JSON(http.StatusOK, servers)
		} else {
			c.String(http.StatusForbidden, "Forbidden")
		}
	})
	r.POST("/api/servers", func(c *gin.Context) {
		if c.ClientIP() == "127.0.0.1" {
			if handlePost(c) {
				c.String(http.StatusOK, "success")
			} else {
				c.String(http.StatusOK, "fail")
			}
		} else {
			c.String(http.StatusForbidden, "Forbidden")
		}
	})
	statikFS, err := fs.New()
	if err != nil {
		logger.Logger.Fatal("Create filesystem failed.",
			zap.Error(err))
	}
	r.StaticFS("/static", statikFS)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "static")
	})

	r.Run(":" + strconv.Itoa(Port))

	// r.Run(":8081")
}
