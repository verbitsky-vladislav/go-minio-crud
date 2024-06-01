package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"minio-gin-crud/internal/common/config"
	"minio-gin-crud/internal/handler"
	"minio-gin-crud/pkg/minio"
)

func main() {
	// Загрузка конфигурации из файла .env
	config.LoadConfig()

	// Инициализация соединения с Minio
	minioClient := minio.NewMinioClient()
	err := minioClient.InitMinio()
	if err != nil {
		log.Fatalf("Ошибка инициализации Minio: %v", err)
	}

	_, s := handler.NewHandler(
		minioClient,
	)

	// Инициализация маршрутизатора Gin
	router := gin.Default()

	s.RegisterRoutes(router)

	// Запуск сервера Gin
	port := config.AppConfig.Port // Мы берем порт из конфига
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера Gin: %v", err)
	}
}
