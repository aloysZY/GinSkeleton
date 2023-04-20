package kube_client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetOneClientset(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := GetConfig(kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

//type Client struct {
//	*kubernetes.Clientset
//}

func GetConfig(kubeconfigPath string) (*rest.Config, error) {

	if kubeconfigPath == "" {
		//// home是家目录，如果能取得家目录的值，就可以用来做默认值
		//if home := homedir.HomeDir(); home != "" {
		//	// 如果输入了kubeconfig参数，该参数的值就是kubeconfig文件的绝对路径，
		//	// 如果没有输入kubeconfig参数，就用默认路径~/.kube/config
		//	kubeconfigPath = filepath.Join(home, ".kube", "config")
		//}
		//从家目录去找
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedConfigDir)
		if err != nil {
			//使得运行在 Kubernetes 集群中的代码可以在集群中获取 Kubernetes client-go 的配置
			//假如我们的代码运行在 Kubernetes 集群中，InClusterConfig 可以帮助我们从集群中创建 config 对象
			config, err := rest.InClusterConfig()
			if err != nil {
				return nil, err
			}
			return config, nil
		}
		return config, nil
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

/*// 获取 client的连接句柄
func GetClient(configPath string, clientType string) (*Client, error) {
	config, err := GetConfig(configPath)
	if err != nil {
		variable.ZapLog.Error("GetConfig(kubeconfigPath) failed error:", zap.Error(err))
		return nil, err
	}
	clientType = strings.Trim(clientType, " ")
	if clientType == "" {
		clientType = variable.ConfigGormv2Yml.GetString("Kubernetes.UseClientType")
	}
	switch strings.ToLower(clientType) {
	// 根据传入参数，创建多重连接类型
	case "clientset":
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			variable.ZapLog.Error("kubernetes.NewForConfig(config) failed error:", zap.Error(err))
			return nil, err
		}
		return &Client{clientset}, nil
		// 初始化其他连接
	default:
		variable.ZapLog.Error(my_errors.ErrorsClientDriverNotExists + clientType)
	}
	return &Client{}, errors.New(my_errors.ErrorsClientDriverNotExists)
}*/
