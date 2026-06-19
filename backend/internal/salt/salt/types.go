package salt

// ---------- 请求参数类型 ----------

// CmdData 封装发往 Salt REST API 的请求参数
// 对应原 Python 项目 salt_https_api.py 中 __init__ 收到的 data 字典
type CmdData struct {
	Fun      string `json:"fun"`                // Salt 执行函数，如 cmd.run / state.sls / grains.items
	Tgt      string `json:"tgt"`                // 目标主机，支持 glob / list / nodegroup
	Arg      string `json:"arg,omitempty"`      // 参数
	ExprForm string `json:"expr_form,omitempty"` // 目标匹配方式，默认 glob
	Client   string `json:"client"`             // 客户端类型：local / local_async / wheel / ssh / runner
}

// ---------- 登录请求 ----------

// LoginData Salt REST API 登录请求参数
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	EAuth    string `json:"eauth"`
}

// ---------- 响应类型 ----------

// SaltResponse Salt API 通用响应结构
// Salt REST API 返回格式：{"return": [{"jid": "...", "ret": {...}}]}
type SaltResponse struct {
	Return []map[string]interface{} `json:"return" yaml:"return"`
}

// JobInfo 作业信息，用于 loadJob 查询返回
type JobInfo struct {
	Result map[string]interface{} `json:"result" yaml:"result"`
}

// ---------- 业务层类型 ----------

// CmdResult 统一命令执行结果
type CmdResult struct {
	Success   bool
	JID       string                 // 异步执行时的作业 ID
	NodeName  string                 // 目标节点名（单节点时有效，多节点时为空）
	NodeNames []string               // 所有目标节点列表（单节点也包含）
	Data      map[string]interface{} // 原始返回数据
	Detail    interface{}            // 执行详情
	Error     string                 // 错误信息
}

// ---------- 主机信息 ----------

// HostInfo 主机信息（通过 grains.items 获取）
type HostInfo struct {
	ID         string `json:"id"`
	OS         string `json:"os"`
	OSRelease  string `json:"osrelease"`
	CPUModel   string `json:"cpu_model"`
	NumCPUs    int    `json:"num_cpus"`
	MemTotal   int    `json:"mem_total"`
	IPv4       string `json:"ipv4"`
	Kernel     string `json:"kernel"`
	SaltVersion string `json:"saltversion"`
	EnableSSH  bool   `json:"enable_ssh"`
}