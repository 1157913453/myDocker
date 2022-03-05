package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// 这里是父进程，也就是当前进程执行的内容
//1. 这里的/proc/self/exe调用中,/proc/self/指的是当前运行进程自己的环境，exec其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始化
//2. 下面的CLONE参数就是去 fork 出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境
func NewParentProcess(tty, interactive bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New Pipe error: %v", err)
		return nil, nil
	}
	//调用init
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC, // 设置命名空间
	}

	// 如果用户指定了-ti参数，就需要把当前进程的输入输出导入到标准输入输出上
	if tty && interactive {
		fmt.Println("ti生效了")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
