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
	attr *RtaAttrConf
}

// 修正结构体定义
type RtaAttrConf struct {
	Version string                `yaml:"version"`
	RtaAttr map[string]RegionConf `yaml:"rta_attr"`
}

type RegionConf map[string]EndpointConf

type EndpointConf struct {
	Url string `yaml:"url"`
	Ak  string `yaml:"ak,omitempty"`
	Sk  string `yaml:"sk,omitempty"`
}

func GetSingleton() *RtaAttrConfig {
	once.Do(func() {
		instance = NewRtaAttrConfig()

	})
	return instance
}

func NewRtaAttrConfig() *RtaAttrConfig {
	return &RtaAttrConfig{
		attr: &RtaAttrConf{},
	}
}

func (r *RtaAttrConfig) Reload(conf string) error {
	var fullConfig RtaAttrConf
	if err := yaml.Unmarshal([]byte(conf), &fullConfig); err != nil {
		zap_loggerex.GetSingleton().Error("bid_stdout_logger", "failed to unmarshal yaml %s, %+v", string(conf), err)
		return err
	}
	r.attr = &fullConfig
	return nil
}

func (r *RtaAttrConfig) GetRtaAttrConf(source string) RegionConf {

	attr := r.attr.RtaAttr
	if attr == nil {
		return nil
	}
	regionConf, ok := attr[source]
	if !ok {
		return nil
	}
	return regionConf
}
