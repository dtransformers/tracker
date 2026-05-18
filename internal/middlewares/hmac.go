package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HMAC returns a middleware that verifies the HMAC signature of the request body.
// It expects the Public API Key in 'X-API-Key' header and the signature in 'X-Signature' header.
func HMAC(getSecretKey func(publicKey string) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		publicKey := c.GetHeader("X-API-Key")
		signature := c.GetHeader("X-Signature")

		if publicKey == "" || signature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API Key or Signature"})
			c.Abort()
			return
		}

		secretKey, err := getSecretKey(publicKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		// Read the body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		// Restore the body for the next handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Calculate HMAC
		h := hmac.New(sha256.New, []byte(secretKey))
		h.Write(bodyBytes)
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		if signature != expectedSignature {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Signature"})
			c.Abort()
			return
		}

		// Store public key in context
		c.Set("publicKey", publicKey)

		c.Next()
	}
}
