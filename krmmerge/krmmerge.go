package krmmerge

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type KrmMerge struct {
	origin   *unstructured.Unstructured
	local    *unstructured.Unstructured
	upstream *unstructured.Unstructured
}

func NewKrmMerge(
	origin *unstructured.Unstructured,
	local *unstructured.Unstructured,
	upstream *unstructured.Unstructured,
) *KrmMerge {
	return &KrmMerge{
		origin,
		local,
		upstream,
	}
}

func (k *KrmMerge) Merge() (*unstructured.Unstructured, error) {
	// A resource missing from origin, local, and upstream will have nil returned.
	if k.origin == nil && k.local == nil && k.upstream == nil {
		return nil, nil
	}

	// A resource present in origin but missing from local and upstream will have nil returned.
	if k.origin != nil && k.local == nil && k.upstream == nil {
		return nil, nil
	}

	// A resource present in local but missing from origin and upstream will have local returned.
	if k.origin == nil && k.local != nil && k.upstream == nil {
		return k.local, nil
	}

	// A resource present in upstream but missing from origin and local will have upstream returned.
	if k.origin == nil && k.local == nil && k.upstream != nil {
		return k.upstream, nil
	}

	if k.origin != nil && k.local == nil && k.upstream != nil {
		return nil, nil
	}

	// A resource present in origin and local but missing from upstream will have nil returned.
	if k.origin != nil && k.local != nil && k.upstream == nil {
		return nil, nil
	}

	if k.origin == nil && k.local != nil && k.upstream != nil {
		return nil, nil // TODO: merge local and upstream
	}

	return k.threeWayMergeRoot()
}

func (k *KrmMerge) threeWayMergeRoot() (*unstructured.Unstructured, error) {
	originGVK := k.origin.GetObjectKind().GroupVersionKind()
	localGVK := k.local.GetObjectKind().GroupVersionKind()
	upstreamGVK := k.upstream.GetObjectKind().GroupVersionKind()

	if originGVK.Kind != localGVK.Kind || originGVK.Kind != upstreamGVK.Kind {
		return nil, fmt.Errorf("kind mismatch: origin=%s, local=%s, upstream=%s", originGVK.Kind, localGVK.Kind, upstreamGVK.Kind)
	}

	if originGVK.Group != localGVK.Group || originGVK.Group != upstreamGVK.Group {
		return nil, fmt.Errorf("group mismatch: origin=%s, local=%s, upstream=%s", originGVK.Group, localGVK.Group, upstreamGVK.Group)
	}

	if originGVK.Version != localGVK.Version || originGVK.Version != upstreamGVK.Version {
		// TODO: handle version mismatch
		return nil, fmt.Errorf("version mismatch: origin=%s, local=%s, upstream=%s", originGVK.Version, localGVK.Version, upstreamGVK.Version)
	}

	m, err := threeWayMergeMap(k.origin.Object, k.local.Object, k.upstream.Object)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{
		Object: m,
	}, nil
}
