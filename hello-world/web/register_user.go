package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func RegisterUser(c *gin.Context) {
	tableName, exists := os.LookupEnv("USER_TABLE")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table does not exist"})
	}
	users := &UserUnencryptedDao{
		tableName: tableName,
	}
	var req registerUserRequest
	e := c.ShouldBindJSON(&req)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}
	e = users.RegisterUser(req.Username, req.Password)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}
}

type registerUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}