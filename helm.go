package k8s_helm

import (
	"bytes"
	"fmt"
	"github.com/hanjin1993/k8s_helm/vo"
	"log"
	"os/exec"
)

type HelmUtil struct {
}

func (s *HelmUtil) CheckName(name string) bool {
	ls, err := s.List(true, "")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	for _, v := range *ls {
		if v.Name == name {
			return false
		}
	}
	return true
}

func (s *HelmUtil) List(all bool, namespace string) (list *[]vo.Helm, err error) {
	// helm list
	list = &[]vo.Helm{}
	args := []string{}
	if all {
		args = append(args, "list", "--all-namespaces")
	} else {
		args = append(args, "list", "-n", namespace)

	}
	//err = bgcmd.ExecStruct("helm", args, list, nil)
	return
}

/**
删除应用
*/
func (h *HelmUtil) Del(name string, namespace string) (err error) {
	cmd := exec.Command("helm", "delete", name, "-n", namespace)
	err = cmd.Run()
	// 执行命令，并返回结果
	//output, err := cmd.Output()
	if err != nil {
		return
	}
	// 因为结果是字节数组，需要转换成string
	fmt.Printf("删除部署应用%s-%s:%s\n", namespace, name)
	return
}

func (h *HelmUtil) Add(name string, namespace string, dsc string) (err error) {
	cmd := exec.Command("helm", "install", name, dsc, "-n", namespace, "-f", dsc+"/values.yaml")
	fmt.Println("helm 添加", cmd.Args)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	// 执行命令，并返回结果
	//output, err := cmd.Output()
	if err != nil {
		return
	}
	// 因为结果是字节数组，需要转换成string
	fmt.Printf("添加部署应用%s-%s:%s\n", namespace, name)
	return
}
