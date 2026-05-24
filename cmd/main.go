package main

import (
	"context"
	"library-rest-api/internal/api"
	"library-rest-api/internal/db"
	"library-rest-api/internal/library"
	"log"
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
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env file")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("database url is not set in .env file")
	}

	pool, err := db.InitDB(ctx, connString)
	if err != nil {
		log.Fatal("unable to connect to database:", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(connString); err != nil {
		log.Fatal("migrations error: ", err)
	}

	// storage := library.NewMemoryStorage()
	storage := library.NewDatabaseStorage(pool)
	handlers := api.NewHTTPHandlers(storage)
	server := api.NewHTTPServer(handlers)
	address := "localhost:8008"

	if err := server.Start(address); err != nil {
		log.Fatal("failed to start server:", err)
	}

}
