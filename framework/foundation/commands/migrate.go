package commands

import (
	"fmt"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/spf13/cobra"
	"github.com/zhaoyang1214/ginco/database/migration"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"strings"
)

func MigrateCommand(a contract.Application) *cobra.Command {
	return migrateCommand(a, "migrate")
}

func MigrateRollbackCommand(a contract.Application) *cobra.Command {
	return migrateCommand(a, "migrate:rollback")
}

func migrateCommand(a contract.Application, use string) *cobra.Command {
	var key, id string
	var keys []string

	migrates := migration.Init(a)

	for k := range migrates {
		keys = append(keys, k)
	}
	desc := "Rollback the last database migration"
	isMigrate := use == "migrate"
	if isMigrate {
		desc = "Run the database migrations"
	}

	var command = &cobra.Command{
		Use:   use,
		Short: desc,
		Long:  desc,
		Run: func(cmd *cobra.Command, args []string) {
			log := a.GetI("log").(contract.LoggerManager).Channel("stderr")
			defer log.Sync()
			ms := make(map[string]*gormigrate.Gormigrate)
			if key == "" {
				ms = migrates
			} else if m, ok := migrates[key]; ok {
				ms[key] = m
			} else {
				log.Fatal("migrations key not found: " + key)
			}

			for k, m := range ms {
				var err error
				if id == "" {
					if isMigrate {
						err = m.Migrate()
					} else {
						err = m.RollbackLast()
					}
				} else {
					if isMigrate {
						err = m.MigrateTo(id)
					} else {
						err = m.RollbackTo(id)
					}
				}

				if err != nil {
					if isMigrate {
						log.Fatal(fmt.Sprintf("Could not migrate: %v", err))
					} else {
						log.Fatal(fmt.Sprintf("Could not rollback: %v", err))
					}
				}
				if isMigrate {
					log.Info("migrations (" + k + ") did run successfully")
				} else {
					log.Info("rollback (" + k + ") did run successfully")
				}
			}
		},
	}

	command.Flags().StringVarP(&key, "key", "k", "", "Gormigrate Key: "+strings.Join(keys, ", "))
	command.Flags().StringVarP(&id, "id", "i", "", "Migrations id")
	return command
}
