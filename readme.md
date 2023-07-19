onTune Kube manager beta v1.0

# Installation
- onTuneKubeManager.zip 파일 압축 해제
- 압축 파일 내용
  - onTuneKubeManager.exe: Manager 실행 파일
  - config/config.env: 환경 설정 파일
  - log/: 로그 저장 디렉토리(설치 해제 시 빈 디렉토리로 생성됨)

# DB 설정 방법
- db_host: Host 주소
- db_name: Database 명칭
- db_port: Port 주소
- db_user: DB Username
  - 입력 후 프로그램 실행 시 db_user 설정값이 사라짐
  - 대신 암호화된 db_stamp 값이 생성됨
  - config 디렉토리에 16진수 .dat 파일도 생성되며, 해당 파일에서 암호화 값 관리
- db_password: DB Password
  - 입력 후 프로그램 실행 시 db_password 설정값이 사라짐
  - 대신 암호화된 db_code 값이 생성됨
  - config 디렉토리에 16진수 .dat 파일도 생성되며, 해당 파일에서 암호화 값 관리

# Kubernetes Cluster 설정 방법
- clustercount
  - 기본값은 1이며, Multi-Cluster의 경우 해당 숫자를 입력해야 함
  - 하기 언급되는 kube_cluster 관련 설정값과 일치하지 않을 경우 에러가 발생하므로 유의 필요

- kube_cluster_xxxx_n: n은 cluster 번호로, clustercount 개수에 맞게 각각 입력되어야 함
  - 예제: kube_cluster_conffile_1, kube_cluster_conffile_2, ...
- kube_cluster_userid_n: Kubernetes Cluster Host 접속을 위한 계정명
- kube_cluster_password_n: Kubernetse Cluster Host 접속을 위한 계정 비밀번호
  - userid, password는 각각 stamp, code로 변환되며, 상기 언급된 db와 동일 방식으로 변경됨
- kube_cluster_conffile_n
  - Kubernetes Cluster 접속을 위한 Config File 경로
  - config/ 경로에 해당 파일이 위치해야 함
  - Config 파일 설정은 하기 기술
- kube_cluster_sshport_n: SSH 접속을 위한 Port이나, 현재는 사용하지 않음

# Interval 설정 방법
- disconnectfailcount: Disconnection이 이루어져서 실패를 체크하기 위한 카운트
- eventloginterval: Event 데이터 수집 주기
- resourceinterval: Info 데이터 수집 주기
- realtimeinterval: Metric 데이터 수집 주기
- resourcechangeinterval: 
  - Metric 데이터 수집 중 Info 데이터의 변경을 감지할 때, Info 데이터 주기를 일시적으로 Metric 데이터 수집 주기와 동일하게 함
  - 이 때 일시적으로 변경된 Info 데이터 수집 주기 유지 시간을 나타냄
- 기타 값: 현재는 사용하지 않으나, 차후 사용 예정

# Kubernetes Config 파일 생성
- Kubernetes Control-Plane 서버에 접속해서 다음의 명령어 실행 필요

```Bash
[root@control-plane ~]# server=https://x.x.x.x:6443
[root@control-plane ~]# name=ks-admin
[root@control-plane ~]# ca=$(kubectl get secret/ks-admin -o jsonpath='{.data.ca\.crt}')
[root@control-plane ~]# token=$(kubectl get secret/ks-admin -o jsonpath='{.data.token}' | base64 --decode)
[root@control-plane ~]# echo "
apiVersion: v1
kind: Config
clusters:
- name: kubernetes
  cluster:
    certificate-authority-data: ${ca}
    server: ${server}
contexts:
- name: ${name}@kubenetes
  context:
    cluster: kubernetes
    user: ${name}
current-context: ${name}@kubenetes
users:
- name: ${name}
  user:
    token: ${token}
" > kubeconfig
```

- 위에서 생성된 config 파일을 Kubernetes Manager의 config/ 경로에 .kube.1.config 등의 이름으로 저장
- 저장된 파일 이름을 config.env 파일의 kube_cluster_conffile_n 부분에 명시해야 함