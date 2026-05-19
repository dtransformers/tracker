package middlewares

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type TrackingInfo struct {
	OriginalIP     string            `json:"original_ip"`
	ForwardedIPs   []string          `json:"forwarded_ips,omitempty"`
	UserAgent      string            `json:"user_agent"`
	DeviceType     string            `json:"device_type"`
	Browser        string            `json:"browser"`
	BrowserVersion string            `json:"browser_version"`
	OS             string            `json:"os"`
	Platform       string            `json:"platform"`
	Mobile         bool              `json:"mobile"`
	Bot            bool              `json:"bot"`
	Cookies        map[string]string `json:"cookies,omitempty"`
	Referer        string            `json:"referer,omitempty"`
	AcceptLanguage string            `json:"accept_language,omitempty"`
	Host           string            `json:"host"`
	Timestamp      time.Time         `json:"timestamp"`
}

const TrackingInfoKey = "tracking_info"

func RequestTracking() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracking := &TrackingInfo{
			Timestamp: time.Now(),
			Host:      c.Request.Host,
		}
		tracking.OriginalIP = extractOriginalIP(c)
		tracking.ForwardedIPs = extractForwardedIPs(c)

		tracking.UserAgent = c.Request.UserAgent()
		parseUserAgent(tracking)

		tracking.Cookies = extractCookies(c)

		tracking.Referer = c.Request.Referer()
		tracking.AcceptLanguage = c.GetHeader("Accept-Language")
		fmt.Println("Extracting tracking info:", tracking)
		c.Set(TrackingInfoKey, tracking)

		c.Next()
	}
}

func extractOriginalIP(c *gin.Context) string {
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return ip
		}
	}
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			if parsedIP := net.ParseIP(clientIP); parsedIP != nil {
				return clientIP
			}
		}
	}
	if ip := c.GetHeader("CF-Connecting-IP"); ip != "" {
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return ip
		}
	}
	if ip := c.GetHeader("True-Client-IP"); ip != "" {
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return ip
		}
	}
	return c.ClientIP()
}

func extractForwardedIPs(c *gin.Context) []string {
	var ips []string

	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		for _, ip := range parts {
			trimmed := strings.TrimSpace(ip)
			if parsedIP := net.ParseIP(trimmed); parsedIP != nil {
				ips = append(ips, trimmed)
			}
		}
	}

	return ips
}

func parseUserAgent(tracking *TrackingInfo) {
	if tracking.UserAgent == "" {
		tracking.DeviceType = "unknown"
		tracking.Browser = "unknown"
		tracking.OS = "unknown"
		return
	}

	ua := user_agent.New(tracking.UserAgent)

	if ua.Mobile() {
		tracking.Mobile = true
		if ua.Mobile() {
			tracking.DeviceType = "mobile"
		} else {
			tracking.DeviceType = "tablet"
		}
	} else {
		tracking.DeviceType = "desktop"
	}

	browserName, browserVersion := ua.Browser()
	tracking.Browser = browserName
	tracking.BrowserVersion = browserVersion
	tracking.OS = ua.OS()
	tracking.Platform = ua.Platform()

	tracking.Bot = ua.Bot()
}

func extractCookies(c *gin.Context) map[string]string {
	cookies := make(map[string]string)

	for _, cookie := range c.Request.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	return cookies
}

func GetTrackingInfo(c *gin.Context) (*TrackingInfo, bool) {
	if info, exists := c.Get(TrackingInfoKey); exists {
		if tracking, ok := info.(*TrackingInfo); ok {
			return tracking, true
		}
	}
	return nil, false
}

func RateLimit(RateLimitWindow time.Duration, RateLimitRequests int) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: RateLimitWindow,
		Limit:  int64(RateLimitRequests),
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	middleware := mgin.NewMiddleware(instance)

	return middleware
}

func Metrics(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"duration_ms", duration.Milliseconds(),
		}

		if tracking, exists := GetTrackingInfo(c); exists {
			attrs = append(attrs,
				"client_ip", tracking.OriginalIP,
				"user_agent", tracking.UserAgent,
				"device_type", tracking.DeviceType,
				"browser", tracking.Browser,
				"os", tracking.OS,
				"mobile", tracking.Mobile,
				"bot", tracking.Bot,
			)

			if tracking.Referer != "" {
				attrs = append(attrs, "referer", tracking.Referer)
			}
		} else {
			attrs = append(attrs,
				"client_ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
			)
		}

		logger.Info("request completed", attrs...)
	}
}

func Validation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPatch {
			contentType := c.GetHeader("Content-Type")
			if contentType != "application/json" && contentType != "" {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
