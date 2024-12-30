package requestidmanager

import (
	"strconv"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequestIdExpiredChecker(t *testing.T) {
	t.Run("input=RequestIdExpiredChecker", func(t *testing.T) {
		requestId, err := NewRequestIdGenerator("1234567890").GenerateRequestId()
		assert.Nil(t, err)
		requestIdExpireChecker := NewRequestIdExpiredChecker(requestId)
		exp := time.Now().Unix() + 10
		expStr := strconv.FormatInt(exp, 10)
		expired, err := requestIdExpireChecker.Check(expStr)
		assert.Nil(t, err)
		assert.False(t, expired)
	})
	t.Run("input=InvalidExp", func(t *testing.T) {
		requestId, err := NewRequestIdGenerator("1234567890").GenerateRequestId()
		assert.Nil(t, err)
		requestIdExpireChecker := NewRequestIdExpiredChecker(requestId)
		expired, err := requestIdExpireChecker.Check("zzzzzzzz")
		assert.NotNil(t, err)
		assert.False(t, expired)
	})
	t.Run("input=InvalidRequestId", func(t *testing.T) {
		requestIdExpireChecker := NewRequestIdExpiredChecker("23456")
		exp := time.Now().Unix() + 10
		expStr := strconv.FormatInt(exp, 10)
		expired, err := requestIdExpireChecker.Check(expStr)
		assert.NotNil(t, err)
		assert.False(t, expired)
	})
	t.Run("input=InvalidRequestIdHexMsTimeError", func(t *testing.T) {
		requestIdExpireChecker := NewRequestIdExpiredChecker("zzzzzzzz-2205077c-8663546")
		exp := time.Now().Unix() + 10
		expStr := strconv.FormatInt(exp, 10)
		expired, err := requestIdExpireChecker.Check(expStr)
		assert.NotNil(t, err)
		assert.False(t, expired)
	})
}
