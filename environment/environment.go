package environment

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blueturbo-ad/go-utils/global"
)

const (
	// os environment KEY list
	KeyEnvironment = "ENVIRONMENT"
	KeyRegion      = "REGION"
	KeyIsPod       = "Pod"
	KeyPodIp       = "POD_IP"
	KeyPodName     = "POD_NAME"
	KeyNameSpace   = "POD_NAMESPACE"
	KeyNodeName    = "NODE_NAME"
	KeyPreVersion  = "PRE_VERSION"
	// error list
	ErrorWorkPathNotFound = "work path error path is %s"
	// fatal list
	FatalPathNotInside = "path not inside %s"
)

const (
	KeyPro  = "Pro"
	KeyPre  = "Pre"
	KeyTest = "Test"
	KeyDev  = "Dev"
)

func Init() {
	GetSingleton().initEnvironment()
}

type Environment struct {
	env          string // 运行环境
	region       string // 运行地区
	workPath     string // 运行路径
	confPath     string // 配置文件路径
	workName     string // 运行名称
	preVersion   string // 区分pre代码版本是新版本还是旧版本
	instanceIp   string
	macAddresses []string
	podInfo      *PodInfo
	uniqueId     string
	namespace    string
}

type PodInfo struct {
	PodIp     string
	PodName   string
	NameSpace string
	NodeName  string
}

var (
	instance *Environment
	once     sync.Once
)

func GetSingleton() *Environment {
	once.Do(func() {
		instance = new(Environment)
	})
	return instance
}

func WorkName() string {
	return GetSingleton().GetWorkName()
}

func GetEnv() string {
	return GetSingleton().GetEnv()
}

func GetRegion() string {
	return GetSingleton().Region()
}

func GetWorkPath() string {
	return GetSingleton().WorkPath()
}

func GetConfPath() string {
	return GetSingleton().ConfPath()
}

func GetPreVersion() string {
	return GetSingleton().GetPreVersion()
}

func GetPodInfo() *PodInfo {
	return GetSingleton().GetPodInfo()
}

func GetMacAddresses() []string {
	return GetSingleton().GetMacAddresses()
}

func GetUniqueId() string {
	return GetSingleton().GetUniqueId()
}

func GetIntanceIp() string {
	return GetSingleton().GetIntanceIp()
}
func GetPodNameInfo() string {
	return GetSingleton().GetPodName()
}

func GetPodNameSpace() string {
	return GetSingleton().GetNamespace()
}

func (e *Environment) GetWorkName() string {
	return e.workName
}

func (e *Environment) GetIntanceIp() string {
	return e.instanceIp
}

func (e *Environment) GetEnv() string {
	return e.env
}
func (e *Environment) Region() string {
	return e.region
}

func (e *Environment) WorkPath() string {
	return e.workPath
}

func (e *Environment) ConfPath() string {
	return e.confPath
}

func (e *Environment) GetPreVersion() string {
	if e.env != KeyPre {
		return ""
	}
	return e.preVersion
}

func (e *Environment) GetPodInfo() *PodInfo {
	return e.podInfo
}

func (e *Environment) GetMacAddresses() []string {
	return e.macAddresses
}

func (e *Environment) GetUniqueId() string {
	return e.uniqueId
}
func (e *Environment) GetPodName() string {
	return e.podInfo.PodName
}

func (e *Environment) GetNamespace() string {
	return e.podInfo.NameSpace
}

func (e *Environment) initEnvironment() {
	e.initEnv()
	e.initRegion()
	e.initWorkName()
	e.initWorkPath()
	e.initConfPath()
	e.initPreVersion()
	e.initPodInfo()
	e.initIntranetIP()
	e.initMacAddresses()
	e.initPodName()
	e.initNameSpace()
	podInfoByte, err := json.Marshal(e.podInfo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("env: %s, region: %s, workPath: %s, confPath: %s workName: %s, instanceIp: %s, macAddresses: %s, podInfo: %s uniqueId: %s\n", e.env, e.region, e.workPath, e.confPath, e.workName, e.instanceIp, strings.Join(e.macAddresses, ","), string(podInfoByte), e.uniqueId)
	if e.env == KeyPre { // 如果为pre环境则输出一下pre version的信息
		fmt.Printf("[INFO] env=[%s] and pre version=[%s]", e.env, e.preVersion)
	}
}

func (e *Environment) initWorkName() {
	e.workName = os.Getenv("WORK_NAME")
}

func (e *Environment) initEnv() {
	env := os.Getenv(KeyEnvironment)
	if env == global.EmptyString {
		env = global.EnvDev
	}
	e.env = env
}

func (e *Environment) initRegion() {
	region := os.Getenv(KeyRegion)
	if region == global.EmptyString {
		region = global.RegionDev
	}
	e.region = region
}

func (e *Environment) initWorkPath() {
	runPath, err := os.Getwd()
	workNameWithSlash := filepath.Join(e.workName)
	idx := strings.LastIndex(runPath, workNameWithSlash)
	if idx == -1 {
		panic(fmt.Sprintf(FatalPathNotInside, workNameWithSlash))
	}
	runPath = runPath[0 : idx+len(workNameWithSlash)]
	if err != nil {
		msg := fmt.Sprintf(ErrorWorkPathNotFound, runPath)
		panic(msg)
	}
	e.workPath = runPath
}

// initConfPath 根据地区不同，初始化conf path
func (e *Environment) initConfPath() {
	var (
		workPath = e.WorkPath()
		region   = e.Region()
	)

	if region == "" {
		panic(fmt.Sprintf("[PANIC] Init confPath failed, please check env:REGION, region: [%s]", region))
	}

	e.confPath = filepath.Join(workPath, "conf", region)
}

func (e *Environment) initPreVersion() {
	if e.GetEnv() != KeyPre {
		return
	}
	e.preVersion = os.Getenv(KeyPreVersion)
	if e.preVersion == "" {
		panic("[PANIC] please check pre version because env=Pre")
	}
}

func (e *Environment) initPodInfo() {
	podInfo := PodInfo{
		PodIp:     os.Getenv(KeyPodIp),
		PodName:   os.Getenv(KeyPodName),
		NameSpace: os.Getenv(KeyNameSpace),
		NodeName:  os.Getenv(KeyNodeName),
	}
	e.podInfo = &podInfo
	e.uniqueId = fmt.Sprintf("%s_%s_%s", podInfo.NodeName, podInfo.NameSpace, podInfo.PodName)
}

func (e *Environment) initPodName() {
	e.podInfo.PodName = os.Getenv(KeyPodName)
}

func (e *Environment) initNameSpace() {
	e.podInfo.NameSpace = os.Getenv(KeyNameSpace)
}

func (e *Environment) initIntranetIP() {
	netInterfaces, err := net.InterfaceAddrs()
	if err != nil {

		panic(err)
	}

	for _, address := range netInterfaces {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				e.instanceIp = ipnet.IP.String()
				return
			}
		}
	}
}

func (e *Environment) initMacAddresses() {
	var macAddresses = make([]string, 0)
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddresses = append(macAddresses, macAddr)
	}
	e.macAddresses = macAddresses
}
