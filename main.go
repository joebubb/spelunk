package main

import (
	"fmt"

	"github.com/joebubb/spelunk/cli"
	"github.com/joebubb/spelunk/util"
)

func main() {
	fmt.Println("running spelunk...")
	fmt.Printf("%v\n", util.UrlIsValidGet("https://google.com/"))
	wp := util.NewWorkerPool(1000, 20)
	wp.Start()

	util.ForEachUrlGen("https://google.com/", cli.Alphabet(), 4, func(s string) {
		wp.SubmitTask(func() {
			if util.UrlIsValidGet(s + "/") {
				fmt.Printf("%s\n", s+"/")
			}
		})
	})

	wp.Stop()
}
