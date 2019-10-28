package tailf

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

type CollectConfig struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

func NewCollectConfig() *CollectConfig {
	return &CollectConfig{}
}

type TailObj struct {
	T        *tail.Tail
	Collect  *CollectConfig
	Status   int
	exitChan chan int
}

type TextMsg struct {
	Msg   string
	Topic string
}

type TailObjMgr struct {
	Tails   []*TailObj
	MsgChan chan *TextMsg
	lock    sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)

func InitTail(collects []*CollectConfig) error {
	tailObjMgr = &TailObjMgr{
		Tails:   make([]*TailObj, 0, 10),
		MsgChan: make(chan *TextMsg, 100),
	}

	if len(collects) == 0 {
		return fmt.Errorf("没有配置: %v", collects)
	}

	for _, co := range collects {
		createNewTask(co)
		// obj := &TailObj{
		//	Collect: co,
		// }
		//
		// t, err := tail.TailFile(co.LogPath, tail.Config{
		//	ReOpen:    true,
		//	Follow:    true,
		//	Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		//	MustExist: false,
		//	Poll:      true,
		// })
		// if err != nil {
		//	return err
		// }
		//
		// obj.T = t
		//
		// tailObjMgr.Tails = append(tailObjMgr.Tails, obj)
		//
		// go readFromTail(obj)
	}

	return nil
}

func readFromTail(obj *TailObj) {
	for {
		select {
		case msg, ok := <-obj.T.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename: %s\n", obj.T.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			tm := &TextMsg{
				Msg:   msg.Text,
				Topic: obj.Collect.Topic,
			}
			tailObjMgr.MsgChan <- tm
		case <-obj.exitChan:
			logs.Warn("关了")
			return
		}
	}
}

func GetOneLine() *TextMsg {
	return <-tailObjMgr.MsgChan
}

func UpdateConfig(cocs []*CollectConfig) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()

	for _, coc := range cocs {
		var isRunning = false
		for _, obj := range tailObjMgr.Tails {
			if coc.LogPath == obj.Collect.LogPath {
				isRunning = true
				break
			}
		}

		if isRunning {
			continue
		}

		createNewTask(coc)
	}

	var tailObjs = []*TailObj{}
	for _, obj := range tailObjMgr.Tails {
		obj.Status = StatusDelete
		for _, coc := range cocs {
			if coc.LogPath == obj.Collect.LogPath {
				obj.Status = StatusNormal
				break
			}
		}

		if obj.Status == StatusDelete {
			obj.exitChan <- 1
		}
		tailObjs = append(tailObjs, obj)
	}

	tailObjMgr.Tails = tailObjs
}

func createNewTask(coc *CollectConfig) {
	obj := &TailObj{
		Collect:  coc,
		exitChan: make(chan int, 1),
	}

	t, err := tail.TailFile(coc.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logs.Error("[error] filename: %v", coc.LogPath, err)
		return
	}

	obj.T = t

	tailObjMgr.Tails = append(tailObjMgr.Tails, obj)

	go readFromTail(obj)
}
