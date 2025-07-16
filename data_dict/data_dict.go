package data_dict

import (
	"sync"
	"sync/atomic"
	"time"

	basetool "github.com/blueturbo-ad/go-utils/base_tool"
	"github.com/blueturbo-ad/go-utils/zap_loggerex"
)

const DefaultCheckDur = 10 * time.Second

type DataDictOption struct {
	CheckDur     time.Duration
	FilePath     string
	FileStatPath string // file timestamp
}

type DataDict[T any] struct {
	// dict
	dict           atomic.Pointer[T]
	lastReloadTime time.Time
	opts           DataDictOption

	// reload
	readFileFunc func() ([]byte, error)
	parseFunc    func([]byte) (*T, error)
	readFileCmdC chan time.Time

	// close
	closeC chan struct{}

	initW      *sync.WaitGroup
	initErrC   chan error
	closeW     *sync.WaitGroup
	loggerName string
}

func NewDataDict[T any](option DataDictOption, initW *sync.WaitGroup, closeW *sync.WaitGroup, initErrC chan error, closeC chan struct{}, logname string) *DataDict[T] {
	d := &DataDict[T]{
		opts:         option,
		readFileCmdC: make(chan time.Time),
		closeC:       closeC,
		initW:        initW,
		initErrC:     initErrC,
		closeW:       closeW,
		loggerName:   logname,
	}

	d.dict.Store(nil)

	closeW.Add(2)
	go d.checkBackground()
	go d.reloadBackground()

	return d
}

func (d *DataDict[T]) GetDict() *T {
	return d.dict.Load()
}

func (d *DataDict[T]) checkFunc() (time.Time, error) {
	data, err := basetool.ReadGCPCloudStorageFile(d.opts.FileStatPath, d.loggerName)
	if err != nil {
		return time.Time{}, err
	}

	unixTime, err := basetool.StringToInt64(basetool.RemoveWhitespace(string(data)))
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(unixTime, 0), nil
}

func (d *DataDict[T]) RegisterFunc(
	readFileFunc func() ([]byte, error), parseFunc func([]byte) (*T, error)) {
	if readFileFunc == nil {
		d.readFileFunc = func() ([]byte, error) {
			return basetool.ReadGCPCloudStorageFile(d.opts.FilePath, d.loggerName)
		}
	} else {
		d.readFileFunc = readFileFunc
	}

	d.parseFunc = parseFunc
	d.readFileCmdC <- time.Now().UTC()
}

func (d *DataDict[T]) checkBackground() {
	defer func() {
		d.closeW.Done()
	}()

	ticker := time.NewTicker(d.opts.CheckDur)

	for {
		zap_loggerex.GetLogger().Debug(d.loggerName, "wait for check trigger")

		select {
		case <-ticker.C:
			zap_loggerex.GetSingleton().Debug(d.loggerName, "check file stat")

			fileModTime, err := d.checkFunc()
			ticker.Reset(d.opts.CheckDur)
			if err != nil {
				zap_loggerex.GetSingleton().Warn(d.loggerName, "failed to get file modify time %s", err)
				continue
			}

			zap_loggerex.GetSingleton().Debug(d.loggerName, "file modify time is %v", fileModTime)

			if fileModTime.After(d.lastReloadTime) {
				zap_loggerex.GetSingleton().Debug(d.loggerName, "trigger reload")
				d.readFileCmdC <- fileModTime
			}

		case <-d.closeC:
			zap_loggerex.GetSingleton().Info(d.loggerName, "close data dict check background")
			return
		}
	}
}

func (d *DataDict[T]) reloadBackground() {
	defer func() {
		d.closeW.Done()
	}()

	for {
		zap_loggerex.GetSingleton().Debug(d.loggerName, "wait for reload trigger")

		select {
		case fileModTime := <-d.readFileCmdC:
			begTime := time.Now().UTC()
			zap_loggerex.GetSingleton().Debug(d.loggerName, "reload dict")

			data, err := d.readFileFunc()
			if err != nil {
				if d.isInitLoad() {
					d.initErrC <- err
					d.initW.Done()
				}

				zap_loggerex.GetSingleton().Warn(d.loggerName, "failed to reload dict %s", err)
				continue
			}

			newDict, err := d.parseFunc(data)
			if err != nil {
				if d.isInitLoad() {
					d.initErrC <- err
					d.initW.Done()
				}

				zap_loggerex.GetSingleton().Warn(d.loggerName, "failed to parse dict %s", err)
				continue
			}

			d.dict.Store(newDict)

			if d.isInitLoad() {
				d.initW.Done()
			}

			d.lastReloadTime = fileModTime

			zap_loggerex.GetSingleton().Info(d.loggerName, "reload dict cost %v", time.Since(begTime))

		case <-d.closeC:
			zap_loggerex.GetSingleton().Info(d.loggerName, "close data dict reload background")
			return
		}
	}
}

func (d *DataDict[T]) isInitLoad() bool {
	return d.lastReloadTime.IsZero()
}
