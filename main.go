package main

import (
	"dangling-tpls/models"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	chartPath       = flag.String("p", "", "path of the chart that has to be scanned for unused tpls")
	dependentCharts = flag.String("dependentCharts", "", "a comma saperated string of absolute paths of dependent charts that make use of tpls from the current '-p' chart")
	exitWithNonZero = flag.Bool("exitWithNonZero", false, "flag indicating whether to exit with a non zero error code in the case if dangling tpls are found.")
)

func main() {
	// TODO: panic if chartPath is not set
	flag.Parse()
	models.InitFileList()
	models.InitTplDefinations()
	models.InitTplUsages()
	models.InitUnusedTpls()
	var chartsToBeScannedForTplUsages []string
	if *dependentCharts != "" {
		chartsToBeScannedForTplUsages = strings.Split(*dependentCharts, ",")
	}
	chartsToBeScannedForTplUsages = append(chartsToBeScannedForTplUsages, *chartPath)
	log.Println("Scanning for tpl usages for the charts below:")
	for _, path := range chartsToBeScannedForTplUsages {
		log.Println(path)
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(chartsToBeScannedForTplUsages))
	for _, path := range chartsToBeScannedForTplUsages {
		go models.FileList.Populate(path, wg)
	}
	wg.Wait()

	models.TplDefs.Populate(*chartPath)
	models.TplUsgs.Populate()
	models.UnusedTpls.Populate()
	// print
	var result [][]string
	result = append(result, []string{"DANGLING TPL NAME", "FILE"})
	for key, value := range models.UnusedTpls.TplUnusedMap {
		result = append(result, []string{key, value})
	}
	printTable(result)
	if len(models.UnusedTpls.TplUnusedMap) > 0 {
		if *exitWithNonZero {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}

func printTable(data [][]string) {
	// Calculate column widths
	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, cell := range row {
			width := len(cell)
			if width > colWidths[i] {
				colWidths[i] = width
			}
		}
	}

	// Print top border
	printBorder(colWidths)

	// Print table content
	for _, row := range data {
		for i, cell := range row {
			fmt.Printf("| %-*s ", colWidths[i], cell)
		}
		fmt.Println("|")
		printBorder(colWidths)
	}
}

func printBorder(colWidths []int) {
	for _, width := range colWidths {
		fmt.Print("+" + strings.Repeat("-", width+2))
	}
	fmt.Println("+")
}
