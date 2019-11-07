# Philips Hue
This plugin is a part of [Casa](https://github.com/getcasa), it's used to interact with Philips Hue ecosystem.

## Downloads
Use the integrated store in casa or [github releases](https://github.com/getcasa/plugin-philipshue/releases).

## Build
```
sudo env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -buildmode=plugin -o philipshue.so *.go
```

## Install
1. Extract `philipshue.zip`
2. Move `philipshue.so` to casa `plugins` folder
3. Restart casa gateway
