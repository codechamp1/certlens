package configs

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

type Config struct {
	Context        string `json:"context,omitempty"`
	KubeConfigPath string `json:"kubeConfigPath,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
	Name           string `json:"name,omitempty"`
}

func Load() *Config {
	config := &Config{}
	flag.StringVar(&config.Context, "context", "", "context to use from kubeconfig, if not set, the current context will be used")
	flag.StringVar(&config.KubeConfigPath, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "path to a kubeconfig")
	flag.StringVar(&config.Namespace, "namespace", "", "namespace to lens, if not set, all namespaces will be used")
	flag.StringVar(&config.Name, "name", "", "name of the secret to lens, if not set, all secrets will be listed")
	flag.Parse()
	return config
}
