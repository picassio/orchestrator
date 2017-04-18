package inst

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	test "github.com/openark/golib/tests"
)

var (
	i710k = InstanceKey{Hostname: "i710", Port: 3306}
	i720k = InstanceKey{Hostname: "i720", Port: 3306}
	i730k = InstanceKey{Hostname: "i730", Port: 3306}
)

var (
	spacesRegexp = regexp.MustCompile(`[ \t\n\r]+`)
)

func normalizeQuery(name string) string {
	name = strings.Replace(name, "`", "", -1)
	name = spacesRegexp.ReplaceAllString(name, " ")
	name = strings.TrimSpace(name)
	return name
}

func mkTestInstances() []*Instance {
	i710 := Instance{Key: i710k, ServerID: 710, ExecBinlogCoordinates: BinlogCoordinates{LogFile: "mysql.000007", LogPos: 10}}
	i720 := Instance{Key: i720k, ServerID: 720, ExecBinlogCoordinates: BinlogCoordinates{LogFile: "mysql.000007", LogPos: 20}}
	i730 := Instance{Key: i730k, ServerID: 730, ExecBinlogCoordinates: BinlogCoordinates{LogFile: "mysql.000007", LogPos: 30}}
	instances := []*Instance{&i710, &i720, &i730}
	for _, instance := range instances {
		instance.Version = "5.6.7"
		instance.VersionComment = "MySQL"
		instance.Binlog_format = "STATEMENT"
		instance.BinlogRowImage = "FULL"
	}
	return instances
}

