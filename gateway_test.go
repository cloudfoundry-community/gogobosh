package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

func TestNewRequest(t *testing.T) {

	gateway := gogobosh.NewDirectorGateway()

	request, apiResponse := gateway.NewRequest("GET", "https://example.com/v2/apps", "admin", "admin", nil)

	assert.True(t, apiResponse.IsSuccessful())
	assert.Equal(t, request.HttpReq.Header.Get("Authorization"), "BEARER my-access-token")
	assert.Equal(t, request.HttpReq.Header.Get("accept"), "application/json")
	assert.Equal(t, request.HttpReq.Header.Get("User-Agent"), "gogobosh "+gogobosh.Version+" / "+runtime.GOOS)
}

func TestNewRequestWithAFileBody(t *testing.T) {

	gateway := gogobosh.NewDirectorGateway()

	body, err := os.Open("../../fixtures/hello_world.txt")
	assert.NoError(t, err)
	request, apiResponse := gateway.NewRequest("GET", "https://example.com/v2/apps", "admin", "admin", body)

	assert.True(t, apiResponse.IsSuccessful())
	assert.Equal(t, request.HttpReq.ContentLength, 12)
}
