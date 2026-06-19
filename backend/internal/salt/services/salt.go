package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/keepsty/go_rds/internal/salt/models"
	"github.com/keepsty/go_rds/internal/salt/salt"

	"github.com/keepsty/go_rds/internal/global"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// ---------- Salt 业务封装 ----------

// SaltService 封装业务层 Salt 调用逻辑
type SaltService struct {
	client *salt.Client
}

// NewSaltService 创建 Salt 业务服务，使用全局配置
func NewSaltService() *SaltService {
	return NewSaltServiceWithClient(
		global.App.Config.Salt.URL,
		global.App.Config.Salt.User,
		global.App.Config.Salt.Password,
		"pam",
	)
}

// NewSaltServiceWithClient 使用自定义参数创建 Salt 业务服务
func NewSaltServiceWithClient(baseURL, username, password, eauth string) *SaltService {
	return &SaltService{
		client: salt.NewClient(baseURL, username, password, eauth),
	}
}

// ---------- 命令执行封装 ----------

// RunCommand 下发 cmd.run 同步命令
func (s *SaltService) RunCommand(host, command string, enableSSH bool) (*salt.CmdResult, error) {
	client := "local"
	if enableSSH {
		client = "ssh"
	}
	resp, err := s.client.CmdRun("cmd.run", host, command, client)
	if err != nil {
		return nil, fmt.Errorf("执行命令失败: %w", err)
	}
	return parseCmdResult(resp, host)
}

// RunState 执行 state.sls，下发 Salt 状态
func (s *SaltService) RunState(host, stateFile string, enableSSH bool) (*salt.CmdResult, error) {
	client := "local"
	if enableSSH {
		client = "ssh"
	}
	resp, err := s.client.StateSLS(host, stateFile, client)
	if err != nil {
		return nil, fmt.Errorf("执行 state.sls 失败: %w", err)
	}
	return parseCmdResult(resp, host)
}

// RunScript 下发脚本执行
func (s *SaltService) RunScript(host, saltPath string, enableSSH bool) (*salt.CmdResult, error) {
	client := "local"
	if enableSSH {
		client = "ssh"
	}
	resp, err := s.client.CmdScript(host, saltPath, client)
	if err != nil {
		return nil, fmt.Errorf("下发脚本失败: %w", err)
	}
	return parseCmdResult(resp, host)
}

// RunAsync 异步执行命令，返回 JID 用于后续查询
func (s *SaltService) RunAsync(fun, tgt, arg string) (string, error) {
	resp, err := s.client.Run(fun, tgt, arg)
	if err != nil {
		return "", fmt.Errorf("异步执行失败: %w", err)
	}
	if len(resp.Return) == 0 {
		return "", fmt.Errorf("异步执行返回为空")
	}
	ret := resp.Return[0]
	jid, ok := ret["jid"].(string)
	if !ok {
		return "", fmt.Errorf("异步执行响应中未找到 jid: %v", ret)
	}
	return jid, nil
}

// LoadJobResult 查询异步作业的执行结果
func (s *SaltService) LoadJobResult(jid string) (*salt.CmdResult, error) {
	resp, err := s.client.LoadJob(jid)
	if err != nil {
		return nil, fmt.Errorf("查询作业失败: %w", err)
	}
	result := &salt.CmdResult{
		Success: true,
		JID:     jid,
		Data:    make(map[string]interface{}),
	}
	if len(resp.Return) > 0 {
		for k, v := range resp.Return[0] {
			result.Data[k] = v
		}
		result.Detail = resp.Return[0]
	}
	return result, nil
}

// ---------- 主机信息采集 ----------

// GetHostGrains 获取主机 grains 信息
func (s *SaltService) GetHostGrains(host string) (*salt.HostInfo, error) {
	resp, err := s.client.GrainsItems(host)
	if err != nil {
		return nil, fmt.Errorf("获取 grains 失败: %w", err)
	}
	if len(resp.Return) == 0 {
		return nil, fmt.Errorf("grains 返回为空")
	}
	ret := resp.Return[0]
	hostData, ok := ret[host].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("grains 响应中未找到主机 %s 的数据", host)
	}
	info := &salt.HostInfo{
		ID:          host,
		OS:          getStringField(hostData, "os"),
		OSRelease:   getStringField(hostData, "osrelease"),
		CPUModel:    getStringField(hostData, "cpu_model"),
		Kernel:      getStringField(hostData, "kernel"),
		SaltVersion: getStringField(hostData, "saltversion"),
	}
	info.IPv4 = getIPString(hostData)
	if nc, ok := hostData["num_cpus"].(int); ok {
		info.NumCPUs = nc
	}
	if mt, ok := hostData["mem_total"].(int); ok {
		info.MemTotal = mt
	}
	return info, nil
}

