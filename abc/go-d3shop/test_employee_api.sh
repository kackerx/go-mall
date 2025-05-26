#!/bin/bash

# 员工API测试脚本

echo "=== 员工API测试 ==="
echo ""

# 基础URL
BASE_URL="http://localhost:8080/api"

# 1. 创建员工（使用MediatorV2）
echo "1. 创建员工（使用MediatorV2）"
curl -X POST ${BASE_URL}/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "李四",
    "email": "lisi@company.com",
    "departmentId": 1,
    "position": "高级工程师",
    "baseSalary": 15000
  }' | jq .

echo ""
echo "等待2秒，让领域事件处理完成..."
sleep 2

# 2. 创建员工（使用服务层）
echo ""
echo "2. 创建员工（使用服务层）"
curl -X POST ${BASE_URL}/employees/service \
  -H "Content-Type: application/json" \
  -d '{
    "name": "王五",
    "email": "wangwu@company.com",
    "departmentId": 1,
    "position": "产品经理",
    "baseSalary": 18000
  }' | jq .

echo ""
echo "等待2秒，让领域事件处理完成..."
sleep 2

# 3. 获取员工信息
echo ""
echo "3. 获取员工信息"
curl -X GET ${BASE_URL}/employees/1 | jq .

# 4. 模拟员工入职集成事件
echo ""
echo "4. 模拟员工入职集成事件"
curl -X POST ${BASE_URL}/employees/1/simulate-event | jq .

echo ""
echo "=== 测试完成 ==="

# 查看服务器日志可以看到：
# - 领域事件处理：更新部门员工数、创建工资记录
# - 集成事件发布：发送到消息队列
# - 通知服务处理：发送邮件通知 