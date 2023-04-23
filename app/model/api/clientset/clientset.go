package clientset

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/kube_client"
)

func CreateClientsetFactory() *Clientset {
	return &Clientset{variable.ControllerClient}
}

type Clientset struct {
	kc *kube_client.KubeControllerClient
}
