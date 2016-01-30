package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"os"
)

type mysqlSetting struct {
	Host string `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port string `short:"p" long:"port" default:"3306" description:"Port"`
	User string `short:"u" long:"user" default:"root" description:"Username"`
	Pass string `short:"P" long:"password" default:"" description:"Password"`
}

type connectionOpts struct {
	mysqlSetting
	Crit int64 `short:"c" long:"critical" description:"critical if uptime seconds is less than this number"`
	Warn int64 `short:"w" long:"warning" description:"warning if uptime seconds is less than this number"`
}

func uptime2str(uptime int64) string {
	day := uptime / 86400
	hour := (uptime % 86400) / 3600
	min := ((uptime % 86400) % 3600) / 60
	sec := ((uptime % 86400) % 3600) % 60
	return fmt.Sprintf("%d days, %02d:%02d:%02d", day, hour, min, sec)
}

func main() {
	ckr := checkUptime()
	ckr.Name = "MySQL Uptime"
	ckr.Exit()
}

func checkUptime() *checkers.Checker {
	opts := connectionOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	_, err := psr.Parse()
	if err != nil {
		os.Exit(1)
	}

	db := mysql.New("tcp", "", fmt.Sprintf("%s:%s", opts.mysqlSetting.Host, opts.mysqlSetting.Port), opts.mysqlSetting.User, opts.mysqlSetting.Pass, "")
	err = db.Connect()
	if err != nil {
		return checkers.Critical("couldn't connect DB")
	}
	defer db.Close()

	rows, res, err := db.Query("SHOW GLOBAL STATUS LIKE 'Uptime'")
	if err != nil {
		return checkers.Critical("couldn't execute query")
	}

	idxValue := res.Map("Value")
	Uptime := rows[0].Int64(idxValue)

	if opts.Crit > 0 && Uptime < opts.Crit {
		return checkers.Critical(fmt.Sprintf("up %s < %s", uptime2str(Uptime), uptime2str(opts.Crit)))
	} else if opts.Warn > 0 && Uptime < opts.Warn {
		return checkers.Warning(fmt.Sprintf("up %s < %s", uptime2str(Uptime), uptime2str(opts.Warn)))
	}
	return checkers.Ok(fmt.Sprintf("up %s", uptime2str(Uptime)))
}

