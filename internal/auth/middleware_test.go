package auth

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker/v3"
)

func TestIsAuthenticated(t *testing.T) {
	// create a new HTTP request with an empty Authorization header
	req := httptest.NewRequest("GET", "/", nil)

	// call the isAuthenticated method with the empty request
	result := isAuthenticated(req)

	// verify that isAuthenticated returns false for an empty Authorization header
	if result != false {
		t.Errorf("isAuthenticated() = %v; want false", result)
	}

	// create a new HTTP request with a valid Authorization header
	req = httptest.NewRequest("GET", "/", nil)
	tk := fmt.Sprintf("Bearer %s", faker.Jwt())
	req.Header.Set("Authorization", tk)

	// // call the isAuthenticated method with the valid request
	// result = isAuthenticated(req)

	// // verify that isAuthenticated returns true for a valid Authorization header
	// if result != true {
	// 	t.Errorf("isAuthenticated() = %v; want true", result)
	// }
}
