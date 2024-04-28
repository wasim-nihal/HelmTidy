package models

import (
	"dangling-tpls/src/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TplDefinations struct {
	DefinationMap map[string]string
}

var TplDefs *TplDefinations

func InitTplDefinations() {
	if TplDefs == nil {
		TplDefs = &TplDefinations{
			DefinationMap: make(map[string]string),
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
		log.Printf("error pupulating the tpl definations list. reason %s\n", err.Error())
	}
}

func (t *TplDefinations) populateTplDefinations(file string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		log.Println("Error opening file:", err)
		return false
	}
	commentless_content := utils.RgxTplComments.ReplaceAll(content, []byte(""))
	allMatches := utils.RgxTplDefinations.FindAll(commentless_content, -1)
	if allMatches == nil {
		return false
	} else {
		for _, v := range allMatches {
			sm := utils.RgxTplDefinations.FindSubmatch(v)
			if len(sm) > 1 {
				t.DefinationMap[string(sm[1])] = string(file)
			}
		}
		return true
	}

}
