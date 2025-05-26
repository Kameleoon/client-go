package network

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// request

type HttpMethod string

const (
	HttpGet  HttpMethod = "GET"
	HttpPost HttpMethod = "POST"
)

type ContentType string

const (
	NoneContentType ContentType = ""
	TextContentType ContentType = "text/plain"
	JsonContentType ContentType = "application/json"
	FormContentType ContentType = "application/x-www-form-urlencoded"
)

type Request struct {
	Method         HttpMethod        // mandatory
	Url            string            // mandatory
	ContentType    ContentType       // optional ("")
	Timeout        time.Duration     // mandatory
	Data           string            // optional ("")
	Headers        map[string]string // optional
	AccessToken    string            // optional ("")
	IsAuthRequired bool              // optional (false)
}

func (r Request) String() string {
	body := "nil"
	if r.Data != "" {
		if strings.HasPrefix(r.Data, "client_id=") {
			body = "****"
		} else {
			body = r.Data
		}
	}
	return fmt.Sprintf("HttpRequest{Method:'%s',Url:'%s',Headers:%v,Body:'%s'}", r.Method, r.Url, r.Headers, body)
}

// response

type Response struct {
	Err         error             // optional (nil)
	Code        int               // optional (0)
	Body        []byte            // optional ([])
	HeadersRead map[string]string // optional ({})
	Request     *Request          // mandatory
}

func (r Response) IsExpectedStatusCode() bool {
	return (r.Code/100 == 2) || (r.Code == 403) || (r.Code == 304)
}

func (r Response) String() string {
	return fmt.Sprintf("HttpResponse{Code:'%d',Reason:'%v',Body:'%s'}", r.Code, r.Err, string(r.Body))
}

// declaration

type NetProvider interface {
	Call(request *Request, headersToRead []string) Response
}

// implementation

const (
	AuthorizationHeader = "Authorization"
)

type NetProviderImpl struct {
	client fasthttp.Client
}

func NewNetProviderImpl(readTimeout time.Duration, writeTimeout time.Duration,
	maxConnsPerHost int, proxyUrl string) *NetProviderImpl {
	np := &NetProviderImpl{
		client: fasthttp.Client{
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

func (np *NetProviderImpl) Call(request *Request, headersToRead []string) Response {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(string(request.Method))
	req.Header.SetRequestURI(request.Url)
	np.setHeaders(req, request)
	if len(request.Data) > 0 {
		req.SetBodyString(request.Data)
	}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := np.client.DoTimeout(req, resp, request.Timeout)
	var response Response
	if err == nil {
		headersRead := np.readHeaders(resp, headersToRead)
		response = Response{Code: resp.StatusCode(), Body: resp.Body(), HeadersRead: headersRead, Request: request}
	} else {
		response = Response{Err: err, Request: request}
	}
	return response
}

func (*NetProviderImpl) setHeaders(req *fasthttp.Request, request *Request) {
	if len(request.Headers) > 0 {
		for key, value := range request.Headers {
			req.Header.Set(key, value)
		}
	}
	if len(request.AccessToken) > 0 {
		req.Header.Set(AuthorizationHeader, "Bearer "+request.AccessToken)
	}
	if len(request.ContentType) > 0 {
		req.Header.SetContentType(string(request.ContentType))
	}
}

func (*NetProviderImpl) readHeaders(resp *fasthttp.Response, headersToRead []string) map[string]string {
	headers := make(map[string]string)
	for _, h := range headersToRead {
		headers[h] = string(resp.Header.Peek(h))
	}
	return headers
}
