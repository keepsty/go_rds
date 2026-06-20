# OWASP Top 10 安全检查清单

## 1. 注入攻击（Injection）
- [ ] SQL 查询是否使用参数化查询或 ORM？
- [ ] 是否存在字符串拼接 SQL 的情况？
- [ ] 用户输入是否经过转义和过滤？

## 2. 身份认证失效（Broken Authentication）
- [ ] 密码是否明文存储？（应使用 bcrypt 等加密）
- [ ] 会话令牌是否使用安全的随机数生成？
- [ ] 是否有登录失败次数限制？

## 3. 敏感数据泄露（Sensitive Data Exposure）
- [ ] API 密钥、数据库密码是否硬编码在代码中？
- [ ] 敏感数据是否通过 HTTPS 传输？
- [ ] 日志中是否记录了敏感信息？

## 4. XSS 跨站脚本攻击
- [ ] 用户输入是否在渲染前经过转义？
- [ ] 是否使用 dangerouslySetInnerHTML 等危险 API？
- [ ] CSP（Content Security Policy）头是否设置？

## 5. 安全配置错误
- [ ] 是否关闭了调试模式？
- [ ] 错误页面是否暴露了堆栈信息？
- [ ] 默认账户密码是否已修改？
