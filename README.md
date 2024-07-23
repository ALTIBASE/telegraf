# How to build Telegraf embedded the Altibase Input Plugin

Telegraf utilizes input plugins to collect metrics from target systems. Among these, the **Altibase input plugin for Telegraf** serves the purpose of collecting metrics from an Altibase database system.

This document describes the sequential steps to compile Telegraf embedded the Altibase input plugin. Following successful compilation, refer to the [Altibase Input Plugin](https://github.com/ALTIBASE/telegraf/blob/main/altibase/README.md) for configuration and execution instructions.


### 1. Intall and configure unixODBC

The Altibase input plugin is connected to an Altibase database system through unixODBC. Therefore, it is mandatory to install a unixODBC and configure environment variables for the Altibase input plugin to reference unixODBC.

It is recommended to install the unixODBC-2.3.12 version for enhanced compatibility. This version can be downloaded from the [official unixODBC homepage](https://www.unixodbc.org/).

After installation, set environment variables as follows:

```bash
export UNIX_ODBC="unixODBC installed path"
export LD_LIBRARY_PATH="$UNIX_ODBC/lib:$LD_LIBRARY_PATH"
```

For more information about unixODBC settings, please refer to [unixODBC Guide for Altibase(Korean)](https://docs.altibase.com/pages/viewpage.action?pageId=11698379).

### 2. Clone the Telegraf repository
Clone the Telegraf repository using the following command:
```bash
git clone https://github.com/influxdata/telegraf.git telegraf
```
The latest Telegraf version compatible with the Altibase input plugin is v1.30.2.

### 3. Clone the Altibase input plugin repository
Clone the Altibase input plugin repository with the following command:
```bash
git clone https://github.com/ALTIBASE/telegraf.git telegraf_altibase
```
The latest version of the Altibase input plugin is v1.0.0.

### 4. Copy the Altibase input plugin into the Telegraf input plugin directory
Copy the Altibase input plugin into the Telegraf input plugin directory using the following command:
```bash
cp telegraf_altibase/all/altibase.go telegraf/plugins/inputs/all
cp -r telegraf_altibase/altibase telegraf/plugins/inputs
```

### 5. Modifiy Makefile

The Makefile is located in the root directory where the Telegraf repository was cloned.

Modify the Makefile by setting CGO_ENABLE=1 since the Altibase input plugin requires cgo.

```makefile
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

Before compiling Telegraf embedded the Altibase input, you need to set some environment variables.

This depends on the SQLLEN size of unixODBC, so check the SQLLEN size by executing the command below.
```bash
${UNIX_ODBC}/bin/odbcinst -j
```

Then, set environment variables according to SQLLEN size.

If SQLLEN=8
```bash
export CGO_CFLAGS="-I$UNIX_ODBC/include"
export CGO_LDFLAGS="-L$UNIX_ODBC/lib"
```

If SQLLEN=4
```bash
export CGO_CFLAGS="-I$UNIX_ODBC/include -DBUILD_LEGACY_64_BIT_MODE=1"
export CGO_LDFLAGS="-L$UNIX_ODBC/lib -lodbc"
```

Finally, navigate to the Telegraf directory and execute the following command to compile Telegraf embedded the Altibase input plugin:
```bash
cd telegraf
go get github.com/alexbrainman/odbc
make build
```
