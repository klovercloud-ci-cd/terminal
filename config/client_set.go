package config

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"sync"
)

var config *rest.Config
var once sync.Once
var KubeConfig *rest.Config

// GetKubeConfig returns rest config for kubernetes.
func GetKubeConfig() *rest.Config {
	var config *rest.Config
	var err error
	if IsK8 == "True" {
		config, err = clientcmd.BuildConfigFromFlags("", "")
	} else {
		if home := homedir.HomeDir(); home != "" {

			kubeConfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeConfig file")
			config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)

		} else {
			config, err = clientcmd.BuildConfigFromFlags("", "")
		}
	}
	if err != nil {
		panic(err)
	}
	return config
}

// GetClientSet returns k8s clientSets
func GetClientSet() *kubernetes.Clientset {
	once.Do(func() {
		config = GetKubeConfig()
	})

	kcs, kcsErr := kubernetes.NewForConfig(config)

	if kcsErr != nil {
		log.Printf("failed to create pipeline clientset: %s", kcsErr)
	}
	return kcs
}
