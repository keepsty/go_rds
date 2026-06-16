package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/keepsty/go_rds/internal/cluster/models"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func FileTemplate(e interface{},files []*models.SaltStateFiles) (err error) {
	for _, f := range files {
		if f == nil {
			continue
		}
		bo, _ := PathExists(f.TargetFilePath)
		if !bo {
			err = os.MkdirAll(f.TargetFilePath, 0777)
			if err != nil {
				panic(err)
				return err
			}
		}
		tpl, err := template.New(f.FileName).Delims("[[", "]]").ParseFiles(fmt.Sprintf("%s/%s", f.FilePath, f.FileName))
		if err != nil {
			panic(err)
		}
		init_file := fmt.Sprintf("%s/%s", f.TargetFilePath, f.TargetFileName)
		_, err = os.Create(init_file)
		os_init_file, _ := os.OpenFile(init_file, os.O_RDWR|os.O_CREATE, os.ModePerm)
		defer os_init_file.Close()
		service_writer := bufio.NewWriter(os_init_file)
		err = tpl.Execute(service_writer, e)
		if err != nil {
			panic(err)
			return err
		}
		service_writer.Flush()
	}
	return
}
