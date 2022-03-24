package k8s_helm

import (
	"fmt"
	"github.com/spf13/pflag"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strconv"
	"time"
)

var K8sClient *rest.Config

/**
获取k8s连接
*/
func GetClientSet() (*kubernetes.Clientset, error) {
	kubeConfig := pflag.String("kubeConfig", "/root/.kube/config", "kubeConfig path")
	if K8sClient == nil {
		client, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
		if err != nil {
			return nil, err
		}
		K8sClient = client
	}
	return kubernetes.NewForConfig(K8sClient)
}

/**
获取
*/
func GetDeployment(namespace string, name string) (info *appv1.Deployment, err error) {
	client, _ := GetClientSet()
	info, err = client.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if info != nil {
		fmt.Println(info.Labels)
	}
	return
}

/**
获取deployment列表
*/
func GetDeploymentList(namespace string, labelSelector string) (list *appv1.DeploymentList, err error) {
	client, err := GetClientSet()
	if err != nil {
		return nil, nil
	}
	getOpt := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	list, err = client.AppsV1().Deployments(namespace).List(getOpt)
	return
}

func GetSvcList(namespace string) (list *v1.ServiceList, err error) {
	client, err := GetClientSet()
	if err != nil {
		return nil, nil
	}
	getOpt := metav1.ListOptions{}
	list, err = client.CoreV1().Services(namespace).List(getOpt)
	return
}

func GetSvcInfo(namespace string, name string) (info *v1.Service, err error) {
	client, err := GetClientSet()
	if err != nil {
		return nil, nil
	}
	info, err = client.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	return
}

func GetNamespaceList() (list *v1.NamespaceList, err error) {
	client, err := GetClientSet()
	list, err = client.CoreV1().Namespaces().List(metav1.ListOptions{})
	return
}

func GetNamespaceInfo(namespace string) (info *v1.Namespace, err error) {
	client, err := GetClientSet()
	info, err = client.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	return
}

func UpdateDeployment(name string, namespace string, labelString string) (err error) {
	client, err := GetClientSet()
	//getDeployment, err := GetDeployment(namespace, name)
	list, err := GetDeploymentList(namespace, labelString)
	if list != nil {
		for _, getDeployment := range list.Items {
			newValue := strconv.FormatInt(time.Now().Unix(), 10)
			if len(getDeployment.Spec.Template.Spec.Containers) > 0 {
				var isSelect bool
				if len(getDeployment.Spec.Template.Spec.Containers[0].Env) > 0 {
					for k, v := range getDeployment.Spec.Template.Spec.Containers[0].Env {
						if v.Name == "create-time" {
							getDeployment.Spec.Template.Spec.Containers[0].Env[k].Value = newValue
							isSelect = true
						}
					}
				}

				if isSelect == false {
					getDeployment.Spec.Template.Spec.Containers[0].Env = append(getDeployment.Spec.Template.Spec.Containers[0].Env, v1.EnvVar{
						Name:  "create-time",
						Value: newValue,
					})
				}
			}
			_, err = client.AppsV1().Deployments(namespace).Update(&getDeployment)
			fmt.Println(err)
		}
	}
	return
}

func GetMetrics() {
	client, err := GetClientSet()
	fmt.Println("get client err ", err)
	raw, err := client.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/namespaces/cloud/pods").DoRaw()
	fmt.Println("get metrics err ", err)
	fmt.Println(string(raw))
}
