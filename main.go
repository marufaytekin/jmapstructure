package main

import (
	"fmt"
	"github.com/marufaytekin/jmapstructure/heap"
	"log"
)

func main() {

	//Sample usage
	javaHome := "/Library/Java/JavaVirtualMachines/openjdk-12.0.2.jdk/Contents/Home/"
	pid := "38602"

	h, err := heap.Get(javaHome, pid)

	log.Println(fmt.Sprintf("%v", h))

	if err != nil {
		log.Println(err)
	}
}
