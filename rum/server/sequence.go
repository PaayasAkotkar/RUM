package rum

type ISequence[In any] struct {
	Name    string // to serach in profile
	Service string // to manipulate Service
	Rank    int    // sequence number whether 1 ,.....
	Input   *In
}
