package kameleoon

import (
	"time"

	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/network"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

const (
	DefaultConfigPath      = "/etc/kameleoon/client-go.yaml"
	DefaultRefreshInterval = time.Hour
	DefaultRequestTimeout  = 10 * time.Second
	DefaultSessionDuration = 30 * time.Minute
	DefaultEnvironment     = ""
)

type KameleoonClientConfig struct {
	defaultsApplied  bool
	Network          NetworkConfig
	Logger           logging.Logger `yml:"-" yaml:"-"`
	ProxyURL         string         `yml:"proxy_url" yaml:"proxy_url"`
	ClientID         string         `yml:"client_id" yaml:"client_id"`
	ClientSecret     string         `yml:"client_secret" yaml:"client_secret"`
	RefreshInterval  time.Duration  `yml:"refresh_interval" yaml:"refresh_interval" default:"1h"`
	DefaultTimeout   time.Duration  `yml:"default_timeout" yaml:"default_timeout" default:"10s"`
	VerboseMode      bool           `yml:"verbose_mode" yaml:"verbose_mode"`
	SessionDuration  time.Duration  `yml:"session_duration" yaml:"session_duration" default:"30m"`
	TopLevelDomain   string         `yml:"top_level_domain" yaml:"top_level_domain"`
	Environment      string         `yml:"environment" yaml:"environment"`
	UserAgentMaxSize int            `yml:"-" yaml:"-"`
	dataApiUrl       string
}

func LoadConfig(path string) (*KameleoonClientConfig, error) {
	c := &KameleoonClientConfig{}
	return c, c.Load(path)
}

func (c *KameleoonClientConfig) defaults() error {
	if c.defaultsApplied {
		return nil
	}
	c.defaultsApplied = true

	if len(c.ClientID) == 0 {
		return errs.NewConfigCredentialsInvalid("Client ID is not specified")
	}
	if len(c.ClientSecret) == 0 {
		return errs.NewConfigCredentialsInvalid("Client secret is not specified")
	}

	if c.Logger == nil {
		c.defaultLogger()
	}
	if len(c.dataApiUrl) == 0 {
		c.dataApiUrl = network.DefaultDataApiUrl
	}
	if c.RefreshInterval < time.Minute {
		if c.RefreshInterval != 0 {
			c.Logger.Printf("Kameloon SDK: Config update interval must not be less than a minute."+
				"Default config update interval (%d minutes) is applied", int(DefaultRefreshInterval.Minutes()))
		}
		c.RefreshInterval = DefaultRefreshInterval
	}
	if c.DefaultTimeout <= 0 {
		if c.DefaultTimeout != 0 {
			c.Logger.Printf("Kameloon SDK: Default timeout must have positive value."+
				"Default default timeout (%d ms) is applied", DefaultRequestTimeout.Milliseconds())
		}
		c.DefaultTimeout = DefaultRequestTimeout
	}
	if c.SessionDuration <= 0 {
		if c.SessionDuration != 0 {
			c.Logger.Printf("Kameloon SDK: Session duration must have positive value."+
				"Default session duration (%d minutes) is applied", int(DefaultSessionDuration.Minutes()))
		}
		c.SessionDuration = DefaultSessionDuration
	}
	if len(c.TopLevelDomain) == 0 {
		c.Logger.Printf("Kameleoon SDK: Setting top level domain is strictly recommended, " +
			"otherwise you may have problems when using subdomains.")
	}
	if len(c.Environment) == 0 {
		c.Environment = DefaultEnvironment
	}
	return c.Network.defaults()
}

func (c *KameleoonClientConfig) defaultLogger() {
	loggerMode := logging.Silent
	if c.VerboseMode {
		loggerMode = logging.Verbose
	}
	c.Logger = logging.NewLogger(loggerMode, logging.DefaultLogger)
}

func (c *KameleoonClientConfig) Load(path string) error {
	if len(path) == 0 {
		path = DefaultConfigPath
	}
	if err := c.loadFile(path); err != nil {
		return err
	}
	return c.defaults()
}

func (c *KameleoonClientConfig) loadFile(configPath string) error {
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
	DoTimeout       time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxConnsPerHost int
}

func (c *NetworkConfig) defaults() error {
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
	return nil
}
