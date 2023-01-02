package kameleoon

import (
	"crypto/sha256"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
	"github.com/valyala/fasthttp"
)

type Cookie struct {
	defaultVisitorCode string
}

const (
	VisitorCodeLength = 16
	CookieKeyJs       = "_js_"
	CookieName        = "kameleoonVisitorCode"
	CookieExpireTime  = 380 * 24 * time.Hour
)

// GetVisitorCode should be called to get the Kameleoon visitorCode for the current visitor.
//
// This is especially important when using Kameleoon in a mixed front-end and back-end environment,
// where user identification consistency must be guaranteed.
//
// The implementation logic is described here:
// First we check if a kameleoonVisitorCode cookie or query parameter associated with the current HTTP request can be
// found. If so, we will use this as the visitor identifier. If no cookie / parameter is found in the current
// request, we either randomly generate a new identifier, or use the defaultVisitorCode argument as identifier if it
// is passed. This allows our customers to use their own identifiers as visitor codes, should they wish to.
// This can have the added benefit of matching Kameleoon visitors with their own users without any additional
// look-ups in a matching table.
func (c *Client) GetVisitorCode(req *fasthttp.Request, defaultVisitorCode ...string) string {
	visitorCode := readVisitorCode(req)
	if len(visitorCode) == 0 {
		if len(defaultVisitorCode) > 0 {
			visitorCode = defaultVisitorCode[0]
		} else {
			visitorCode = utils.GetRandomString(VisitorCodeLength)
		}
	}
	return visitorCode
}

// SetVisitorCode should be called to set the Kameleoon visitorCode in response cookie.
//
// The server-side (via HTTP header) kameleoonVisitorCode cookie is set with the value.
func (c *Client) SetVisitorCode(resp *fasthttp.Response, visitorCode, domain string) error {
	cookie := newVisitorCodeCookie(visitorCode, domain)
	resp.Header.SetCookie(cookie)
	fasthttp.ReleaseCookie(cookie)
	return nil
}

func (c *Client) ObtainVisitorCode(req *fasthttp.Request, resp *fasthttp.Response, domain string, defaultVisitorCode ...string) (string, error) {
	visitorCode := c.GetVisitorCode(req, defaultVisitorCode...)
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return visitorCode, err
	}
	c.SetVisitorCode(resp, visitorCode, domain)
	return visitorCode, nil
}

func (c *Client) validateVisitorCode(visitorCode string) (bool, error) {
	if visitorCode == "" {
		return false, newErrVisitorCodeNotValid("empty visitor code")
	} else if len(visitorCode) > KAMELEOON_VISITOR_CODE_LENGTH {
		return false, newErrVisitorCodeNotValid("is longer than 255 chars")
	}
	return true, nil
}

func readVisitorCode(req *fasthttp.Request) string {
	cookie := string(req.Header.Cookie(CookieName))
	if strings.HasPrefix(cookie, CookieKeyJs) {
		cookie = cookie[len(CookieKeyJs):]
	}
	if len(cookie) < VisitorCodeLength {
		return ""
	}
	return cookie[:VisitorCodeLength]
}

func newVisitorCodeCookie(visitorCode, domain string) *fasthttp.Cookie {
	c := fasthttp.AcquireCookie()
	c.SetKey(CookieName)
	c.SetValue(visitorCode)
	c.SetExpire(time.Now().Add(CookieExpireTime))
	c.SetPath("/")
	c.SetDomain(domain)
	return c
}

func getHashDouble(visitorCode string, containerID int, respoolTime []types.RespoolTime) float64 {
	return getHashDoubleSuffix(visitorCode, containerID, respoolTime, "")
}

func getHashDoubleV2(visitorCode string, containerID int, suffix string) float64 {
	return getHashDoubleSuffix(visitorCode, containerID, nil, suffix)
}

func getHashDoubleSuffix(visitorCode string, containerID int, respoolTime []types.RespoolTime, suffix string) float64 {
	var b []byte
	b = append(b, visitorCode...)
	b = append(b, utils.WriteUint(containerID)...)
	b = append(b, suffix...)

	vals := make([]float64, len(respoolTime))
	i := 0
	for _, v := range respoolTime {
		vals[i] = v.Value
		i++
	}
	sort.Float64s(vals)
	for _, v := range vals {
		b = append(b, strconv.FormatFloat(v, 'f', -1, 64)...)
	}

	h := sha256.New()
	h.Write(b)

	z := new(big.Int).SetBytes(h.Sum(nil))
	n1 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

	f1 := new(big.Float).SetInt(z)
	f2 := new(big.Float).SetInt(n1)
	f, _ := new(big.Float).Quo(f1, f2).Float64()
	return f
}
