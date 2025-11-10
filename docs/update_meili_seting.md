# 更新meili 索引配置
## 1. 测试命令

创建测试脚本 `test_settings.sh`：

```bash
#!/bin/bash

# 设置API基础URL
BASE_URL="http://localhost:8081/api/v1"

echo "=== MeiliSearch 设置接口测试 ==="
echo "基础URL: $BASE_URL"
echo ""

# 1. 获取当前设置
echo "1. 获取当前设置:"
curl -s -X GET "$BASE_URL/settings" | python -m json.tool
echo ""

# 2. 设置可搜索字段及其权重
echo "2. 设置可搜索字段及其权重:"
curl -s -X PUT "$BASE_URL/settings/searchable-attributes" \
  -H "Content-Type: application/json" \
  -d '{
    "searchable_attributes": ["title", "description", "genres"],
    "weights": {
      "title": 10,
      "description": 5
    }
  }' | python -m json.tool
echo ""

# 3. 设置可过滤字段
echo "3. 设置可过滤字段:"
curl -s -X PUT "$BASE_URL/settings/filterable-attributes" \
  -H "Content-Type: application/json" \
  -d '{
    "filterable_attributes": ["id", "genres", "release_date"]
  }' | python -m json.tool
echo ""

# 4. 设置可排序字段
echo "4. 设置可排序字段:"
curl -s -X PUT "$BASE_URL/settings/sortable-attributes" \
  -H "Content-Type: application/json" \
  -d '{
    "sortable_attributes": ["id", "release_date", "rating"]
  }' | python -m json.tool
echo ""

# 5. 更新排序规则
echo "5. 更新排序规则:"
curl -s -X PUT "$BASE_URL/settings/ranking-rules" \
  -H "Content-Type: application/json" \
  -d '{
    "ranking_rules": [
      "words",
      "typo",
      "proximity",
      "attribute",
      "sort",
      "exactness"
    ]
  }' | python -m json.tool
echo ""

# 6. 重置设置
echo "6. 重置设置:"
curl -s -X POST "$BASE_URL/settings/reset" | python -m json.tool
echo ""

echo "=== 设置接口测试完成 ==="
```

## 2. 使用 HTTPie 测试

```bash
# 获取设置
http GET http://localhost:8081/api/v1/settings

# 设置可搜索字段
http PUT http://localhost:8081/api/v1/settings/searchable-attributes \
  searchable_attributes:='["title", "description", "genres"]' \
  weights:='{"title": 10, "description": 5}'

# 设置可过滤字段
http PUT http://localhost:8081/api/v1/settings/filterable-attributes \
  filterable_attributes:='["id", "genres", "release_date"]'

# 设置可排序字段
http PUT http://localhost:8081/api/v1/settings/sortable-attributes \
  sortable_attributes:='["id", "release_date", "rating"]'

# 更新排序规则
http PUT http://localhost:8081/api/v1/settings/ranking-rules \
  ranking_rules:='["words", "typo", "proximity", "attribute", "sort", "exactness"]'

# 重置设置
http POST http://localhost:8081/api/v1/settings/reset
```

## 3. 完整的测试脚本（包含搜索和设置）

创建 `test_all.sh`：

```bash
#!/bin/bash

BASE_URL="http://localhost:8081/api/v1"

echo "=== 完整功能测试 ==="

# 健康检查
echo "1. 健康检查:"
curl -s -X GET "$BASE_URL/health" | python -m json.tool

# 索引信息
echo "2. 索引信息:"
curl -s -X GET "$BASE_URL/index" | python -m json.tool

# 获取当前设置
echo "3. 当前设置:"
curl -s -X GET "$BASE_URL/settings" | python -m json.tool

# 配置搜索设置
echo "4. 配置搜索设置:"
curl -s -X PUT "$BASE_URL/settings/searchable-attributes" \
  -H "Content-Type: application/json" \
  -d '{
    "searchable_attributes": ["title", "genres"],
    "weights": {"title": 10}
  }' | python -m json.tool

# 等待设置生效
sleep 2

# 测试搜索
echo "5. 测试搜索:"
curl -s -X GET "$BASE_URL/search?query=action" | python -m json.tool

echo "=== 测试完成 ==="
```

## 主要功能说明

1. **获取设置** (`GET /api/v1/settings`) - 查看当前索引的所有配置
2. **设置可搜索字段** (`PUT /api/v1/settings/searchable-attributes`) - 配置哪些字段参与搜索及其权重
3. **设置可过滤字段** (`PUT /api/v1/settings/filterable-attributes`) - 配置哪些字段可用于过滤
4. **设置可排序字段** (`PUT /api/v1/settings/sortable-attributes`) - 配置哪些字段可用于排序
5. **更新排序规则** (`PUT /api/v1/settings/ranking-rules`) - 配置搜索结果的排序优先级
6. **重置设置** (`POST /api/v1/settings/reset`) - 恢复所有设置为默认值
