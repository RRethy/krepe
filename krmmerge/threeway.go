package krmmerge

func threeWayMergeAny(origin, local, upstream any) any {

	return nil
}

func threeWayMergeMap(origin, local, upstream map[string]any) (map[string]any, error) {
	keys := make(map[string]struct{})
	for k := range local {
		keys[k] = struct{}{}
	}
	for k := range upstream {
		keys[k] = struct{}{}
	}

	res := make(map[string]any)
	for k := range keys {
		originVal, originOk := origin[k]
		localVal, localOk := local[k]
		upstreamVal, upstreamOk := upstream[k]

		if originOk && !localOk && !upstreamOk {
		} else if !originOk && localOk && !upstreamOk {
			res[k] = localVal
		} else if !originOk && !localOk && upstreamOk {
			res[k] = upstreamVal
		} else if originOk && localOk && !upstreamOk {
			val := delta(localVal, originVal)
			if val != nil {
				res[k] = val
			}
		} else if originOk && !localOk && upstreamOk {
			res[k] = upstreamVal
		} else if !originOk && localOk && upstreamOk {
			res[k] = twoWayMergeAny(localVal, upstreamVal)
		} else if originOk && localOk && upstreamOk {
			res[k] = threeWayMergeAny(originVal, localVal, upstreamVal)
		}
	}

	return res, nil
}
