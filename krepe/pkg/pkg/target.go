package pkg

import (
	"fmt"

	"github.com/RRethy/krepe/krepe/pkg/types"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Target struct {
	APIVersion string
	Group      string
	Version    string
	Kind       string
	Name       string
	Namespace  string
}

func NewTarget(t types.Target) (Target, error) {
	if t.APIVersion != "" && (t.Group != "" || t.Version != "") {
		return Target{}, fmt.Errorf("cannot specify both apiVersion and group/version")
	}

	if t.APIVersion != "" {
		gv, err := schema.ParseGroupVersion(t.APIVersion)
		if err != nil {
			return Target{}, fmt.Errorf("parsing apiVersion: %w", err)
		}
		t.Group = gv.Group
		t.Version = gv.Version
	}

	return Target{
		APIVersion: t.APIVersion,
		Group:      t.Group,
		Version:    t.Version,
		Kind:       t.Kind,
		Name:       t.Name,
		Namespace:  t.Namespace,
	}, nil
}

func (t Target) ToTypesTarget() types.Target {
	if t.APIVersion != "" {
		return types.Target{
			APIVersion: t.APIVersion,
			Kind:       t.Kind,
			Name:       t.Name,
			Namespace:  t.Namespace,
		}
	} else {
		return types.Target{
			Group:     t.Group,
			Version:   t.Version,
			Kind:      t.Kind,
			Name:      t.Name,
			Namespace: t.Namespace,
		}
	}
}

func (t Target) Matches(res *Resource) bool {
	if t.APIVersion != "" && t.APIVersion != res.GetAPIVersion() {
		return false
	}
	gvk := res.GetObjectKind().GroupVersionKind()
	if t.Group != "" && t.Group != gvk.Group {
		return false
	}
	if t.Version != "" && t.Version != gvk.Version {
		return false
	}
	if t.Kind != "" && t.Kind != res.GetKind() {
		return false
	}
	if t.Name != "" && t.Name != res.GetName() {
		return false
	}
	if t.Namespace != "" && t.Namespace != res.GetNamespace() {
		return false
	}
	return true
}
