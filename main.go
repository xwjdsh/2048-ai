package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("./2048/index.html")
	router.Static("/style", "./2048/style")
	router.Static("/js", "./2048/js")

	router.GET("/", index)

	router.Run(":8080")
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
