package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pheelee/dyndns/pkg/dns"
	"github.com/pheelee/dyndns/pkg/model"
	"gopkg.in/yaml.v2"
)

// Format defines the config serializer
var Format string = ".yaml"

// Config represents the apps config file
type Config struct {
	Global struct {
		TimeZone     string
		Listen       int
		RealIPHeader string
	}
	HTTPReqAuth struct {
		Username string
		Password string
	}
	DNSServer struct {
		Active dns.Server `yaml:"-" toml:"-"`
		Bind   dns.Bind
	}

	Domains map[string]model.Domain
}

// New initializes an example config
func (c *Config) New() {
	c.Global.TimeZone = "Europe/Zurich"
	c.Global.Listen = 8090
	c.Domains = make(map[string]model.Domain)
	c.Global.RealIPHeader = "X-Forwarded-For"
	c.Domains["dyn.example.com"] = model.Domain{Credentials: "2DE70679C33C9B652F0FCCE1DCD62124020BDCB5676E113C40C969050EA60074"} // dynuser / dynpass
	c.DNSServer.Bind = dns.Bind{Enabled: true, Host: "localhost"}
}

// Save writes the config at path
func (c *Config) Save(path string) {
	var (
		err error
		f   *os.File
	)
	defer f.Close()
	switch Format {
	case ".yaml":
		if f, err = os.Create(path); err != nil {
			panic(fmt.Sprintf("Could not open file, %s", err))
		}
		enc := yaml.NewEncoder(f)
		err = enc.Encode(c)
		if err != nil {
			panic(fmt.Sprintf("Could not write config file, %s", err))
		}
	case ".toml":
		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		if err := enc.Encode(c); err != nil {
			panic(fmt.Sprintf("Could not write config file, %s", err))
		}
		ioutil.WriteFile(path, buf.Bytes(), 0644)
	default:
		panic(fmt.Sprintf("Encoder unknown %s", Format))
	}

}

// Load reads the config from path
func (c *Config) Load(path string) {
	var (
		err error
		f   *os.File
	)
	defer f.Close()
	Format = filepath.Ext(path)
	// config exists read it
	if _, err = os.Stat(path); err == nil {
		switch Format {
		case ".yaml":
			if f, err = os.Open(path); err != nil {
				panic(fmt.Sprintf("Could not open file, %s", err))
			}
			d := yaml.NewDecoder(f)
			if d.Decode(c) != nil {
				panic("Could not decode config using yaml decoder")
			}
		case ".toml":
			if _, err := toml.DecodeFile(path, c); err != nil {
				panic(fmt.Sprintf("Could parse config file %s", path))
			}
		default:
			panic("Config format not recognized, must be either yaml or toml")
		}

	} else {
		// config does not exist create default and write it
		c.New()
	}

	c.HTTPReqAuth.Username = os.Getenv("HTTPREQ_USERNAME")
	c.HTTPReqAuth.Password = os.Getenv("HTTPREQ_PASSWORD")
	c.Save(path)

	// set the active DNS server
	if c.DNSServer.Bind.Enabled {
		c.DNSServer.Active = &c.DNSServer.Bind
	}
}

// TimeNow gets the current time based on the location specified in the config or env
func (c *Config) TimeNow() time.Time {
	var (
		loc *time.Location
		err error
	)

	if loc, err = time.LoadLocation(c.Global.TimeZone); err == nil {
		return time.Now().In(loc)
	}
	if loc, err = time.LoadLocation(os.Getenv("TZ")); err == nil {
		return time.Now().In(loc)
	}
	return time.Now()
}
