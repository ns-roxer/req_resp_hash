package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	parallelWorkersNum := flag.Uint("parallel", 10, "Limit of parallel requests. Must be above zero")
	flag.Parse()
	if parallelWorkersNum == nil || *parallelWorkersNum == 0 {
		flag.Usage()
		os.Exit(1)
	}
	urls := flag.Args()

	app := NewApp(*parallelWorkersNum, NewSimpleHttpClient(http.DefaultClient))
	urlsWithHashes, err := app.CalcUrlHashes(urls)
	if err != nil {
		fmt.Println("error occurred: " + err.Error())
		os.Exit(1)
	}
	for _, u := range urlsWithHashes {
		fmt.Println(u)
	}
}
