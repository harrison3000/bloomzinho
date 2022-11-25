package bloomzinho

import "errors"

func parsEq(a, b *Filter) bool {
	return a.filterParams == b.filterParams
}

func (f *Filter) Intersects(b *Filter) bool {
	if !parsEq(f, b) {
		return false
	}

	if len(f.state) != len(b.state) {
		//this helps the bound checker
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

func copyFilterFunc(a, b *Filter, cb func(_, _ bucket_t) bucket_t) (*Filter, error) {
	if a == nil || b == nil {
		return nil, errors.New("one of the input filters is nil")
	}
	if !parsEq(a, b) {
		return nil, errors.New("both filters must have the same dimensions")
	}
	ret := &Filter{
		filterParams: a.filterParams,
		state:        make([]bucket_t, len(a.state)),
	}

	for i := range ret.state {
		ret.state[i] = cb(a.state[i], b.state[i])
	}

	return ret, nil
}

func NewUnion(a, b *Filter) (*Filter, error) {
	return copyFilterFunc(a, b, func(ab, bb bucket_t) bucket_t {
		return ab | bb
	})
}

func NewIntersection(a, b *Filter) (*Filter, error) {
	return copyFilterFunc(a, b, func(ab, bb bucket_t) bucket_t {
		return ab & bb
	})
}
