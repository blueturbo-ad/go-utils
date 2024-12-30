package requestidmanager

import (
	"fmt"
	"strconv"
	"strings"

	"git.domob-inc.cn/blueturbo/go-utils.git/global"
)

const request_sub_id_formatter_split_count int = 4

type RequestSubIdParser struct {
	RequestSubId string
}

func NewRequestSubIdParser(RequestSubId string) *RequestSubIdParser {
	return &RequestSubIdParser{RequestSubId: RequestSubId}
}

func (c *RequestSubIdParser) Parse() (string, string, error) {
	splits := strings.Split(c.RequestSubId, "-")
	if splits == nil || len(splits) != request_sub_id_formatter_split_count {
		return global.EmptyString, global.EmptyString, fmt.Errorf("request sub format error %s", c.RequestSubId)
	}
	curCounterPerMicroSecond, err := strconv.ParseInt(splits[2], 10, 64)
	if err != nil {
		return global.EmptyString, global.EmptyString, err
	}
	requestId := fmt.Sprintf(request_id_format, splits[0], splits[1], curCounterPerMicroSecond)
	requestSubId := splits[3]
	return requestId, requestSubId, nil
}
