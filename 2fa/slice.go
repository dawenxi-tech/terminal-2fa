package main

func sliceMoveElement[S ~[]E, E any](s S, index int, offset int) S {
	if index < 0 || index >= len(s) {
		return s
	}
	end := index + offset
	switch {
	case offset == 0:
		return s
	case offset < 0:
		for i := index; i > end; i-- {
			if i-1 < 0 {
				break
			}
			s[i], s[i-1] = s[i-1], s[i]
		}
	default:
		// offset > 0
		for i := index; i < end; i++ {
			if i+1 >= len(s) {
				break
			}
			s[i], s[i+1] = s[i+1], s[i]
		}
	}
	//fmt.Println(s)
	return s
}
