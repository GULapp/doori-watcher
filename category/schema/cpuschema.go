package schema

type CpuSchema struct {
	interestedClients []interestedClient
	Usage             int
}

func NewCpuSchema() *CpuSchema {
	return &CpuSchema{Usage: -1}
}

