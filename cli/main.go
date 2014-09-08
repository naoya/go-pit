package main

import (
	"fmt"
	"github.com/naoya/go-pit"
)

func main() {
	config := pit.Get("kaizenplatform.in")
	fmt.Println(config["aws_access_key_id"])

	pit.Switch("development")
	config = pit.Get("kaizenplatform.in")
	fmt.Println(config["aws_access_key_id"])

	pit.Switch("default")
	config = pit.Get("kaizenplatform.in")
	fmt.Println(config["aws_access_key_id"])
}
