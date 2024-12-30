package requestidmanager

import (
	"fmt"
	"testing"
)

func TestNewSnowflake(t *testing.T) {
	t.Run("input=valid", func(t *testing.T) {
		workerID := int64(1)
		snowflake := NewSnowflake(workerID)
		if snowflake.workerID != workerID {
			t.Errorf("workerID not match")
		}

		// 生成一系列唯一 ID
		for i := 0; i < 10; i++ {
			id := snowflake.NextID()
			fmt.Println(id)
		}
	})

}
