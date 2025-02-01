package main

import (
	"fmt"
	"log"
	"wallet-test/internal/entities"

	"wallet-test/internal/config"
	"wallet-test/internal/controllers"
	"wallet-test/internal/db"
	"wallet-test/internal/repositories"
	"wallet-test/internal/server"
	"wallet-test/internal/services"
)

func main() {
	//Получаем конфиг
	cfg, err := config.LoadDbConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	// Подключаемся к БД
	pgDB, err := db.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка соединения: %v", err)
	}

	// Выполняем автоматическую миграцию для структуры Wallet
	if err := pgDB.AutoMigrate(&entities.Wallet{}); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	// Инициализация repository + service + handler
	repo := repositories.NewWalletRepository(pgDB)
	svc := services.NewWalletService(repo)
	controller := controllers.NewWalletController(svc)

	// Создаём роутер
	r := server.NewRouter(controller)

	// Запускаем сервер
	addr := fmt.Sprintf(":%d", cfg.AppPort)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
