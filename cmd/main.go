package main

import (
	"context"
	"library-rest-api/internal/api"
	"library-rest-api/internal/db"
	"library-rest-api/internal/library"
	"library-rest-api/internal/logger"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

/*
— В этой библиотеке хранятся книги, каждая книга имеет:
	1. Название
	2. Автора
	З. Количество страниц
	4. Информацию о том, прочитана эта книга нами, либо же нет
	5. Время добавления в библиотеку
	6. Время окончательного прочтения (когда книга была дочитана нами до конца)
— Приложение предоставляет возможность:
	1. Добавлять новые книги в нашу личную библиотеку
	2. Отмечать отдельные книги как прочитанные
	З. Получать информацию о какой-то конкретной книге
	4. Получать список всех книг, с учётом возможной фильтрации по:
	   автору, прочитано/не прочитано
	5. Удалять книги из нашей библиотеки
— Взаимодействие с приложением происходит по протоколу НТТР,
  контракт приложения должен удовлетворять принципам REST API
*/

func main() {
	appLogger, file, err := logger.New()
	if err != nil {
		log.Fatal("failed to create logger:", err)
	}
	defer file.Close()
	slog.SetDefault(appLogger)

	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system env file")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		slog.Error("database url is not set in .env file")
		os.Exit(1)
	}
	pool, err := db.InitDB(ctx, connString)
	if err != nil {
		slog.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.RunMigrations(connString); err != nil {
		slog.Error("migrations error", "error", err)
		os.Exit(1)
	}

	// storage := library.NewMemoryStorage()
	storage := library.NewDatabaseStorage(pool)
	handlers := api.NewHTTPHandlers(storage)
	server := api.NewHTTPServer(handlers, appLogger)
	address := "localhost:8008"
	slog.Info("starting server", "address", address)

	if err := server.Start(address); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
