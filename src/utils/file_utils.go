package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ChartDir        string
	once            sync.Once
	LibChartAbsPath string
)

// initialize the working directory
func initWd() {
	if ChartDir == "" {
		currDir, err := os.Getwd()
		if err != nil {
			log.Fatalln("cannot initialize working directory. reason ", err.Error())
		}
		ChartDir = filepath.Join(currDir, "charts")
	}
}

// GetChartHttp downloads and extracts the chart from the http url
func GetChartHttp(chartUrl string) string {
	initWd()
	log.Printf("starting downloading the of the file %s\n", chartUrl)
	fileURL, err := url.Parse(chartUrl)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("cannot create file %s. reason %s", fileName, err.Error())
	}
	defer file.Close()
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(chartUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, _ := io.Copy(file, resp.Body)

	log.Printf("Downloaded a file %s with size %d", fileName, size)
	err = untarTgz(fileName)
	if err != nil {
		log.Fatalf("cannot untar %s file. reason %s", fileName, err.Error())
	}
	// remove the library charts from the scan to prevent reporting of false positives
	err = os.RemoveAll(filepath.Join(LibChartAbsPath, "charts"))
	if err != nil {
		log.Fatalf("error removing the library chart directory. reason %s\n", err.Error())
	}
	return ChartDir
}

// untar the tarball
func untarTgz(tgzPath string) error {
	tgzFile, err := os.Open(tgzPath)
	if err != nil {
		return fmt.Errorf("failed to open tgz file: %w", err)
	}
	defer tgzFile.Close()
	gzipReader, err := gzip.NewReader(tgzFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}
		targetPath := filepath.Join(ChartDir, header.Name)

		if err := extractFile(tarReader, header, targetPath); err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	return nil
}

// extract the file
func extractFile(reader *tar.Reader, header *tar.Header, targetPath string) error {
	// Create the directory structure if it doesn't exist
	if header.Typeflag != tar.TypeDir {
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		once.Do(func() {
			LibChartAbsPath = targetDir
		})
	}

	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, reader)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}
