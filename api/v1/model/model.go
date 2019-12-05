package model

var (
	IDLen = MinMax{
		Min: 4,
		Max: 16,
	}
	ContentLen = MinMax{
		Min: 0,
		Max: 65535,
	}
	ProofLen = 32
)

type Pad struct {
	ID, Content, Proof, NewProof string
}

type MinMax struct {
	Min, Max int
}

func (minmax MinMax) Check(val string) bool {
	len := len(val)

	if len < minmax.Min || len > minmax.Max {
		return false
	}

	return true
}
