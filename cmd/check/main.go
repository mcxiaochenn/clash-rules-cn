package main

import (
	"log"

	"clash-rules-cn/internal/artifacts"
	"clash-rules-cn/internal/config"
)

func main() {
	cfg, err := config.Load("sources.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := artifacts.Validate(cfg, "rules", "geoip"); err != nil {
		log.Fatal(err)
	}
	log.Println("发布产物校验通过")
}
