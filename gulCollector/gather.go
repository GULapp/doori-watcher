package gulCollector

var gathers []Gather

/**
Gathering 리눅스 자원 모니터링, 결과는 []byte 형식으로
Done 그 다음 행위(보통은, Monitoring SERVER에 []byte를 보냄)
PrettyPrint human friendly print
 */
type Gather interface {
	Gathering() []byte
	Done([]byte) error
	PrettyPrint()
}
