package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/blueturbo-ad/go-utils/config_manage"
)

// 小于等于0 是 warning error fail 大于等于1 error fail  大于等于2 fail
type FeishuManage struct {
	Config [2]*config_manage.FeishuConfig
	index  int
	lock   sync.RWMutex
}

var (
	instance *FeishuManage
	once     sync.Once
)

type CardMessage struct {
	MsgType string `json:"msg_type"`
	Card    struct {
		Config struct {
			WideScreenMode bool `json:"wide_screen_mode"`
		} `json:"config"`
		Header struct {
			Title struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"title"`
			Template string `json:"template"`
		} `json:"header"`
		Elements []struct {
			Tag  string `json:"tag"`
			Text struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"text"`
		} `json:"elements"`
	} `json:"card"`
}

const (
	TitleContent = "Dsp告警级别：%s"
	Tag          = "plain_text"
	MsgType      = "interactive"
)

var ColorMap = map[string]string{
	"warning": "orange",
	"error":   "red",
	"fail":    "carmine",
}

type FeishuResponse struct {
	StatusCode    int      `json:"StatusCode"`
	StatusMessage string   `json:"StatusMessage"`
	Code          int      `json:"code"`
	Data          struct{} `json:"data"`
	Msg           string   `json:"msg"`
}

func GetInstance() *FeishuManage {

	once.Do(func() {
		instance = &FeishuManage{
			Config: [2]*config_manage.FeishuConfig{new(config_manage.FeishuConfig), new(config_manage.FeishuConfig)},
			index:  -1,
			lock:   sync.RWMutex{},
		}
	})

	return instance
}

// 函数用于内存更新etcd配置
func (l *FeishuManage) UpdateFromEtcd(env string, eventType string, key string, value string) error {
	fmt.Printf("Event Type: %s, Key: %s, Value: %s\n", eventType, key, value)

	var err error
	switch key {
	case "logger":
		var e = new(config_manage.FeishuConfig)
		err = e.LoadMemoryConfig([]byte(value), env)
		if err != nil {
			return err
		}
		return l.UpdateLogger(e)
	default:
		return fmt.Errorf("unknown UpdateFromEtcd: %s", key)
	}
}

func (f *FeishuManage) UpdateFromFile(confPath string, env string) error {
	var err error
	var e = new(config_manage.FeishuConfig)
	err = e.LoadConfig(confPath, env)
	if err != nil {
		return err
	}

	return f.UpdateLogger(e)
}

func (l *FeishuManage) UpdateLogger(config *config_manage.FeishuConfig) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.index = (l.index + 1) % 2
	l.Config[l.index] = config
	return nil
}
func (l *FeishuManage) GetConfig() *config_manage.FeishuConfig {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.Config[l.index]
}

func (l *FeishuManage) Send(errType, message string) error {
	conf := l.GetConfig()
	messageData := l.BuildFeishuCare(errType, message)

	msg, err := json.Marshal(messageData)
	if err != nil {
		return fmt.Errorf("failed to marshal card message: %v", err)
	}

	response, err := http.Post(conf.Config.Url, "application/json", bytes.NewBuffer(msg))
	if err != nil {
		return fmt.Errorf("failed to send card message: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", response.StatusCode)
	}
	body := new(bytes.Buffer)
	_, err = body.ReadFrom(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	// 解析响应
	var feishuResp FeishuResponse
	err = json.Unmarshal(body.Bytes(), &feishuResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %v", err)
	}
	if feishuResp.Code != 0 {
		return fmt.Errorf("received non-0 code: %d, msg: %s", feishuResp.Code, feishuResp.Msg)
	}
	return nil
}

func (l *FeishuManage) BuildFeishuCare(errType, message string) *CardMessage {
	cardMessage := new(CardMessage)
	title := fmt.Sprintf(TitleContent, errType)
	color, ok := ColorMap[errType]
	if !ok {
		color = "red"
	}
	cardMessage.MsgType = MsgType
	cardMessage.Card.Config.WideScreenMode = true
	cardMessage.Card.Header.Title.Tag = Tag
	cardMessage.Card.Header.Title.Content = title
	cardMessage.Card.Header.Template = color // 设置卡片头部颜色
	cardMessage.Card.Elements = []struct {
		Tag  string `json:"tag"`
		Text struct {
			Tag     string `json:"tag"`
			Content string `json:"content"`
		} `json:"text"`
	}{
		{
			Tag: "div",
			Text: struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			}{
				Tag:     "lark_md",
				Content: message,
			},
		},
	}
	return cardMessage
}
