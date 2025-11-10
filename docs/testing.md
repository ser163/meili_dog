以下是测试命令行，用于测试您的 MeiliSearch 搜索服务：

## 1. 健康检查测试

```bash
# 健康检查
curl -X GET "http://localhost:8081/api/v1/health"

# 带详细输出的健康检查
curl -v -X GET "http://localhost:8081/api/v1/health"
```

## 2. 索引信息测试

```bash
# 获取当前配置的索引信息
curl -X GET "http://localhost:8081/api/v1/index"

# 格式化输出索引信息
curl -s -X GET "http://localhost:8081/api/v1/index" | python -m json.tool
```

## 3. 基本搜索测试

```bash
# 简单搜索
curl -X GET "http://localhost:8081/api/v1/search?query=wonder"

# 带分页的搜索
curl -X GET "http://localhost:8081/api/v1/search?query=action&page=1&limit=5"

# 搜索特定字段
curl -X GET "http://localhost:8081/api/v1/search?query=drama"
```

## 4. 高级搜索测试

```bash
# 带过滤条件的搜索
curl -X GET "http://localhost:8081/api/v1/search?query=wonder&filters[id]=2"

# 多条件过滤（需要URL编码）
curl -G "http://localhost:8081/api/v1/search" \
  --data-urlencode "query=action" \
  --data-urlencode "filters[genres]=Action" \
  --data-urlencode "filters[id]=2"

# 排序搜索
curl -X GET "http://localhost:8081/api/v1/search?query=action&sort=id:desc"
```

## 5. 完整的测试脚本

创建一个测试脚本 `test_search.sh`：

```bash
#!/bin/bash

# 设置API基础URL
BASE_URL="http://localhost:8081/api/v1"

echo "=== MeiliSearch 搜索服务测试 ==="
echo "基础URL: $BASE_URL"
echo ""

# 1. 健康检查
echo "1. 健康检查:"
curl -s -X GET "$BASE_URL/health" | python -m json.tool
echo ""

# 2. 索引信息
echo "2. 索引信息:"
curl -s -X GET "$BASE_URL/index" | python -m json.tool
echo ""

# 3. 基本搜索测试
echo "3. 基本搜索测试:"
echo "搜索 'wonder':"
curl -s -X GET "$BASE_URL/search?query=wonder" | python -m json.tool
echo ""

echo "搜索 'action' (第1页，每页3条):"
curl -s -X GET "$BASE_URL/search?query=action&page=1&limit=3" | python -m json.tool
echo ""

# 4. 过滤搜索测试
echo "4. 过滤搜索测试:"
echo "搜索 'action' 并过滤类型为 'Action':"
curl -s -G "$BASE_URL/search" \
  --data-urlencode "query=action" \
  --data-urlencode "filters[genres]=Action" | python -m json.tool
echo ""

echo "=== 测试完成 ==="
```

给脚本执行权限并运行：
```bash
chmod +x test_search.sh
./test_search.sh
```

## 6. 使用 HTTPie 进行测试（如果安装了）

```bash
# 安装 HTTPie: pip install httpie

# 健康检查
http GET http://localhost:8081/api/v1/health

# 搜索测试
http GET http://localhost:8081/api/v1/search query=="wonder"

# 带过滤的搜索
http GET http://localhost:8081/api/v1/search \
  query=="action" \
  filters[genres]=="Action" \
  page==1 \
  limit==5
```

## 7. 使用 Postman 或 curl 的复杂测试

```bash
# 复杂过滤条件（JSON格式）
curl -X GET "http://localhost:8081/api/v1/search" \
  -G \
  --data-urlencode "query=action" \
  --data-urlencode 'filters={"genres": "Action", "id": 2}'

# 测试错误情况
echo "测试错误情况:"
curl -X GET "http://localhost:8081/api/v1/search"  # 缺少query参数
curl -X GET "http://localhost:8081/api/v1/search?query="  # 空查询
```

## 8. 性能测试脚本

创建性能测试脚本 `test_performance.sh`：

```bash
#!/bin/bash

echo "=== 性能测试 ==="

# 测试搜索响应时间
for i in {1..5}; do
    echo "测试 $i:"
    time curl -s -o /dev/null -w "HTTP状态: %{http_code}\n总时间: %{time_total}秒\n" \
        "http://localhost:8081/api/v1/search?query=test&limit=10"
    echo "---"
done
```

## 9. 环境变量设置

如果您的服务使用不同的端口或主机，可以设置环境变量：

```bash
export SEARCH_HOST="http://localhost:8081"
export API_BASE="$SEARCH_HOST/api/v1"

# 然后使用变量进行测试
curl -X GET "$API_BASE/health"
```

## 使用说明

1. **首先确保服务正在运行**：
   ```bash
   go run main.go
   ```

2. **测试前确保 MeiliSearch 服务可用**：
   ```bash
   curl -X GET "http://localhost:7700/health"
   ```