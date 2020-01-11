package oauth

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("about to start oauth tests")
	os.Exit(m.Run())
}

func TestOauthConstants(t *testing.T) {
	assert.EqualValues(t, headerXPublic, "X-Public")
	assert.EqualValues(t, headerXClientID, "X-Client-ID")
	assert.EqualValues(t, headerXCallerID, "X-Caller-ID")
	assert.EqualValues(t, paramAccessToken, "access_token")
}

func TestIsPublicNilRequest(t *testing.T) {
	assert.True(t, IsPublic(nil))
}

func TestIsPublicNoError(t *testing.T) {
	request := http.Request{
		Header: make(http.Header),
	}
	assert.False(t, IsPublic(&request))

	request.Header.Add("X-Public", "true")
	assert.True(t, IsPublic(&request))
}

func TestGetCallerIDNilRequest(t *testing.T) {

}

func TestGetAccessTokenInvalidRestclientResponse(t *testing.T) {
	accessToken, err := getAccessToken("AbC123")
	assert.Nil(t, accessToken)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "invalid restclient response when trying to get access token", err.Message())
}
