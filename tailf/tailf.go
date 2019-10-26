package tailf

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"time"
)

type CollectConfig struct {
	LogPath string
	Topic   string
}

func NewCollectConfig() *CollectConfig {
	return &CollectConfig{}
}

type TailObj struct {
	T       *tail.Tail
	Collect *CollectConfig
}

type TextMsg struct {
	Msg   string
	Topic string
}

type TailObjMgr struct {
	Tails   []*TailObj
	MsgChan chan *TextMsg
}

var (
	tailObjMgr *TailObjMgr
)

func InitTail(collects []*CollectConfig) error {
	if len(collects) == 0 {
		return fmt.Errorf("没有配置: %v", collects)
	}

	tailObjMgr = &TailObjMgr{
		Tails:   make([]*TailObj, 0, 10),
		MsgChan: make(chan *TextMsg, 100),
	}

	for _, co := range collects {
		obj := &TailObj{
			Collect: co,
		}

		t, err := tail.TailFile(co.LogPath, tail.Config{
			ReOpen:    true,
			Follow:    true,
			Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
			MustExist: false,
			Poll:      true,
		})
		if err != nil {
			return err
		}

		obj.T = t

		tailObjMgr.Tails = append(tailObjMgr.Tails, obj)

		go readFromTail(obj)
	}

	return nil
}

func readFromTail(obj *TailObj) {
	for {
		msg, ok := <-obj.T.Lines
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
	}
}

func GetOneLine() *TextMsg {
	return <-tailObjMgr.MsgChan
}
