package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/errors"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type Dialer interface {
	Dial(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)
}

// APITestCase represents the data needed to describe an API test case.
type APITestCase[T any] struct {
	Client       *http.Client
	Router       *gin.Engine
	URL          *url.URL
	Method, Body *string
	Header       *http.Header
	WantStatus   *int
	WantResponse *string
	WSDialer     Dialer
	WSChannel    chan T
	Logger       log.Logger
}

func (tc APITestCase[T]) CopyWith(new APITestCase[T]) APITestCase[T] {
	tc.Router = utils.Or(new.Router, tc.Router)
	tc.Method = utils.Or(new.Method, tc.Method)
	tc.URL = utils.Or(new.URL, tc.URL)
	tc.Body = utils.Or(new.Body, tc.Body)
	tc.Header = utils.Or(new.Header, tc.Header)
	tc.WantStatus = utils.Or(new.WantStatus, tc.WantStatus)
	tc.WantResponse = utils.Or(new.WantResponse, tc.WantResponse)
	return tc
}

func (tc APITestCase[T]) WithClient(client *http.Client) APITestCase[T] {
	tc.Client = client
	return tc
}

func (tc APITestCase[T]) WithRouter(router *gin.Engine) APITestCase[T] {
	tc.Router = router
	return tc
}

func (tc APITestCase[T]) WithMethod(method string) APITestCase[T] {
	tc.Method = &method
	return tc
}

func (tc APITestCase[T]) WithURL(URL url.URL) APITestCase[T] {
	tc.URL = &URL
	return tc
}

func (tc APITestCase[T]) WithHost(host string) APITestCase[T] {
	tc.URL.Host = host
	return tc
}

func (tc APITestCase[T]) WithBody(body string) APITestCase[T] {
	tc.Body = &body
	return tc
}

func (tc APITestCase[T]) WithRequest(req any) APITestCase[T] {
	b, err := json.Marshal(req)
	if err != nil {
		panic(fmt.Sprint("failed to marshal request ", err))
	}
	tc.Body = utils.String(string(b))
	return tc
}

func (tc APITestCase[T]) WithHeader(header http.Header) APITestCase[T] {
	tc.Header = &header
	return tc
}

func (tc APITestCase[T]) WithWantStatus(status int) APITestCase[T] {
	tc.WantStatus = &status
	return tc
}

func (tc APITestCase[T]) WithWantResponse(response string) APITestCase[T] {
	tc.WantResponse = &response
	return tc
}

func (tc APITestCase[T]) WithLogger(logger log.Logger) APITestCase[T] {
	tc.Logger = logger
	return tc
}

// CheckEndpoint tests an HTTP endpoint using the given APITestCase spec without using t.Run.
// The check is made through a client or directly on the router.
func (tc APITestCase[T]) CheckEndpoint(t *testing.T) T {
	if tc.Client != nil {
		tc.Logger.Debugf("testing the api with a client")
		return tc.checkEndpoint(t, tc.Client.Do, func(res *http.Response) {
			res.Body.Close()
		})
	}

	tc.Logger.Debugf("Unit test of the api")
	return tc.checkEndpoint(t, func(req *http.Request) (*http.Response, error) {
		res := httptest.NewRecorder()

		if tc.Router == nil {
			t.Fatal("router is required")
		}

		defer func() (*http.Response, error) {
			if err := recover(); err != nil {
				t.Fatalf("panic: %v", err)
			}
			return res.Result(), nil
		}() // nolint:errcheck

		tc.Router.ServeHTTP(res, req)
		return res.Result(), nil
	}, func(res *http.Response) {})
}

func (tc APITestCase[T]) Dial(t *testing.T) *websocket.Conn {
	if tc.URL.Scheme == "ws" && tc.WSDialer != nil {
		req := tc.prepareRequest(t)
		tc.Logger.Debugf("testing the api with a websocket dialer")

		tc.Logger.Infof("connecting to %s", tc.URL.String())

		c, response, err := tc.WSDialer.Dial(
			req.URL.String(),
			nil)
		if err != nil {
			if tc.WantResponse != nil {
				assert.Regexp(t, *tc.WantResponse, err.Error(), "response mismatch")
			}
		}
		httpStatusEqual(t, *tc.WantStatus, response.StatusCode, "status mismatch")
		return c
	}
	t.Fatalf("websocket dialer is required")
	return nil
}

