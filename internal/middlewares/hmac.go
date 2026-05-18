package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HMACConfig struct {
	GetSecretKey func(publicKey string) (string, error)

	MaxClockSkew time.Duration

	NonceStore NonceStore

	RequireContentInSignature bool
}

type NonceStore interface {
	CheckAndStore(publicKey, nonce string, expiry time.Time) bool
}

func HMAC(config HMACConfig) gin.HandlerFunc {
	if config.MaxClockSkew == 0 {
		config.MaxClockSkew = 5 * time.Minute
	}

	return func(c *gin.Context) {
		publicKey, signature, timestamp, nonce, ok := validateHeaders(c)
		if !ok {
			abortWithError(c, http.StatusUnauthorized, "Missing required headers: X-API-Key, X-Signature, X-Timestamp, X-Nonce")
			return
		}

		requestTimestamp, ok := validateTimestamp(timestamp, config.MaxClockSkew)
		if !ok {
			abortWithError(c, http.StatusUnauthorized, "Invalid timestamp or request expired")
			return
		}

		if !validateNonce(config.NonceStore, publicKey, nonce, requestTimestamp, config.MaxClockSkew) {
			abortWithError(c, http.StatusUnauthorized, "Invalid nonce or replay attack detected")
			return
		}

		secretKey, err := config.GetSecretKey(publicKey)
		if err != nil {
			abortWithError(c, http.StatusNotFound, "Invalid API Key")
			return
		}

		bodyBytes, err := readRequestBody(c)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, "Failed to read request body")
			return
		}

		canonicalRequest := buildCanonicalRequest(c, timestamp, nonce, bodyBytes, config.RequireContentInSignature)

		h := hmac.New(sha256.New, []byte(secretKey))
		h.Write([]byte(canonicalRequest))
		expectedSignature := hex.EncodeToString(h.Sum(nil))

		if !verifySignature(signature, expectedSignature) {
			abortWithError(c, http.StatusForbidden, "Invalid Signature")
			return
		}

		c.Set("publicKey", publicKey)
		c.Set("authenticatedAt", time.Now())

		c.Next()
	}
}

func abortWithError(c *gin.Context, status int, message string) {
	requestID := c.GetHeader("X-Request-Id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	c.AbortWithStatusJSON(status, gin.H{
		"code":       status,
		"message":    message,
		"version":    "v1",
		"request_id": requestID,
	})
}

func validateHeaders(c *gin.Context) (publicKey, signature, timestamp, nonce string, ok bool) {
	publicKey = c.GetHeader("X-API-Key")
	signature = c.GetHeader("X-Signature")
	timestamp = c.GetHeader("X-Timestamp")
	nonce = c.GetHeader("X-Nonce")

	if publicKey == "" || signature == "" || timestamp == "" || nonce == "" {
		return "", "", "", "", false
	}
	return publicKey, signature, timestamp, nonce, true
}

func validateTimestamp(timestamp string, maxClockSkew time.Duration) (time.Time, bool) {
	requestTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, false
	}

	requestTimestamp := time.Unix(requestTime, 0)
	timeDiff := time.Since(requestTimestamp).Abs()

	if timeDiff > maxClockSkew {
		return time.Time{}, false
	}
	return requestTimestamp, true
}

func validateNonce(store NonceStore, publicKey, nonce string, requestTimestamp time.Time, maxClockSkew time.Duration) bool {
	if store == nil {
		return true
	}
	if nonce == "" {
		return false
	}

	expiry := requestTimestamp.Add(maxClockSkew * 2)
	return store.CheckAndStore(publicKey, nonce, expiry)
}

func readRequestBody(c *gin.Context) ([]byte, error) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return bodyBytes, nil
}

func verifySignature(signature, expectedSignature string) bool {
	providedSignature, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	expectedSignatureBytes, _ := hex.DecodeString(expectedSignature)
	return hmac.Equal(providedSignature, expectedSignatureBytes)
}

// buildCanonicalRequest creates the standardized string that gets signed
// Format: METHOD\nPATH\nQUERY\nTIMESTAMP\nNONCE\n[CONTENT-TYPE\nCONTENT-LENGTH\n]BODY
func buildCanonicalRequest(c *gin.Context, timestamp, nonce string, body []byte, includeContent bool) string {
	var parts []string

	// Core request elements (always included)
	parts = append(parts,
		c.Request.Method,
		c.Request.URL.Path,
		c.Request.URL.RawQuery, // Include query parameters
		timestamp,
	)

	// Include nonce if provided
	if nonce != "" {
		parts = append(parts, nonce)
	}

	// Optional: include content headers for additional security
	if includeContent {
		parts = append(parts,
			c.GetHeader("Content-Type"),
			c.GetHeader("Content-Length"),
		)
	}

	// Add body (can be empty)
	parts = append(parts, string(body))

	return strings.Join(parts, "\n")
}

// Utility: Generate signature (for client-side reference)
func GenerateSignature(method, path, query, timestamp, nonce string, body []byte, secretKey string) string {
	parts := []string{method, path, query, timestamp}
	if nonce != "" {
		parts = append(parts, nonce)
	}
	parts = append(parts, string(body))

	canonical := strings.Join(parts, "\n")

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(canonical))
	return hex.EncodeToString(h.Sum(nil))
}