// ---------- 主机发现 ----------

// MinionList minion 列表
type MinionList struct {
	Minions    []string // 已接受
	Unaccepted []string // 未接受
	Rejected   []string // 已拒绝
	Denied     []string // 已拒绝
}

// ListMinions 列出所有已认证的 Minion
func (s *SaltService) ListMinions() (*MinionList, error) {
	resp, err := s.client.KeyListAll()
	if err != nil {
		return nil, fmt.Errorf("获取 minion 列表失败: %w", err)
	}
	list := &MinionList{}
	if len(resp.Return) > 0 {
		if data, ok := resp.Return[0]["data"].(map[string]interface{}); ok {
			if local, ok := data["local"].(map[string]interface{}); ok {
				list.Minions = toStringSlice(local["minions"])
				list.Unaccepted = toStringSlice(local["minions_pre"])
				list.Rejected = toStringSlice(local["minions_rejected"])
				list.Denied = toStringSlice(local["minions_denied"])
			}
		}
		if list.Minions == nil {
			list.Minions = toStringSlice(resp.Return[0]["minions"])
			list.Unaccepted = toStringSlice(resp.Return[0]["minions_pre"])
			list.Rejected = toStringSlice(resp.Return[0]["minions_rejected"])
		}
	}
	return list, nil
}

// ListOnlineMinions 获取在线 minion 列表
func (s *SaltService) ListOnlineMinions() ([]string, error) {
	resp, err := s.client.ManageStatus()
	if err != nil {
		return nil, fmt.Errorf("获取 minion 状态失败: %w", err)
	}
	if len(resp.Return) > 0 {
		if up, ok := resp.Return[0]["up"].([]interface{}); ok {
			return toStringSlice(up), nil
		}
	}
	return []string{}, nil
}

// ---------- 通用执行 ----------

// RawExecute 通用 Salt 函数执行
func (s *SaltService) RawExecute(fun, tgt, arg, client string) (*salt.CmdResult, error) {
	var resp *salt.SaltResponse
	var err error
	switch client {
	case "wheel":
		resp, err = s.client.WheelRun(fun, arg)
	case "runner":
		resp, err = s.client.RunnerRun(fun, arg)
	case "ssh":
		resp, err = s.client.SSHExecute(fun, tgt, arg)
	case "local_async":
		resp, err = s.client.Run(fun, tgt, arg)
	default:
		resp, err = s.client.CmdRun(fun, tgt, arg, client)
	}
	if err != nil {
		return nil, fmt.Errorf("RawExecute 失败: %w", err)
	}
	return parseCmdResult(resp, tgt)
}

// ---------- 模板管理 ----------

// GetTemplates 从数据库读取模板列表
func GetTemplates() ([]models.InsightSaltTemplates, error) {
	var templates []models.InsightSaltTemplates
	result := global.App.DB.Model(&models.InsightSaltTemplates{}).Find(&templates)
	return templates, result.Error
}

// DeployFromTemplate 根据模板名和配置执行部署，同步等待结果
func DeployFromTemplate(name string, config map[string]interface{}, targets []string) (interface{}, error) {
	var tmpl models.InsightSaltTemplates
	if err := global.App.DB.Where("name = ?", name).First(&tmpl).Error; err != nil {
		return nil, fmt.Errorf("模板 '%s' 不存在", name)
	}

	svc := NewSaltService()

	switch name {
	case "cmd-run":
		cmd, _ := config["command"].(string)
		if cmd == "" {
			return nil, fmt.Errorf("命令不能为空")
		}
		result, err := svc.RunCommand(targets[0], cmd, false)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": result.Success,
			"message": result.Error,
			"detail":  result.Detail,
		}, nil

	case "state-apply":
		stateFile, _ := config["state_file"].(string)
		if stateFile == "" {
			return nil, fmt.Errorf("state 文件不能为空")
		}
		result, err := svc.RunState(targets[0], stateFile, false)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": result.Success,
			"message": result.Error,
			"detail":  result.Detail,
		}, nil

	case "mysql-deploy":
		port := int64(3306)
		if p, ok := config["port"].(float64); ok {
			port = int64(p)
		}
		version := "8.0"
		if v, ok := config["version"].(string); ok {
			version = v
		}
		hp := []*models.SaltMysqlHostPost{
			{Port: port, Host: targets[0], Version: version},
		}
		si := &models.SaltMysqlServerInfo{HostPort: hp}
		data, err := InstallMySQLHandler("prod", si, &global.App.Config.Salt)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": true,
			"detail":  data,
		}, nil

	default:
		return nil, fmt.Errorf("不支持模板类型: %s", name)
	}
}

