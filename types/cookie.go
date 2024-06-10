package types

type Cookie struct {
	cookies map[string]string
}

func NewCookie(cookies map[string]string) *Cookie {
	return &Cookie{cookies: cookies}
}

func (c *Cookie) dataRestriction() {}

func (c *Cookie) Cookies() map[string]string {
	return c.cookies
}

func (c *Cookie) DataType() DataType {
	return DataTypeCookie
}
