package models

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type TplDefinations struct {
	DefinationMap map[string]string
	lock          sync.Mutex
}

var TplDefs *TplDefinations

func InitTplDefinations() {
	if TplDefs == nil {
		TplDefs = &TplDefinations{
			DefinationMap: make(map[string]string),
			lock:          sync.Mutex{},
		}
	}
}

func GetTplDefinations() *TplDefinations {
	return TplDefs
}

func (t *TplDefinations) Populate(dir string) {
	// Populate TplDefinations
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".tpl") {
			t.populateTplDefinations(path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

func (t *TplDefinations) populateTplDefinations(file string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		log.Println("Error opening file:", err)
		return false
	}
	commentless_content := regex_comments.ReplaceAll(content, []byte(""))
	allMatches := regex_defination.FindAll(commentless_content, -1)
	if allMatches == nil {
		return false
	} else {
		for _, v := range allMatches {
			x := regex_defination.FindSubmatch(v)
			if len(x) > 1 {
				t.DefinationMap[string(x[1])] = string(file)
			}
		}
		return true
	}

}
