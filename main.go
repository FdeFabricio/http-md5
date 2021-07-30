package main

import (
	"flag"
	"github.com/FdeFabricio/http-md5/myhttp"
)

func main() {
	parallel := flag.Int("parallel", 10, "limit of parallel requests")
	flag.Parse()
	urls := flag.Args()
	myhttp.Execute(*parallel, urls)
}
