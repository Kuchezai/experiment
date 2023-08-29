package handlers

type Validator struct {
}

//Own validator, since the built-in GIN does not perform its functions,
//and i do not want to drag a dependency in the form of another validator
func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) IsIntersect(a, b []string) bool {
	m := make(map[string]struct{})
	for _, el := range a {
		m[el] = struct{}{}
	}

	for _, el := range b {
		if _, ok := m[el]; ok {
			return true
		}
	}
	return false
}

func (v *Validator) checkAddedSegmentsIsValid(segments []AddSegments) bool {
	for _, seg := range segments {
		if seg.Slug == "" {
			return false
		}
		if seg.TTL < 0 || seg.TTL > 366 {
			return false
		}
	}
	return true
}
