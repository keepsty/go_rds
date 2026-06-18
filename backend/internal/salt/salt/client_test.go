package salt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// ---------- 模拟 Salt REST API 服务器 ----------

// mockSaltServer 创建一个模拟的 Salt API 服务器
// 它会打印收到的请求，并返回预设的模拟响应
func mockSaltServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ===== 打印收到的请求信息 =====
		fmt.Println("\n  ┌──────────────────────────────────────────────")
		fmt.Printf("  │ 收到请求: %s %s\n", r.Method, r.URL.Path)
		fmt.Println("  │──────────────────────────────────────────────")
		fmt.Printf("  │ Headers:\n")
		for k, v := range r.Header {
			// 隐藏 token 值
			if k == "X-Auth-Token" {
				fmt.Printf("  │   %s: %s\n", k, v[0][:8]+"...(截断)")
			} else {
				fmt.Printf("  │   %s: %s\n", k, strings.Join(v, ", "))
			}
		}

		// 如果有请求体，打印出来
		if r.Body != nil && r.ContentLength > 0 {
			buf := make([]byte, r.ContentLength)
			r.Body.Read(buf)
			fmt.Printf("  │ Body:\n")
			// 尝试格式化 JSON
			var pretty map[string]interface{}
			if err := json.Unmarshal(buf, &pretty); err == nil {
				prettyJSON, _ := json.MarshalIndent(pretty, "  │   ", "  ")
				fmt.Printf("  │   %s\n", string(prettyJSON))
			} else {
				fmt.Printf("  │   %s\n", string(buf))
			}
		}
		fmt.Println("  └──────────────────────────────────────────────")

		// ===== 根据请求路径返回模拟响应 =====
		w.Header().Set("Content-Type", "application/x-yaml")

		switch {
		case strings.HasSuffix(r.URL.Path, "/login"):
			// 登录响应
			loginResp := map[string]interface{}{
				"return": []interface{}{
					map[string]interface{}{
						"token":  "mock-salt-token-abc123def456",
						"expire": 43200,
						"start":  0,
						"user":   "saltops",
						"eauth":  "pam",
					},
				},
			}
			writeYAML(w, loginResp)
			fmt.Println("  >> 返回模拟登录响应 ✓")

		case strings.Contains(r.URL.Path, "/jobs/"):
			// 作业查询响应
			jobResp := map[string]interface{}{
				"return": []interface{}{
					map[string]interface{}{
						"jid": strings.TrimPrefix(r.URL.Path, "/jobs/"),
						"result": map[string]interface{}{
							"minion-01": map[string]interface{}{
								"success": true,
								"return":  "hello world",
							},
						},
					},
				},
			}
			writeYAML(w, jobResp)
			fmt.Println("  >> 返回模拟作业查询响应 ✓")

		default:
			// 普通命令执行响应
			cmdResp := map[string]interface{}{
				"return": []interface{}{
					map[string]interface{}{
						"jid": "20250617010101010101",
						"minion-01": map[string]interface{}{
							"success": true,
							"return":  `模拟执行结果: 命令已执行`,
						},
					},
				},
			}
			writeYAML(w, cmdResp)
			fmt.Println("  >> 返回模拟命令执行响应 ✓")
		}
	}))
}

func writeYAML(w http.ResponseWriter, data interface{}) {
	body, _ := yaml.Marshal(data)
	w.Write(body)
}

// ---------- 测试用例 ----------

// TestClient_Login 测试登录流程
func TestClient_Login(t *testing.T) {
	server := mockSaltServer(t)
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试1: Login — 登录 Salt API")
	fmt.Println("========================================")

	token, err := Login(server.URL+"/", "saltops", "saltops", "pam")
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	fmt.Printf("\n  ✦ 获取到 Token: %s\n", token)
	if token != "mock-salt-token-abc123def456" {
		t.Errorf("Token 不匹配，期望 mock-salt-token-abc123def456，得到 %s", token)
	}
	fmt.Println("  ✦ 登录测试通过 ✓")
}

// TestClient_CmdRun 测试同步执行 cmd.run
func TestClient_CmdRun(t *testing.T) {
	server := mockSaltServer(t)
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试2: CmdRun — 同步执行 cmd.run")
	fmt.Println("========================================")

	client := NewClient(server.URL+"/", "saltops", "saltops", "pam")
	resp, err := client.CmdRun("cmd.run", "*", "uptime", "local")
	if err != nil {
		t.Fatalf("CmdRun 失败: %v", err)
	}

	fmt.Printf("\n  ✦ 响应 Return 条目数: %d\n", len(resp.Return))
	for k, v := range resp.Return[0] {
		fmt.Printf("  ✦ %s: %v\n", k, v)
	}
	fmt.Println("  ✦ CmdRun 测试通过 ✓")
}

// TestClient_StateSLS 测试执行 state.sls
func TestClient_StateSLS(t *testing.T) {
	server := mockSaltServer(t)
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试3: StateSLS — 执行 state.sls")
	fmt.Println("========================================")

	client := NewClient(server.URL+"/", "saltops", "saltops", "pam")
	resp, err := client.StateSLS("minion-01", "my_state", "local")
	if err != nil {
		t.Fatalf("StateSLS 失败: %v", err)
	}

	fmt.Printf("\n  ✦ 响应 Return 条目数: %d\n", len(resp.Return))
	for k, v := range resp.Return[0] {
		fmt.Printf("  ✦ %s: %v\n", k, v)
	}
	fmt.Println("  ✦ StateSLS 测试通过 ✓")
}

