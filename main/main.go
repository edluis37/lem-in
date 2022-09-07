package main

import (
	"flag"
	"fmt"
	"leminmod"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Program takes argument (fileName or --flags)!")
		os.Exit(1)
	}
	argument := os.Args[1]
	// Set Flags
	filename := flag.String("file", "", "--file=filename\n")
	flag.Parse()
	// Start Program
	if strings.HasPrefix(argument, "--") { // if has flags
		if *filename != "" {
			leminmod.RunProgramWithFile(*filename, true)
		} else { // Default = Help
			flag.Usage()
			os.Exit(1)
		}
	} else { // Default
		leminmod.RunProgramWithFile(argument, false)
	}
}
