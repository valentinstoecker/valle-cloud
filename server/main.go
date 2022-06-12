package main

import (
	"github.com/gin-gonic/gin"
	"github.com/valentinstoecker/valle-cloud/server/files"
)

func main() {
	engine := gin.New()
	engine.Use(gin.Logger())
	api := engine.Group("/api")
	api.GET("/files", files.GetFiles)
	api.GET("/files/:hash", files.GetImage)
	api.GET("/files/:hash/thumbnail", files.GetThumbnail)
	api.POST("/files", files.UploadFiles)
	engine.Static("/static", "../dist/valle-cloud")
	engine.NoRoute(func(c *gin.Context) {
		c.File("../dist/valle-cloud/index.html")
	})
	engine.Run(":8080")
}
