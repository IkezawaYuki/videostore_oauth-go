package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/IkezawaYuki/videostore_utils-go/rest_errors"
)

const (
	headerXPublic    = "X-Public"
	headerXClientID  = "X-Client-ID"
	headerXCallerID  = "X-Caller-ID"
	paramAccessToken = "access_token"
)

var (
	// oauthRestClient = rest.RequestBuilder{
	// 	BaseURL: "http://localhost:8080",
	// 	Timeout: 200 * time.Millisecond,
	// }
	httpClient = &http.Client{}
)

type accessToken struct {
	ID       string `json:"id"`
	UserID   int64  `json:"user_id"`
	ClientID int64  `json:"client_id"`
}

func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}
	return request.Header.Get(headerXPublic) == "true"
}

func GetCallerID(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	callerID, err := strconv.ParseInt(request.Header.Get(headerXCallerID), 10, 64)
	if err != nil {
		return 0
	}
	return callerID
}

func GetClientID(request *http.Request) int64 {
	if request == nil {
		return 0
	}
	clientID, err := strconv.ParseInt(request.Header.Get(headerXClientID), 10, 64)
	if err != nil {
		return 0
	}
	return clientID
}

func AuthenticateRequest(request *http.Request) rest_errors.RestErr {
	if request == nil {
		return nil
	}
	cleanRequest(request)

	accessTokenID := strings.TrimSpace(request.URL.Query().Get(paramAccessToken))
	if accessTokenID == "" {
		return nil
	}

	at, err := getAccessToken(accessTokenID)
	if err != nil {
		if err.Status() == http.StatusNotFound {
			return nil
		}
		return err
	}
	request.Header.Add(headerXClientID, fmt.Sprintf("%v", at.ClientID))
	request.Header.Add(headerXCallerID, fmt.Sprintf("%v", at.UserID))

	return nil
}

func cleanRequest(request *http.Request) {
	if request == nil {
		return
	}
	request.Header.Del(headerXClientID)
	request.Header.Del(headerXCallerID)
}

func getAccessToken(accessTokenID string) (*accessToken, rest_errors.RestErr) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/oauth/access_token/%s", accessTokenID), nil)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to get access token", errors.New("api error"))
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to get access token", errors.New("api error"))
	}
	// response := oauthRestClient.Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenID))
	// if response == nil || response.Response == nil {
	// 	return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to get access token", errors.New("api error"))
	// }
	byte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("invalid error interface when trying to get access token", errors.New("api error"))
	}
	if response.StatusCode > 299 {
		var restErr rest_errors.RestErr
		err = json.Unmarshal(byte, &restErr)
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to get access token", errors.New("api error"))
		}
		return nil, restErr
	}
	var at accessToken
	if err := json.Unmarshal(byte, &at); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal access token response", errors.New("api error"))
	}
	return &at, nil
}