func (tc APITestCase[T]) prepareRequest(t *testing.T) *http.Request {
	body := ""
	if tc.Body != nil {
		body = *tc.Body
	}

	var buffer *bytes.Buffer
	tc.Logger.Debugf("Request sent with body %s", body)
	buffer = bytes.NewBufferString(body)

	if tc.Method == nil {
		t.Fatal("method is required")
	}

	tc.Logger.Debugf("Request sent with method %s", *tc.Method)

	URL := tc.URL.String()
	if tc.URL != nil {
		tc.Logger.Debugf("Request sent with URL %s", URL)
	}

	req, err := http.NewRequest(*tc.Method, URL, buffer)
	if err != nil {
		if t != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		t.Fatalf("failed to create request: %v", err)
	}

	req.Close = true

	if tc.Header != nil {
		if *tc.Header != nil {
			req.Header = *tc.Header
		}
	}

	tc.Logger.Debugf("Request sent with header %s", req.Header)

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func (tc APITestCase[T]) checkEndpoint(
	t *testing.T,
	Do func(*http.Request) (*http.Response, error),
	postProc func(*http.Response),
) T {
	req := tc.prepareRequest(t)

	res, err := Do(req)
	if err != nil {
		t.Fatalf("failed to run the request: %v", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read the body: %v", err)
	}

	bs := string(bodyBytes)

	tc.Logger.Debugf("Response received with status %d", res.StatusCode)
	tc.Logger.Debugf("Response received with body %s", bs)

	var response T
	err = json.Unmarshal([]byte(bs), &response)

	if err != nil {
		tc.Logger.Error("failed to unmarshal response: ", err)
	}

	httpStatusEqual(t, *tc.WantStatus, res.StatusCode, "status mismatch")
	if tc.WantResponse != nil {
		assert.Regexp(t, *tc.WantResponse, bs, "response mismatch")
	}

	if t.Failed() {
		t.Logf("\n\tEndpoint: %s", tc.URL)
	}

	postProc(res)
	return response
}

func httpStatusEqual(t *testing.T, expected int, actual int, msgAndArgs ...interface{}) bool {
	if expected != actual {
		return assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %d (%s)\n"+
			"actual  : %d (%s)", expected, http.StatusText(expected), actual, http.StatusText(actual)), msgAndArgs...)
	}

	return true
}

func (tc APITestCase[T]) Call() (T, error) {
	var output T
	if tc.Client == nil {
		return output, fmt.Errorf("client is required")
	}

	return tc.call()
}

func (tc APITestCase[T]) call() (T, error) {
	var output T

	if tc.URL == nil {
		return output, fmt.Errorf("url is required")
	}

	URL := tc.URL.String()

	tc.Logger.Debugf("Request sent with URL %s", URL)

	body := ""
	if tc.Body != nil {
		body = *tc.Body
	}

	tc.Logger.Debugf("Request sent with body %s", body)

	buffer := bytes.NewBufferString(body)

	if tc.Method == nil {
		return output, fmt.Errorf("method is required")
	}
	tc.Logger.Debugf("Request sent with method %s", *tc.Method)

	req, err := http.NewRequest(*tc.Method, URL, buffer)
	if err != nil {
		return output, fmt.Errorf("failed to create request: %v", err)
	}

	req.Close = true

	if tc.Header != nil {
		if *tc.Header != nil {
			req.Header = *tc.Header
		}
	}

	tc.Logger.Debugf("Request sent with header %s", req.Header)

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := tc.Client.Do(req)

	if err != nil {
		return output, fmt.Errorf("failed to run the request: %v", err)
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return output, fmt.Errorf("failed to read the body: %v", err)
	}

	bs := string(bodyBytes)

	tc.Logger.Debugf("Response received with status %d", res.StatusCode)
	tc.Logger.Debugf("Response received with body %s", bs)

	b, err := regexp.Match("2[0-9][0-9]", []byte(res.Status))
	if err != nil {
		return output, err
	} else if b {
		err = json.Unmarshal([]byte(bs), &output)
		if err != nil {
			return output, err
		}
		return output, nil
	}
	var errResponse errors.ErrorResponse
	err = json.Unmarshal([]byte(bs), &errResponse)
	if err != nil {
		return output, err
	}
	return output, errResponse
}
