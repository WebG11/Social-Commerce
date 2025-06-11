1. 启动安全测试环境：
```bash
cd tests/security
docker-compose up -d
```

2. 运行不同类型的扫描：

a. 基础扫描：
```bash
docker-compose run zap-baseline
```

b. API扫描：
```bash
docker-compose run zap-api-scan
```

c. 自定义扫描：
```bash
# 进入ZAP容器
docker exec -it goby_zap_1 bash

# 运行自定义扫描脚本
python /zap/wrk/custom_scan.py http://host.docker.internal:8080
```

3. 查看测试结果：
- 所有报告都会保存在 `reports` 目录下
- 可以打开 `reports/security_report.html` 查看详细报告
- 报告包含：
  - 发现的安全漏洞
  - 风险等级
  - 问题描述
  - 修复建议

4. 测试范围包括：
- XSS（跨站脚本）攻击
- SQL注入
- CSRF（跨站请求伪造）
- 目录遍历
- 命令注入
- 文件包含漏洞
- 其他常见Web安全漏洞

5. 安全测试建议：
- 在测试环境中运行，避免影响生产系统
- 定期运行安全扫描
- 及时修复发现的安全问题
- 保持ZAP工具和规则库更新

6. 自定义测试：
- 可以修改 `gen.conf` 文件调整扫描规则
- 可以编辑 `custom_scan.py` 添加特定的测试场景
- 可以配置不同的扫描策略

7. 注意事项：
- 确保目标应用已经启动
- 检查网络连接是否正常
- 注意扫描可能对系统造成的影响
- 保存测试报告以便后续分析

8. 停止测试环境：
```bash
docker-compose down
```

这个安全测试环境提供了：
- 自动化安全扫描
- 多种扫描模式
- 详细的测试报告
- 可定制的测试规则

您可以根据需要：
1. 调整扫描规则
2. 添加特定的测试场景
3. 自定义报告格式
4. 集成到CI/CD流程
