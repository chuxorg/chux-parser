package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwtverifier "github.com/okta/okta-jwt-verifier-golang"

	cfg "idc-okta-api/internal/config"
)

func isAuthenticated(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		log.Printf("Access token not found")
		return false
	}

	tokenParts := strings.Split(authHeader, "Bearer ")
	bearerToken := tokenParts[1]

	toValidate := map[string]string{}
	toValidate["aud"] = "api://sandbox.acadian.am"

	url := cfg.GetOktaOAuth2Issuer()
	verifier := jwtverifier.JwtVerifier{
		Issuer:           url,
		ClaimsToValidate: toValidate,
	}
	_, err := verifier.New().VerifyAccessToken(bearerToken)

	if err != nil {
		log.Printf("Validation failed: %s", err.Error())
		return false
	}
	return true
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthenticated(c.Request) {
			path := c.Request.URL.Path
			log.Printf("Unauthorized route: %s", path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized route"})
			c.AbortWithError(403, fmt.Errorf("You are Unauthorized for %s", path))
			return
		}

		c.Next()
	}
}
