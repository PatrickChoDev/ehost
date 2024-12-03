package main

import (
	"fmt"
	"os"

	"github.com/PatrickChoDev/ehost/utils"
	"github.com/urfave/cli/v2"
)

const hostsFilePath = "/etc/hosts"

func main() {
	app := &cli.App{
		Name:  "ehost",
		Usage: "Manage your /etc/hosts at ease",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add an entry to the hosts file",
				Action: func(c *cli.Context) error {
					ip := c.Args().Get(0)
					hostname := c.Args().Get(1)
					if ip == "" || hostname == "" {
						return fmt.Errorf("IP and Hostname are required")
					}
					err := utils.AddEntry(hostsFilePath, ip, hostname)
					if err != nil {
						return err
					}
					fmt.Println("Entry added successfully")
					return nil
				},
			},
			{
				Name:  "remove",
				Usage: "Remove entries from the hosts file",
				Action: func(c *cli.Context) error {
					arg1 := c.Args().Get(0)
					arg2 := c.Args().Get(1)
					if arg1 == "" {
						return fmt.Errorf("IP or Hostname is required")
					}
					switch utils.IsIPorHostEntry(arg1) {
					case utils.IP:
						if arg2 != "" {
							err := utils.RemoveEntry(hostsFilePath, arg1, arg2)
							if err != nil {
								return err
							}
							fmt.Println("Entry removed successfully")
						} else {
							err := utils.RemoveAllIP(hostsFilePath, arg1)
							if err != nil {
								return err
							}
							fmt.Println("All entries with the specified IP removed successfully")
						}
					case utils.HOST:
						err := utils.RemoveAllHostname(hostsFilePath, arg1)
						if err != nil {
							return err
						}
						fmt.Println("All entries with the specified Hostname removed successfully")
					default:
						return fmt.Errorf("Invalid IP or Hostname")
					}
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List entries from the hosts file",
				Subcommands: []*cli.Command{
					{
						Name:  "ip",
						Usage: "List all entries with the specified IP",
						Action: func(c *cli.Context) error {
							ip := c.Args().Get(0)
							if ip == "" {
								return fmt.Errorf("IP is required")
							}
							entries, err := utils.GetEntriesByIP(hostsFilePath, ip)
							if err != nil {
								return err
							}
							for _, entry := range entries {
								fmt.Println(entry)
							}
							return nil
						},
					},
					{
						Name:  "host",
						Usage: "List all entries with the specified Hostname",
						Action: func(c *cli.Context) error {
							hostname := c.Args().Get(0)
							if hostname == "" {
								return fmt.Errorf("Hostname is required")
							}
							entries, err := utils.GetEntriesByHostname(hostsFilePath, hostname)
							if err != nil {
								return err
							}
							for _, entry := range entries {
								fmt.Println(entry)
							}
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
