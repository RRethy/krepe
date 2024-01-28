package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Krepe struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata,omitempty"`

	Imports   Imports    `yaml:"imports,omitempty"`
	Pipelines []Pipeline `yaml:"pipelines,omitempty"`
}

type Imports struct {
	Files    []string        `yaml:"files,omitempty"`
	Packages []PackageImport `yaml:"packages,omitempty"`
}

type PackageImport struct {
	Ref  string `yaml:"ref,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type Pipeline struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps,omitempty"`
}

type Step struct {
	Function  string         `yaml:"function"`
	Target    Target         `yaml:"target,omitempty"`
	ConfigMap map[string]any `yaml:"configMap,omitempty"`
}

type Target struct {
	APIVersion string `yaml:"apiVersion,omitempty"`
	Group      string `yaml:"group,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Kind       string `yaml:"kind,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Namespace  string `yaml:"namespace,omitempty"`
}