// CreateTemplate 创建模板
func CreateTemplate(name, title, description, fieldsSchema, defaults string) error {
	tmpl := models.InsightSaltTemplates{
		Name:         name,
		Title:        title,
		Description:  description,
		FieldsSchema: datatypes.JSON(fieldsSchema),
		Defaults:     datatypes.JSON(defaults),
	}
	return global.App.DB.Model(&models.InsightSaltTemplates{}).Create(&tmpl).Error
}

// UpdateTemplate 更新模板
func UpdateTemplate(id uint64, title, description, fieldsSchema, defaults string) error {
	updates := map[string]interface{}{
		"title":        title,
		"description":  description,
		"fields_schema": datatypes.JSON(fieldsSchema),
		"defaults":     datatypes.JSON(defaults),
	}
	return global.App.DB.Model(&models.InsightSaltTemplates{}).Where("id=?", id).Updates(updates).Error
}

// DeleteTemplate 删除模板
func DeleteTemplate(id uint64) error {
	return global.App.DB.Where("id=?", id).Delete(&models.InsightSaltTemplates{}).Error
}

// ---------- 主机配置管理 ----------

func GetHostConfigs() ([]models.SaltHostConfig, error) {
	var list []models.SaltHostConfig
	err := global.App.DB.Model(&models.SaltHostConfig{}).Find(&list).Error
	return list, err
}

func CreateHostConfig(name, hosts, description string) error {
	record := models.SaltHostConfig{
		Name: name, Hosts: datatypes.JSON(hosts), Description: description,
	}
	return global.App.DB.Model(&models.SaltHostConfig{}).Create(&record).Error
}

func UpdateHostConfig(id uint64, name, hosts, description string) error {
	return global.App.DB.Model(&models.SaltHostConfig{}).Where("id=?", id).Updates(map[string]interface{}{
		"name": name, "hosts": datatypes.JSON(hosts), "description": description,
	}).Error
}

func DeleteHostConfig(id uint64) error {
	return global.App.DB.Where("id=?", id).Delete(&models.SaltHostConfig{}).Error
}

// ---------- 部署任务管理 ----------

func GetTasks(status string) ([]models.SaltTask, error) {
	var list []models.SaltTask
	tx := global.App.DB.Model(&models.SaltTask{}).Order("id desc")
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	err := tx.Find(&list).Error
	return list, err
}

func CreateTask(name, templateName string, hostConfigID uint64, configParams map[string]interface{}, createdBy string) (uint64, error) {
	paramsJSON, _ := json.Marshal(configParams)
	record := models.SaltTask{
		Name: name, TemplateName: templateName,
		HostConfigID: hostConfigID, Status: "pending",
		ConfigParams: datatypes.JSON(paramsJSON), CreatedBy: createdBy,
	}
	err := global.App.DB.Model(&models.SaltTask{}).Create(&record).Error
	return record.ID, err
}

func ApproveTask(id uint64, action, username string) error {
	if action == "approve" {
		return global.App.DB.Model(&models.SaltTask{}).Where("id=? AND status=?", id, "pending").
			Updates(map[string]interface{}{"status": "approved", "approved_by": username}).Error
	}
	return global.App.DB.Model(&models.SaltTask{}).Where("id=? AND status=?", id, "pending").
		Update("status", "rejected").Error
}

