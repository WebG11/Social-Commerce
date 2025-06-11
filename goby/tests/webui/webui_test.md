# GOBY Web应用 WebUI测试

## 1. 测试概述

本文档详细介绍了GOBY Web应用的WebUI自动化测试实现，使用Selenium WebDriver框架对前端用户界面进行自动化测试，确保用户交互功能的正确性和稳定性。

## 2. 测试环境配置

### 2.1 技术栈
- **测试框架**: Python unittest
- **WebDriver**: Selenium WebDriver
- **浏览器**: Chrome WebDriver
- **等待机制**: WebDriverWait + Expected Conditions

## 2. 测试用例设计
测试代码：python tests/webui/webui_test.py

### 2.1 测试用例1：首页标题验证
**测试目的**: 验证应用首页能够正确加载并显示预期标题

**测试步骤**:
1. 启动Chrome浏览器
2. 访问应用首页 `http://localhost:3000`
3. 获取页面标题
4. 验证标题是否为"GoBy"

**预期结果**: 页面标题应显示为"GoBy"

### 2.2 测试用例2：用户登录功能测试
**测试目的**: 验证用户登录功能的完整流程

**测试步骤**:
1. 访问登录页面 `http://localhost:3000/seller/login`
2. 等待页面完全加载
3. 定位邮箱输入框并输入测试邮箱: `22009201315@stu.xidian.edu.cn`
4. 定位密码输入框并输入测试密码: `123456`
5. 点击"Continue"登录按钮
6. 验证页面URL是否发生跳转

**预期结果**: 登录成功后页面URL应发生改变，跳转到其他页面