func TestMkInsertOdku(t *testing.T) {
	instances := mkTestInstances()

	sql, args := mkInsertOdkuForInstances(nil, true, true)
	test.S(t).ExpectEquals(sql, "")
	test.S(t).ExpectEquals(len(args), 0)

	// one instance
	s1 := `INSERT ignore INTO database_instance
                (hostname, port, last_checked, last_attempted_check, uptime, server_id, server_uuid, version, major_version, version_comment, binlog_server, read_only, binlog_format, binlog_row_image, log_bin, log_slave_updates, binary_log_file, binary_log_pos, master_host, master_port, slave_sql_running, slave_io_running, has_replication_filters, supports_oracle_gtid, oracle_gtid, executed_gtid_set, gtid_purged, mariadb_gtid, pseudo_gtid, master_log_file, read_master_log_pos, relay_master_log_file, exec_master_log_pos, relay_log_file, relay_log_pos, last_sql_error, last_io_error, seconds_behind_master, slave_lag_seconds, sql_delay, num_slave_hosts, slave_hosts, cluster_name, suggested_cluster_alias, data_center, physical_environment, replication_depth, is_co_master, replication_credentials_available, has_replication_credentials, allow_tls, semi_sync_enforced, instance_alias, last_discovery_latency, last_seen)
        VALUES
                (?, ?, NOW(), NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
        ON DUPLICATE KEY UPDATE
                hostname=VALUES(hostname), port=VALUES(port), last_checked=VALUES(last_checked), last_attempted_check=VALUES(last_attempted_check), uptime=VALUES(uptime), server_id=VALUES(server_id), server_uuid=VALUES(server_uuid), version=VALUES(version), major_version=VALUES(major_version), version_comment=VALUES(version_comment), binlog_server=VALUES(binlog_server), read_only=VALUES(read_only), binlog_format=VALUES(binlog_format), binlog_row_image=VALUES(binlog_row_image), log_bin=VALUES(log_bin), log_slave_updates=VALUES(log_slave_updates), binary_log_file=VALUES(binary_log_file), binary_log_pos=VALUES(binary_log_pos), master_host=VALUES(master_host), master_port=VALUES(master_port), slave_sql_running=VALUES(slave_sql_running), slave_io_running=VALUES(slave_io_running), has_replication_filters=VALUES(has_replication_filters), supports_oracle_gtid=VALUES(supports_oracle_gtid), oracle_gtid=VALUES(oracle_gtid), executed_gtid_set=VALUES(executed_gtid_set), gtid_purged=VALUES(gtid_purged), mariadb_gtid=VALUES(mariadb_gtid), pseudo_gtid=VALUES(pseudo_gtid), master_log_file=VALUES(master_log_file), read_master_log_pos=VALUES(read_master_log_pos), relay_master_log_file=VALUES(relay_master_log_file), exec_master_log_pos=VALUES(exec_master_log_pos), relay_log_file=VALUES(relay_log_file), relay_log_pos=VALUES(relay_log_pos), last_sql_error=VALUES(last_sql_error), last_io_error=VALUES(last_io_error), seconds_behind_master=VALUES(seconds_behind_master), slave_lag_seconds=VALUES(slave_lag_seconds), sql_delay=VALUES(sql_delay), num_slave_hosts=VALUES(num_slave_hosts), slave_hosts=VALUES(slave_hosts), cluster_name=VALUES(cluster_name), suggested_cluster_alias=VALUES(suggested_cluster_alias), data_center=VALUES(data_center), physical_environment=VALUES(physical_environment), replication_depth=VALUES(replication_depth), is_co_master=VALUES(is_co_master), replication_credentials_available=VALUES(replication_credentials_available), has_replication_credentials=VALUES(has_replication_credentials), allow_tls=VALUES(allow_tls), semi_sync_enforced=VALUES(semi_sync_enforced), instance_alias=VALUES(instance_alias), last_discovery_latency=VALUES(last_discovery_latency), last_seen=VALUES(last_seen)
        `
	a1 := `i710, 3306, 0, 710, , 5.6.7, 5.6, MySQL, false, false, STATEMENT, FULL, false, false, , 0, , 0, false, false, false, false, false, , , false, false, , 0, mysql.000007, 10, , 0, , , {0 false}, {0 false}, 0, 0, [], , , , , 0, false, false, false, false, false, , 0, `

	sql1, args1 := mkInsertOdkuForInstances(instances[:1], false, true)

	test.S(t).ExpectEquals(normalizeQuery(sql1), normalizeQuery(s1))
	test.S(t).ExpectEquals(strings.TrimSpace(fmtArgs(args1)), strings.TrimSpace(a1))

	// three instances
	s3 := `INSERT  INTO database_instance
                (hostname, port, last_checked, last_attempted_check, uptime, server_id, server_uuid, version, major_version, version_comment, binlog_server, read_only, binlog_format, binlog_row_image, log_bin, log_slave_updates, binary_log_file, binary_log_pos, master_host, master_port, slave_sql_running, slave_io_running, has_replication_filters, supports_oracle_gtid, oracle_gtid, executed_gtid_set, gtid_purged, mariadb_gtid, pseudo_gtid, master_log_file, read_master_log_pos, relay_master_log_file, exec_master_log_pos, relay_log_file, relay_log_pos, last_sql_error, last_io_error, seconds_behind_master, slave_lag_seconds, sql_delay, num_slave_hosts, slave_hosts, cluster_name, suggested_cluster_alias, data_center, physical_environment, replication_depth, is_co_master, replication_credentials_available, has_replication_credentials, allow_tls, semi_sync_enforced, instance_alias, last_discovery_latency, last_seen)
        VALUES
                (?, ?, NOW(), NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW()),
                (?, ?, NOW(), NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW()),
                (?, ?, NOW(), NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
        ON DUPLICATE KEY UPDATE
                hostname=VALUES(hostname), port=VALUES(port), last_checked=VALUES(last_checked), last_attempted_check=VALUES(last_attempted_check), uptime=VALUES(uptime), server_id=VALUES(server_id), server_uuid=VALUES(server_uuid), version=VALUES(version), major_version=VALUES(major_version), version_comment=VALUES(version_comment), binlog_server=VALUES(binlog_server), read_only=VALUES(read_only), binlog_format=VALUES(binlog_format), binlog_row_image=VALUES(binlog_row_image), log_bin=VALUES(log_bin), log_slave_updates=VALUES(log_slave_updates), binary_log_file=VALUES(binary_log_file), binary_log_pos=VALUES(binary_log_pos), master_host=VALUES(master_host), master_port=VALUES(master_port), slave_sql_running=VALUES(slave_sql_running), slave_io_running=VALUES(slave_io_running), has_replication_filters=VALUES(has_replication_filters), supports_oracle_gtid=VALUES(supports_oracle_gtid), oracle_gtid=VALUES(oracle_gtid), executed_gtid_set=VALUES(executed_gtid_set), gtid_purged=VALUES(gtid_purged), mariadb_gtid=VALUES(mariadb_gtid), pseudo_gtid=VALUES(pseudo_gtid), master_log_file=VALUES(master_log_file), read_master_log_pos=VALUES(read_master_log_pos), relay_master_log_file=VALUES(relay_master_log_file), exec_master_log_pos=VALUES(exec_master_log_pos), relay_log_file=VALUES(relay_log_file), relay_log_pos=VALUES(relay_log_pos), last_sql_error=VALUES(last_sql_error), last_io_error=VALUES(last_io_error), seconds_behind_master=VALUES(seconds_behind_master), slave_lag_seconds=VALUES(slave_lag_seconds), sql_delay=VALUES(sql_delay), num_slave_hosts=VALUES(num_slave_hosts), slave_hosts=VALUES(slave_hosts), cluster_name=VALUES(cluster_name), suggested_cluster_alias=VALUES(suggested_cluster_alias), data_center=VALUES(data_center), physical_environment=VALUES(physical_environment), replication_depth=VALUES(replication_depth), is_co_master=VALUES(is_co_master), replication_credentials_available=VALUES(replication_credentials_available), has_replication_credentials=VALUES(has_replication_credentials), allow_tls=VALUES(allow_tls), semi_sync_enforced=VALUES(semi_sync_enforced), instance_alias=VALUES(instance_alias), last_discovery_latency=VALUES(last_discovery_latency), last_seen=VALUES(last_seen)
        `
	a3 := `i710, 3306, 0, 710, , 5.6.7, 5.6, MySQL, false, false, STATEMENT, FULL, false, false, , 0, , 0, false, false, false, false, false, , , false, false, , 0, mysql.000007, 10, , 0, , , {0 false}, {0 false}, 0, 0, [], , , , , 0, false, false, false, false, false, , 0, i720, 3306, 0, 720, , 5.6.7, 5.6, MySQL, false, false, STATEMENT, FULL, false, false, , 0, , 0, false, false, false, false, false, , , false, false, , 0, mysql.000007, 20, , 0, , , {0 false}, {0 false}, 0, 0, [], , , , , 0, false, false, false, false, false, , 0, i730, 3306, 0, 730, , 5.6.7, 5.6, MySQL, false, false, STATEMENT, FULL, false, false, , 0, , 0, false, false, false, false, false, , , false, false, , 0, mysql.000007, 30, , 0, , , {0 false}, {0 false}, 0, 0, [], , , , , 0, false, false, false, false, false, , 0, `

	sql3, args3 := mkInsertOdkuForInstances(instances[:3], true, true)

	test.S(t).ExpectEquals(normalizeQuery(sql3), normalizeQuery(s3))
	test.S(t).ExpectEquals(strings.TrimSpace(fmtArgs(args3)), strings.TrimSpace(a3))
}

func fmtArgs(args []interface{}) string {
	b := &bytes.Buffer{}
	for _, a := range args {
		fmt.Fprint(b, a)
		fmt.Fprint(b, ", ")
	}
	return b.String()
}
