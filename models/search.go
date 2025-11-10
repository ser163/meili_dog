package models

import (
	"time"
)

// SettingRequest 设置请求参数
type SettingRequest struct {
	SearchableAttributes []string       `json:"searchable_attributes,omitempty"`
	FilterableAttributes []string       `json:"filterable_attributes,omitempty"`
	SortableAttributes   []string       `json:"sortable_attributes,omitempty"`
	RankingRules         []string       `json:"ranking_rules,omitempty"`
	Weights              map[string]int `json:"weights,omitempty"` // 字段权重
}

// SettingResponse 设置响应
type SettingResponse struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message,omitempty"`
	TaskUID  int64       `json:"task_uid,omitempty"`
	Settings interface{} `json:"settings,omitempty"`
	Error    string      `json:"error,omitempty"`
}

// CurrentSettingsResponse 当前设置响应
type CurrentSettingsResponse struct {
	SearchableAttributes []string            `json:"searchable_attributes"`
	FilterableAttributes []string            `json:"filterable_attributes"`
	SortableAttributes   []string            `json:"sortable_attributes"`
	RankingRules         []string            `json:"ranking_rules"`
	DisplayedAttributes  []string            `json:"displayed_attributes"`
	StopWords            []string            `json:"stop_words"`
	Synonyms             map[string][]string `json:"synonyms"`
	TypoTolerance        interface{}         `json:"typo_tolerance"`
}

// SearchRequest 搜索请求参数
type SearchRequest struct {
	Query   string                 `form:"query" binding:"required"`
	Page    int                    `form:"page,default=1"`
	Limit   int                    `form:"limit,default=20"`
	Offset  int                    `form:"-"` // 不绑定查询参数，由程序计算
	Filters map[string]interface{} `form:"filters"`
	Sort    []string               `form:"sort"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Hits               []map[string]interface{} `json:"hits"`
	EstimatedTotalHits int64                    `json:"estimatedTotalHits"`
	Query              string                   `json:"query"`
	Limit              int                      `json:"limit"`
	Offset             int                      `json:"offset"`
	Page               int                      `json:"page"`
	TotalPages         int                      `json:"totalPages,omitempty"`
	ProcessingTimeMs   int                      `json:"processingTimeMs"`
	IndexUID           string                   `json:"indexUID"` // 返回使用的索引UID
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Address string        `toml:"address"`
	APIKey  string        `toml:"api_key"`
	Timeout time.Duration `toml:"timeout"`
}

// SearchConfig 搜索配置
type SearchConfig struct {
	Default      SearchDefault      `toml:"default"`
	Optimization SearchOptimization `toml:"optimization"`
	Fields       FieldConfig        `toml:"-"`
}

// SearchDefault 默认搜索参数
type SearchDefault struct {
	Limit  int `toml:"limit"`
	Offset int `toml:"offset"`
}

// SearchOptimization 搜索优化参数
type SearchOptimization struct {
	AttributesToCrop      []string `toml:"attributes_to_crop"`
	CropLength            int      `toml:"crop_length"`
	AttributesToHighlight []string `toml:"attributes_to_highlight"`
	HighlightPreTag       string   `toml:"highlight_pre_tag"`
	HighlightPostTag      string   `toml:"highlight_post_tag"`
	AttributesToRetrieve  []string `toml:"attributes_to_retrieve"`
	AttributesToSearchOn  []string `toml:"attributes_to_search_on"`
}

// FieldConfig 字段类型配置
type FieldConfig struct {
	StringFields  []string `toml:"string"`
	NumberFields  []string `toml:"number"`
	BooleanFields []string `toml:"boolean"`
	ArrayFields   []string `toml:"array"`
	ObjectFields  []string `toml:"object"`
}

// AppConfig 应用配置
type AppConfig struct {
	Server struct {
		Address   string `toml:"address"`
		APIKey    string `toml:"api_key"`
		LocalPort int64  `toml:"local_port"`
	} `toml:"server"`
	Search struct {
		IndexUID     string `toml:"index_uid"` // 添加索引UID配置
		Optimization struct {
			AttributesToCrop      []string `toml:"attributes_to_crop"`
			AttributesToHighlight []string `toml:"attributes_to_highlight"`
			AttributesToRetrieve  []string `toml:"attributes_to_retrieve"`
			AttributesToSearchOn  []string `toml:"attributes_to_search_on"`
			CropLength            int      `toml:"crop_length"`
			HighlightPreTag       string   `toml:"highlight_pre_tag"`
			HighlightPostTag      string   `toml:"highlight_post_tag"`
		} `toml:"optimization"`
	} `toml:"search"`
}
