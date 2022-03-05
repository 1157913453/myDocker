package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// 通过/proc/self/mountinfo找出挂载了某个subsystem的hierarchy cgroup根节点所在的目录,返回结果示例：/sys/fs/cgroup/memory/mydocker-cgroup
func GetCgroupPath(subsystem, cgroupPath string) (string, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fileds := strings.Split(txt, " ")
		for _, opt := range strings.Split(fileds[len(fileds)-1], ",") {
			if opt == subsystem {
				subsysCgroupPath := path.Join(fileds[4], cgroupPath)
				if _, err = os.Stat(subsysCgroupPath); err != nil { // 如果路径不存在
					return subsysCgroupPath, err
				}
				return subsysCgroupPath, nil
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("在/proc/self/mountinfo中没找到对应的%s", subsystem)
}
