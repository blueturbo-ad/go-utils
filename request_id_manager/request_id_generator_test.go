package requestidmanager

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestIdGenerator(t *testing.T) {
	t.Run("input=valid", func(t *testing.T) {
		requestIdGenerator := NewRequestIdGenerator("1234567890")
		requestId, err := requestIdGenerator.GenerateRequestId()
		assert.Nil(t, err)
		requestSubIdFirst, err := requestIdGenerator.GenerateRequestSubId(requestId)
		assert.Nil(t, err)
		requestIdFirst, countFirst, err := NewRequestSubIdParser(requestSubIdFirst).Parse()
		assert.Nil(t, err)
		assert.Equal(t, requestId, requestIdFirst)
		requestSubIdSecond, err := requestIdGenerator.GenerateRequestSubId(requestId)
		assert.Nil(t, err)
		requestIdSecond, countSecond, err := NewRequestSubIdParser(requestSubIdSecond).Parse()
		assert.Nil(t, err)
		assert.Equal(t, requestId, requestIdSecond)
		first, err := strconv.ParseInt(countFirst, 10, 64)
		assert.Nil(t, err)
		second, err := strconv.ParseInt(countSecond, 10, 64)
		assert.Nil(t, err)
		assert.Equal(t, first+1, second)
	})
}
