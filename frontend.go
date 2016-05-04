package main

import (
	"github.com/gin-gonic/gin"
)

// Frontend serves the application
func Frontend(c *gin.Context) {
	c.HTML(200, "page.html", gin.H{
		"test": "test",
	})
}
