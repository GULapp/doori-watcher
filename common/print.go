package common

import (
	"reflect"
	LOG "watcher/common/log"
)

// 구조체 변수입력 받으면,
// 구조체의 저장된 형(type)대로 출력하는 함수
func PrintAsStruct(st interface{}) {
	ref := reflect.ValueOf(st)
	elements := ref.Elem()

	for i := 0; i < elements.NumField(); i++ {
		value := elements.Field(i)
		field := elements.Type().Field(i)
		LOG.Debug("%-10s : %v", field.Name, value.Interface())
	}
}

// 구조체 변수입력 받으면,
// Json Tag가 정의되어 있는 변수만 출력하는 함수
func PrintAsStructForJson(st interface{}) {
	ref := reflect.ValueOf(st)
	elements := ref.Elem()

	for i := 0; i < elements.NumField(); i++ {
		value := elements.Field(i)
		field := elements.Type().Field(i)
		switch field.Tag.Get("json") {
		case "-":
		default:
			LOG.Debug("%-10s : %v", field.Name, value.Interface())
		}
	}
}
