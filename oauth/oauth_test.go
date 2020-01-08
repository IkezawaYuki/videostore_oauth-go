package oauth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOauthConstants(t *testing.T) {
	assert.EqualValues(t, headerXPublic, "X-Public")
	assert.EqualValues(t, headerXClientID, "X-Client-ID")
	assert.EqualValues(t, headerXCallerID, "X-Caller-ID")
	assert.EqualValues(t, paramAccessToken, "access_token")
}
