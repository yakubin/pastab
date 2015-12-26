package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
)

func main() {
	hashfun, keywords := loadCmdLine()
	hash := flag.Arg(0)

	c := genPermCombs(keywords)
	total := totalPerms(len(keywords))

	uiprogress.Start()
	bar := uiprogress.AddBar(total)
	bar.AppendCompleted()
	bar.PrependElapsed()

	for v := range c {
		bar.Incr()

		if hashfun(v) == hash {
			fmt.Println(v)
			os.Exit(0)
		}
	}

	fmt.Println("no match")
}
