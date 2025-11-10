package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"meili_dog/config"
	"meili_dog/handlers"
	"meili_dog/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 检查索引配置
	if cfg.Search.IndexUID == "" {
		log.Fatal("配置文件中未设置索引UID (search.index_uid)")
	}

	log.Printf("配置索引: %s", cfg.Search.IndexUID)

	// 初始化搜索处理器
	searchHandler := handlers.NewSearchHandler(*cfg)

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.Default()

	// 添加CORS中间件（可选）
	router.Use(corsMiddleware())

	// 路由定义
	api := router.Group("/api/v1")
	{
		api.GET("/health", searchHandler.HealthCheck)
		api.GET("/index", searchHandler.GetIndexInfo) // 改为单数，获取当前索引信息
		api.GET("/search", searchHandler.Search)

		settings := api.Group("/settings")
		{
			settings.GET("/", searchHandler.GetSettings)                                     // 获取所有设置
			settings.PUT("/searchable-attributes", searchHandler.UpdateSearchableAttributes) // 设置可搜索字段
			settings.PUT("/filterable-attributes", searchHandler.UpdateFilterableAttributes) // 设置可过滤字段
			settings.PUT("/sortable-attributes", searchHandler.UpdateSortableAttributes)     // 设置可排序字段
			settings.PUT("/ranking-rules", searchHandler.UpdateRankingRules)                 // 更新排序规则
			settings.POST("/reset", searchHandler.ResetSettings)                             // 重置所有设置
		}
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.FormatInt(cfg.Server.LocalPort, 10)
	}

	log.Printf("服务器启动在端口 %s", port)
	log.Printf("MeiliSearch 地址: %s", cfg.Server.Address)
	log.Printf("目标索引: %s", cfg.Search.IndexUID)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// loadConfig 加载配置
func loadConfig() (*models.AppConfig, error) {
	// 获取配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.toml"
	}

	// 确保路径是绝对路径
	if !filepath.IsAbs(configPath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(wd, configPath)
	}

	return config.LoadConfig(configPath)
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
