package main

import (
	"fmt"
	"ipsetsv/handlers"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nadoo/ipset"
	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, cmdServe())
}

func cmdServe() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "accept ips from client",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "port",
				Usage:    "listen port",
				Value:    9090,
				Required: true,
			},
			&cli.StringFlag{
				Name:     "token",
				Usage:    "token for auth",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			if err := ipset.Init(); err != nil {
				log.Printf("error in ipset Init: %s", err)
				return err
			}

			e := echo.New()
			e.Use(middleware.Logger())
			e.Use(middleware.Recover())
			ipsetController := handlers.IPSet{
				Token:   c.String("token"),
				Timeout: 3600,
			}
			e.POST("/", ipsetController.SyncIPSet)
			return e.Start(fmt.Sprintf(":%d", c.Int("port")))
		},
	}
}
