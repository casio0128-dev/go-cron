package conf

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
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

const (
	MINUTE_MIN = 0
	MINUTE_MAX = 59
	HOUR_MIN = 0
	HOUR_MAX = 23
	DAY_MIN = 1
	DAY_MAX = 31
	MONTH_MIN = 1
	MONTH_MAX = 12
	WEEK_MIN = 0	// sun mon tue wed thu fri sat sun
	WEEK_MAX = 7 	//  0   1   2   3   4   5   6   7

	ASTAH = "*"
	SEPARATER_SLASH = "/"
	SEPARATER_CONMA = ","
)

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
	Minute cronDuration
	Hour cronDuration
	Day cronDuration
	Month cronDuration
	Week cronDuration
}

type cronDuration struct {
	min int				// 範囲指定（最小）
	max int				// 範囲指定（最大）
	duration int		// min - max間での実行期間
}

func NewDuration(min, max, duration int) *cronDuration{
	return &cronDuration{
		min: min,
		max: max,
		duration: duration,
	}
}

func (file FileInfo) GetCron() (error, *CronInfo) {

	var min cronDuration
	var hour cronDuration
	var day cronDuration
	var mon cronDuration
	var week cronDuration

	cronInfo := strings.Split(file.Cron, " ")
	if !strings.EqualFold(file.Cron, "") && len(cronInfo) <= 4 {
		return fmt.Errorf("invalid format."), nil
	}

	minCron := strings.Split(cronInfo[0], SEPARATER_SLASH)
	if strings.EqualFold(cronInfo[0], SEPARATER_SLASH) {
		min = *NewDuration(MINUTE_MIN, MINUTE_MAX, 1)
	} else if len(minCron) == 2 {
		var num int

		if strings.EqualFold(strings.Split(cronInfo[0], SEPARATER_SLASH)[0], ASTAH) {
			num = MINUTE_MAX
		} else {
			num, _ = strconv.Atoi(minCron[0])
		}
		denom, err := strconv.Atoi(minCron[1])

		if err != nil {
			return err, nil
		}
		min = *NewDuration(MINUTE_MIN, num, denom)
	}

	return nil, &CronInfo{
		Minute: min,
		Hour:   hour,
		Day:   day,
		Month:  mon,
		Week:   week,
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