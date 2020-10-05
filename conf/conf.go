package conf

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	User User
}

type User struct {
	Name  string        `toml:"name"`
	Email  string       `toml:"email"`
	Date time.Time		`toml:"date"`
	Files []FileInfo	`toml:"fileInfo"`
}

type FileInfo struct {
	Id int				`toml:"id"`
	Path string			`toml:"path"`
	When string			`toml:"when"`
	Cron string			`toml:"cron"`
}

func (conf *Config) Parse () error {
	reg, err := regexp.Compile(`({.{1,4}})`)
	if err != nil {
		return err
	}

	if conf.User.Date.Format("2006") == "0001" {
		conf.User.Date = time.Now()
	}

	for index, file := range conf.User.Files {
		for _, match := range reg.FindAllString(file.Path, -1) {
			if reg.MatchString(file.Path) {
				file.Path = conf.User.parseDate(file.Path, match)
			}
		}
		conf.User.Files[index] = file
	}


	return nil
}

func (user User) parseDate (src, format string) (result string) {
	switch format {
	case "{YYYY}":
		result = strings.ReplaceAll(src, "{YYYY}", user.Date.Format("2006"))
	case "{MM}":
		result = strings.ReplaceAll(src, "{MM}", user.Date.Format("01"))
	case "{M}":
		result = strings.ReplaceAll(src, "{M}", user.Date.Format("1"))
	case "{DD}":
		result = strings.ReplaceAll(src, "{DD}", user.Date.Format("02"))
	case "{D}":
		result = strings.ReplaceAll(src, "{D}", user.Date.Format("2"))
	}
	return
}

type CronInfo struct {
	Minute string
	Hour string
	Day string
	Month string
	Week string
}

func (file FileInfo) GetCron() *CronInfo {
	cronInfo := strings.Split(file.Cron, " ")
	if len(cronInfo) <= 4 {
		return nil
	}

/*	min := cronInfo[0]
	hour  := cronInfo[1]
	day  := cronInfo[2]
	mon  := cronInfo[3]
	week  := cronInfo[4]*/

/*	if strings.EqualFold(min, "*") {
		min = "1"
	} else if len(strings.Split(min, "/")) == 2 {
		num, err := strings.Atoi(strings.Split(min, "/")[0])
		denom, err := strings.Atoi(strings.Split(min, "/")[1])
		if err != nil {
			panic(err)
		}
		min = fmt.Sprintf("%s", int(num / denom))
	}*/


	return &CronInfo{
		Minute: cronInfo[0],
		Hour:   cronInfo[1],
		Day:   cronInfo[2],
		Month:  cronInfo[3],
		Week:   cronInfo[4],
	}
}

func (file FileInfo) Run () error {
	if strings.EqualFold(file.Cron, "") {
		if err := open(file.Path); err != nil {
			return err
		}

		fmt.Println("open to ",file.Path)
	}
	return nil
}

func open(cmd string) error{
	var arg1Cmd string

	switch runtime.GOOS {
	case "darwin":
		arg1Cmd = "open"
	case "windows":
		arg1Cmd = "start"
	}

	if err := exec.Command(arg1Cmd, cmd).Run(); err != nil {
		return err
	}
	return nil
}