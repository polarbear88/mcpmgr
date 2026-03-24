package main

import (
	"context"
)

type App struct {
	ctx     context.Context
	service *AppService
}

func NewApp() *App {
	return &App{
		service: NewAppService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetAppState() (AppState, error) {
	return a.service.GetState()
}

func (a *App) SaveServer(input ServerInput) (AppState, error) {
	return a.service.SaveServer(input)
}

func (a *App) DeleteServer(id string) (AppState, error) {
	return a.service.DeleteServer(id)
}

func (a *App) ApplyToAllClients() (ApplyResult, error) {
	return a.service.ApplyToAllClients()
}

func (a *App) PreviewClientConfig(clientID string) (ClientConfigPreview, error) {
	return a.service.PreviewClientConfig(clientID)
}

func (a *App) PreviewAppConfig() (ClientConfigPreview, error) {
	return a.service.PreviewAppConfig()
}

func (a *App) EnableClient(clientID string) (AppState, error) {
	return a.service.EnableClient(clientID)
}

func (a *App) DisableClient(clientID string, restoreBackup bool) (AppState, error) {
	return a.service.DisableClient(clientID, restoreBackup)
}
