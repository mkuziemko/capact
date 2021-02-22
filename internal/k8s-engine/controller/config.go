package controller

import "time"

type Config struct {
	BuiltinRunner BuiltinRunnerConfig
	ClusterPolicy ClusterPolicyConfig
}

type BuiltinRunnerConfig struct {
	Timeout time.Duration `envconfig:"default=30m"`
	Image   string
}

type ClusterPolicyConfig struct {
	Name      string `envconfig:"default=voltron-engine-cluster-policy"`
	Namespace string `envconfig:"default=voltron-system"`
}