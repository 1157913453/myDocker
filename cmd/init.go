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
	"path"
	"path/filepath"
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

		setUpMount()
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

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("获取当前目录失败：%v", err)
		return
	}
	log.Infof("当前的目录是：%v", pwd)
	pivotRoot(pwd)

	// 在容器内启动并初始化进程
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		log.Errorf("挂载proc失败：%v", err)
		return
	}

	if err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755"); err != nil {
		log.Errorf("挂载tmpfs失败：%v", err)
		return
	}
}

func pivotRoot(root string) {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		log.Errorf("Mount rootfs to itself error: %v", err)
		return
	}
	// 创建 rootfs/.pivot_root存储old_root
	pivotDir := path.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		log.Errorf("创建pivotDir失败：%v", err)
		return
	}

	// pivot_root 到新的rootfs, 现在老的old_root是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		log.Errorf("挂载就文件系统失败：%v", err)
		return
	}

	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		log.Errorf("Chdir错误：%v", err)
		return
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		log.Errorf("unmount pivot_root dir %v", err)
		return
	}
	// 删除临时文件夹
	if err := os.Remove(pivotDir); err != nil {
		log.Errorf("删除临时文件失败：%v", err)
		return
	}

}
