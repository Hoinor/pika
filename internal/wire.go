//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/dushixiang/pika/internal/config"
	"github.com/dushixiang/pika/internal/handler"
	"github.com/dushixiang/pika/internal/repo"
	"github.com/dushixiang/pika/internal/service"
	"github.com/dushixiang/pika/internal/websocket"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitializeApp 初始化应用
func InitializeApp(logger *zap.Logger, db *gorm.DB, cfg *config.AppConfig) (*AppComponents, error) {
	wire.Build(
		// Repositories
		repo.NewAgentRepo,
		repo.NewMetricRepo,
		repo.NewUserRepo,
		provideApiKeyRepo,
		provideAlertRepo,
		providePropertyRepo,

		// Services
		provideApiKeyService,
		provideAgentService,
		service.NewUserService,
		provideNotifier,
		providePropertyService,
		provideAlertService,

		// Providers for services with config
		provideAccountService,
		provideAgentHandler,
		provideApiKeyHandler,
		provideAlertHandler,
		providePropertyHandler,

		// WebSocket Manager
		websocket.NewManager,

		// Handlers
		handler.NewAccountHandler,
		handler.NewUserHandler,

		// App Components
		wire.Struct(new(AppComponents), "*"),
	)
	return nil, nil
}

// AppComponents 应用组件
type AppComponents struct {
	AccountHandler  *handler.AccountHandler
	AgentHandler    *handler.AgentHandler
	UserHandler     *handler.UserHandler
	ApiKeyHandler   *handler.ApiKeyHandler
	AlertHandler    *handler.AlertHandler
	PropertyHandler *handler.PropertyHandler
	AgentService    *service.AgentService
	UserService     *service.UserService
	AlertService    *service.AlertService
	PropertyService *service.PropertyService
	WSManager       *websocket.Manager
}
