package responses

import "github.com/gin-gonic/gin"

func ResponseWithError(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

func ResponseWithData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"data": data,
	})
	c.Done()
}
