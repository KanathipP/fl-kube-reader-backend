package config

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func LoadKubeRestConfig() (*rest.Config, error) {
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}
	kubeconfig := clientcmd.RecommendedHomeFile

	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
