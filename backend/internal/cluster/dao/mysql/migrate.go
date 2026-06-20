package mysql

import "fmt"

// AutoMigrate 根据模型结构自动创建集群管理模块所需的数据库表
func AutoMigrate() error {
	tables := []struct {
		name string
		sql  string
	}{
		{
			name: "production_line",
			sql: `CREATE TABLE IF NOT EXISTS production_line (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '产品线ID',
				name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '产品线名称',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品线表'`,
		},
		{
			name: "database_servicegroup",
			sql: `CREATE TABLE IF NOT EXISTS database_servicegroup (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '集群ID',
				prod_id BIGINT NOT NULL DEFAULT 0 COMMENT '产品线ID',
				name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '集群名称',
				dba_user VARCHAR(64) NOT NULL DEFAULT '' COMMENT '负责DBA',
				rd_owner VARCHAR(64) NOT NULL DEFAULT '' COMMENT '业务负责人',
				service_level TINYINT NOT NULL DEFAULT 0 COMMENT '集群等级',
				environment TINYINT NOT NULL DEFAULT 0 COMMENT '环境 0:prod 1:rc 2:k8s 3:press',
				ha_type TINYINT NOT NULL DEFAULT 0 COMMENT 'HA类型 0:mha 1:orc 2:mgr',
				middleware TINYINT NOT NULL DEFAULT 0 COMMENT '中间件 0:proxysql 1:zebra 2:mgw',
				dns_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'DNS名称',
				vip VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'VIP地址',
				cluster_description VARCHAR(512) NOT NULL DEFAULT '' COMMENT '集群描述',
				peak_time VARCHAR(64) NOT NULL DEFAULT '' COMMENT '高峰期',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
				INDEX idx_prod_id (prod_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库服务集群表'`,
		},
		{
			name: "hosts",
			sql: `CREATE TABLE IF NOT EXISTS hosts (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主机ID',
				name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '主机名',
				ip VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'IP地址',
				cpu_number INT NOT NULL DEFAULT 0 COMMENT 'CPU核数',
				memory INT NOT NULL DEFAULT 0 COMMENT '内存大小(G)',
				disk_size INT NOT NULL DEFAULT 0 COMMENT '磁盘大小(G)',
				disk_type INT NOT NULL DEFAULT 0 COMMENT '磁盘类型',
				raid_type INT NOT NULL DEFAULT 0 COMMENT 'RAID类型',
				cpu_platform VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'CPU平台',
				system_type VARCHAR(64) NOT NULL DEFAULT '' COMMENT '系统类型',
				system_filename VARCHAR(64) NOT NULL DEFAULT '' COMMENT '系统镜像',
				idc_location VARCHAR(128) NOT NULL DEFAULT '' COMMENT '机房位置',
				host_status INT NOT NULL DEFAULT 0 COMMENT '主机状态',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机信息表'`,
		},
		{
			name: "database_instance",
			sql: `CREATE TABLE IF NOT EXISTS database_instance (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '实例ID',
				host_id BIGINT NOT NULL DEFAULT 0 COMMENT '主机ID',
				servicegroup_id BIGINT NOT NULL DEFAULT 0 COMMENT '集群ID',
				port INT NOT NULL DEFAULT 3306 COMMENT '端口',
				status INT NOT NULL DEFAULT 0 COMMENT '状态 0:offline 1:online 2:mantain 3:problem',
				role INT NOT NULL DEFAULT 0 COMMENT '角色 0:从库 1:主库',
				purpose VARCHAR(64) NOT NULL DEFAULT '' COMMENT '实例用途',
				mysql_version VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'MySQL版本',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
				INDEX idx_host_id (host_id),
				INDEX idx_sg_id (servicegroup_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库实例表'`,
		},
		{
			name: "database_database",
			sql: `CREATE TABLE IF NOT EXISTS database_database (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '数据库ID',
				servicegroup_id BIGINT NOT NULL DEFAULT 0 COMMENT '集群ID',
				name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '数据库名',
				rd_user VARCHAR(256) NOT NULL DEFAULT '' COMMENT '业务RD(多值逗号分隔)',
				database_size DECIMAL(20,4) NOT NULL DEFAULT 0 COMMENT '数据空间大小',
				database_charset VARCHAR(32) NOT NULL DEFAULT '' COMMENT '字符集',
				db_rd_leader VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'RD负责人',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
				INDEX idx_sg_id (servicegroup_id),
				INDEX idx_name (name)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库信息表'`,
		},
		{
			name: "database_tables",
			sql: `CREATE TABLE IF NOT EXISTS database_tables (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '表ID',
				db_id BIGINT NOT NULL DEFAULT 0 COMMENT '数据库ID',
				table_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '表名',
				table_size DECIMAL(20,4) NOT NULL DEFAULT 0 COMMENT '表大小(G)',
				free_size DECIMAL(20,4) NOT NULL DEFAULT 0 COMMENT '空洞大小(G)',
				table_rows BIGINT NOT NULL DEFAULT 0 COMMENT '行数',
				auto_increase BIGINT NOT NULL DEFAULT 0 COMMENT '自增值',
				table_collation VARCHAR(64) NOT NULL DEFAULT '' COMMENT '字符集校验规则',
				tb_schema TEXT COMMENT '建表语句',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
				INDEX idx_db_id (db_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库表信息表'`,
		},
		{
			name: "database_proxysql",
			sql: `CREATE TABLE IF NOT EXISTS database_proxysql (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'ProxySQL ID',
				host_id BIGINT NOT NULL DEFAULT 0 COMMENT '主机ID',
				servicegroup_id BIGINT NOT NULL DEFAULT 0 COMMENT '集群ID',
				admin_port INT NOT NULL DEFAULT 6032 COMMENT '管理端口',
				app_port INT NOT NULL DEFAULT 6033 COMMENT '应用端口',
				proxy_weight INT NOT NULL DEFAULT 0 COMMENT '权重',
				proxy_version VARCHAR(32) NOT NULL DEFAULT '' COMMENT 'ProxySQL版本',
				rule_info TEXT COMMENT '路由规则',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
				INDEX idx_sg_id (servicegroup_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ProxySQL节点表'`,
		},
		{
			name: "database_proxysql_details",
			sql: `CREATE TABLE IF NOT EXISTS database_proxysql_details (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '规则ID',
				proxysql_id BIGINT NOT NULL DEFAULT 0 COMMENT 'ProxySQL ID',
				username VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
				flag_out INT NOT NULL DEFAULT 0 COMMENT 'flag_out',
				match_pattern VARCHAR(256) NOT NULL DEFAULT '' COMMENT '匹配规则',
				destination_hostgroup INT NOT NULL DEFAULT 0 COMMENT '目标hostgroup',
				hostname VARCHAR(128) NOT NULL DEFAULT '' COMMENT '主机名',
				port INT NOT NULL DEFAULT 3306 COMMENT '端口',
				status VARCHAR(32) NOT NULL DEFAULT 'ONLINE' COMMENT '状态',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
				INDEX idx_proxysql_id (proxysql_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ProxySQL路由规则表'`,
		},
		{
			name: "database_proxysql_runtime_mysql_server",
			sql: `CREATE TABLE IF NOT EXISTS database_proxysql_runtime_mysql_server (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
				proxysql_id BIGINT NOT NULL DEFAULT 0 COMMENT 'ProxySQL ID',
				hostgroup_id INT NOT NULL DEFAULT 0 COMMENT 'Hostgroup ID',
				mysql_hostname VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'MySQL主机名',
				port INT NOT NULL DEFAULT 3306 COMMENT '端口',
				status VARCHAR(32) NOT NULL DEFAULT 'ONLINE' COMMENT '状态',
				weight INT NOT NULL DEFAULT 1 COMMENT '权重',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
				INDEX idx_proxysql_id (proxysql_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ProxySQL运行时MySQL服务器表'`,
		},
		{
			name: "database_dms_query_record",
			sql: `CREATE TABLE IF NOT EXISTS database_dms_query_record (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '记录ID',
				username VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户名',
				host VARCHAR(64) NOT NULL DEFAULT '' COMMENT '目标主机',
				db_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '数据库名',
				tables VARCHAR(512) NOT NULL DEFAULT '' COMMENT '相关表',
				execute_query TEXT COMMENT '执行的SQL',
				query_consume_time VARCHAR(32) NOT NULL DEFAULT '' COMMENT '查询耗时(秒)',
				query_status VARCHAR(32) NOT NULL DEFAULT '' COMMENT '查询状态',
				affected_rows INT NOT NULL DEFAULT 0 COMMENT '影响行数',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				INDEX idx_username (username)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='DMS查询记录表'`,
		},
		{
			name: "salt_deployment",
			sql: `CREATE TABLE IF NOT EXISTS salt_deployment (
				id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '部署ID',
				env VARCHAR(32) NOT NULL DEFAULT '' COMMENT '部署环境',
				servicegroup_id BIGINT NOT NULL DEFAULT 0 COMMENT '集群ID',
				deploy_type VARCHAR(32) NOT NULL DEFAULT '' COMMENT '部署类型 mysql/proxysql/mgr',
				mysql_version VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'MySQL版本',
				config_json TEXT COMMENT '部署配置JSON',
				status VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT '状态 pending/running/success/failed',
				result TEXT COMMENT '部署结果',
				created_by VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建人',
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
				INDEX idx_env (env),
				INDEX idx_sg_id (servicegroup_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='SaltStack部署记录表'`,
		},
		{
			name: "database_backup_config",
			sql: `CREATE TABLE IF NOT EXISTS database_backup_config (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(128) NOT NULL,
				db_type VARCHAR(32) NOT NULL,
				instance_id VARCHAR(128) NOT NULL,
				backup_type VARCHAR(32) NOT NULL DEFAULT 'full',
				schedule_cron VARCHAR(64) DEFAULT '',
				retention_days INT NOT NULL DEFAULT 7,
				storage_path VARCHAR(256) DEFAULT '',
				status VARCHAR(32) NOT NULL DEFAULT 'enabled',
				remark TEXT,
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		},
		{
			name: "database_backup_task",
			sql: `CREATE TABLE IF NOT EXISTS database_backup_task (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				config_id BIGINT NOT NULL,
				name VARCHAR(128) NOT NULL,
				db_type VARCHAR(32) NOT NULL,
				backup_type VARCHAR(32) NOT NULL DEFAULT 'full',
				instance_id VARCHAR(128) NOT NULL,
				status VARCHAR(32) NOT NULL DEFAULT 'pending',
				started_at DATETIME,
				finished_at DATETIME,
				created_by VARCHAR(64) NOT NULL,
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		},
		{
			name: "database_backup_record",
			sql: `CREATE TABLE IF NOT EXISTS database_backup_record (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				task_id BIGINT NOT NULL,
				file_name VARCHAR(256) DEFAULT '',
				file_size BIGINT DEFAULT 0,
				file_path VARCHAR(512) DEFAULT '',
				status VARCHAR(32) NOT NULL DEFAULT 'running',
				output TEXT,
				started_at DATETIME,
				finished_at DATETIME,
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		},
		{
			name: "database_backup_template",
			sql: `CREATE TABLE IF NOT EXISTS database_backup_template (
				id BIGINT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(128) NOT NULL,
				db_type VARCHAR(32) NOT NULL,
				description TEXT,
				config_schema TEXT NOT NULL,
				default_config TEXT NOT NULL,
				create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
				update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		},
	}

	for _, t := range tables {
		if _, err := db.Exec(t.sql); err != nil {
			return fmt.Errorf("自动建表失败 [%s]: %w", t.name, err)
		}
	}

	return nil
}
