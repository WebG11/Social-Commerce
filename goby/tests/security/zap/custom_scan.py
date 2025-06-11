#!/usr/bin/env python3
from zapv2 import ZAPv2
import time
import sys

# ZAP API配置
zap = ZAPv2(apikey='your-api-key', proxies={'http': 'http://localhost:8080', 'https': 'http://localhost:8080'})

def run_scan(target_url):
    print('开始扫描目标: ' + target_url)
    
    # 启动爬虫
    print('启动爬虫...')
    zap.spider.scan(target_url)
    
    # 等待爬虫完成
    while (int(zap.spider.status()) < 100):
        print('爬虫进度: ' + zap.spider.status() + '%')
        time.sleep(5)
    
    print('爬虫完成')
    
    # 启动主动扫描
    print('启动主动扫描...')
    zap.ascan.scan(target_url)
    
    # 等待主动扫描完成
    while (int(zap.ascan.status()) < 100):
        print('主动扫描进度: ' + zap.ascan.status() + '%')
        time.sleep(5)
    
    print('主动扫描完成')
    
    # 生成报告
    print('生成报告...')
    report = zap.core.htmlreport()
    
    # 保存报告
    with open('/zap/reports/security_report.html', 'w') as f:
        f.write(report)
    
    print('报告已保存到 /zap/reports/security_report.html')
    
    # 输出警报
    print('\n发现的安全问题:')
    for alert in zap.core.alerts():
        print('风险等级: ' + alert['risk'])
        print('URL: ' + alert['url'])
        print('问题: ' + alert['name'])
        print('描述: ' + alert['description'])
        print('解决方案: ' + alert['solution'])
        print('---')

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print('使用方法: python custom_scan.py <target_url>')
        sys.exit(1)
    
    run_scan(sys.argv[1]) 