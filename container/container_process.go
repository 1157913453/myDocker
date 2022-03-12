package container

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io/fs"
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
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	rootUrl := "/root/"
	mergeUrl := "/root/merge"
	NewWorkSpace(rootUrl, mergeUrl)
	cmd.Dir = mergeUrl
	return cmd, writePipe
}

func NewWorkSpace(rootUrl, mergeUrl string) {
	//CreateLower(rootUrl)
	CreateUpper(rootUrl)
	CreateWork(rootUrl)
	CreateMerge(rootUrl, mergeUrl)
}

func CreateLower(rootUrl string) {
	lowerUrl := rootUrl + "lower"
	if err := os.Mkdir(lowerUrl, 0777); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Errorf("创建lowerUrl失败：%v", err)
		}
	}
}

func CreateUpper(rootUrl string) {
	upperUrl := rootUrl + "upper"
	if err := os.Mkdir(upperUrl, 0777); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Errorf("创建upperUrl失败：%v", err)
		}
	}
}

func CreateWork(rootUrl string) {
	workUrl := rootUrl + "work"
	if err := os.Mkdir(workUrl, 0777); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Errorf("创建work失败：%v", err)
		}
	}
}

func CreateMerge(rootUrl, mergeUrl string) {
	if err := os.Mkdir(mergeUrl, 0777); err != nil {
		if !errors.Is(err, fs.ErrExist) {
			log.Errorf("创建mergeUrl失败：%v", err)
			return
		}
	}

	dirs := "lowerdir=" + rootUrl + "busybox," + "upperdir=" + rootUrl + "upper," + "workdir=" + rootUrl + "work"
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mergeUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("cmd.Run()错误：%v", err)
	}
}

func PathExits(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func DeleteWorkSpace(rootUrl, mntUrl string) {
	DeleteMountPoint(mntUrl)
	DeleteWriteLayer(rootUrl)
}

func DeleteMountPoint(mntUrl string) {
	cmd := exec.Command("umount", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("卸载%s失败：%v", mntUrl, err)
	}

	if err := os.RemoveAll(mntUrl); err != nil {
		log.Errorf("删除%s失败：%v", mntUrl, err)
	}
}

func DeleteWriteLayer(rootUrl string) {
	writeLayerUrl := rootUrl + "writeLayer"
	if err := os.RemoveAll(writeLayerUrl); err != nil {
		log.Errorf("删除读写层失败：%v", err)
	}

}
