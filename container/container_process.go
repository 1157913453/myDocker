package container

import (
	"os"
	"os/exec"
	"syscall"
)

// 这里是父进程，也就是当前进程执行的内容
//1. 这里的/proc/self/exe调用中,/proc/self/指的是当前运行进程自己的环境，exec其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始化
//2. 下面的CLONE参数就是去 fork 出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境
func NewParentProcess(tty, interactive bool, command string) *exec.Cmd {
	args := []string{"init", command}
	//调用myDocker自己，传递参数，创建namespace隔离的进程，command一般为/bin/bash
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC, // 设置命名空间
	}

	// 如果用户指定了-ti参数，就需要把当前进程的输入输出导入到标准输入输出上
	if tty && interactive {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