func RunTask(id uint64) error {
	var task models.SaltTask
	if err := global.App.DB.First(&task, id).Error; err != nil {
		return fmt.Errorf("任务不存在")
	}
	if task.Status != "approved" {
		return fmt.Errorf("任务未审批，无法执行")
	}

	global.App.DB.Model(&task).Updates(map[string]interface{}{
		"status": "running", "started_at": time.Now(),
	})

	var config map[string]interface{}
	if task.ConfigParams != nil {
		json.Unmarshal(task.ConfigParams, &config)
	}

	result, err := DeployFromTemplate(task.TemplateName, config, []string{})
	output, _ := json.Marshal(result)
	now := time.Now()

	if err != nil {
		global.App.DB.Model(&task).Updates(map[string]interface{}{
			"status": "failed", "run_output": datatypes.JSON(output), "finished_at": now,
		})
		return err
	}

	global.App.DB.Model(&task).Updates(map[string]interface{}{
		"status": "success", "run_output": datatypes.JSON(output), "finished_at": now,
	})
	return nil
}

// ---------- GoRDS 风格的 Service 结构体 ----------

// SaltCmdService 执行 Salt 命令（GoRDS 模块模式）
type SaltCmdService struct {
	C       *gin.Context
	Host    string
	Command string
}

func (s *SaltCmdService) Run() (interface{}, error) {
	svc := NewSaltService()
	return svc.RunCommand(s.Host, s.Command, false)
}

// ---------- 辅助函数 ----------

var systemKeys = map[string]bool{
	"jid": true, "result": true, "return": true,
	"success": true, "data": true, "output": true,
	"minions": true, "failure": true,
}

func parseCmdResult(resp *salt.SaltResponse, host string) (*salt.CmdResult, error) {
	result := &salt.CmdResult{
		Success: true,
		Data:    make(map[string]interface{}),
	}
	if len(resp.Return) == 0 {
		result.Success = false
		result.Error = "响应中 return 为空"
		return result, nil
	}
	ret := resp.Return[0]
	result.Detail = ret
	for k, v := range ret {
		result.Data[k] = v
	}
	if jid, ok := ret["jid"].(string); ok {
		result.JID = jid
	}
	nodeNames := extractNodeNames(ret)
	if len(nodeNames) > 0 {
		result.NodeNames = nodeNames
		if len(nodeNames) == 1 {
			result.NodeName = nodeNames[0]
		}
	} else if host != "" && host != "*" {
		result.NodeName = host
		result.NodeNames = []string{host}
	}
	if errStr := checkResultError(ret); errStr != "" {
		result.Success = false
		result.Error = errStr
	}
	return result, nil
}

func extractNodeNames(ret map[string]interface{}) []string {
	var names []string
	for k, v := range ret {
		if systemKeys[k] {
			continue
		}
		if nodeData, ok := v.(map[string]interface{}); ok {
			if _, hasSuccess := nodeData["success"]; hasSuccess {
				names = append(names, k)
				continue
			}
			if _, hasReturn := nodeData["return"]; hasReturn {
				names = append(names, k)
			}
		}
	}
	return names
}

func checkResultError(ret map[string]interface{}) string {
	for _, v := range ret {
		switch val := v.(type) {
		case string:
			if strings.Contains(val, "Failed to execute") {
				return val
			}
		case map[string]interface{}:
			if errStr, ok := val["error"].(string); ok {
				return errStr
			}
			if resultVal, ok := val["result"]; ok {
				if b, ok := resultVal.(bool); ok && !b {
					if comment, ok := val["comment"].(string); ok {
						return comment
					}
					return "执行失败"
				}
			}
		}
	}
	return ""
}

func getStringField(data map[string]interface{}, field string) string {
	if v, ok := data[field].(string); ok {
		return v
	}
	return ""
}

func getIPString(data map[string]interface{}) string {
	raw, ok := data["ipv4"].([]interface{})
	if !ok || len(raw) == 0 {
		return ""
	}
	ips := make([]string, 0, len(raw))
	for _, ip := range raw {
		if s, ok := ip.(string); ok && s != "127.0.0.1" {
			ips = append(ips, s)
		}
	}
	return strings.Join(ips, ",")
}

func toStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []string:
		return val
	case []interface{}:
		s := make([]string, len(val))
		for i, item := range val {
			s[i] = fmt.Sprintf("%v", item)
		}
		return s
	default:
		return nil
	}
}
