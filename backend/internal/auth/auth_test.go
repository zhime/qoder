package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager(t *testing.T) {
	jwtManager := NewJWTManager("test-secret", 3600, 7200)

	t.Run("GenerateToken", func(t *testing.T) {
		token, err := jwtManager.GenerateToken(1, "testuser", "admin")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("ValidateToken", func(t *testing.T) {
		token, err := jwtManager.GenerateToken(1, "testuser", "admin")
		assert.NoError(t, err)

		claims, err := jwtManager.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), claims.UserID)
		assert.Equal(t, "testuser", claims.Username)
		assert.Equal(t, "admin", claims.Role)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		expiredManager := NewJWTManager("test-secret", -1, 7200)
		token, err := expiredManager.GenerateToken(1, "testuser", "admin")
		assert.NoError(t, err)

		time.Sleep(time.Second * 2)

		_, err = jwtManager.ValidateToken(token)
		assert.Error(t, err)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, err := jwtManager.ValidateToken("invalid-token")
		assert.Error(t, err)
	})
}

func TestPasswordHash(t *testing.T) {
	password := "testpassword123"

	t.Run("HashPassword", func(t *testing.T) {
		hash, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("CheckPassword", func(t *testing.T) {
		hash, err := HashPassword(password)
		assert.NoError(t, err)

		// 正确密码
		assert.True(t, CheckPassword(password, hash))

		// 错误密码
		assert.False(t, CheckPassword("wrongpassword", hash))
	})
}