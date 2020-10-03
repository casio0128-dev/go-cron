package conf

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	User User
}

type User struct {
	Name  string        `toml:"name"`
	Email  string        `toml:"email"`
	Date time.Time		`toml:"date"`
	Files []FileInfo	 `toml:"fileInfo"`
}

type FileInfo struct {
	Id int				`toml:"id"`
	Path string			`toml:"path"`
	Filename string		`toml:"filename"`
	When string			`toml:"when"`
}

func (conf *Config) Parse () error {
	reg, err := regexp.Compile(`({{.{1,4}}})`)
	if err != nil {
		return err
	}
	for _, v := range conf.User.Files {
		for _, match := range reg.FindAllString(v.Path, -1) {
			//if reg.Match([]byte(v.Path)) {
				v.Path = conf.User.parseYear(v.Path)
			//	v.Path = conf.User.parseMonth(v.Path, match)
				fmt.Println(conf.User.parseDay(v.Path, match))
			//}ß
		}
	}
	return nil
}

func (user User) parseYear (src string) string {
	return strings.ReplaceAll(src, "{{YYYY}}", user.Date.Format("2006"))
}

func (user User) parseMonth (src, format string) string {
	var result string

	switch format {
	case "{{MM}}":
		result = strings.ReplaceAll(src, "{{MM}}", user.Date.Format("01"))
	case "{{M}}":
		result = strings.ReplaceAll(src, "{{M}}", user.Date.Format("1"))
	default:
		fmt.Println("can't find format")
	}
	return result
}

func (user User) parseDay (src, format string) string {
	var result string

	switch format {
	case "{{DD}}":
		result = strings.ReplaceAll(src, "{{DD}}", user.Date.Format("02"))
	case "{{D}}":
		result = strings.ReplaceAll(src, "{{D}}", user.Date.Format("2"))
	default:
		fmt.Println("can't find format")
	}

	return result
}