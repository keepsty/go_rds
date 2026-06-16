package mysql

import "fmt"

// InitializeSeedData 初始化 cluster 模块的元数据种子数据（幂等安全）
func InitializeSeedData() error {
	seedData := []struct {
		name string
		sql  string
	}{
		// ====== 产品线 ======
		{
			name: "production_line",
			sql: `INSERT IGNORE INTO production_line (id, name) VALUES
				(1, '电商平台'),
				(2, '支付系统'),
				(3, '数据平台')`,
		},
		// ====== 主机 ======
		{
			name: "hosts",
			sql: `INSERT IGNORE INTO hosts (id, name, ip, cpu_number, memory, disk_size, disk_type, cpu_platform, system_type, idc_location, host_status) VALUES
				(1, 'db-master-01', '192.168.1.101', 32, 128, 2048, 0, 'Intel Xeon Platinum', 'CentOS 7.9', '北京IDC-A', 1),
				(2, 'db-slave-01',  '192.168.1.102', 32, 128, 2048, 0, 'Intel Xeon Platinum', 'CentOS 7.9', '北京IDC-A', 1),
				(3, 'db-slave-02',  '192.168.1.103', 32, 128, 2048, 0, 'Intel Xeon Platinum', 'CentOS 7.9', '上海IDC-B', 1),
				(4, 'proxy-01',     '192.168.2.101', 16,  64,  512, 0, 'Intel Xeon Gold',   'CentOS 7.9', '北京IDC-A', 1),
				(5, 'proxy-02',     '192.168.2.102', 16,  64,  512, 0, 'Intel Xeon Gold',   'CentOS 7.9', '上海IDC-B', 1)`,
		},
		// ====== 集群 ======
		{
			name: "database_servicegroup",
			sql: `INSERT IGNORE INTO database_servicegroup (id, prod_id, name, dba_user, rd_owner, service_level, environment, ha_type, middleware, dns_name, vip, cluster_description, peak_time) VALUES
				(1, 1, '电商-core-db',   '张工', '李经理', 1, 0, 0, 0, 'core.ecommerce.xxx.com',  '10.0.1.100', '电商平台核心数据库集群',  '20:00-22:00'),
				(2, 1, '电商-order-db',  '张工', '王磊',   2, 0, 0, 0, 'order.ecommerce.xxx.com', '10.0.1.101', '电商平台订单数据库集群',  '20:00-23:00'),
				(3, 2, '支付-pay-db',    '赵工', '刘总',   1, 0, 0, 0, 'pay.payment.xxx.com',     '10.0.2.100', '支付系统核心数据库集群',  '00:00-01:00')`,
		},
		// ====== 实例 ======
		{
			name: "database_instance",
			sql: `INSERT IGNORE INTO database_instance (id, host_id, servicegroup_id, port, status, role, purpose, mysql_version) VALUES
				(1, 1, 1, 3306, 1, 1, '读写',   '8.0.32'),
				(2, 2, 1, 3306, 1, 0, '只读',   '8.0.32'),
				(3, 3, 1, 3306, 1, 0, '只读',   '8.0.32'),
				(4, 1, 2, 3307, 1, 1, '读写',   '8.0.32'),
				(5, 2, 2, 3307, 1, 0, '只读',   '8.0.32'),
				(6, 3, 3, 3306, 1, 1, '读写',   '8.0.28')`,
		},
		// ====== 数据库 ======
		{
			name: "database_database",
			sql: `INSERT IGNORE INTO database_database (id, servicegroup_id, name, rd_user, database_size, database_charset, db_rd_leader) VALUES
				(1, 1, 'user_db',    ',张三,李四,', 128.5000, 'utf8mb4', '张三'),
				(2, 1, 'product_db', ',王五,',       256.0000, 'utf8mb4', '王五'),
				(3, 1, 'order_db',   ',李四,赵六,',  512.7500, 'utf8mb4', '李四'),
				(4, 2, 'order_2024', ',王磊,',       1024.0000,'utf8mb4', '王磊'),
				(5, 2, 'order_2025', ',王磊,周七,',  64.2500,  'utf8mb4', '王磊'),
				(6, 3, 'pay_db',     ',刘总,',       512.0000, 'utf8mb4', '刘总')`,
		},
		// ====== ProxySQL ======
		{
			name: "database_proxysql",
			sql: `INSERT IGNORE INTO database_proxysql (id, host_id, servicegroup_id, admin_port, app_port, proxy_weight, proxy_version, rule_info) VALUES
				(1, 4, 1, 6032, 6033, 100, '2.5.5', '电商路由规则: 读写分离, 从库负载均衡'),
				(2, 5, 1, 6032, 6033, 90,  '2.5.5', '电商路由规则: 读写分离, 从库负载均衡'),
				(3, 4, 3, 6032, 6033, 100, '2.5.5', '支付路由规则: 读写分离')`,
		},
		// ====== ProxySQL 运行时 MySQL 服务器 ======
		{
			name: "database_proxysql_runtime_mysql_server",
			sql: `INSERT IGNORE INTO database_proxysql_runtime_mysql_server (proxysql_id, hostgroup_id, mysql_hostname, port, status, weight) VALUES
				(1, 10, '192.168.1.101', 3306, 'ONLINE', 1),
				(1, 11, '192.168.1.102', 3306, 'ONLINE', 1000),
				(1, 11, '192.168.1.103', 3306, 'ONLINE', 1000),
				(2, 10, '192.168.1.101', 3306, 'ONLINE', 1),
				(2, 11, '192.168.1.102', 3306, 'ONLINE', 1000),
				(2, 11, '192.168.1.103', 3306, 'ONLINE', 1000),
				(3, 10, '192.168.1.103', 3306, 'ONLINE', 1)`,
		},
	}

	for _, s := range seedData {
		if _, err := db.Exec(s.sql); err != nil {
			return fmt.Errorf("种子数据插入失败 [%s]: %w", s.name, err)
		}
	}

	return nil
}
