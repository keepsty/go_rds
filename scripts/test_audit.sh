#!/bin/bash
# =============================================================================
# GoRDS 审核功能测试脚本
# 用法: bash scripts/test_audit.sh
# 前提: 后端服务已运行 (go run cmd/main.go)
# =============================================================================

set -euo pipefail

BASE="http://localhost:8083"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
PASS=0; FAIL=0
P() { echo -e "${GREEN}[PASS]${NC} $1"; PASS=$((PASS+1)); }
F() { echo -e "${RED}[FAIL]${NC} $1"; FAIL=$((FAIL+1)); }
I() { echo -e "${BLUE}[INFO]${NC} $1"; }
W() { echo -e "${YELLOW}[WARN]${NC} $1"; }

# ---- 登录 ----
I "Step 1: 登录..."
RESP=$(curl -s -X POST "$BASE/api/v1/user/login" \
    -H 'Content-Type: application/json' \
    -d '{"username":"admin","password":"1234.Com!"}')
TOKEN=$(echo "$RESP" | python3 -c "import sys,json;print(json.load(sys.stdin)['data']['token'])" 2>/dev/null)
[ -n "$TOKEN" ] && [ "$TOKEN" != "null" ] && P "登录成功" || { F "登录失败"; exit 1; }
AUTH="Authorization: JWT $TOKEN"

# ---- 准备环境 ----
I "Step 2: 准备测试环境..."
ENV_ID=$(curl -s "$BASE/api/v1/admin/environments" -H "$AUTH" | python3 -c "import sys,json;d=json.load(sys.stdin)['data'];print(d[0]['id'] if d else '')")
if [ -z "$ENV_ID" ]; then
    curl -s -X POST "$BASE/api/v1/admin/environment" -H "$AUTH" -H 'Content-Type: application/json' -d '{"name":"测试环境"}' >/dev/null
    ENV_ID=$(curl -s "$BASE/api/v1/admin/environments" -H "$AUTH" | python3 -c "import sys,json;print(json.load(sys.stdin)['data'][0]['id'])")
fi
P "环境ID: $ENV_ID"
ORG=$(curl -s "$BASE/api/v1/admin/organizations" -H "$AUTH" | python3 -c "import sys,json;d=json.load(sys.stdin)['data'];print(d[0]['organization_key'] if d else '')")
if [ -z "$ORG" ]; then
    curl -s -X POST "$BASE/api/v1/admin/organization" -H "$AUTH" -H 'Content-Type: application/json' -d '{"name":"测试组织","organization_key":"test_org"}' >/dev/null
    ORG="test_org"; fi
P "组织: $ORG"
INSTANCE_ID=$(curl -s "$BASE/api/v1/admin/instances" -H "$AUTH" | python3 -c "import sys,json;d=json.load(sys.stdin)['data'];print(d[0]['instance_id'] if d else '')")
if [ -z "$INSTANCE_ID" ]; then
    curl -s -X POST "$BASE/api/v1/admin/instances" -H "$AUTH" -H 'Content-Type: application/json' \
        -d "{\"hostname\":\"127.0.0.1\",\"port\":3306,\"user\":\"go_rds_rw\",\"password\":\"1234.Com!\",\"use_type\":\"工单\",\"db_type\":\"MySQL\",\"environment\":$ENV_ID,\"organization_key\":[\"$ORG\"],\"remark\":\"测试实例\"}" >/dev/null
    INSTANCE_ID=$(curl -s "$BASE/api/v1/admin/instances" -H "$AUTH" | python3 -c "import sys,json;print(json.load(sys.stdin)['data'][0]['instance_id'])")
fi
[ -n "$INSTANCE_ID" ] && P "实例ID: $INSTANCE_ID" || { F "无法获得实例ID"; exit 1; }

# ---- 审核测试 ----
I ""
I "============================================="
I "Step 3: SQL 审核功能测试"
I "============================================="

