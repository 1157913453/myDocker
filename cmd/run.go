/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/1157913453/myDocker/container"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	containerName        string // --name=容器名字
	detach               bool   // -d 后台运行
	interactive          bool   // -i 交互
	tty                  bool   // -t 分配伪终端
	defaultContainerName string // 默认容器名字
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long:  `Create a new container with namespace and cgroups limit mydocker run -it [command]`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Errorf("Missing container command")
			return
		}
		Run(args)
	},
}

func Run(args []string) {
	command := args[0]                                              // 获取执行的命令，如/bin/bash
	parent := container.NewParentProcess(tty, interactive, command) // 创建进程
	if err := parent.Start(); err != nil {
		log.Error(err)
		return
	}
	parent.Wait()
	os.Exit(-1)
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "是否提供交互")
	runCmd.Flags().BoolVarP(&tty, "tty", "t", false, "是否提供伪终端")
	runCmd.Flags().BoolVarP(&detach, "detach", "d", false, "容器是否在后台运行")
	runCmd.Flags().StringVarP(&containerName, "name", "n", defaultContainerName, "容器名字")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
