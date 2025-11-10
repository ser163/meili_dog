package handlers

import (
	"log"
	"net/http"
	"strconv"

	"meili_dog/models"

	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	client   *meilisearch.Client
	config   models.AppConfig
	indexUID string
}

// NewSearchHandler 创建新的搜索处理器
func NewSearchHandler(config models.AppConfig) *SearchHandler {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   config.Server.Address,
		APIKey: config.Server.APIKey,
	})

	return &SearchHandler{
		client:   client,
		config:   config,
		indexUID: config.Search.IndexUID,
	}
}

// Search 执行搜索
func (h *SearchHandler) Search(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查索引UID是否配置
	if h.indexUID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "索引未配置"})
		return
	}

	// 计算偏移量
	if req.Page < 1 {
		req.Page = 1
	}
	req.Offset = (req.Page - 1) * req.Limit

	// 构建搜索参数
	searchRequest := &meilisearch.SearchRequest{
		Query:  req.Query,
		Limit:  int64(req.Limit),
		Offset: int64(req.Offset),
	}

	// 应用优化参数
	h.applyOptimizationParams(searchRequest)

	// 应用过滤条件
	h.applyFilters(searchRequest, req.Filters)

	// 应用排序
	if len(req.Sort) > 0 {
		searchRequest.Sort = req.Sort
	}

	// 执行搜索
	result, err := h.client.Index(h.indexUID).Search(req.Query, searchRequest)
	if err != nil {
		log.Printf("搜索错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败: " + err.Error()})
		return
	}

	// 构建响应
	response := models.SearchResponse{
		Hits:               h.convertHits(result.Hits),
		EstimatedTotalHits: result.EstimatedTotalHits,
		Query:              req.Query,
		Limit:              req.Limit,
		Offset:             req.Offset,
		Page:               req.Page,
		ProcessingTimeMs:   int(result.ProcessingTimeMs),
		IndexUID:           h.indexUID,
	}

	// 计算总页数
	if result.EstimatedTotalHits > 0 {
		response.TotalPages = int((result.EstimatedTotalHits + int64(req.Limit) - 1) / int64(req.Limit))
	}

	c.JSON(http.StatusOK, response)
}

// applyOptimizationParams 应用优化参数
func (h *SearchHandler) applyOptimizationParams(req *meilisearch.SearchRequest) {
	opt := h.config.Search.Optimization

	if len(opt.AttributesToCrop) > 0 {
		req.AttributesToCrop = opt.AttributesToCrop
		req.CropLength = int64(opt.CropLength)
	}

	if len(opt.AttributesToHighlight) > 0 {
		req.AttributesToHighlight = opt.AttributesToHighlight
		req.HighlightPreTag = opt.HighlightPreTag
		req.HighlightPostTag = opt.HighlightPostTag
	}

	if len(opt.AttributesToRetrieve) > 0 && opt.AttributesToRetrieve[0] != "*" {
		req.AttributesToRetrieve = opt.AttributesToRetrieve
	}

	if len(opt.AttributesToSearchOn) > 0 {
		req.AttributesToSearchOn = opt.AttributesToSearchOn
	}
}

// applyFilters 应用过滤条件
func (h *SearchHandler) applyFilters(req *meilisearch.SearchRequest, filters map[string]interface{}) {
	if filters == nil || len(filters) == 0 {
		return
	}

	filterStrings := h.buildFilterStrings(filters)
	if len(filterStrings) > 0 {
		req.Filter = filterStrings
	}
}

// buildFilterStrings 构建过滤字符串
func (h *SearchHandler) buildFilterStrings(filters map[string]interface{}) []string {
	var filterStrings []string

	for field, value := range filters {
		switch v := value.(type) {
		case string:
			filterStrings = append(filterStrings, field+" = '"+v+"'")
		case int, int32, int64, float32, float64:
			filterStrings = append(filterStrings, field+" = "+toString(v))
		case bool:
			filterStrings = append(filterStrings, field+" = "+strconv.FormatBool(v))
		case []interface{}:
			if len(v) > 0 {
				filterStrings = append(filterStrings, h.buildArrayFilter(field, v))
			}
		}
	}

	return filterStrings
}

