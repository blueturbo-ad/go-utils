package requestidmanager

import (
	"strconv"
	"time"
)

type RequestIdExpiredChecker struct {
	RequestId string
}

func NewRequestIdExpiredChecker(RequestId string) *RequestIdExpiredChecker {
	return &RequestIdExpiredChecker{RequestId: RequestId}
}

func (r *RequestIdExpiredChecker) Check(exp string) (bool, error) {
	hexMsTime, _, _, err := NewRequestIdParser(r.RequestId).Parse()
	if err != nil {
		return false, err
	}
	n, err := strconv.ParseUint(hexMsTime, 16, 64)
	if err != nil {
		return false, err
	}
	exptime, err := strconv.Atoi(exp)
	if err != nil {
		return false, err
	}
	n /= 1e6
	nowstamp := uint64(time.Now().Unix())
	return nowstamp >= n+uint64(exptime), nil
}
