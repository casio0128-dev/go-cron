package main

import (
	"fileOpener/conf"
	"fmt"
	"github.com/BurntSushi/toml"
)

func main() {
	var config conf.Config

	if _, err := toml.DecodeFile("config.tml", &config); err != nil {
		panic(err)
	}

	if err := config.Parse(); err != nil {
		panic(err)
	}

	fmt.Printf("Host is :%s\n", config.User.Name)
	fmt.Printf("Port is :%s\n", config.User.Email)
	fmt.Printf("Date is :%v\n", config.User.Date)
	for k, v := range config.User.Files {
		fmt.Printf("File %d\n", k)
		fmt.Printf("  id is %d\n", v.Id)
		fmt.Printf("  path is %s\n", v.Path)
		fmt.Printf("  filename is %s\n", v.Name)
		fmt.Printf("  when is %s\n", v.When)
	}
}