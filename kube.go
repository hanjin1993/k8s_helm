package k8s_helm

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

//k8s初始化 采用service_account 权限 操作k8s api
func InitKubectl() {
	k8sApiServer := "https://kubernetes.default:443"
	k8sToken := ""
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	if k8sToken == "" {
		b, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			panic("k8s RBAC is not exsit or error:" + err.Error())
		} else {
			k8sToken = string(b)
		}
	}
	cmd := exec.Command("kubectl", "config", "set-cluster", "default", "--insecure-skip-tls-verify=true", "--server="+k8sApiServer)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		//logrus.Info(err)
	}
	cmd = exec.Command("kubectl", "config", "set-credentials", "default", "--token="+k8sToken)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		//logrus.Info(err)
	}
	cmd = exec.Command("kubectl", "config", "set-context", "default", "--cluster=default", "--user=default")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		//logrus.Info(err)
	}
	cmd = exec.Command("kubectl", "config", "use-context", "default")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		//logrus.Info(err)
	}
}

func KubectlApply(filePath string, namespace string) (err error) {
	cmd := exec.Command("kubectl", "apply", "-f", filePath, "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func KubectlDelete(filePath string, namespace string) (err error) {
	cmd := exec.Command("kubectl", "delete", "-f", filePath, "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

//探针注入
func AddIstioLabel(namespace string) (err error) {
	istioLabel := "istio.io/rev=canary"
	cmd := exec.Command("kubectl", "label", "ns", namespace, istioLabel)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

//去掉探针
func DelIstioLabel(namespace string) (err error) {
	istioLabel := "istio.io/rev=canary"
	split := strings.Split(istioLabel, "=")
	cmd := exec.Command("kubectl", "label", "ns", namespace, split[0]+"-")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}
