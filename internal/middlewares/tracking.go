package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"net"
	"net/http"
	"strings"
	"time"

	"tracking/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"

	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RequestTracking() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracking := &TrackingInfo{
			Timestamp: time.Now(),
			Host:      c.Request.Host,
		}
		tracking.OriginalIP = extractOriginalIP(c)
		tracking.ForwardedIPs = extractForwardedIPs(c)

		if tracking.OriginalIP != "" {
			extractGeolocation(tracking.OriginalIP, tracking)
		}

		tracking.UserAgent = c.Request.UserAgent()
		parseUserAgent(tracking)

		tracking.Cookies = extractCookies(c)
		if tracking.Cookies != nil {
			tracking.VisitorUUID = tracking.Cookies[constants.CookieVisitorUUID]
			tracking.SessionToken = tracking.Cookies[constants.CookieSessionToken]
		}

		tracking.Referer = c.Request.Referer()
		tracking.AcceptLanguage = c.GetHeader("Accept-Language")

		c.Set(TrackingInfoKey, tracking)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), TrackingInfoKey, tracking))

		c.Next()
	}
}

func extractOriginalIP(c *gin.Context) string {
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

func extractGeolocation(ip string, tracking *TrackingInfo) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.IsLoopback() || parsedIP.IsPrivate() || parsedIP.IsUnspecified() {
		return
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var data struct {
		Status     string `json:"status"`
		Country    string `json:"country"`
		RegionName string `json:"regionName"`
		City       string `json:"city"`
		Timezone   string `json:"timezone"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	if data.Status != "success" {
		return
	}

	tracking.Country = data.Country
	tracking.Region = data.RegionName
	tracking.City = data.City
	tracking.Timezone = data.Timezone
}

func parseUserAgent(tracking *TrackingInfo) {
	if tracking.UserAgent == "" {
		tracking.Browser = "unknown"
		tracking.OS = "unknown"
		return
	}

	ua := user_agent.New(tracking.UserAgent)

	tracking.Mobile = ua.Mobile()

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

func GetTrackingInfo(ctx context.Context) (*TrackingInfo, bool) {
	if ctx == nil {
		return nil, false
	}
	if info := ctx.Value(TrackingInfoKey); info != nil {
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
	if logger == nil {
		logger = slog.Default()
	}
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
			if contentType != "application/json" {
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

func Cors(allowedOrigins []string) gin.HandlerFunc {
	allowedMap := make(map[string]bool)
	for _, o := range allowedOrigins {
		allowedMap[o] = true
	}
	
	// Allow any origin if allowedOrigins slice contains "*"
	allowAll := allowedMap["*"]

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		if origin != "" {
			if allowAll || allowedMap[origin] {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			} else {
				// If not allowed, we don't set CORS headers so browser blocks it
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key, X-Signature, X-Timestamp, X-Nonce, X-Request-Id, X-Real-IP, X-Forwarded-For, CF-Connecting-IP, True-Client-IP")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
