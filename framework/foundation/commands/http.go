package commands

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/router"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func HttpCommand(a contract.Application) *cobra.Command {
	configServer, err := a.Get("config")
	if err != nil {
		panic(err)
	}

	var defaultPort = configServer.(contract.Config).GetInt("http.port")
	if defaultPort == 0 {
		defaultPort = 8080
	}
	var port int
	var isDaemon bool
	var httpCommand = &cobra.Command{
		Use:     "start",
		Aliases: []string{"server"},
		Short:   "Start http server (alias: server)",
		Long:    "Run attaches the router to a http.Server and starts listening and serving HTTP requests.",
		Run: func(cmd *cobra.Command, args []string) {
			router.Register(a)
			routerServer, err := a.Get("router")
			if err != nil {
				panic(err)
			}
			r := routerServer.(*gin.Engine)
			portStr := strconv.Itoa(port)
			if isDaemon {
				pidFile := filepath.Join(a.RuntimePath(), "ginco-"+portStr+".pid")
				logPath := filepath.Join(a.RuntimePath(), "logs")
				if _, err := os.Stat(logPath); err != nil {
					if err := os.MkdirAll(logPath, 0755); err != nil {
						panic(err)
					}
				}
				logFile := filepath.Join(logPath, "ginco-"+portStr+".log")
				context := &daemon.Context{
					PidFileName: pidFile,
					PidFilePerm: 0644,
					LogFileName: logFile,
					LogFilePerm: 0640,
					WorkDir:     a.BasePath(""),
					Umask:       027,
					Args:        []string{"", "start", "-p", portStr, "-d"},
				}

				child, err := context.Reborn()
				if err != nil {
					fmt.Println(err)
					return
				}
				if child != nil {
					fmt.Println("pid: ", child.Pid)
					fmt.Println("log: ", logFile)
					return
				}
				defer context.Release()
				if err := r.Run(":" + portStr); err != nil {
					panic(err)
				} else {
					fmt.Println("http started")
				}
				return
			}

			if err := r.Run(":" + portStr); err != nil {
				panic(err)
			}
		},
	}

	httpCommand.Flags().IntVarP(&port, "port", "p", defaultPort, "Listening and serving HTTP port")
	httpCommand.Flags().BoolVarP(&isDaemon, "daemon", "d", false, "Start http.Server daemon")
	return httpCommand
}

func HttpStopCommand(a contract.Application) *cobra.Command {
	configServer, err := a.Get("config")
	if err != nil {
		panic(err)
	}

	var defaultPort = configServer.(contract.Config).GetInt("http.port")
	if defaultPort == 0 {
		defaultPort = 8080
	}
	var port int

	var httpStopCommand = &cobra.Command{
		Use:   "stop",
		Short: "Stop http server",
		Long:  "Stop http server.",
		Run: func(cmd *cobra.Command, args []string) {
			pidFile := filepath.Join(a.RuntimePath(), "ginco-"+strconv.Itoa(port)+".pid")
			pid, err := ioutil.ReadFile(pidFile)
			if err != nil {
				fmt.Println(pidFile, ":", err)
				return
			}
			command := exec.Command("kill", string(pid))
			_ = command.Start()
			fmt.Println("http stop: ", port)
			_ = os.Remove(pidFile)
		},
	}

	httpStopCommand.Flags().IntVarP(&port, "port", "p", defaultPort, "Stop HTTP port")
	return httpStopCommand
}

func HttpRestartCommand(a contract.Application) *cobra.Command {
	configServer, err := a.Get("config")
	if err != nil {
		panic(err)
	}

	var defaultPort = configServer.(contract.Config).GetInt("http.port")
	if defaultPort == 0 {
		defaultPort = 8080
	}
	var port int

	var httpRestartCommand = &cobra.Command{
		Use:   "restart",
		Short: "Restart http server",
		Long:  "Restart http server.",
		Run: func(cmd *cobra.Command, args []string) {
			portStr := strconv.Itoa(port)
			cmdServer, _ := a.Get("cmd")
			c := cmdServer.(*cobra.Command)
			c.SetArgs([]string{"stop", "-p", portStr})

			fmt.Println("Stop Http server: ", portStr)
			_, err = c.ExecuteC()
			if err != nil {
				fmt.Println(err)
			}

			c.SetArgs([]string{"start", "-p", portStr, "-d"})
			fmt.Println("Start Http server: ", portStr)
			if _, err = c.ExecuteC(); err != nil {
				fmt.Println(err)
			}
		},
	}

	httpRestartCommand.Flags().IntVarP(&port, "port", "p", defaultPort, "Listening and serving HTTP port")
	return httpRestartCommand
}
