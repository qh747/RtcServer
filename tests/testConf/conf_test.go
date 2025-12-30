package testconf

import (
	"fmt"
	"testing"

	"github.com/unknwon/goconfig"
)

func TestLoadFile(t *testing.T) {
	conf, err := goconfig.LoadConfigFile("../../conf/ini/config.ini")
	if nil != err {
		fmt.Printf("Load config file error. err: %e, path: ../../conf/ini/config.ini", err)
		return
	}

	dir, err := conf.GetValue("LOG", "log_dir")
	fmt.Printf("Load config dir: %s\n", dir)

	name, err := conf.GetValue("LOG", "log_name")
	fmt.Printf("Load config name: %s\n", name)

	level, err := conf.GetValue("LOG", "log_level")
	fmt.Printf("Load config level: %s\n", level)

	size, err := conf.GetValue("LOG", "log_size")
	fmt.Printf("Load config size: %s\n", size)
}
