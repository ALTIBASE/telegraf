# Altibase Input Plugin

The Altibase input plugin gathers metrics from an Altibase database system.

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

[CONFIGURATION.md]: https://github.com/influxdata/telegraf/blob/master/docs/CONFIGURATION.md#plugins

## Configuration

```toml @sample.conf
# Read metrics from one altibase server
[[inputs.altibase]]
  ## Specify DSN
  dsn = "Altiodbc"

  ## Setting the address is optional, and you can set only some of the items.
  ## If you set the address, the connection configuration in dsn will be overwritten.
  # address = "SERVER=127.0.0.1;PORT=20300;USER=sys;PASSWORD=manager"

  ## Specify metric query filename
  query_file = "query.toml"
```

## Metrics

The Altibase input plugin gathers the execution results of SQL in query.toml, primarily querying Altibase's performance views. Refer to the [General Reference-2.The Data Dictionary](https://github.com/ALTIBASE/Documents/blob/master/Manuals/Altibase_7.3/eng/General%20Reference-2.The%20Data%20Dictionary.md) for details.


Certain metrics can only be collected when TIMED_STATISTICS=1 is set in the Altibase server. Refer to the [General Reference-1.Data Types & Altibase Properties](https://github.com/ALTIBASE/Documents/blob/master/Manuals/Altibase_7.3/eng/General%20Reference-1.Data%20Types%20%26%20Altibase%20Properties.md) for details.


* altibase - Metrics related to Altibase overview. This measurement encompasses fields such as:
  * ALTIBASE_VERSION
  * ARCHIVE_MODE
  * DISK_TBL_USAGE      
  * FULLSCAN_QUERY_COUNT
  * HIT_RATIO           
  * LF_PREPARE_WAIT_COUNT
  * LOCK_HOLD_COUNT     
  * LOCK_WAIT_COUNT     
  * LOGFILE_CURRENT     
  * LOGFILE_GAP         
  * LOGFILE_OLDEST      
  * LONG_RUN_QUERY_COUNT
  * MEMORY_TBL_USAGE    
  * MEMSTAT_ALLOC       
  * MEMSTAT_MAX_TOTAL   
  * MEM_DELTHR          
  * MEM_LOGICAL_AGER    
  * REP_RECEIVER_COUNT  
  * REP_SENDER_COUNT    
  * SESSION_COUNT       
  * STATEMENT_COUNT     
  * UTRANS_QUERY_COUNT  
  * VICTIM_FAILS        
  * WORKING_TIME        
* Altibase usage - Statistics related to Altibase's memory or disk usage. There are measurements such as:
  * altibase_memstat_max
  * altibase_memstat_alloc
  * altibase_tbs_total
  * altibase_tbs_usage
  * altibase_mem_tbl_usage
  * altibase_disk_tbl_usage
* ServiceThread stats - Statistics related to the service threads of Altibase. There is a measurement such as:
  * altibase_srv_thr
* Query stats - Query statistics that can affect Altibase's performance. There are measurements such as:
  * altibase_fullscan_query
  * altibase_utrans_query
  * altibase_long_run_query
* Lock stats - Statistics on the queries holding or waiting for locks. There are measurements such as:
  * altibase_lock_hold_info
  * altibase_lock_wait_info
  * altibase_tx_of_memory_view_scn
* File I/O stats - Statistics about I/O on individual disk files of Altibase. There are measurements such as:
  * altibase_file_io_reads
  * altibase_file_io_wait
  * altibase_file_io_writes
* Waits stats -Statistics on waits classified by wait event. There are measurements such as:
  * altibase_session_event
  * altibase_system_event
* Replication gap - Replication gaps by replication name. There is a measurement such as:
  * altibase_rep_gap
* Overall system stats - Overall statistics for the system. There is a measurement such as:
  * altibase_sysstat

## Tags

* All measurements have following tags.
  * dsn - the DSN from which the metrics are gathered
  * host - the host name from which the metrics are gathered
* Measurement 'altibase_tbs_usage_per_tbsname' additionally has the following tag.
  * TBS_NAME - tablespace name

## Example Output

When run with:

```sh
./telegraf --config telegraf.conf --input-filter altibase --test
```

The output is as follows:

```text
> altibase,dsn=Altiodbc,host=ux-ubuntu WORKING_TIME=3641304i 1713242733000000000
> altibase,dsn=Altiodbc,host=ux-ubuntu ALTIBASE_VERSION="7.1.0.9.2" 1713242733000000000
> altibase,dsn=Altiodbc,host=ux-ubuntu ARCHIVE_MODE=0i 1713242733000000000
> altibase,dsn=Altiodbc,host=ux-ubuntu SESSION_COUNT=6i 1713242733000000000
> altibase,dsn=Altiodbc,host=ux-ubuntu SESSION_COUNT=6i,STATEMENT_COUNT=4i 1713242733000000000
> altibase_srv_thr,dsn=Altiodbc,host=ux-ubuntu SOCKET(MULTIPLEXING)=8i 1713242733000000000
> altibase_srv_thr,dsn=Altiodbc,host=ux-ubuntu IPC=1i,SOCKET(MULTIPLEXING)=8i 1713242733000000000
> altibase_srv_thr,dsn=Altiodbc,host=ux-ubuntu IPC=1i,POLL=8i,SOCKET(MULTIPLEXING)=8i 1713242733000000000
> altibase_srv_thr,dsn=Altiodbc,host=ux-ubuntu EXECUTE=1i,IPC=1i,POLL=8i,SOCKET(MULTIPLEXING)=8i 1713242733000000000
> altibase_srv_thr,dsn=Altiodbc,host=ux-ubuntu EXECUTE=1i,IPC=1i,POLL=8i,SHARED=8i,SOCKET(MULTIPLEXING)=8i 1713242733000000000
> altibase_srv_thr,dsn=Altiodbc,host=ux-ubuntu DEDICATED=1i,EXECUTE=1i,IPC=1i,POLL=8i,SHARED=8i,SOCKET(MULTIPLEXING)=8i 1713242733000000000
> altibase_srv_thr,dsn=Altiodbc,host=ux-ubuntu DEDICATED=1i,EXECUTE=1i,IPC=1i,POLL=8i,SHARED=8i,SOCKET(MULTIPLEXING)=8i,TOTAL_COUNT=9i 1713242733000000000
> altibase,dsn=Altiodbc,host=ux-ubuntu MEMSTAT_ALLOC=3997000000i,MEMSTAT_MAX_TOTAL=388959798128i 1713242733000000000
...
```
