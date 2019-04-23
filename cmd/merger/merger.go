package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	flag.Parse()
	for _, dir := range flag.Args() {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			fmt.Println(file.Name())
		}
	}
}
