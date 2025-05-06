# macl
## 개요
* 네트워크 ACL 점검을 위한 간단한 애플리케이션입니다.
* controller와 agent로 구성되어 있습니다.
* controller는 ACL을 점검할 수 있는 명령어를 agent에 전달하고, agent는 해당 명령어를 수행하여 결과를 controller에 전달합니다.
* agent 는 controller 의 요청에 따라 패킷을 리스닝하거니 전송합니다. 


## Agent 
* Agent 는 controller 의 요청에 따라 패킷을 리스닝하거나 전송합니다.
* Agent 모드는 기본 옵션이기 때문에 별도로 설정하지 않아도 됩니다.

### Agent 실행
```shell
$ nohup ./macl & 
```

## Controller
* Controller 는 agent 에게 명령어를 전달하고, agent 의 응답을 수신합니다.

### 설정
* 확인하고 싶은 ACL을 config.csv 또는 config-{profile}.csv 파일에 작성합니다.
    * 헤더 없이 작성하세요.
    * config.csv
```csv
1,172.30.1.96,172.30.1.97,10001,tcp
2,172.30.1.96,172.30.1.97,10002,udp
```

### 실행
* controller 모드로 실행합니다.
```shell
$ ./macl -type controller
```
* 디버그 모드로 실행하고 싶으면 -profile k8s -debug true 옵션을 추가합니다.
```shell
$ ./macl -type controller -debug true
```
* 기본 control 통신 포트는 10000입니다. 만약, 변경하고 싶다면 -controlPort {port} 옵션을 추가합니다.
```shell
$ ./macl -type controller -controlPort 10001
```

* 설정파일을 분리하고 싶으면 설정파일을 config-{profile}.csv 형식으로 작성하고 -profile 옵션을 사용합니다. 
```shell
$ ./macl -type controller -profile k8s
```