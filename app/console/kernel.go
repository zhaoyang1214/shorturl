package console

import (
	"github.com/spf13/cobra"
	"github.com/zhaoyang1214/ginco/app/console/command/version"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

func Register(a contract.Application) {
	cmdServer, err := a.Get("cmd")
	if err != nil {
		panic(err)
	}

	cmd := cmdServer.(*cobra.Command)

	cmd.AddCommand(version.Command(a))
}
