package salt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gopkg.in/yaml.v3"
)

// ---------- Salt REST API 客户端 ----------

// Client 封装所有对 Salt REST API 的 HTTP 调用
// 每个实例持有自己的 Token 缓存，互不干扰
// 对应原 Python 项目 saltrest/salt_https_api.py 中的 salt_api_token 类
type Client struct {
	baseURL    string
	username   string
	password   string
	eauth      string
	httpClient *http.Client
	tokenCache *tokenCache
}

// NewClient 创建新的 Salt API 客户端
func NewClient(baseURL, username, password, eauth string) *Client {
	if eauth == "" {
		eauth = "pam"
	}
	return &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
		eauth:    eauth,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		tokenCache: &tokenCache{},
	}
}

// getToken 从缓存获取有效 token，缓存失效则重新登录
func (c *Client) getToken() (string, error) {
	// 先查实例缓存
	if token := c.tokenCache.get(); token != "" {
		return token, nil
	}

	// 缓存失效，重新登录
	token, err := Login(c.baseURL, c.username, c.password, c.eauth)
	if err != nil {
		return "", err
	}

	// 存入缓存
	c.tokenCache.set(token)
	return token, nil
}

// ForceLogin 强制重新登录（忽略缓存），用于 token 失效后的重试
func (c *Client) ForceLogin() (string, error) {
	c.tokenCache.clear()
	return c.getToken()
}

// buildHeaders 设置请求头
func (c *Client) buildHeaders(req *http.Request, token string) {
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Accept", "application/x-yaml")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GoSaltOps/1.0")
}

// ---------- 核心请求方法 ----------

// doPost 核心 POST 请求方法
func (c *Client) doPost(data *CmdData) (*SaltResponse, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("获取 Token 失败: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建 POST 请求失败: %w", err)
	}
	c.buildHeaders(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Salt API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Salt API 返回非 200 (状态码=%d): %s", resp.StatusCode, string(respBody))
	}

	var result SaltResponse
	if err := yaml.Unmarshal(respBody, &result); err != nil {
		// 兼容 JSON 格式
		if jsonErr := json.Unmarshal(respBody, &result); jsonErr != nil {
			return nil, fmt.Errorf("解析响应失败 (YAML/JSON): %w | 原始内容: %s", err, string(respBody))
		}
	}

	return &result, nil
}

// doGet GET 请求方法（用于查询作业）
func (c *Client) doGet(path string) (*SaltResponse, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("获取 Token 失败: %w", err)
	}

	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 GET 请求失败: %w", err)
	}
	c.buildHeaders(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Salt API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result SaltResponse
	if err := yaml.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w | 原始内容: %s", err, string(respBody))
	}

	return &result, nil
}

// ---------- 公开方法：对应原项目的 5 种调用方式 ----------

// Run 异步执行任务（client=local_async）
// 对应原项目 salt_api_token.run()
func (c *Client) Run(fun, tgt, arg string) (*SaltResponse, error) {
	data := &CmdData{
		Fun:    fun,
		Tgt:    tgt,
		Arg:    arg,
		Client: "local_async",
	}
	return c.doPost(data)
}

// CmdRun 同步执行任务
// 对应原项目 salt_api_token.CmdRun(client='local')
//
// client 参数: "local" 或 "ssh"
func (c *Client) CmdRun(fun, tgt, arg, client string) (*SaltResponse, error) {
	if client == "" {
		client = "local"
	}
	data := &CmdData{
		Fun:    fun,
		Tgt:    tgt,
		Arg:    arg,
		Client: client,
	}
	return c.doPost(data)
}

// CmdRunWithExpr 同步执行（带 expr_form 匹配方式参数）
func (c *Client) CmdRunWithExpr(fun, tgt, arg, client, exprForm string) (*SaltResponse, error) {
	if client == "" {
		client = "local"
	}
	if exprForm == "" {
		exprForm = "glob"
	}
	data := &CmdData{
		Fun:      fun,
		Tgt:      tgt,
		Arg:      arg,
		Client:   client,
		ExprForm: exprForm,
	}
	return c.doPost(data)
}

// WheelRun 使用 Wheel 客户端执行（用于 Master 管理操作，如 key 管理）
// 对应原项目 salt_api_token.wheelRun()
func (c *Client) WheelRun(fun, arg string) (*SaltResponse, error) {
	data := &CmdData{
		Fun:    fun,
		Arg:    arg,
		Client: "wheel",
	}
	return c.doPost(data)
}

// SSHExecute 使用 Salt-SSH 模式执行
// 对应原项目 salt_api_token.sshRun()
func (c *Client) SSHExecute(fun, tgt, arg string) (*SaltResponse, error) {
	data := &CmdData{
		Fun:    fun,
		Tgt:    tgt,
		Arg:    arg,
		Client: "ssh",
	}
	return c.doPost(data)
}

// RunnerRun 使用 Runner 客户端执行（Master 端模块，如 manage.status）
// 对应原项目 salt_api_token.runnerRun()
func (c *Client) RunnerRun(fun, arg string) (*SaltResponse, error) {
	data := &CmdData{
		Fun:    fun,
		Arg:    arg,
		Client: "runner",
	}
	return c.doPost(data)
}

// LoadJob 查询作业执行状态
// 对应原项目 salt_api_token.loadJob(jid)
//
// GET {baseURL}/jobs/{jid}
func (c *Client) LoadJob(jid string) (*SaltResponse, error) {
	return c.doGet("jobs/" + jid)
}

// ---------- 便捷方法 ----------

// StateApply 执行 state.apply
func (c *Client) StateApply(tgt, client string) (*SaltResponse, error) {
	return c.CmdRun("state.apply", tgt, "", client)
}

// StateSLS 执行 state.sls
func (c *Client) StateSLS(tgt, slsName, client string) (*SaltResponse, error) {
	return c.CmdRun("state.sls", tgt, slsName, client)
}

// CmdScript 下发脚本执行
func (c *Client) CmdScript(tgt, saltPath, client string) (*SaltResponse, error) {
	return c.CmdRun("cmd.script", tgt, saltPath, client)
}

// GrainsItems 获取主机 grains 信息
func (c *Client) GrainsItems(tgt string) (*SaltResponse, error) {
	return c.CmdRun("grains.items", tgt, "", "local")
}

// ManageStatus 获取 minion 在线状态
func (c *Client) ManageStatus() (*SaltResponse, error) {
	return c.RunnerRun("manage.status", "")
}

// KeyListAll 列出所有 Salt Key
func (c *Client) KeyListAll() (*SaltResponse, error) {
	return c.RunnerRun("key.list_all", "")
}
