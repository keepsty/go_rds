package models

type HostsDetail struct {
	ID             int64  `json:"id" db:"id"`
	HostStatus     int64  `json:"host_status" db:"host_status"`
	DiskType       int64  `json:"disk_type" db:"disk_type"`
	DiskSize       int64  `json:"disk_size" db:"disk_size"`
	Memory         int64  `json:"memory" db:"memory"`
	RaidType       int64  `json:"raid_type" db:"raid_type"`
	CpuNumber      int64  `json:"cpu_number" db:"cpu_number"`
	Name           string `json:"name" db:"name"`
	IP             string `json:"ip" db:"ip"`
	CpuPlatform    string `json:"cpu_platform" db:"cpu_platform"`
	SystemType     string `json:"system_type" db:"system_type"`
	SystemFilename string `json:"system_filename" db:"system_filename"`
	IdcLocation    string `json:"idc_location" db:"idc_location"`
}
