package kameleoon

import (
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v2/logging"
	"github.com/Kameleoon/client-go/v2/network"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

const (
	DefaultConfigPath           = "/etc/kameleoon/client-go.yaml"
	DefaultConfigUpdateInterval = time.Hour
	DefaultRequestTimeout       = 2 * time.Second
	DefaultVisitorDataMaxSize   = 500 // 500 mb
	ClientHeader                = "sdk/" + SdkLanguage + "/"
	DefaultEnvironment          = ""
	DefaultUserAgentMaxSize     = 100_000
)

type Config struct {
	Network              NetworkConfig
	Logger               logging.Logger `yml:"-" yaml:"-"`
	SiteCode             string         `yml:"site_code" yaml:"site_code"`
	ProxyURL             string         `yml:"proxy_url" yaml:"proxy_url"`
	ClientID             string         `yml:"client_id" yaml:"client_id"`
	ClientSecret         string         `yml:"client_secret" yaml:"client_secret"`
	ConfigUpdateInterval time.Duration  `yml:"config_update_interval" yaml:"config_update_interval" default:"1h"`
	Timeout              time.Duration  `yml:"timeout" yaml:"timeout" default:"2s"`
	VisitorDataMaxSize   int            `yml:"visitor_data_max_size" yaml:"visitor_data_max_size"`
	VerboseMode          bool           `yml:"verbose_mode" yaml:"verbose_mode"`
	Environment          string         `yml:"environment" yaml:"environment"`
	UserAgentMaxSize     int            `yml:"-" yaml:"-"`
	dataApiUrl           string
}

func LoadConfig(path string) (*Config, error) {
	c := &Config{}
	return c, c.Load(path)
}

func (c *Config) defaults() {
	if len(c.dataApiUrl) == 0 {
		c.dataApiUrl = network.DefaultDataApiUrl
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

	if len(c.Environment) == 0 {
		c.Environment = DefaultEnvironment
	}
	if c.UserAgentMaxSize == 0 {
		c.UserAgentMaxSize = DefaultUserAgentMaxSize
	}
	c.Network.defaults(SdkVersion)
	if c.Logger == nil {
		c.defaultLogger()
	}
}

func (c *Config) defaultLogger() {
	loggerMode := logging.Silent
	if c.VerboseMode {
		loggerMode = logging.Verbose
	}
	c.Logger = logging.NewLogger(loggerMode, logging.DefaultLogger)
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

type NetworkConfig struct {
	ProxyURL        string
	KameleoonClient string
	DoTimeout       time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxConnsPerHost int
}

func (c *NetworkConfig) defaults(version string) {
	var b strings.Builder
	b.WriteString(ClientHeader)
	b.WriteString(version)
	c.KameleoonClient = b.String()

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
