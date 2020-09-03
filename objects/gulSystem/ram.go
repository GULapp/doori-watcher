package gulSystem

type Ram struct {
	Usermode   uint32 `json:"usermode"`
	Systemmode uint32 `json:"systemmode"`
	Nice       uint32 `json:"nice"`
	Idle       uint32 `json:"idle"`
	Wait       uint32 `json:"wait"`
	Irq        uint32 `json:"irq"`
	Srq        uint32 `json:"srq"`
}
