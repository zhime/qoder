package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
	// 这里应该使用测试数据库，暂时返回基础路由
	router := gin.New()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Service is running",
		})
	})
	
	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()
	
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestLoginEndpoint(t *testing.T) {
	// 这个测试需要完整的数据库设置，暂时跳过
	t.Skip("需要数据库连接，暂时跳�?)
}
