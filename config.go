package kameleoon

import (
	"strings"
	"time"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

const (
	DefaultConfigPath           = "/etc/kameleoon/client-go.yaml"
	DefaultConfigUpdateInterval = time.Hour
	DefaultRequestTimeout       = 2 * time.Second
	DefaultVisitorDataMaxSize   = 500 // 500 mb
	DefaultTrackingVersion      = "sdk/go/2.0.1"
	UserAgent                   = "kameleoon-client-go/"
)

type Config struct {
	REST                 RestConfig
	Logger               Logger        `yml:"-" yaml:"-"`
	SiteCode             string        `yml:"site_code" yaml:"site_code"`
	TrackingURL          string        `yml:"tracking_url" yaml:"tracking_url" default:"https://api-ssx.kameleoon.com"`
	TrackingVersion      string        `yml:"tracking_version" yaml:"tracking_version"`
	ProxyURL             string        `yml:"proxy_url" yaml:"proxy_url"`
	ClientID             string        `yml:"client_id" yaml:"client_id"`
	ClientSecret         string        `yml:"client_secret" yaml:"client_secret"`
	Version              string        `yml:"version" yaml:"version"`
	ConfigUpdateInterval time.Duration `yml:"config_update_interval" yaml:"config_update_interval" default:"1h"`
	Timeout              time.Duration `yml:"timeout" yaml:"timeout" default:"2s"`
	VisitorDataMaxSize   int           `yml:"visitor_data_max_size" yaml:"visitor_data_max_size"`
	BlockingMode         bool          `yml:"blocking_mode" yaml:"blocking_mode"`
	VerboseMode          bool          `yml:"verbose_mode" yaml:"verbose_mode"`
}

func LoadConfig(path string) (*Config, error) {
	c := &Config{}
	return c, c.Load(path)
}

func (c *Config) defaults() {
	if c.Logger == nil {
		c.Logger = defaultLogger
	}
	if len(c.TrackingURL) == 0 {
		c.TrackingURL = API_SSX_URL
	}
	if c.ConfigUpdateInterval == 0 {
		c.ConfigUpdateInterval = DefaultConfigUpdateInterval
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultRequestTimeout
	}
	if c.VisitorDataMaxSize == 0 {
		c.VisitorDataMaxSize = DefaultVisitorDataMaxSize
	}
	if len(c.Version) == 0 {
		c.Version = sdkVersion
	}
	if len(c.TrackingVersion) == 0 {
		c.TrackingVersion = DefaultTrackingVersion
	}
	c.REST.defaults(c.Version)
}

func (c *Config) Load(path string) error {
	if len(path) == 0 {
		path = DefaultConfigPath
	}
	err := c.loadFile(path)
	c.defaults()
	return err
}

func (c *Config) loadFile(configPath string) error {
	yml := aconfigyaml.New()
	loader := aconfig.LoaderFor(c, aconfig.Config{
		SkipFlags:          true,
		SkipEnv:            true,
		FailOnFileNotFound: true,
		AllowUnknownFields: true,
		Files:              []string{configPath},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": yml,
			".yml":  yml,
		},
	})

	return loader.Load()
}

const (
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 5 * time.Second
	DefaultDoTimeout       = 10 * time.Second
	DefaultMaxConnsPerHost = 10000
)

type RestConfig struct {
	ProxyURL        string
	UserAgent       string
	DoTimeout       time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxConnsPerHost int
}

func (c *RestConfig) defaults(version string) {
	var b strings.Builder
	b.WriteString(UserAgent)
	b.WriteString(version)
	c.UserAgent = b.String()

	if c.ReadTimeout == 0 {
		c.ReadTimeout = DefaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = DefaultWriteTimeout
	}
	if c.DoTimeout == 0 {
		c.DoTimeout = DefaultDoTimeout
	}
	if c.MaxConnsPerHost == 0 {
		c.MaxConnsPerHost = DefaultMaxConnsPerHost
	}
}
