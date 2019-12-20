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
	ID       string `json:",omitempty"`
	Content  string `json:",omitempty"`
	Proof    string `json:",omitempty"`
	NewProof string `json:",omitempty"`
}

// MinMax is a simple struct representing a minimum and maximum length for a string.
type MinMax struct {
	Min, Max int
}

// Check checks if a string is within the length requirements.
func (minmax MinMax) Check(val string) bool {
	len := len(val)
	return !(len < minmax.Min || len > minmax.Max)
}
