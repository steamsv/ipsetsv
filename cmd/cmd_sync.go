package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ipsetsv/common"
	"ipsetsv/sync"

	"github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, cmdSync())
}

func cmdSync() *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "sync ips from db to server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Usage:    "path to config file",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			content, err := ioutil.ReadFile(c.String("config"))
			if err != nil {
				return err
			}
			if err = json.Unmarshal(content, &common.Conf); err != nil {
				return err
			}
			dbConfig := mysql.NewConfig()
			{
				config := common.Conf
				dbConfig.User = config.MySQL.User
				dbConfig.Passwd = config.MySQL.Password
				dbConfig.Net = "tcp"
				dbConfig.Addr = fmt.Sprintf("%s:%d", config.MySQL.Host, config.MySQL.Port)
				dbConfig.DBName = config.MySQL.Database
			}
			db, err := sql.Open("mysql", dbConfig.FormatDSN())
			if err != nil {
				return err
			}
			defer db.Close()
			syncer := sync.NewSyncer(db)
			return syncer.Sync()
		},
	}
}
