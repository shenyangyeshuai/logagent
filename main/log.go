package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

var (
	logLevels = map[string]int{
		"debug": logs.LevelDebug,
		"warn":  logs.LevelWarn,
		"info":  logs.LevelInfo,
		"trace": logs.LevelTrace,
	}
)

func initLogger() error {
	config := make(map[string]interface{})
	config["filename"] = appConfig.LogPath
	level, _ := logLevels[appConfig.LogLevel]
	config["level"] = level

	s, err := json.Marshal(config) // s's type is []byte
	if err != nil {
		fmt.Printf("initLogger() failed: %v", err)
		return err
	}

	logs.SetLogger(logs.AdapterFile, string(s))

	return nil
}
