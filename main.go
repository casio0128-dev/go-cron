package main

import (
	"fileOpener/conf"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"sync"
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

	wg := &sync.WaitGroup{}

	for k, v := range config.User.Files {
		fmt.Printf("File %d\n", k)
		fmt.Printf("  id is %d\n", v.Id)
		fmt.Printf("  path is %s\n", v.Path)
		fmt.Printf("  cron is %s\n", v.Cron)
		fmt.Printf("  when is %s\n", v.When)

		wg.Add(1)
		go func(wg *sync.WaitGroup, file conf.FileInfo) {
			if err := file.Run(); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(wg, v)
	}

	wg.Wait()
}