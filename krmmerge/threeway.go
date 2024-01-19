package krmmerge

func threeWayMergeAny(origin, local, upstream any) any {

	return nil
}

func threeWayMergeMap(origin, local, upstream map[string]any) (map[string]any, error) {
	keepKeys := make(map[string]struct{})
	for k := range local {
		keepKeys[k] = struct{}{}
	}
	for k := range upstream {
		keepKeys[k] = struct{}{}
	}

	res := make(map[string]any)
	for key := range keepKeys {
		originVal, originOk := origin[key]
		localVal, localOk := local[key]
		upstreamVal, upstreamOk := upstream[key]

		if originOk && !localOk && !upstreamOk {
		} else if !originOk && localOk && !upstreamOk {
			res[key] = localVal
		} else if !originOk && !localOk && upstreamOk {
			res[key] = upstreamVal
		} else if originOk && localOk && !upstreamOk {
			val := delta(localVal, originVal)
			if val != nil {
				res[key] = val
			}
		} else if originOk && !localOk && upstreamOk {
			res[key] = upstreamVal
		} else if !originOk && localOk && upstreamOk {
			res[key] = twoWayMerge(localVal, upstreamVal)
		} else if originOk && localOk && upstreamOk {
			res[key] = threeWayMergeAny(originVal, localVal, upstreamVal)
		}
	}

	return res, nil
}
