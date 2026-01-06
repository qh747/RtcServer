# WebRTC Server

## 1. environment

ubuntu 24.04 server

## 2.build

### 2.1 build signal server

```bash
cd scripts
./buildProto.sh
./buildSigServ.sh
```

### 2.2 build librtcbase

```bash
cd cpp/librtcbase
cmake -B build 
cmake --build build
```

## 3.run

```bash
cd scripts
./runSigServ.sh
```