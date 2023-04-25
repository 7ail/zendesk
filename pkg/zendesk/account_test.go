package zendesk

import (
	"errors"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrl(t *testing.T) {
	subdomain := "subdomain"
	host := "host"
	requestUrl := fmt.Sprintf("http://%s.%s.com/api/v2/account.json", subdomain, host)

	act := func(doer doer, subdomain, host string) (string, error) {
		zendesk, _ := NewGateway(doer, subdomain, host)

		return zendesk.Url()
	}

	testCases := []struct{
		name string
		arrange func(*testing.T) (doer, string, string)
		assert func(*testing.T, string, error)
	}{
		{
			name: "happy path",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(`{"url": "https://test.zendesk.com"}`)),
					StatusCode: http.StatusOK,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualUrl string, actualErr error) {
				assert.Equal(t, "https://test.zendesk.com", actualUrl)
				assert.NoError(t, actualErr)
			},
		},
		{
			name: "api call failed",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(nil, fmt.Errorf("something went wrong"))

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualUrl string, actualErr error) {
				assert.Equal(t, "", actualUrl)
				assert.True(t, errors.Is(actualErr, ErrUnexpected))
				assert.True(t, strings.Contains(actualErr.Error(), "g.client.Do:"))
			},
		},
		{
			name: "response code 404",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualUrl string, actualErr error) {
				assert.Equal(t, "", actualUrl)
				assert.False(t, errors.Is(actualErr, ErrUnexpected))
				assert.True(t, errors.Is(actualErr, ErrNotFound))
			},
		},
		{
			name: "response code not 200",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					StatusCode: http.StatusLocked,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualUrl string, actualErr error) {
				assert.Equal(t, "", actualUrl)
				assert.True(t, errors.Is(actualErr, ErrUnexpected))
			},
		},
		{
			name: "invalid response body",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("not a json payload")),
					StatusCode: http.StatusOK,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualUrl string, actualErr error) {
				assert.Equal(t, "", actualUrl)
				assert.True(t, errors.Is(actualErr, ErrUnexpected))
				assert.True(t, strings.Contains(actualErr.Error(), "json.Unmarshal:"))
			},
		},
		{
			name: "empty url value",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
					StatusCode: http.StatusOK,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualUrl string, actualErr error) {
				assert.Equal(t, "", actualUrl)
				assert.True(t, errors.Is(actualErr, ErrUnexpected))
				assert.True(t, errors.Is(actualErr, ErrNotFound))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doer, subdomain, host := tc.arrange(t)
			actualUrl, actualErr := act(doer, subdomain, host)
			tc.assert(t, actualUrl, actualErr)
		})
	}
}

func TestSubdomain(t *testing.T) {
	subdomain := "subdomain"
	host := "host"
	requestUrl := fmt.Sprintf("http://%s.%s.com/api/v2/account.json", subdomain, host)

	act := func(doer doer, subdomain, host string) (string, error) {
		zendesk, _ := NewGateway(doer, subdomain, host)

		return zendesk.Subdomain()
	}

	testCases := []struct{
		name string
		arrange func(*testing.T) (doer, string, string)
		assert func(*testing.T, string, error)
	}{
		{
			name: "happy path",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(`{"subdomain": "subdomain"}`)),
					StatusCode: http.StatusOK,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualSubdomain string, actualErr error) {
				assert.Equal(t, "subdomain", actualSubdomain)
				assert.NoError(t, actualErr)
			},
		},
		{
			name: "api call failed",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(nil, fmt.Errorf("something went wrong"))

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualSubdomain string, actualErr error) {
				assert.Equal(t, "", actualSubdomain)
				assert.True(t, errors.Is(actualErr, ErrUnexpected))
				assert.True(t, strings.Contains(actualErr.Error(), "g.client.Do:"))
			},
		},
		{
			name: "response code 404",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					StatusCode: http.StatusNotFound,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualSubdomain string, actualErr error) {
				assert.Equal(t, "", actualSubdomain)
				assert.False(t, errors.Is(actualErr, ErrUnexpected))
				assert.True(t, errors.Is(actualErr, ErrNotFound))
			},
		},
		{
			name: "response code not 200",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					StatusCode: http.StatusLocked,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualSubdomain string, actualErr error) {
				assert.Equal(t, "", actualSubdomain)
				assert.True(t, errors.Is(actualErr, ErrUnexpected))
			},
		},
		{
			name: "invalid response body",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("not a json payload")),
					StatusCode: http.StatusOK,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualSubdomain string, actualErr error) {
				assert.Equal(t, "", actualSubdomain)
				assert.True(t, errors.Is(actualErr, ErrUnexpected))
				assert.True(t, strings.Contains(actualErr.Error(), "json.Unmarshal:"))
			},
		},
		{
			name: "empty url value",
			arrange: func(t *testing.T) (doer, string, string) {
				req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
				assert.NoError(t, err)

				doer := newMockDoer(t)
				doer.On("Do", req).Return(&http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
					StatusCode: http.StatusOK,
				}, nil)

				return doer, subdomain, host
			},
			assert: func(t *testing.T, actualSubdomain string, actualErr error) {
				assert.Equal(t, "", actualSubdomain)
				assert.True(t, errors.Is(actualErr, ErrUnexpected))
				assert.True(t, errors.Is(actualErr, ErrNotFound))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doer, subdomain, host := tc.arrange(t)
			actualSubdomain, actualErr := act(doer, subdomain, host)
			tc.assert(t, actualSubdomain, actualErr)
		})
	}
}