package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

type mysqlSetting struct {
	Host    string        `short:"H" long:"host" default:"localhost" description:"Hostname"`
	Port    string        `short:"p" long:"port" default:"3306" description:"Port"`
	User    string        `short:"u" long:"user" default:"root" description:"Username"`
	Pass    string        `short:"P" long:"password" default:"" description:"Password"`
	Timeout time.Duration `long:"timeout" default:"5s" description:"Timeout to connect mysql"`
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

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/",
			opts.mysqlSetting.User,
			opts.mysqlSetting.Pass,
			opts.mysqlSetting.Host,
			opts.mysqlSetting.Port,
		),
	)
	if err != nil {
		return checkers.Critical(fmt.Sprintf("couldn't connect DB: %v", err))
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()
	ch := make(chan error, 1)
	var uptime int64

	go func() {
		ch <- db.QueryRow("SELECT VARIABLE_VALUE FROM information_schema.GLOBAL_STATUS WHERE VARIABLE_NAME='Uptime'").Scan(&uptime)
	}()

	select {
	case err = <-ch:
		// nothing
	case <-ctx.Done():
		err = fmt.Errorf("connection or query timeout")
	}

	if err != nil {
		return checkers.Critical(fmt.Sprintf("couldn't execute query: %v", err))
	}

	if opts.Crit > 0 && uptime < opts.Crit {
		return checkers.Critical(fmt.Sprintf("up %s < %s", uptime2str(uptime), uptime2str(opts.Crit)))
	} else if opts.Warn > 0 && uptime < opts.Warn {
		return checkers.Warning(fmt.Sprintf("up %s < %s", uptime2str(uptime), uptime2str(opts.Warn)))
	}
	return checkers.Ok(fmt.Sprintf("up %s", uptime2str(uptime)))
}
