package krmmerge

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// KrmMerge is a 3-tuple of origin, local, and upstream resources.
// These can be used to merge them into a single resource.
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

// Merge returns the result of performing a 3-way merge on the origin, local, and upstream resources.
// The specifics of the 3-way merge algorithm are described in the README.
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

	// A resource present in local and upstream but missing from origin will have local merged with upstream returned.
	if k.origin == nil && k.local != nil && k.upstream != nil {
		return nil, nil // TODO: merge local and upstream
	}

	return k.threeWayMergeRoot()
}

// threeWayMergeRoot performs a 3-way merge on the root of the origin, local, and upstream resources.
// It will take into account the GVK of the resources.
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

	return &unstructured.Unstructured{
		Object: threeWayMergeMap(k.origin.Object, k.local.Object, k.upstream.Object),
	}, nil
}
