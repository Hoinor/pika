package internal

import (
	"github.com/dushixiang/pika/internal/config"
	"github.com/dushixiang/pika/internal/handler"
	"github.com/dushixiang/pika/internal/repo"
	"github.com/dushixiang/pika/internal/service"
	"github.com/dushixiang/pika/internal/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// provideAccountService 提供AccountService
func provideAccountService(logger *zap.Logger, userService *service.UserService, cfg *config.AppConfig) *service.AccountService {
	return service.NewAccountService(logger, userService, cfg.JWT.Secret, cfg.JWT.ExpiresHours)
}

// provideApiKeyRepo 提供ApiKeyRepo
func provideApiKeyRepo(db *gorm.DB) *repo.ApiKeyRepo {
	return repo.NewApiKeyRepo(db)
}

// provideApiKeyService 提供ApiKeyService
func provideApiKeyService(logger *zap.Logger, apiKeyRepo *repo.ApiKeyRepo) *service.ApiKeyService {
	return service.NewApiKeyService(logger, apiKeyRepo)
}

// provideApiKeyHandler 提供ApiKeyHandler
func provideApiKeyHandler(logger *zap.Logger, apiKeyService *service.ApiKeyService) *handler.ApiKeyHandler {
	return handler.NewApiKeyHandler(logger, apiKeyService)
}

// provideAgentService 提供AgentService
func provideAgentService(logger *zap.Logger, agentRepo *repo.AgentRepo, metricRepo *repo.MetricRepo, apiKeyService *service.ApiKeyService) *service.AgentService {
	return service.NewAgentService(logger, agentRepo, metricRepo, apiKeyService)
}

// provideAgentHandler 提供AgentHandler
func provideAgentHandler(logger *zap.Logger, agentService *service.AgentService, wsManager *websocket.Manager, cfg *config.AppConfig) *handler.AgentHandler {
	return handler.NewAgentHandler(
		logger,
		agentService,
		wsManager,
	)
}

// provideAlertRepo 提供AlertRepo
func provideAlertRepo(db *gorm.DB) *repo.AlertRepo {
	return repo.NewAlertRepo(db)
}

// provideNotifier 提供Notifier
func provideNotifier(logger *zap.Logger) *service.Notifier {
	return service.NewNotifier(logger)
}

// providePropertyRepo 提供PropertyRepo
func providePropertyRepo(db *gorm.DB) *repo.PropertyRepo {
	return repo.NewPropertyRepo(db)
}

// providePropertyService 提供PropertyService
func providePropertyService(propertyRepo *repo.PropertyRepo, logger *zap.Logger) *service.PropertyService {
	return service.NewPropertyService(propertyRepo, logger)
}

// providePropertyHandler 提供PropertyHandler
func providePropertyHandler(logger *zap.Logger, propertyService *service.PropertyService, notifier *service.Notifier) *handler.PropertyHandler {
	return handler.NewPropertyHandler(logger, propertyService, notifier)
}

// provideAlertService 提供AlertService
func provideAlertService(alertRepo *repo.AlertRepo, agentRepo *repo.AgentRepo, propertyService *service.PropertyService, notifier *service.Notifier, logger *zap.Logger) *service.AlertService {
	return service.NewAlertService(alertRepo, agentRepo, propertyService, notifier, logger)
}

// provideAlertHandler 提供AlertHandler
func provideAlertHandler(logger *zap.Logger, alertService *service.AlertService) *handler.AlertHandler {
	return handler.NewAlertHandler(logger, alertService)
}
