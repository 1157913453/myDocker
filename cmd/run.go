/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/1157913453/myDocker/cgroups"
	"github.com/1157913453/myDocker/cgroups/subsystems"
	"github.com/1157913453/myDocker/container"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	containerName string // --name=容器名字
	detach        bool   // -d 后台运行
	interactive   bool   // -i 交互
	tty           bool   // -t 分配伪终端
	memoryLimit   string
	cpuShare      string
	cpuSet        string
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "是否提供交互")
	runCmd.Flags().BoolVarP(&tty, "tty", "t", false, "是否提供伪终端")
	runCmd.Flags().BoolVarP(&detach, "detach", "d", false, "容器是否在后台运行")
	runCmd.Flags().StringVarP(&containerName, "name", "n", "", "容器名字")
	runCmd.Flags().StringVarP(&memoryLimit, "memoryLimit", "m", "", "内存限制")
	runCmd.Flags().StringVar(&cpuShare, "cpuShare", "", "cpu核心限制")
	runCmd.Flags().StringVar(&cpuSet, "cpuSet", "", "cpu限制")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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

func Run(args []string) { // 获取执行的命令，如/bin/bash
	res := &subsystems.ResourceConfig{
		MemoryLimit: memoryLimit,
		CpuShare:    cpuShare,
		CpuSet:      cpuSet,
	}
	parent, writePipe := container.NewParentProcess(tty, interactive) // 创建进程
	if err := parent.Start(); err != nil {
		log.Error(err)
		return
	}
	// 给父进程进行cgroup限制
	cgroupManager := cgroups.NewCgroupManager("mydockerCgroup")
	defer cgroupManager.Destroy() //删除cgroup
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	fmt.Println("当前的进程id：", parent.Process.Pid)
	sendInitcommand(args, writePipe)
	err := parent.Wait()
	if err != nil {
		log.Errorf("父进程执行出错：%v", err)
		return
	}
	os.Exit(-1)
}

func sendInitcommand(cmdArgs []string, writePipe *os.File) {
	command := strings.Join(cmdArgs, " ")
	log.Infof("所有的参数是 %s", command)
	_, err := writePipe.WriteString(command)
	if err != nil {
		log.Errorf("写管道错误：%v", err)
		return
	}
	err = writePipe.Close()
	if err != nil {
		log.Errorf("管道关闭错误：%v", err)
		return
	}
}
