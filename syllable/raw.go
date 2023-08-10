package syllable

// a raw phoneme
type Raw struct {
	Phoneme string
}

func (r *Raw) componentTag() {}
