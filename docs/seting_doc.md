# MeiliDog 配置文件文档

## 配置文件位置
默认配置文件路径：`config/config.toml`

## 配置文件示例

```toml
# MeiliDog 配置文件
# 服务器配置
[server]
address = "http://localhost:7700"  # MeiliSearch 服务器地址
api_key = "your-master-key-here"   # MeiliSearch API 密钥
timeout = "30s"                    # 请求超时时间

# 搜索配置
[search]

# 默认搜索参数
[search.default]
limit = 20                         # 默认每页显示数量
offset = 0                         # 默认偏移量

# 搜索优化配置
[search.optimization]
attributes_to_crop = ["content", "description"]  # 需要裁剪的属性
crop_length = 200                               # 裁剪长度
attributes_to_highlight = ["title", "content"]   # 需要高亮的属性
highlight_pre_tag = "<mark>"                    # 高亮开始标签
highlight_post_tag = "</mark>"                   # 高亮结束标签
attributes_to_retrieve = ["*"]                  # 需要返回的属性
attributes_to_search_on = ["title", "content"]   # 指定搜索字段
```

## 配置项详解

### 服务器配置 (`[server]`)

| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `address` | string | 是 | - | MeiliSearch 服务器地址，格式：`http://host:port` |
| `api_key` | string | 否 | - | MeiliSearch API 密钥，如果服务器设置了认证则需要 |
| `timeout` | duration | 否 | "30s" | HTTP 请求超时时间，支持单位：s(秒)、m(分钟)、h(小时) |

### 默认搜索参数 (`[search.default]`)

| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `limit` | int | 否 | 20 | 每页显示的搜索结果数量 |
| `offset` | int | 否 | 0 | 搜索结果偏移量 |

### 搜索优化配置 (`[search.optimization]`)

| 配置项 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| `attributes_to_crop` | []string | 否 | [] | 需要裁剪长文本的属性列表 |
| `crop_length` | int | 否 | 200 | 裁剪后的文本长度 |
| `attributes_to_highlight` | []string | 否 | [] | 需要高亮匹配关键词的属性列表 |
| `highlight_pre_tag` | string | 否 | "<mark>" | 高亮开始标签 |
| `highlight_post_tag` | string | 否 | "</mark>" | 高亮结束标签 |
| `attributes_to_retrieve` | []string | 否 | ["*"] | 指定返回的字段列表，["*"] 表示返回所有字段 |
| `attributes_to_search_on` | []string | 否 | [] | 指定搜索的字段列表，空数组表示搜索所有字段 |

## 完整配置示例

```toml
# MeiliDog 完整配置示例

[server]
# MeiliSearch 服务器配置
address = "http://localhost:7700"
api_key = "masterKey"
timeout = "30s"

[search]

[search.default]
# 分页配置
limit = 20
offset = 0

[search.optimization]
# 文本裁剪配置（适用于长文本内容）
attributes_to_crop = [
    "content",
    "description",
    "body"
]
crop_length = 200

# 高亮配置
attributes_to_highlight = [
    "title",
    "content",
    "description"
]
highlight_pre_tag = "<em class='highlight'>"
highlight_post_tag = "</em>"

# 字段控制
attributes_to_retrieve = [
    "id",
    "title",
    "content",
    "created_at",
    "updated_at"
]

# 搜索字段限制
attributes_to_search_on = [
    "title",
    "content",
    "tags"
]
```

## 环境变量覆盖

MeiliDog 支持通过环境变量覆盖配置文件中的设置：

| 环境变量 | 对应配置项 | 示例 |
|----------|------------|------|
| `MEILI_DOG_SERVER_ADDRESS` | `server.address` | `export MEILI_DOG_SERVER_ADDRESS=http://192.168.1.100:7700` |
| `MEILI_DOG_SERVER_API_KEY` | `server.api_key` | `export MEILI_DOG_SERVER_API_KEY=your-secret-key` |
| `MEILI_DOG_SERVER_TIMEOUT` | `server.timeout` | `export MEILI_DOG_SERVER_TIMEOUT=60s` |
| `MEILI_DOG_SEARCH_DEFAULT_LIMIT` | `search.default.limit` | `export MEILI_DOG_SEARCH_DEFAULT_LIMIT=50` |

## 配置验证

启动应用时会自动验证配置文件的正确性，如果配置有误会显示具体的错误信息。

## 注意事项

1. **路径配置**：确保 `address` 配置的 MeiliSearch 服务器可访问
2. **权限配置**：确保 `api_key` 具有足够的权限执行搜索操作
3. **性能优化**：根据数据量调整 `limit` 和超时时间
4. **内存考虑**：返回大量数据时注意内存使用情况

## 故障排除

如果遇到配置问题，可以：

1. 检查配置文件语法是否正确
2. 确认 MeiliSearch 服务器状态
3. 查看应用启动日志中的错误信息
4. 验证环境变量是否正确设置

这个配置文件文档提供了完整的配置说明和示例，可以根据实际需求进行调整。