// TestClient_AsyncRun 测试异步执行
func TestClient_AsyncRun(t *testing.T) {
	server := mockSaltServer(t)
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试4: Run — 异步执行 (local_async)")
	fmt.Println("========================================")

	client := NewClient(server.URL+"/", "saltops", "saltops", "pam")
	resp, err := client.Run("cmd.run", "*", "echo hello")
	if err != nil {
		t.Fatalf("AsyncRun 失败: %v", err)
	}

	jid := ""
	if len(resp.Return) > 0 {
		if j, ok := resp.Return[0]["jid"].(string); ok {
			jid = j
		}
	}
	fmt.Printf("\n  ✦ 获取到 JID: %s\n", jid)
	fmt.Println("  ✦ 异步执行测试通过 ✓")
}

// TestClient_LoadJob 测试查询作业结果
func TestClient_LoadJob(t *testing.T) {
	server := mockSaltServer(t)
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试5: LoadJob — 查询作业结果")
	fmt.Println("========================================")

	client := NewClient(server.URL+"/", "saltops", "saltops", "pam")
	resp, err := client.LoadJob("20250617010101010101")
	if err != nil {
		t.Fatalf("LoadJob 失败: %v", err)
	}

	fmt.Printf("\n  ✦ 响应 Return 条目数: %d\n", len(resp.Return))
	for k, v := range resp.Return[0] {
		fmt.Printf("  ✦ %s: %v\n", k, v)
	}
	fmt.Println("  ✦ LoadJob 测试通过 ✓")
}

// TestClient_WheelRunner 测试 Wheel 和 Runner 调用
func TestClient_WheelRunner(t *testing.T) {
	server := mockSaltServer(t)
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试6: WheelRun + RunnerRun")
	fmt.Println("========================================")

	client := NewClient(server.URL+"/", "saltops", "saltops", "pam")

	// Wheel
	resp, err := client.WheelRun("key.list_all", "")
	if err != nil {
		t.Fatalf("WheelRun 失败: %v", err)
	}
	fmt.Printf("\n  ✦ WheelRun 返回: %d 条\n", len(resp.Return))

	// Runner
	resp, err = client.RunnerRun("manage.status", "")
	if err != nil {
		t.Fatalf("RunnerRun 失败: %v", err)
	}
	fmt.Printf("  ✦ RunnerRun 返回: %d 条\n", len(resp.Return))

	fmt.Println("  ✦ Wheel/Runner 测试通过 ✓")
}

// TestClient_SSH 测试 Salt-SSH 模式
func TestClient_SSH(t *testing.T) {
	server := mockSaltServer(t)
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试7: SSHExecute — Salt-SSH 模式")
	fmt.Println("========================================")

	client := NewClient(server.URL+"/", "saltops", "saltops", "pam")
	resp, err := client.SSHExecute("cmd.run", "host-without-minion", "whoami")
	if err != nil {
		t.Fatalf("SSHExecute 失败: %v", err)
	}

	fmt.Printf("\n  ✦ SSHExecute 返回: %d 条\n", len(resp.Return))
	fmt.Println("  ✦ SSHExecute 测试通过 ✓")
}

// TestClient_Methods 测试便捷方法合集
func TestClient_Methods(t *testing.T) {
	server := mockSaltServer(t)
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试8: 便捷方法 — StateApply / CmdScript / GrainsItems / ManageStatus / KeyListAll")
	fmt.Println("========================================")

	client := NewClient(server.URL+"/", "saltops", "saltops", "pam")

	tests := []struct {
		name string
		fn   func() (*SaltResponse, error)
	}{
		{"StateApply", func() (*SaltResponse, error) { return client.StateApply("*", "local") }},
		{"CmdScript", func() (*SaltResponse, error) { return client.CmdScript("*", "salt://deploy.sh", "local") }},
		{"GrainsItems", func() (*SaltResponse, error) { return client.GrainsItems("minion-01") }},
		{"ManageStatus", func() (*SaltResponse, error) { return client.ManageStatus() }},
		{"KeyListAll", func() (*SaltResponse, error) { return client.KeyListAll() }},
	}

	for _, tt := range tests {
		resp, err := tt.fn()
		if err != nil {
			t.Errorf("%s 失败: %v", tt.name, err)
			continue
		}
		fmt.Printf("  ✦ %-15s → %d 条 return\n", tt.name, len(resp.Return))
	}

	fmt.Println("  ✦ 所有便捷方法测试通过 ✓")
}

// TestClient_TokenCache 测试 Token 缓存机制
func TestClient_TokenCache(t *testing.T) {
	loginCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")

		if strings.HasSuffix(r.URL.Path, "/login") {
			loginCount++
			writeYAML(w, map[string]interface{}{
				"return": []interface{}{map[string]interface{}{
					"token":  "cached-token",
					"expire": 43200,
					"user":   "saltops",
					"eauth":  "pam",
				}},
			})
			return
		}

		writeYAML(w, map[string]interface{}{
			"return": []interface{}{map[string]interface{}{
				"jid": "test-jid",
				"*": map[string]interface{}{
					"success": true,
					"return":  "ok",
				},
			}},
		})
	}))
	defer server.Close()

	fmt.Println("\n========================================")
	fmt.Println("  测试9: Token 缓存机制")
	fmt.Println("========================================")

	client := NewClient(server.URL+"/", "saltops", "saltops", "pam")

	// 第一次调用应触发登录
	client.CmdRun("test.fun", "*", "", "local")
	fmt.Printf("  第一次调用后 login 次数: %d\n", loginCount)

	// 第二次调用应使用缓存，不再重复登录
	client.CmdRun("test.fun", "*", "", "local")
	fmt.Printf("  第二次调用后 login 次数: %d (应该还是 1)\n", loginCount)

	if loginCount != 1 {
		t.Errorf("Token 缓存失效: 期望 1 次登录，实际 %d 次", loginCount)
	} else {
		fmt.Println("  ✦ Token 缓存机制正常 ✓")
	}
}
