package rtaattr

import (
	"sync"

	"github.com/blueturbo-ad/go-utils/zap_loggerex"
	"github.com/stretchr/testify/assert/yaml"
)

/*
curused: Dev
Dev:

	version: "1.0.0"
	rta_attr:
	  -
	    tiktok:
	      us:
	        url: "https://growth-rta.tiktokv-us.com"
	      eu:
	        url: "https://growth-rta.tiktokv-eu.com"
	      row:
	        url: "https://growth-rta.byteintl.com"

Pro:

	version: "1.0.0"
	rta_attr:
	  -
	    tiktok:
	      us:
	        url: "https://growth-rta.tiktokv-us.com"
	      eu:
	        url: "https://growth-rta.tiktokv-eu.com"
	      row:
	        url: "https://growth-rta.byteintl.com"
*/

var (
	instance    *RtaAttrConfig
	once        sync.Once
	EmptyString = ""
)

type RtaAttrConfig struct {
	attr *FullConfig
}

// 修正结构体定义
type RtaAttrConf struct {
	Version string                `yaml:"version"`
	RtaAttr map[string]RegionConf `yaml:"rta_attr"`
}

type RegionConf map[string]EndpointConf

type EndpointConf struct {
	Url string `yaml:"url"`
}

// 完整的配置结构，包含环境配置
type FullConfig struct {
	CurUsed string      `yaml:"curused"`
	Dev     RtaAttrConf `yaml:"Dev"`
	Pro     RtaAttrConf `yaml:"Pro"`
}

func GetSingleton() *RtaAttrConfig {
	once.Do(func() {
		instance = NewRtaAttrConfig()

	})
	return instance
}

func NewRtaAttrConfig() *RtaAttrConfig {
	return &RtaAttrConfig{
		attr: &FullConfig{},
	}
}

func (r *RtaAttrConfig) Reload(conf string) error {
	var fullConfig FullConfig
	if err := yaml.Unmarshal([]byte(conf), &fullConfig); err != nil {
		zap_loggerex.GetSingleton().Error("bid_stdout_logger", "failed to unmarshal yaml %s, %+v", string(conf), err)
		return err
	}
	r.attr = &fullConfig
	return nil
}

func (r *RtaAttrConfig) GetRtaAttrConf(source string, env string) RegionConf {
	var conf *RtaAttrConf
	if env == "" {
		env = r.attr.CurUsed
	}
	switch env {
	case "Dev":
		conf = &r.attr.Dev
	case "Pro":
		conf = &r.attr.Pro
	default:
		zap_loggerex.GetSingleton().Error("bid_stdout_logger", "invalid environment: %s", env)
		return nil
	}

	if conf == nil {
		zap_loggerex.GetSingleton().Error("bid_stdout_logger", "configuration for environment %s is nil", env)
		return nil
	}

	endpointConf, exists := conf.RtaAttr[source]
	if !exists {
		zap_loggerex.GetSingleton().Error("bid_stdout_logger", "source %s not found in environment %s", source, env)
		return nil
	}

	return endpointConf
}
