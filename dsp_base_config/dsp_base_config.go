package dspbaseconfig

import (
	"fmt"
	"sync"

	"github.com/blueturbo-ad/go-utils/config_manage"
)

type DspBaseConfigManage struct {
	Hooks []func(config string) error
}

var (
	instance    *DspBaseConfigManage
	once        sync.Once
	EmptyString = ""
)

func GetSingleton() *DspBaseConfigManage {
	once.Do(func() {
		instance = &DspBaseConfigManage{}

	})
	return instance
}

func (d *DspBaseConfigManage) RegistHookFunc(f func(config string) error) {
	d.Hooks = append(d.Hooks, f)
}

func (d *DspBaseConfigManage) LoadK8sConfigMap(configMapName, env string) error {
	var e = new(config_manage.GeneralConfigManage)
	err := e.LoadK8sConfigMap(configMapName, env)
	if err != nil {
		return fmt.Errorf("dsp base config  LoadK8sConfigMap is error %s", err.Error())
	}
	if err := d.RefreshConfig(e); err != nil {
		return fmt.Errorf("dsp base config  LoadConfig is error %s", err.Error())
	}
	return nil
}

func (d *DspBaseConfigManage) LoadFileConfig(filePath, env string) error {
	var e = new(config_manage.GeneralConfigManage)
	err := e.LoadConfig(filePath, env)
	if err != nil {
		return fmt.Errorf("dsp base config  LoadFileConfig is error %s", err.Error())
	}
	if err := d.RefreshConfig(e); err != nil {
		return fmt.Errorf("dsp base config  LoadConfig is error %s", err.Error())
	}
	return nil
}

func (d *DspBaseConfigManage) RefreshConfig(e *config_manage.GeneralConfigManage) error {
	for _, hook := range d.Hooks {
		err := hook(e.Config)
		if err != nil {
			return fmt.Errorf("hook func error %s", err.Error())
		}
	}
	return nil
}
