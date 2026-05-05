package main

import (
	"fmt"
	"md2html/internal/config"
	"md2html/internal/converter"
	"md2html/internal/utils"
	"os"
	"time"
)

func main() {
	start := time.Now()
	cfg , err := utils.LoadConfig()
	if err != nil {
		fmt.Println("failed to load config file.")
		cfg = &config.Config{}
	}

	utils.GetFlags(cfg)

	converter := converter.NewConverter(cfg)
	if err := converter.ConvertDirectory(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("cost:", time.Since(start))
}
