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
	Minute []cronDuration
	Hour []cronDuration
	Day []cronDuration
	Month []cronDuration
	Week []cronDuration
}

func (ci *CronInfo) addMinute(cd cronDuration) error {
	if &cd == nil {
		return fmt.Errorf("cronDuration is empty.")
	}

	ci.Minute = append(ci.Minute, cd)
	return nil
}

func (ci *CronInfo) addHour(cd cronDuration) error {
	if &cd == nil {
		return fmt.Errorf("cronDuration is empty.")
	}

	ci.Hour = append(ci.Hour, cd)
	return nil
}

func (ci *CronInfo) addDay(cd cronDuration) error {
	if &cd == nil {
		return fmt.Errorf("cronDuration is empty.")
	}

	ci.Day = append(ci.Day, cd)
	return nil
}

func (ci *CronInfo) addMonth(cd cronDuration) error {
	if &cd == nil {
		return fmt.Errorf("cronDuration is empty.")
	}

	ci.Month = append(ci.Month, cd)
	return nil
}

func (ci *CronInfo) addWeek(cd cronDuration) error {
	if &cd == nil {
		return fmt.Errorf("cronDuration is empty.")
	}

	ci.Week = append(ci.Week, cd)
	return nil
}

type cronDuration struct {
	min int				// 範囲指定（最小）
	max int				// 範囲指定（最大）
	duration int		// min - max間での実行期間
}

func NewDuration(min, max, duration int) *cronDuration {
	return &cronDuration{
		min: min,
		max: max,
		duration: duration,
	}
}

func NewDurationWithCron(min, max int, cronString string) *cronDuration {
	var result cronDuration

	if strings.EqualFold(cronString, "") {
		return nil
	}

	cronNum := strings.Split(cronString, SEPARATER_SLASH)
	if strings.EqualFold(cronString, ASTAH) {
		result = *NewDuration(min, max, 1)
	} else if len(cronNum) == 2 {
		var num int

		if strings.EqualFold(strings.Split(cronString, SEPARATER_SLASH)[0], ASTAH) {
			num = MINUTE_MAX
		} else {
			num, _ = strconv.Atoi(cronNum[0])
		}
		denom, err := strconv.Atoi(cronNum[1])

		if err != nil {
			return nil
		}
		result = *NewDuration(MINUTE_MIN, num, denom)
	}
	return &result
}

func (file FileInfo) GetCron() (*CronInfo, error) {

	var min []cronDuration
	var hour []cronDuration
	var day []cronDuration
	var mon []cronDuration
	var week []cronDuration

	cronInfo := strings.Split(file.Cron, " ")
	if !strings.EqualFold(file.Cron, "") && len(cronInfo) <= 4 {
		return nil, fmt.Errorf("invalid format.")
	} else if strings.EqualFold(file.Cron, "") {
		return nil, nil
	}

	if c := NewDurationWithCron(MINUTE_MIN, MINUTE_MAX, cronInfo[0]); c != nil {
		min = append(min, *c)
	}
	if c := NewDurationWithCron(HOUR_MIN, HOUR_MAX, cronInfo[1]); c != nil {
		hour = append(hour, *c)
	}
	if c := NewDurationWithCron(DAY_MIN, DAY_MAX, cronInfo[2]); c != nil {
		day = append(day, *c)
	}
	if c := NewDurationWithCron(MONTH_MIN, MONTH_MAX, cronInfo[3]); c != nil {
		mon = append(mon, *c)
	}
	if c := NewDurationWithCron(WEEK_MIN, WEEK_MAX, cronInfo[4]); c != nil {
		week = append(week, *c)
	}

	return &CronInfo{
		Minute: min,
		Hour:   hour,
		Day:   day,
		Month:  mon,
		Week:   week,
	}, nil
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