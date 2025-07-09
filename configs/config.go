package configs

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

type Config struct {
	KubeConfigPath string `json:"kubeConfigPath,omitempty"`
}

func Load() *Config {
	config := &Config{}
	flag.StringVar(&config.KubeConfigPath, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "path to a kubeconfig")
	flag.Parse()
	return config
}
