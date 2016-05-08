package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Frontend serves the application
func Frontend(c *gin.Context) {
	q := c.Request.URL.Query()
	var elements []int
	if vals, exists := q["score_id"]; exists {
		for _, val := range vals {
			newEl, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			var c bool
			for _, existingVal := range elements {
				if existingVal == newEl {
					c = true
					continue
				}
			}
			if c {
				continue
			}
			elements = append(elements, newEl)
		}
	}

	c.HTML(200, "page.html", gin.H{
		"elements": elements,
	})
}
