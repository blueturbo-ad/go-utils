package data_dict

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestReload(t *testing.T) {
	initW := sync.WaitGroup{}
	closeW := sync.WaitGroup{}
	initErrC := make(chan error, 1)
	closeC := make(chan struct{})

	d := NewDataDict[map[string]int](DataDictOption{CheckDur: 5 * time.Second}, &initW, &closeW, initErrC, closeC)
	if d == nil {
		t.Error("NewDataDict returned nil")
		return
	}

	d.initW.Add(1)
	d.RegisterFunc(
		func() ([]byte, error) {
			return []byte{1, 2, 3}, nil
		},
		func(data []byte) (*map[string]int, error) {
			// generate random data
			dict := make(map[string]int)
			dict["a"] = time.Now().Nanosecond()
			dict["b"] = time.Now().Nanosecond()
			fmt.Println("reload at", time.Now())
			return &dict, nil
		},
	)

	d.initW.Wait()

	dict := d.dict.Load()
	if dict == nil {
		t.Error("GetDict returned nil")
	}
	// fmt.Println(utils.PrettyStringify(dict))

	d.readFileCmdC <- time.Now()
	time.Sleep(1 * time.Second)
	dict = d.dict.Load()
	if dict == nil {
		t.Error("GetDict returned nil")
	}
	// fmt.Println(utils.PrettyStringify(dict))

	fmt.Println("send close")
	close(d.closeC)

	fmt.Println("close wait")
	d.closeW.Wait()
	fmt.Println("close done")
}
