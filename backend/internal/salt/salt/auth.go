package salt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// ---------- Token 管理 ----------

// tokenCache 缓存 X-Auth-Token，避免每次调用都重新登录
// 每个 Client 实例持有自己的缓存，互不干扰
type tokenCache struct {
	mu        sync.RWMutex
	token     string
	expiresAt time.Time
}

// get 返回缓存的 token，过期返回空
func (tc *tokenCache) get() string {
	if tc == nil {
		return ""
	}
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if tc.token != "" && time.Now().Before(tc.expiresAt) {
		return tc.token
	}
	return ""
}

// set 缓存 token（Salt token 默认有效期 12 小时，这里保守设为 10 小时）
func (tc *tokenCache) set(token string) {
	if tc == nil {
		return
	}
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.token = token
	tc.expiresAt = time.Now().Add(10 * time.Hour)
}

// clear 清除缓存的 token
func (tc *tokenCache) clear() {
	if tc == nil {
		return
	}
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.token = ""
	tc.expiresAt = time.Time{}
}

// Login 登录 Salt REST API 获取 X-Auth-Token（无缓存，每次真实请求）
//
// 对应原 Python 项目 saltrest/salt_token_id.py 的 token_id() 函数
// 客户端应使用 Client.Login() 或 Client.ForceLogin() 以获得缓存能力
//
// 请求 URL:   {baseURL}/login
// 请求参数:   username, password, eauth
// 响应格式:   {"return": [{"token": "..."}]}
func Login(baseURL, username, password, eauth string) (string, error) {
	payload := LoginData{
		Username: username,
		Password: password,
		EAuth:    eauth,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("序列化登录参数失败: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+"login", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("创建登录请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/x-yaml")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Salt API 登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取登录响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Salt API 登录返回非 200 (状态码=%d): %s", resp.StatusCode, string(respBody))
	}

	// Salt API 默认返回 YAML 格式
	var result map[string]interface{}
	if err := yaml.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析登录响应(YAML)失败: %w", err)
	}

	// 提取 token：响应结构为 {"return": [{"token": "xxx", ...}]}
	returnList, ok := result["return"].([]interface{})
	if !ok || len(returnList) == 0 {
		return "", fmt.Errorf("登录响应中缺少 return 字段: %s", string(respBody))
	}

	returnItem, ok := returnList[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("登录响应 return[0] 格式异常: %s", string(respBody))
	}

	token, ok := returnItem["token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("登录响应中未找到 token: %s", string(respBody))
	}

	return token, nil
}