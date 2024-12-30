package requestidmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestIdParser(t *testing.T) {
	t.Run("input=valid", func(t *testing.T) {
		requestId := "5fcd0f24633ca-2205077c-8663546"
		requestIdParser := NewRequestIdParser(requestId)
		hexMsTime, sha1, counterPerMicroSecondStep, err := requestIdParser.Parse()
		assert.Nil(t, err)
		assert.Equal(t, hexMsTime, "5fcd0f24633ca")
		assert.Equal(t, sha1, "2205077c")
		assert.Equal(t, counterPerMicroSecondStep, "8663546")
	})
	t.Run("input=invalidLessArg", func(t *testing.T) {
		requestId := "5fcd0f24633ca2205077c8663546-0"
		requestIdParser := NewRequestIdParser(requestId)
		_, _, _, err := requestIdParser.Parse()
		assert.NotNil(t, err)
	})

	t.Run("input=invalidMoreArg", func(t *testing.T) {
		requestId := "5fcd0f246-33ca220507-7c866-3546-0"
		requestIdParser := NewRequestIdParser(requestId)
		_, _, _, err := requestIdParser.Parse()
		assert.NotNil(t, err)
	})

}