# 用 python3 构建 JSON body，避免 shell 转义问题
do_audit() {
    local sql_type="$1" schema="$2" content="$3"
    python3 -c "
import json, sys
body = json.dumps({
    'db_type': 'MySQL',
    'sql_type': '$sql_type',
    'instance_id': '$INSTANCE_ID',
    'schema': '$schema',
    'content': '''$content'''
})
sys.stdout.write(body)
" | curl -s --max-time 8 -X POST "$BASE/api/v1/orders/inspect-syntax" \
    -H "$AUTH" -H 'Content-Type: application/json' -d @-
}

check_audit() {
    local label="$1" sql_type="$2" schema="$3" content="$4" expect="$5"
    RESP=$(do_audit "$sql_type" "$schema" "$content")
    CODE=$(echo "$RESP" | python3 -c "import sys,json;print(json.load(sys.stdin).get('code','ERROR'))" 2>/dev/null)

    if [ "$expect" = "api_error" ]; then
        if [ "$CODE" != "0000" ]; then
            P "$label → API拒绝(符合预期)"
        else
            F "$label → 期望API报错但成功了"
        fi
        return
    fi

    if [ "$CODE" != "0000" ]; then
        MSG=$(echo "$RESP" | python3 -c "import sys,json;print(json.load(sys.stdin).get('message',''))")
        F "$label → API错误: $MSG"
        return
    fi

    STATUS=$(echo "$RESP" | python3 -c "import sys,json;print(json.load(sys.stdin)['data']['status'])")
    DETAILS=$(echo "$RESP" | python3 -c "
import sys,json
for item in json.load(sys.stdin)['data']['data']:
    for s in item['summary']:
        if s['level'] != 'INFO':
            print(f\" [{s['level']}] {s['message'][:40]}\", end='')
" 2>/dev/null)

    if [ "$STATUS" = "$expect" ]; then
        P "$label$DETAILS"
    else
        F "$label → 期望=$expect 实际=$STATUS$DETAILS"
    fi
}

# ========== DDL ==========
I ""; I "--- DDL 审核 ---"
check_audit \
    "规范建表(完整)" \
    "DDL" "go_rds" \
    "CREATE TABLE test_new_table (id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键', name varchar(64) NOT NULL DEFAULT '' COMMENT '名称', created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间', updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间', PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='测试用表';" \
    "0"

check_audit \
    "建表(有建议)" \
    "DDL" "go_rds" \
    "CREATE TABLE test_ok_table (id bigint NOT NULL AUTO_INCREMENT COMMENT 'ID', name varchar(64) NOT NULL COMMENT '名称', PRIMARY KEY (id)) ENGINE=InnoDB CHARSET=utf8mb4;" \
    "1"

# ========== DML ==========
I ""; I "--- DML 审核 ---"
check_audit \
    "UPDATE 有WHERE" \
    "DML" "go_rds" \
    "UPDATE insight_users SET nick_name='test' WHERE username='admin';" \
    "0"

check_audit \
    "INSERT 带字段列表" \
    "DML" "go_rds" \
    "INSERT INTO insight_das_records (user_name, instance_id, db_name, sql_text, status) VALUES ('admin','abc-def','go_rds','SELECT 1','FINISH');" \
    "1"

check_audit \
    "DELETE 有WHERE" \
    "DML" "go_rds" \
    "DELETE FROM insight_users WHERE username='nonexist';" \
    "0"

# ========== 边界 ==========
I ""; I "--- 边界情况 ---"
check_audit \
    "语法错误(不完整)" \
    "DDL" "go_rds" \
    "CREATE TABLE bad (id bigint" \
    "api_error"

# ========== 参数 ==========
I ""; I "--- 审核参数接口 ---"
PARAMS=$(curl -s "$BASE/api/v1/admin/inspect/params" -H "$AUTH" | python3 -c "import sys,json;d=json.load(sys.stdin);print(d.get('total',0))")
[ "$PARAMS" -gt 0 ] && P "审核参数接口正常 ($PARAMS 条)" || W "审核参数为空"

# ---- 结果 ----
echo ""
echo "============================================="
echo "  通过: $PASS  失败: $FAIL"
echo "============================================="
[ "$FAIL" -eq 0 ] && echo -e "${GREEN} 全部通过！${NC}" || echo -e "${RED} 有 $FAIL 项失败${NC}"
exit $FAIL
