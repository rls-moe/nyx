package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var configFileName = "./config.yml"

func init() {
	flag.StringVar(&configFileName, "config", "./config.yml", "Config File Location")
}

type Config struct {
	Site              SiteConfig    `yaml:"site"`             // Site/HTML Configuration
	DB                DBConfig      `yaml:"db"`               // Database Configuration
	HostnameWhiteList []string      `yaml:"hosts"`            // List of valid hostnames, ignored if empty
	ListenOn          string        `yaml:"listen_on"`        // Address & Port to use
	MasterSecret      string        `yaml:"secret"`           // Master Secret for keychain
	DisableSecurity   bool          `yaml:"disable_security"` // Disables various flags to ensure non-HTTPS requests work
	Captcha           CaptchaConfig `yaml:"captcha"`
}

const (
	CaptchaRecaptcha = "recaptcha"
	CaptchaInternal  = "internal"
	CaptchaHybrid    = "hybrid"
	CaptchaDisabled  = "disabled"
)

type CaptchaConfig struct {
	Mode     string            `yaml:"mode"` // Captcha Mode
	Settings map[string]string `yaml:"settings,inline"`
}

type SiteConfig struct {
	Title        string `yaml:"title"`       // Site Title
	Description  string `yaml:"description"` // Site Description
	PrimaryColor string `yaml:"color"`       // Primary Color for Size
}

type DBConfig struct {
	File     string `yaml:"file"`
	ReadOnly bool   `yaml:"read_only"`
}

func Load() (*Config, error) {
	var config = &Config{
		Site: SiteConfig{
			Title:        "NyxChan",
			PrimaryColor: "#78909c",
			Description:  "NyxChan Default Configuration",
		},
		DB: DBConfig{
			File:     ":memory:",
			ReadOnly: false,
		},
		HostnameWhiteList: []string{},
		ListenOn:          ":8080",
		MasterSecret:      "changeme",
		DisableSecurity:   true,
		Captcha: CaptchaConfig{
			Mode: CaptchaDisabled,
		},
	}
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return config, nil
	}
	dat, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(dat, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c Config) IsHostNameValid(hostname string) bool {
	if c.HostnameWhiteList == nil {
		return true
	}
	if len(c.HostnameWhiteList) == 0 {
		return true
	}
	for _, v := range c.HostnameWhiteList {
		if v == hostname {
			return true
		}
	}
	return false
}
