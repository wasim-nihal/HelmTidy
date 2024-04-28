package models

import "sync"

type TplUnused struct {
	TplUnusedMap map[string]string
	lock         sync.Mutex
}

var UnusedTpls *TplUnused

func InitUnusedTpls() {
	if UnusedTpls == nil {
		UnusedTpls = &TplUnused{
			TplUnusedMap: make(map[string]string),
			lock:         sync.Mutex{},
		}
	}
}

func (u *TplUnused) Calculate() {
	for key, val := range GetTplDefinations().DefinationMap {
		if _, ok := TplUsgs.TplUsageMap[key]; !ok {
			UnusedTpls.TplUnusedMap[key] = val
		}
	}
}
