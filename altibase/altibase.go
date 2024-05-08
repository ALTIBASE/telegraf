//go:generate ../../../tools/readme_config_includer/generator
package altibase

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"sync"

	_ "github.com/alexbrainman/odbc"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/toml"
)

//go:embed sample.conf
var sampleConfig string

// Altibase struct
type Altibase struct {
	Dsn              string `toml:"dsn"`
	ConnectionString string `toml:"address"`
	QueryFile        string `toml:"query_file"`

	conn     *sql.DB
	queryMap QueryMap
}

type Query struct {
	Measurement string   `toml:"measurement"`
	Sql         string   `toml:"sql"`
	Tags        []string `toml:"tags"`
	Fields      []string `toml:"fields"`
	Pivot       bool     `toml:"pivot"`
	PivotKey    string   `toml:"pivot_key"`
	Enable      bool     `toml:"enable"`
}

type QueryMap struct {
	Title   string  `toml:"Title"`
	Queries []Query `toml:"sql_metric"`
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (a *Altibase) Init() error {
	if len(a.Dsn) == 0 {
		return fmt.Errorf("dsn must be set.")
	}

	if len(a.QueryFile) == 0 {
		return fmt.Errorf("query_file must be set.")
	}

	return nil
}

// Start starts the ServiceInput's service, whatever that may be
func (a *Altibase) Start(acc telegraf.Accumulator) error {
	buf, err := os.ReadFile(a.QueryFile)
	if err != nil {
		acc.AddError(err)
		return err
	}

	err = toml.Unmarshal(buf, &a.queryMap)
	if err != nil {
		acc.AddError(err)
		return err
	}

	a.conn, err = sql.Open("odbc", a.getConnectionString())
	if err != nil {
		acc.AddError(err)
		return err
	}

	_, err = a.conn.Exec("exec set_client_info('telegraf')")
	if err != nil {
		acc.AddError(err)
		return err
	}

	return nil
}

// Stop stops the services and closes connections
func (a *Altibase) Stop() {
	a.conn.Close()
}

func init() {
	inputs.Add("altibase", func() telegraf.Input {
		return &Altibase{}
	})
}

func (a *Altibase) SampleConfig() string {
	return sampleConfig
}

func (a *Altibase) getConnectionString() string {
	return fmt.Sprintf("DSN=%s;%s", a.Dsn, a.ConnectionString)
}

func (a *Altibase) Gather(acc telegraf.Accumulator) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		acc.AddError(a.runSQL(acc))
	}()
	wg.Wait()

	return nil
}

func (a *Altibase) runSQL(acc telegraf.Accumulator) error {
	for _, element := range a.queryMap.Queries {
		if !element.Enable {
			continue
		}

		rows, err := a.conn.Query(element.Sql)
		if err != nil {
			return err
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return err
		}

		result_data := make([]map[string]interface{}, 0)

		for rows.Next() {
			entry, err := a.getResult(rows, columns)
			if err != nil {
				return err
			}
			result_data = append(result_data, entry)
		}

		//fmt.Println("result_data : ", result_data)
		a.accRow(element, acc, result_data)
	}

	return nil
}

func (a *Altibase) getResult(row scanner, columns []string) (map[string]interface{}, error) {
	column_count := len(columns)

	value_data := make([]interface{}, column_count)
	value_ptrs := make([]interface{}, column_count)

	for i := 0; i < column_count; i++ {
		value_ptrs[i] = &value_data[i]
	}

	err := row.Scan(value_ptrs...)
	if err != nil {
		return nil, err
	}

	entry := make(map[string]interface{})

	for i, col := range columns {
		var v interface{}
		val := value_data[i]
		//fmt.Println("val", val)

		b, ok := val.([]byte)

		if ok {
			v = string(b)
		} else {
			v = val
		}
		//fmt.Println("v", v)
		entry[col] = v
	}

	return entry, nil
}

func (a *Altibase) accRow(element Query, acc telegraf.Accumulator, result_data []map[string]interface{}) {
	tags := make(map[string]string)
	fields := make(map[string]interface{})

	tags["dsn"] = a.Dsn

	for _, v := range result_data {
		for _, v2 := range element.Tags {
			tags[v2] = v[v2].(string)
		}

		if element.Pivot {
			key := v[element.PivotKey].(string)
			data := v[element.Fields[0]]
			fields[key] = data
		} else {
			for _, v2 := range element.Fields {
				fields[v2] = v[v2]
			}
		}
		//fmt.Println("accRow : ", fields)
		acc.AddFields(element.Measurement, fields, tags)
	}
}
