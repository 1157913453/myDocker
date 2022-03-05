package subsystems

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path"
	"strconv"
)

type CpuSubsystem struct{}

func (m *CpuSubsystem) Name() string {
	return "cpu"
}

func (m *CpuSubsystem) Set(cgroupPath string, res *ResourceConfig) {
	subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath)
	if subsysCgroupPath == "" {
		log.Errorf("获取cgropPath失败：%v", err)
		return
	}
	if res.CpuShare != "" {
		// 在宿主机的对应目录创建cgroup
		if err = os.Mkdir(subsysCgroupPath, 0755); err != nil {
			if !errors.Is(err, fs.ErrExist) {
				log.Errorf("创建%s失败：%v", cgroupPath, err)
				return
			}
		}
		// 设置这个cgroup的cpuShare限制，即将限制写入到cgroup对应目录的cpu.shares文件中。
		if err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuSet), 0644); err != nil {
			log.Errorf("写入cgroup的cpuShare失败：%v", err)
			return
		}
		log.Infof("设置cpuShare限制成功")
	}
}

func (m *CpuSubsystem) Apply(cgroupPath string, pid int) {
	if subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath); err == nil {
		if err = os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			log.Errorf("写入的tasks文件失败:%v", err)
			return
		}
		log.Infof("应用cpuShare成功")
	}
}

func (m *CpuSubsystem) Remove(cgroupPath string) {
	if subsysCgroupPath, err := GetCgroupPath(m.Name(), cgroupPath); err == nil {
		err = os.RemoveAll(subsysCgroupPath)
		if err != nil {
			log.Warnf("删除cgroup失败:%v", err)
		}
	}
}
