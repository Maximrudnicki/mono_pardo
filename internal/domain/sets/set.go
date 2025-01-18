package sets

type WordSet struct{}

func NewWordSet() (*WordSet, error) {
	return &WordSet{}, nil
}
