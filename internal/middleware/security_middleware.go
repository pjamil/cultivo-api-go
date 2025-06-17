package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders sets security-related HTTP headers
// @Summary Set security headers
// @Description Sets security headers to protect against common vulnerabilities
// @Tags middleware
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Security headers set"
// @Failure 500 {object} map[string]interface{}
// @Router /middleware/security [get]
// SecurityHeaders is a middleware that sets security headers to protect against common vulnerabilities.
// It adds headers like X-Content-Type-Options, X-Frame-Options, and X-XSS-Protection.
// This middleware should be applied globally to all routes to enhance security.
// It is recommended to use this middleware in production environments to mitigate risks such as MIME type sniffing, clickjacking, and cross-site scripting (XSS) attacks.
// It is a good practice to include this middleware in your application to ensure that security headers are consistently applied across all responses.

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
