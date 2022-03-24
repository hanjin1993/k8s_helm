# k8s_helm
k8s helm 操作示例合计

1. helm.go
* 采用部署helm chart 包方式部署应用 可以批量部署 批量卸载 回滚等操作

2. k8s.go
* 使用client-go serviceAccount方式 连接k8s api执行k8s的deployment service等操作

3. kube.go
* 初始化kubectl命令权限

4. yaml.go
* yaml的一些操作
