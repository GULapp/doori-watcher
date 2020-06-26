# go_monitoring

## 기능
리눅스전용 통합 모니터링 어플리케이션.

## 개발환경
* Ubuntu 20.04
* golang 1.15
* GoLand/Bash

## 방향
* 시스템모니터링
  * syslog
  * secure
  * cpu
  * ram
* 커스텀 어플리케이션 모니터링
  * user application 상태모니터링
  * user application log모니터링
* 3rd party 어플리케이션 모니터링
  * DB모니터링
  * Middleware모니터링

## 용어
* Agent : 모니터링 대상서버에 데이터 수집프로세스
* Server : 모니터링 대상서버로 부터 수집데이터를 가공처리 및 Client Notification
* Client : 모니터링 정보를 효율적으로 보여주는 UI어플리케이션

## todo
* log level에 따라 로그로직 구현부
* leveldb 에 실시간데이터 저장.
* Client prototype 개발.
