package cgroups

import (
	"github.com/1157913453/myDocker/cgroups/subsystems"
)

type CgroupManager struct {
	// cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	Path string

	// 资源配置
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 将进程加入到每个cgroup中
func (c *CgroupManager) Apply(pid int) {
	for _, subsysIns := range subsystems.SubsystemsIns {
		subsysIns.Apply(c.Path, pid)
	}
}

// 设置各个subsystem挂载中的cgroup资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Set(c.Path, res)
	}
}

// 释放个个挂载在subsystem中的cgroup
func (c *CgroupManager) Destroy() {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Remove(c.Path)
	}
}
