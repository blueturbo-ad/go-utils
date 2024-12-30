package requestidmanager

/*
获取request id的功能
1. GenerateRequestId 获取request 的主id
1.1 GenerateRequestId 当RequestBody 为空时就不会携带RequestBody 进行计算requestid 如果RequestBody 不为空那么就会携带RequestBody 进行request id的计算
1.2 GenerateRequestSubId 获取到request 的子id， 该id 会继承父id进行递增操作。
1.3 调用GenerateRequestSubId的时候必须保证RequestTool 结构体在同一个生命周期内
*/
import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"git.domob-inc.cn/blueturbo/go-utils.git/global"
)

const request_id_format string = "%s-%s-%04d"
const request_sub_id_format string = "%s-%d"

var counterPerMicroSecond uint64

const counterStep int = 1
const sha1MinLength int = 32
const time2IntBase int = 16
const counterPerMicroSecondStep uint64 = 1

type RequestIdGenerator struct {
	hashfrom  string
	requestId string
	counter   int
}

func NewRequestIdGenerator(hashFrom string) *RequestIdGenerator {
	return &RequestIdGenerator{hashfrom: hashFrom, counter: 0}
}

func (r *RequestIdGenerator) GenerateRequestId() (string, error) {
	now := time.Now().UnixMicro()
	hexTime := strconv.FormatInt(now, time2IntBase)
	sha1 := Sha1Encode(r.hashfrom)
	if len(sha1) <= sha1MinLength {
		return global.EmptyString, fmt.Errorf("sha1 encode error")
	}
	sha1 = sha1[sha1MinLength:]
	curCounterPerMicroSecond := atomic.AddUint64(&counterPerMicroSecond, counterPerMicroSecondStep)
	r.requestId = fmt.Sprintf(request_id_format, hexTime, sha1, curCounterPerMicroSecond)
	return r.requestId, nil
}

func (r *RequestIdGenerator) GenerateRequestSubId(requestId string) (string, error) {
	if requestId == global.EmptyString {
		return global.EmptyString, fmt.Errorf("request id is empty")
	}
	now := time.Now().UnixMicro()
	subId := fmt.Sprintf(request_sub_id_format, requestId, now)
	return subId, nil
}

// todo 要提出去
func Sha1Encode(input string) string {
	c := sha1.New()
	c.Write([]byte(input))
	bytes := c.Sum(nil)
	return hex.EncodeToString(bytes)
}
