package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSelfPro(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}
