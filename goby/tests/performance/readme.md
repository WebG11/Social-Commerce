1. 首先，确保您的项目服务已经启动：
```bash
# 在项目根目录下
make
```

2. 启动性能测试环境：
```bash
cd tests/performance
docker-compose up -d
```

3. 等待所有容器启动完成，可以通过以下命令查看容器状态：
```bash
docker-compose ps
```

4. 运行性能测试：
```bash
# 进入JMeter容器
docker exec -it goby_jmeter_1 bash

# 在容器内执行测试
jmeter -n -t /jmeter/performance_test.jmx -l /results/results.jtl -e -o /results/report
```

5. 查看测试结果：
- 打开浏览器访问以下地址：
  - JMeter测试报告：`http://localhost:8080/report`
  - Grafana仪表板：`http://localhost:3100`（默认账号：admin/admin）
  - Prometheus指标：`http://localhost:9090`

6. 在Grafana中配置数据源：
   - 登录Grafana（http://localhost:3000）
   - 添加Prometheus数据源（URL: http://prometheus:9090）
   - 创建新的仪表板，添加以下面板：
     - 响应时间
     - 请求数/秒
     - 错误率
     - 并发用户数

7. 分析测试结果：
   - 查看响应时间分布
   - 检查错误率
   - 分析系统瓶颈
   - 观察资源使用情况

8. 如果需要停止测试环境：
```bash
docker-compose down
```

9. 如果需要修改测试参数，可以编辑 `jmeter/performance_test.jmx` 文件：
   - 调整并发用户数（ThreadGroup.num_threads）
   - 修改循环次数（LoopController.loops）
   - 更改测试接口和参数

10. 保存测试结果：
    - 测试结果会保存在 `results` 目录下
    - 可以导出Grafana仪表板配置
    - 保存Prometheus查询配置

建议：
1. 首次运行时，建议使用较小的并发用户数（如10-20）进行测试
2. 逐步增加并发用户数，观察系统表现
3. 记录每次测试的参数和结果，便于对比分析
4. 如果发现性能问题，可以：
   - 检查服务器资源使用情况
   - 分析数据库性能
   - 查看应用日志
   - 优化代码和配置
