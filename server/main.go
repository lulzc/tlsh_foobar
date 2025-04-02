package main

import (
	"flag"
	"log"
	"os"
	"tlsh_foobar/server/internal/tlsh_foo"
	"tlsh_foobar/server/internal/transport"
)

func main() {

	var csvFile string

	flag.StringVar(&csvFile, "csv", "", "path to CSV file")

	flag.Usage = func() {
		log.Println("Usage: ./tlshServer -csv csvFile ")
	}
	flag.Parse()

	if csvFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	svc := tlsh_foo.NewService()
	server := transport.NewServer(svc)

	//var csvFile = "./server/internal/db/full.csv"
	svc.ReadCsv(csvFile)

	if err := server.Serve(); err != nil {
		log.Println(err)
	}

}
