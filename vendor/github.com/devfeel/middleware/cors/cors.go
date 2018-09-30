package cors

import (
	"github.com/devfeel/dotweb"
	"strconv"
)

//CROS配置
type Config struct {
	enabledCROS      bool
	allowedOrigins   string
	allowedMethods   string
	allowedHeaders   string
	allowCredentials bool
	exposeHeaders    string
	allowedP3P       string
	maxAge           int
}

func (c *Config) UseDefault() *Config {
	c.enabledCROS = true
	c.allowedOrigins = "*"
	c.allowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	c.allowedHeaders = "Content-Type,Authorization"
	c.allowedP3P = "CP=\"CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR\""
	return c
}

func (c *Config) Enabled() *Config {
	c.enabledCROS = true
	return c
}

func (c *Config) SetOrigin(origins string) *Config {
	c.allowedOrigins = origins
	return c
}

func (c *Config) SetMethod(methods string) *Config {
	c.allowedMethods = methods
	return c
}

func (c *Config) SetHeader(headers string) *Config {
	c.allowedHeaders = headers
	return c
}

func (c *Config) SetExposeHeaders(headers string) *Config {
	c.exposeHeaders = headers
	return c
}

func (c *Config) SetAllowCredentials(flag bool) *Config {
	c.allowCredentials = flag
	return c
}

func (c *Config) SetMaxAge(maxAge int) *Config {
	c.maxAge = maxAge
	return c
}

func (c *Config) SetP3P(p3p string) *Config {
	c.allowedP3P = p3p
	return c
}

func NewConfig() *Config {
	return &Config{}
}

//jwt中间件
type CORSMiddleware struct {
	dotweb.BaseMiddlware
	config *Config
}

func (m *CORSMiddleware) Handle(ctx dotweb.Context) error {
	if m.config.enabledCROS {
		ctx.Response().SetHeader(dotweb.HeaderAccessControlAllowOrigin, m.config.allowedOrigins)
		ctx.Response().SetHeader(dotweb.HeaderAccessControlAllowMethods, m.config.allowedMethods)
		ctx.Response().SetHeader(dotweb.HeaderAccessControlAllowHeaders, m.config.allowedHeaders)
		ctx.Response().SetHeader(dotweb.HeaderAccessControlExposeHeaders, m.config.exposeHeaders)
		ctx.Response().SetHeader(dotweb.HeaderAccessControlAllowCredentials, strconv.FormatBool(m.config.allowCredentials))
		ctx.Response().SetHeader(dotweb.HeaderAccessControlMaxAge, strconv.Itoa(m.config.maxAge))
		ctx.Response().SetHeader(dotweb.HeaderP3P, m.config.allowedP3P)
	}
	return m.Next(ctx)
}

// Middleware create new CORS Middleware
func Middleware(config *Config) *CORSMiddleware {
	return &CORSMiddleware{config: config}
}

// DefaultMiddleware create new CORS Middleware with default config
func DefaultMiddleware() *CORSMiddleware {
	option := NewConfig().UseDefault()
	return &CORSMiddleware{config: option}
}
