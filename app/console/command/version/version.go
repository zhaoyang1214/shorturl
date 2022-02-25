package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

func Command(a contract.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get Application version",
		Long:  "Get Application version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("The Application version is v%s\n", a.Version())
		},
	}
}
