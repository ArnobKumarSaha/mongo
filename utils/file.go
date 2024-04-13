package utils

import (
	"fmt"
	"log"
	"os"
)

func MakeDir(dir string) {
	_, err := os.Stat(dir)
	if err == nil {
		fmt.Printf("Directory %s already exists; cleaning up.\n", dir)
		err = os.RemoveAll(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteFile(dir, fileName string, data []byte) {
	fileName = dir + "/" + fileName + ".json"
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