// buildArrayFilter 构建数组过滤条件
func (h *SearchHandler) buildArrayFilter(field string, values []interface{}) string {
	var conditions []string
	for _, val := range values {
		switch v := val.(type) {
		case string:
			conditions = append(conditions, field+" = '"+v+"'")
		default:
			conditions = append(conditions, field+" = "+toString(v))
		}
	}
	return "(" + joinStrings(conditions, " OR ") + ")"
}

// convertHits 转换命中结果
func (h *SearchHandler) convertHits(hits []interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, hit := range hits {
		if hitMap, ok := hit.(map[string]interface{}); ok {
			result = append(result, hitMap)
		}
	}
	return result
}

// 辅助函数
func toString(v interface{}) string {
	switch val := v.(type) {
	case int:
		return strconv.Itoa(val)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return ""
	}
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}

// HealthCheck 健康检查端点
func (h *SearchHandler) HealthCheck(c *gin.Context) {
	_, err := h.client.Health()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

// GetIndexInfo 获取当前配置的索引信息
func (h *SearchHandler) GetIndexInfo(c *gin.Context) {
	if h.indexUID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "索引未配置"})
		return
	}

	index := h.client.Index(h.indexUID)
	stats, err := index.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取索引统计信息失败: " + err.Error()})
		return
	}

	// 获取索引设置
	settings, err := index.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取索引设置失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"index_uid": h.indexUID,
		"stats":     stats,
		"settings":  settings,
	})
}

// GetSettings 获取当前索引的所有设置
func (h *SearchHandler) GetSettings(c *gin.Context) {
	if h.indexUID == "" {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "索引未配置",
		})
		return
	}

	index := h.client.Index(h.indexUID)

	// 获取各种设置
	searchableAttrs, err := index.GetSearchableAttributes()
	if err != nil {
		log.Printf("获取可搜索字段错误: %v", err)
	}

	filterableAttrs, err := index.GetFilterableAttributes()
	if err != nil {
		log.Printf("获取可过滤字段错误: %v", err)
	}

	sortableAttrs, err := index.GetSortableAttributes()
	if err != nil {
		log.Printf("获取可排序字段错误: %v", err)
	}

	rankingRules, err := index.GetRankingRules()
	if err != nil {
		log.Printf("获取排序规则错误: %v", err)
	}

	displayedAttrs, err := index.GetDisplayedAttributes()
	if err != nil {
		log.Printf("获取显示字段错误: %v", err)
	}

	stopWords, err := index.GetStopWords()
	if err != nil {
		log.Printf("获取停用词错误: %v", err)
	}

	synonyms, err := index.GetSynonyms()
	if err != nil {
		log.Printf("获取同义词错误: %v", err)
	}

	typoTolerance, err := index.GetTypoTolerance()
	if err != nil {
		log.Printf("获取拼写容错设置错误: %v", err)
	}

	// 创建响应对象，处理指针类型
	response := models.CurrentSettingsResponse{
		SearchableAttributes: derefStringSlice(searchableAttrs),
		FilterableAttributes: derefStringSlice(filterableAttrs),
		SortableAttributes:   derefStringSlice(sortableAttrs),
		RankingRules:         derefStringSlice(rankingRules),
		DisplayedAttributes:  derefStringSlice(displayedAttrs),
		StopWords:            derefStringSlice(stopWords),
		Synonyms:             derefStringMap(synonyms),
		TypoTolerance:        typoTolerance,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSearchableAttributes 设置可搜索字段及其权重
func (h *SearchHandler) UpdateSearchableAttributes(c *gin.Context) {
	var req models.SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "请求参数错误: " + err.Error(),
		})
		return
	}

	if h.indexUID == "" {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "索引未配置",
		})
		return
	}

	index := h.client.Index(h.indexUID)

	// 如果有权重信息，需要处理字段权重
	searchableAttrs := req.SearchableAttributes
	if req.Weights != nil && len(req.Weights) > 0 {
		searchableAttrs = h.applyWeightsToAttributes(req.SearchableAttributes, req.Weights)
	}

	task, err := index.UpdateSearchableAttributes(&searchableAttrs)
	if err != nil {
		log.Printf("更新可搜索字段错误: %v", err)
		c.JSON(http.StatusInternalServerError, models.SettingResponse{
			Success: false,
			Error:   "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SettingResponse{
		Success: true,
		Message: "可搜索字段更新成功",
		TaskUID: task.TaskUID,
	})
}

