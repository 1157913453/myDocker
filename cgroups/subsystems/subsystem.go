package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	// 返回subsystem名字
	Name() string

	// 设置某个cgroup在该subsystem的资源限制
	Set(path string, res *ResourceConfig)

	// 将进程添加进cgroup中
	Apply(path string, pid int)

	// 将进程从cgroup中移除
	Remove(path string)
}

var (
	SubsystemsIns = []Subsystem{
		&CpuSubsystem{},
		&CpusetSubsystem{},
		&MemorySubsystem{},
	}
)
