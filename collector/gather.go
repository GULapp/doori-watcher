package collector

var gathers []Gather

type Gather interface {
	Gathering() error
	Done([]byte) error
	PrettyPrint()
}
