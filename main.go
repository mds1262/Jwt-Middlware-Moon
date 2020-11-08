package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sample-middleware/controller"
	"sample-middleware/utils"
)

type Result struct {
	status map[string]interface{}
}

type SessionClass struct {
	*utils.SessionConf
}

func main() {
	result := &Result{}
	sc := &SessionClass{
		&utils.SessionConf{
			RootStoreName: utils.STORNAME,
			StoreNames: []string{utils.KUSUPLOADKEYS,utils.KUSTOKEN},
		},
	}
	g := gin.Default()

	g.Use(gin.Logger())

	g.Use(gin.Recovery())

	g.Use(utils.InitSession(sc))

	g.Use(controller.Passport(sc.SessionConf))

	gg := g.Group("/")

	gg.GET("/1", func(c *gin.Context) {
		msgCode := map[string]interface{}{
			"code" : 200,
			"msg" : "success",
		}
		result.status = msgCode
		c.JSON(http.StatusOK,msgCode)
	})

	g.Run(":8081")
}