// UpdateFilterableAttributes 设置可过滤字段
func (h *SearchHandler) UpdateFilterableAttributes(c *gin.Context) {
	var req models.SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "请求参数错误: " + err.Error(),
		})
		return
	}

	if h.indexUID == "" {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "索引未配置",
		})
		return
	}

	index := h.client.Index(h.indexUID)

	task, err := index.UpdateFilterableAttributes(&req.FilterableAttributes)
	if err != nil {
		log.Printf("更新可过滤字段错误: %v", err)
		c.JSON(http.StatusInternalServerError, models.SettingResponse{
			Success: false,
			Error:   "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SettingResponse{
		Success: true,
		Message: "可过滤字段更新成功",
		TaskUID: task.TaskUID,
	})
}

// UpdateSortableAttributes 设置可排序字段
func (h *SearchHandler) UpdateSortableAttributes(c *gin.Context) {
	var req models.SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "请求参数错误: " + err.Error(),
		})
		return
	}

	if h.indexUID == "" {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "索引未配置",
		})
		return
	}

	index := h.client.Index(h.indexUID)

	task, err := index.UpdateSortableAttributes(&req.SortableAttributes)
	if err != nil {
		log.Printf("更新可排序字段错误: %v", err)
		c.JSON(http.StatusInternalServerError, models.SettingResponse{
			Success: false,
			Error:   "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SettingResponse{
		Success: true,
		Message: "可排序字段更新成功",
		TaskUID: task.TaskUID,
	})
}

// UpdateRankingRules 更新排序规则
func (h *SearchHandler) UpdateRankingRules(c *gin.Context) {
	var req models.SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "请求参数错误: " + err.Error(),
		})
		return
	}

	if h.indexUID == "" {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "索引未配置",
		})
		return
	}

	index := h.client.Index(h.indexUID)

	task, err := index.UpdateRankingRules(&req.RankingRules)
	if err != nil {
		log.Printf("更新排序规则错误: %v", err)
		c.JSON(http.StatusInternalServerError, models.SettingResponse{
			Success: false,
			Error:   "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SettingResponse{
		Success: true,
		Message: "排序规则更新成功",
		TaskUID: task.TaskUID,
	})
}

// ResetSettings 重置所有设置为默认值
func (h *SearchHandler) ResetSettings(c *gin.Context) {
	if h.indexUID == "" {
		c.JSON(http.StatusBadRequest, models.SettingResponse{
			Success: false,
			Error:   "索引未配置",
		})
		return
	}

	index := h.client.Index(h.indexUID)

	task, err := index.ResetSettings()
	if err != nil {
		log.Printf("重置设置错误: %v", err)
		c.JSON(http.StatusInternalServerError, models.SettingResponse{
			Success: false,
			Error:   "重置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SettingResponse{
		Success: true,
		Message: "设置重置成功",
		TaskUID: task.TaskUID,
	})
}

// 辅助方法：应用权重到字段
func (h *SearchHandler) applyWeightsToAttributes(attributes []string, weights map[string]int) []string {
	var weightedAttributes []string

	for _, attr := range attributes {
		if weight, exists := weights[attr]; exists {
			weightedAttributes = append(weightedAttributes, attr+":"+strconv.Itoa(weight))
		} else {
			weightedAttributes = append(weightedAttributes, attr)
		}
	}

	return weightedAttributes
}

// 辅助函数：解引用字符串切片指针
func derefStringSlice(ptr *[]string) []string {
	if ptr == nil {
		return []string{}
	}
	return *ptr
}

// 辅助函数：解引用字符串映射指针
func derefStringMap(ptr *map[string][]string) map[string][]string {
	if ptr == nil {
		return map[string][]string{}
	}
	return *ptr
}
