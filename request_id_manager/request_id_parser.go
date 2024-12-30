package requestidmanager

import (
	"fmt"
	"strings"

	"git.domob-inc.cn/blueturbo/go-utils.git/global"
)

const request_id_formatter_split_count int = 3

type RequestIdParser struct {
	RequestId string
}

func NewRequestIdParser(RequestId string) *RequestIdParser {
	return &RequestIdParser{RequestId: RequestId}
}

func (c *RequestIdParser) Parse() (string, string, string, error) {
	splits := strings.Split(c.RequestId, "-")
	if splits == nil || len(splits) != request_id_formatter_split_count {
		return global.EmptyString, global.EmptyString, global.EmptyString, fmt.Errorf("request sub format error %s", c.RequestId)
	}
	hexMsTime := splits[0]
	sha1 := splits[1]
	counterPerMicroSecondStep := splits[2]
	return hexMsTime, sha1, counterPerMicroSecondStep, nil
}
