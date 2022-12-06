package bloomzinho

import "errors"

func parsEq(a, b *Filter) bool {
	//we only compare the number of hashes and number of bits
	//the other values are derived from these 2 values anyway
	return a.filterParams == b.filterParams
}

// Intersects tests if the filter intersects with another given filter "b"
// 2 filters intersecting means that a value that was added to one filter was probably also added to the other
//
// this function is useful if you have a set of values that you want to test against a bunch of filters
// it's more efficient to make a new filter and test for intersection than to do multiple lookups
//
// if the 2 filters being tested have different size or number of hashes, this method gives a false positive
func (f *Filter) Intersects(b *Filter) bool {
	if !parsEq(f, b) {
		return true
	}

	if len(f.state) != len(b.state) {
		//this helps the bound checker
		return true
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
