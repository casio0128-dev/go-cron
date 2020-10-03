package conf

import (
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
	Name string		`toml:"name"`
	When string			`toml:"when"`
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

		for _, match := range reg.FindAllString(file.Name, -1) {
			if reg.MatchString(file.Name) {
				file.Name = conf.User.parseDate(file.Name, match)
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

func (user User) parseYear (src string) string {
	if !strings.Contains(src, "{{YYYY}}") {
		return ""
	}

	return strings.ReplaceAll(src, "{{YYYY}}", user.Date.Format("2006"))
}

func (user User) parseMonth (src, format string) string {
	var result string

	if !strings.Contains(src, format) {
		return ""
	}

	switch format {
	case "{{MM}}":
		result = strings.ReplaceAll(src, "{{MM}}", user.Date.Format("01"))
	case "{{M}}":
		result = strings.ReplaceAll(src, "{{M}}", user.Date.Format("1"))
	}
	return result
}

func (user User) parseDay (src, format string) string {
	var result string

	if !strings.Contains(src, format) {
		return ""
	}

	switch format {
	case "{{DD}}":
		result = strings.ReplaceAll(src, "{{DD}}", user.Date.Format("02"))
	case "{{D}}":
		result = strings.ReplaceAll(src, "{{D}}", user.Date.Format("2"))
	}

	return result
}