package main

import (
	"dangling-tpls/src/models"
	"dangling-tpls/src/utils"
	"flag"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	chartPath       = flag.String("p", "", `path of the chart that has to be scanned for the unused tpls. It can be either relative/absolute/http url. If the value starts with 'http://' or 'https://', then a http request is made to get the chart.`)
	dependentCharts = flag.String("dependentCharts", "", `a list comma saperated paths of dependent charts that make use of tpls from the current '-p' chart. Values can be either relative/absolute/http url. If the value starts with 'http://' or 'https://', then a http request is made to get the chart.`)
	exitWithNonZero = flag.Bool("exitWithNonZero", false, "flag indicating whether to exit with a non zero error code in the case if dangling tpls are found.")
)

func main() {
	flag.Parse()
	if *chartPath == "" {
		panic("chart path '-p' cannot be empty")
	}
	if strings.HasPrefix(*chartPath, "http://") || strings.HasPrefix(*chartPath, "https://") {
		*chartPath = utils.GetChartHttp(*chartPath)
	}
	// initialize model objs
	models.InitFileList()
	models.InitTplDefinations()
	models.InitTplUsages()
	models.InitUnusedTpls()
	var chartsToBeScannedForTplUsages []string
	var dependentChartList []string
	if *dependentCharts != "" {
		dependentChartList = strings.Split(*dependentCharts, ",")
	}
	dependentChartList = append(dependentChartList, *chartPath)
	log.Println("Scanning for tpl usages for the charts below:")
	for _, path := range dependentChartList {
		localPath := path
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			localPath = utils.GetChartHttp(path)
		}
		log.Println(path)
		chartsToBeScannedForTplUsages = append(chartsToBeScannedForTplUsages, localPath)
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(chartsToBeScannedForTplUsages))
	for _, path := range chartsToBeScannedForTplUsages {
		go models.FileList.Populate(path, wg)
	}
	wg.Wait()

	models.TplDefs.Populate(*chartPath)
	models.TplUsgs.Populate()
	models.UnusedTpls.Calculate()
	// print
	var result [][]string
	result = append(result, []string{"DANGLING TPL NAME", "FILE"})

	for key, value := range models.UnusedTpls.TplUnusedMap {
		result = append(result, []string{key, value})
	}

	utils.PrintTable(result)

	if len(models.UnusedTpls.TplUnusedMap) > 0 {
		if *exitWithNonZero {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

}
