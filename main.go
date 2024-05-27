package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/helloworlddan/tortuneai/tortuneai"
)

func main() {
	versionToggle := flag.Bool("version", false, "show version info")
	flag.Parse()

	if *versionToggle {
		fmt.Printf("tortuneai version %s\n", tortuneai.Version)
		return
	}

	joke, err := tortuneai.HitMe("", os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Panicf("error: %v\n", err)
	}

	fmt.Println(joke)
}
