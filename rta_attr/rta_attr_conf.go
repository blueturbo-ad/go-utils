package rtaattr

import (
	"fmt"
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
	instance    *RtaAttrConf
	once        sync.Once
	EmptyString = ""
)

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

func GetSingleton() *RtaAttrConf {
	once.Do(func() {
		instance = NewRtaAttrConf()

	})
	return instance
}

func NewRtaAttrConf() *RtaAttrConf {
	return &RtaAttrConf{
		RtaAttr: make(map[string]RegionConf, 0),
	}
}

func (r *RtaAttrConf) Reload(conf string, env string) error {
	var fullConfig FullConfig
	if err := yaml.Unmarshal([]byte(conf), &fullConfig); err != nil {
		zap_loggerex.GetSingleton().Error("bid_stdout_logger", "failed to unmarshal yaml %s, %+v", string(conf), err)
		return err
	}
	switch env {
	case "Dev":
		*r = fullConfig.Dev
	case "Pro":
		*r = fullConfig.Pro
	default:
		zap_loggerex.GetSingleton().Error("bid_stdout_logger", "unknown environment: %s", env)
		return fmt.Errorf("unknown environment: %s", env)
	}
	return nil
}

func (r *RtaAttrConf) GetRtaAttrConf(source string) RegionConf {
	if r.RtaAttr == nil {
		return nil
	}
	if val, ok := r.RtaAttr[source]; ok {
		return val
	}
	return nil
}
