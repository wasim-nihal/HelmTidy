package models

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

type ListOfFiles struct {
	List []string
	lock sync.Mutex
}

var FileList *ListOfFiles

func InitFileList() {
	if FileList == nil {
		FileList = &ListOfFiles{
			List: make([]string, 0),
			lock: sync.Mutex{},
		}
	}
}

func GetFileList() []string {
	return FileList.List
}

func (t *ListOfFiles) Populate(dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	var visit func(path string, info os.FileInfo, err error) error
	var visitedDirs = make(map[string]bool)
	visit = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if visitedDirs[path] {
				return nil
			}
			visitedDirs[path] = true
			err := filepath.Walk(path, visit)
			if err != nil {
				log.Printf("Error walking into directory %q: %v\n", path, err)
				return err
			}
			return nil
		}
		t.lock.Lock()
		t.List = append(t.List, path)
		t.lock.Unlock()
		return nil
	}
	err := filepath.Walk(dir, visit)
	if err != nil {
		log.Println(err)
	}
}
