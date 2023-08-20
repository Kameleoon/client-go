package network

import (
	"crypto/tls"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// request

type HttpMethod string

const (
	HttpGet  = "GET"
	HttpPost = "POST"
)

type ContentType string

const (
	NoneContentType ContentType = ""
	TextContentType ContentType = "text/plain"
	JsonContentType ContentType = "application/json"
	FormContentType ContentType = "application/x-www-form-urlencoded"
)

type Request struct {
	Method      HttpMethod    // mandatory
	Url         string        // mandatory
	ContentType ContentType   // optional ("")
	AuthToken   string        // optional ("")
	Timeout     time.Duration // mandatory
	UserAgent   string        // optional ("")
	Data        string        // optional ("")
}

// response

type Response struct {
	Err     error    // optional (nil)
	Code    int      // optional (0)
	Body    []byte   // optional ([])
	Request *Request // mandatory
}

// declaration

type NetProvider interface {
	Call(request Request) Response
}

// implementation

const (
	AuthorizationHeader = "Authorization"
)

type NetProviderImpl struct {
	client fasthttp.Client
}

func NewNetProviderImpl(name string, readTimeout time.Duration, writeTimeout time.Duration,
	maxConnsPerHost int, proxyUrl string) *NetProviderImpl {
	np := &NetProviderImpl{
		client: fasthttp.Client{
			Name:            name,
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			MaxConnsPerHost: maxConnsPerHost,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	if len(proxyUrl) > 0 {
		np.client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxyUrl)
	}
	return np
}

func (np *NetProviderImpl) Call(request Request) Response {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(string(request.Method))
	req.Header.SetRequestURI(request.Url)
	if len(request.UserAgent) > 0 {
		req.Header.SetUserAgent(request.UserAgent)
	}
	if len(request.AuthToken) > 0 {
		req.Header.Set(AuthorizationHeader, request.AuthToken)
	}
	if len(request.ContentType) > 0 {
		req.Header.SetContentType(string(request.ContentType))
	}
	if len(request.Data) > 0 {
		req.SetBodyString(request.Data)
	}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := np.client.DoTimeout(req, resp, request.Timeout)
	var response Response
	if err == nil {
		response = Response{Code: resp.StatusCode(), Body: resp.Body(), Request: &request}
	} else {
		response = Response{Err: err, Request: &request}
	}
	return response
}
