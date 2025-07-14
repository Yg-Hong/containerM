# containerM 🐳

`containerM`는 Go 언어로 구현한 Linux 환경 경량 컨테이너 런타임입니다.
컨테이너의 핵심 개념인 **네임스페이스, cgroups, pivot_root, 가상 네트워크** 등을 직접 다루며, 
컨테이너가 어떻게 작동하는지를 학습하고 실험하기 위한 목적으로 제작되었습니다.

---

## ✨ 주요 기능

- ✅ UTS, PID, Mount, Network, User 네임스페이스 격리
- ✅ `pivot_root`를 통한 루트 파일 시스템 분리
- ✅ `/proc`, `/dev`, `/sys` 마운트
- ✅ `veth` 기반 가상 네트워크 + 브릿지 설정
- ✅ `/dev/null`, `/dev/zero`, `/dev/random` 등 기본 장치 노드 포함
- ✅ 리소스 정리를 위한 안전한 종료 처리 (`defer`, trap 등)

---

## 📁 디렉토리 구조
```
containerM/
├── main.go # 프로그램 엔트리포인트
├── go.mod # Go 모듈 설정
├── cmd/
│    └── run.go # run, child 명령 처리
├── container/
│    └── process.go # 컨테이너 프로세스 실행, 초기화
├── fs/
│    └── pivot.go # pivot_root 및 파일 시스템 마운트
├── namespace/
│    ├── uts.go # hostname 설정
│    ├── pid.go # (확장 가능)
│    ├── mount.go # /proc, /sys 마운트
│    ├── net.go # (확장 가능)
│    └── dev.go # /dev 마운트 및 디바이스 생성
├── scripts/
│    ├── host_setup.sh # 호스트 측 veth/bridge 설정
│    └── container_setup.sh # 컨테이너 측 eth0 설정
├── utils/
│    └── must.go # 에러 핸들 헬퍼
└── README.md
```

---

## ⚙️ 빌드 및 실행 방법

### 0. 테스트 환경
* ubuntu 24.04.2
  
### 1. 빌드

```bash
go build -o containerM .
```

### 2. 컨테이너 루트 파일 시스템 준비
```bash
mkdir -p /tmp/ubuntufs
debootstrap --variant=minbase focal /tmp/ubuntufs http://archive.ubuntu.com/ubuntu # 또는 busybox 환경 구성
```

### 3. 실행 예시
```bash
sudo ./containerM run /bin/bash
```
