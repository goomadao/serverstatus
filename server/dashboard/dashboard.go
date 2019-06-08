package dashboard

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/goomadao/serverstatus/server/status"
	"github.com/goomadao/serverstatus/util/logger"
	"github.com/gin-gonic/gin"

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
	r := gin.Default()
	r.GET("/api/servers", func(c *gin.Context) {
		servers := new(status.Servers)
		status.GetServers(servers)
		c.JSON(200, servers)
	})
	r.POST("/api/servers", func(c *gin.Context) {
		if handlePost(c) {
			c.String(200, "success")
		} else {
			c.String(200, "fail")
		}
	})
	r.LoadHTMLGlob("../web/dist/*.html")
	r.LoadHTMLFiles("../web/static/*/*")
	r.Static("/static", "../web/dist/static")
	r.StaticFile("/", "../web/dist/index.html")

	r.Run(":" + strconv.Itoa(Port))

	// r.Run(":8081")
}
