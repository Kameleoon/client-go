package kameleoon

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

const (
	HeaderContentTypeJson = "application/json"
	HeaderContentTypeForm = "application/x-www-form-urlencoded"
	HeaderContentTypeText = "text/plain"
	HeaderAuthorization   = "Authorization"
	HeaderPaginationCount = "X-Pagination-Page-Count"
	HeaderTracking        = "Kameleoon-Client"

	MethodGet  = fasthttp.MethodGet
	MethodPost = fasthttp.MethodPost
)

type restClient interface {
	Do(r request, callback respCallback) error
}

type rest struct {
	cfg *RestConfig
	c   *fasthttp.Client
}

type respCallback func(resp *fasthttp.Response, err error) error

func newRESTClient(cfg *RestConfig) restClient {
	c := &fasthttp.Client{
		Name:            cfg.UserAgent,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		MaxConnsPerHost: cfg.MaxConnsPerHost,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		NoDefaultUserAgentHeader: true, // Don't send: User-Agent: fasthttp
	}
	if len(cfg.ProxyURL) > 0 {
		c.Dial = fasthttpproxy.FasthttpHTTPDialer(cfg.ProxyURL)
	}
	return &rest{
		cfg: cfg,
		c:   c,
	}
}

func (c *rest) Do(r request, callback respCallback) error {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(r.Method)
	req.Header.SetUserAgent(c.cfg.UserAgent)
	req.Header.SetRequestURI(r.URL)
	if len(r.AuthToken) > 0 {
		req.Header.Set(HeaderAuthorization, r.AuthToken)
	}
	if len(r.ContentType) > 0 {
		req.Header.SetContentType(r.ContentType)
	}
	if len(r.ClientHeader) > 0 {
		req.Header.Set(HeaderTracking, r.ClientHeader)
	}
	if r.Body != nil {
		req.SetBody(r.Body)
	} else if len(r.BodyString) > 0 {
		req.SetBodyString(r.BodyString)
	}
	timeout := r.Timeout
	if timeout == 0 {
		timeout = c.cfg.DoTimeout
	}
	resp := fasthttp.AcquireResponse()
	doErr := c.c.DoTimeout(req, resp, timeout)

	if callback == nil {
		callback = defaultRespCallback
	}
	err := callback(resp, doErr)

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return err
}

type request struct {
	Method       string
	URL          string
	AuthToken    string
	ContentType  string
	BodyString   string
	Body         []byte
	Timeout      time.Duration
	ClientHeader string
}

func (r request) String() string {
	var s strings.Builder
	s.WriteString("method=")
	s.WriteString(r.Method)
	s.WriteString(", url=")
	s.WriteString(r.URL)
	if len(r.AuthToken) > 0 {
		s.WriteString(", auth_token=")
		s.WriteString(r.AuthToken)
	}
	if len(r.ContentType) > 0 {
		s.WriteString(", content_type=")
		s.WriteString(r.ContentType)
	}
	if r.Timeout > 0 {
		s.WriteString(", timeout=")
		s.WriteString(r.Timeout.String())
	}
	if len(r.ClientHeader) > 0 {
		s.WriteString(", client_header=")
		s.WriteString(r.ClientHeader)
	}
	if r.Body != nil {
		s.WriteString(", body=")
		s.Write(r.Body)
	} else if len(r.BodyString) > 0 {
		s.WriteString(", body=")
		s.WriteString(r.BodyString)
	}
	return s.String()
}

func respCallbackJson(i interface{}) respCallback {
	return func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		if resp.StatusCode() >= fasthttp.StatusBadRequest {
			return ErrBadStatus
		}
		return json.Unmarshal(resp.Body(), i)
	}
}

var defaultRespCallback = func(resp *fasthttp.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode() >= fasthttp.StatusBadRequest {
		return ErrBadStatus
	}
	return nil
}
