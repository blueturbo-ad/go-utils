package requestidmanager

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blueturbo-ad/go-utils/global"

	"github.com/stretchr/testify/assert"
)

func TestRequestSubId(t *testing.T) {
	requestSubId := fmt.Sprintf(request_sub_id_format, request_id_format, 0)
	assert.Equal(t, len(strings.Split(requestSubId, "-")), request_sub_id_formatter_split_count)

	t.Run("input=valid", func(t *testing.T) {
		requestSubId := "5ff416a12ddd7-ca803e9a-001-0"
		requestId := "5ff416a12ddd7-ca803e9a-001"
		requestSubIdParser := NewRequestSubIdParser(requestSubId)
		requestId, requestSubId, err := requestSubIdParser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, requestId, requestId)
		assert.Equal(t, requestSubId, "0")
	})
	t.Run("ParseError", func(t *testing.T) {
		requestSubId := "5fcd0f24633ca-2205077c-asdasd-1"
		requestSubIdParser := NewRequestSubIdParser(requestSubId)
		requestId, requestSubId, err := requestSubIdParser.Parse()
		assert.Equal(t, err != nil, true)
		assert.Equal(t, requestId, global.EmptyString)
		assert.Equal(t, requestSubId, global.EmptyString)
	})
	t.Run("InvalidParse", func(t *testing.T) {
		requestSubId := "5fcd0f24633ca-2205077c-8663546"
		requestSubIdParser := NewRequestSubIdParser(requestSubId)
		requestId, requestSubId, err := requestSubIdParser.Parse()
		assert.Equal(t, err != nil, true)
		assert.Equal(t, requestId, global.EmptyString)
		assert.Equal(t, requestSubId, global.EmptyString)
	})
}
