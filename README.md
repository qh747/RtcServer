# WebRTC Server

## 1. environment

ubuntu 24.04 server

## 2. signal server

### 2.1 build signal server

```bash
cd scripts
./buildProto.sh
./buildSigServ.sh
```

## 2.2 run signal server

```bash
cd scripts
./runSigServ.sh
```

## 3. media server

### 3.1 build librtcbase

```bash
cd cpp/librtcbase
cmake -B build 
cmake --build build
```

