package common

import "encoding/json"

// Agent --> Server 프로토콜 정의함
// Header와 Body 구조를 가짐
type Protocol struct {
	Header Header `json:"header"`
	Body   Body   `json:"body"`
}

// Header는 데이터를 보낸 hostname, hostname Ip, Body 데이터의 Category
// Category는 system | process | ...
// Category가 system이면 cpu, ram 모니터링, process는 프로세스모니터링
type Header struct {
	Server   string `json:"server"`
	Ip       string `json:"ip"`
	Category string `json:"category"`
}

// Body는 Tr, 그리고, 어떤 형태인지 모르는 json.RawMessage([]byte) 구성되어 있음
// Server 쪽에서, Tr를 확인후, 그 에 맞게, Data(json.RawMessage) 변수로 가공함
type Body struct {
	Tr   string          `json:"tr"`
	Data json.RawMessage `json:"data"`
}

// Protocol Header 정보를 셋팅함. 이 Header부분은 자주 바뀌지 않으므로
// Init 이라는 함수명 부여
func (protocolSet *Protocol) Init(server string, ip string, category string) {
	protocolSet.Header.Server = server
	protocolSet.Header.Ip = ip
	protocolSet.Header.Category = category
}

// Protocol Body 정보를 셋팅함.
func (protocolSet *Protocol) Set(trName string, data json.RawMessage) {
	protocolSet.Body.Tr = trName
	protocolSet.Body.Data = data
}

// Protocol 셋팅이 완료되면, json.RawMessage로 변환하여 리턴함
func (protocolSet Protocol) Marshaling() (json.RawMessage, error) {
	return json.Marshal(protocolSet)
}