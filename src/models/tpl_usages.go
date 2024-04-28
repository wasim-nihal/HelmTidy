package models

import (
	"dangling-tpls/src/utils"
	"flag"
	"log"
	"os"
	"sync"
)

var numberofWorkers = flag.Int("workerThreads", 4, "number of worker threads to be used")

type TplUsages struct {
	TplUsageMap map[string]string
	lock        sync.Mutex
}

var TplUsgs *TplUsages

func InitTplUsages() {
	if TplUsgs == nil {
		TplUsgs = &TplUsages{
			TplUsageMap: make(map[string]string),
			lock:        sync.Mutex{},
		}
	}
}

func (t *TplUsages) Populate() {
	fileList := GetFileList()
	fileChan := make(chan string, len(fileList))
	var wg sync.WaitGroup
	for i := 0; i < *numberofWorkers; i++ {
		wg.Add(1)
		go t.populateTplUsages(fileChan, &wg)
	}
	for _, file := range fileList {
		fileChan <- file
	}
	close(fileChan)
	wg.Wait()
}

func (t *TplUsages) populateTplUsages(files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Println("Error opening file:", err)
		}
		yaml_commentless_content := utils.RgxYamlComments.ReplaceAll(content, []byte(""))
		commentless_content := utils.RgxTplComments.ReplaceAll(yaml_commentless_content, []byte(""))
		allMatches := utils.RgxTplUsages.FindAll(commentless_content, -1)
		if allMatches == nil {
		} else {
			for _, v := range allMatches {
				sm := utils.RgxTplUsages.FindSubmatch(v)
				if len(sm) > 2 {
					t.lock.Lock()
					t.TplUsageMap[string(sm[2])] = string(file)
					t.lock.Unlock()
				}
			}
		}
	}
}
