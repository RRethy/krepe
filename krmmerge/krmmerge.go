package krmmerge

type KrmMerge struct {
	origin   map[string]any
	local    map[string]any
	upstream map[string]any
}

func NewKrmMerge(
	origin map[string]any,
	local map[string]any,
	upstream map[string]any,
) *KrmMerge {
	return &KrmMerge{
		origin,
		local,
		upstream,
	}
}

func (k *KrmMerge) Merge() map[string]any {
	return nil
}
