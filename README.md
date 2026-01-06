# WebRTC Server

## 1.build

### 1.1 build signal server

```bash
cd scripts
./buildProto.sh
./buildSigServ.sh
```

### 1.2 build librtcbase

```bash
cd cpp/librtcbase
cmake -B build 
cmake --build build
```

## 2.run

```bash
cd scripts
./runSigServ.sh
```