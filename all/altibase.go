//go:build !custom || inputs || inputs.altibase

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/altibase" // register plugin
