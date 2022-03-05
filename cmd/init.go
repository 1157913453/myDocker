/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "",
	Long:  `Init myDocker container`,

	// 这里的 init 函数是在容器内部执行的，也就是说,代码执行到这里后,容器所在的进程其实就已经创建出来了，
	// 这是本容器执行的第一个进程/bin/bash。
	Run: func(cmd *cobra.Command, args []string) {
		cmdArray := readUserCommand()
		if cmdArray == nil || len(cmdArray) == 0 {
			log.Errorf("获取用户命令错误，命令为空")
			return
		}
		// 在容器内启动并初始化进程
		defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

		// 使用 mount 先去挂载 proc 文件系统，以便后面通过 ps 等系统命令去查看当前进程资源的情况。
		syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
		path, err := exec.LookPath(cmdArray[0])
		if err != nil {
			log.Errorf("Exec loop path error %v", err)
			return
		}
		log.Infof("当前命令路径为： %s", path)
		if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
			log.Errorf(err.Error())
			return
		}
	},
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Errorf("初始化读管道失败:%v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
