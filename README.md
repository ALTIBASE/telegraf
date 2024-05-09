# How to build Telegraf embedded the Altibase Input Plugin

Telegraf input plugins are used to collect metrics from target systems. The Altibase input plugin for Telegraf is provided to gathers metrics from an Altibase database system as one of Telegraf input plugins.

This document shows how to compile Telegraf embedded the Altibase input plugin step by step. Then, consult [Altibase Input Plugin][1] to configure and run the successfully compiled binary executable.

[1]: https://github.com/jiee-altibase/telegraf/blob/main/altibase/README.md

### 1. Intall and configure unixODBC

The Altibase input plugin is connected to an Altibase database system through unixODBC. So, it is mandatory to install a unixODBC and configure environment variables for the Altibase input plugin to reference it.

unixodbc-2.3.12 version is recommended for higher compatibility and download it at the official [unixODBC homepage](https://www.unixodbc.org).

After installing it, set environment variables as follows.

```
LD_LIBRARY_PATH = "$HOME/unixodbc/lib:$LD_LIBRARY_PATH"
CGO_CFLAGS = "-I$HOME/unixodbc/include"
CGO_LDFLAGS = "-L$HOME/unixodbc/lib"
```

### 2. Clone the Telegraf repository:

```
git clone https://github.com/influxdata/telegraf.git telegraf
```
The last verified version of Telegraf for the Altibase input plugin is v1.30.2.

### 3. Clone the Altibase input plugin repository:

```
git clone https://github.com/ALTIBASE/telegraf.git telegraf_altibase
```
The last verified version of the Altibase input plugin is v1.0.0.

### 4. Copy the Altibase input plugin into the Telegraf input plugin directory

```
cp -r telegraf_altibase/* telegraf/plugins/inputs
```

### 5. Modifiy Makefile

The Makefile is located in the root directory where you cloned the Telegraf repository.

You must set CGO_ENABLE=1 in the Makefile since the Altibase input plugin requires cgo.

```
...
.PHONY: build
build:
        CGO_ENABLED=1 go build -tags "$(BUILDTAGS)" -ldflags "$(LDFLAGS)" ./cmd/telegraf
...
$(buildbin):
        echo $(GOOS)
        @mkdir -pv $(dir $@)
        CGO_ENABLED=1 go build -o $(dir $@) -tags "$(BUILDTAGS)" -ldflags "$(LDFLAGS)" ./cmd/telegraf
```

### 6. Compile Telegraf embedded the Altibase input plugin
```
cd telegraf
make build
```

