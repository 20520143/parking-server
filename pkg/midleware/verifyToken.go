package midleware

import "github.com/gin-gonic/gin"

func verifyToken(c *gin.Context) {
	//get access token
	token := c.Request.Header.Get("x-access-token")
	if token != "" {

	}
}
