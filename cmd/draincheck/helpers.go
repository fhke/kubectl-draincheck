package draincheck

import (
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func mustNewLogger() *zap.SugaredLogger {
	unsugared, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return unsugared.Sugar()
}

func defaultKubeConfig() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Could not get user home directory: " + err.Error())
	}

	return path.Join(homeDir, ".kube", "config")
}

func getKubeconfigPath(kubeconfig string) string {
	if k := os.Getenv("KUBECONFIG"); k != "" {
		return k
	} else {
		return kubeconfig
	}
}

func newClientset(kubeconfig string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func mustMarshalWrite(log *zap.SugaredLogger, m func() ([]byte, error)) {
	data, err := m()
	if err != nil {
		log.Panicw("Error marshalling data", "error", err)
	}

	fmt.Print(string(data))
}
