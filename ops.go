package bloomzinho

func parsEq(a, b *Filter) bool {
	if a == b {
		return true
	}

	return a.filterParams == b.filterParams && len(a.state) == len(b.state)
}

//TODO add union and intersection

func (f *Filter) Intersects(b *Filter) bool {
	if !parsEq(f, b) {
		return false
	}

	for i := range f.state {
		res := f.state[i] & b.state[i]
		if res != 0 {
			return true
		}
	}

	return false
}
