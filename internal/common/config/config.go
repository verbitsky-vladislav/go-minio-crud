package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

// Config структура, обозначающая структуру .env файла
type Config struct {
	Port               string // Порт, на котором запускается сервер
	MinioEndpoint      string // Адрес конечной точки Minio
	BucketName         string // Название конкретного бакета в Minio
	MinioRootUser      string // Имя пользователя для доступа к Minio
	MinioRootPassword  string // Пароль для доступа к Minio
	MinioUseSSL        bool   // Переменная, отвечающая за доступ к Minio с SSL или без
	FileTimeExpiration int    // Время в часах
}

var AppConfig *Config

// LoadConfig загружает конфигурацию из файла .env
func LoadConfig() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Устанавливаем конфигурационные параметры
	AppConfig = &Config{
		Port: getEnv("PORT", "8080"),

		MinioEndpoint:      getEnv("MINIO_ENDPOINT", "localhost:9000"),
		BucketName:         getEnv("MINIO_BUCKET_NAME", "defaultBucket"),
		MinioRootUser:      getEnv("MINIO_ROOT_USER", "root"),
		MinioRootPassword:  getEnv("MINIO_ROOT_PASSWORD", "minio_password"),
		MinioUseSSL:        getEnvAsBool("MINIO_USE_SSL", false),
		FileTimeExpiration: getEnvAsInt("FILE_TIME_EXPIRATION", int(time.Second*60*24)),
	}
}

// getEnv считывает значение переменной окружения или возвращает значение по умолчанию, если переменная не установлена
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt считывает значение переменной окружения как целое число или возвращает значение по умолчанию, если переменная не установлена или не может быть преобразована в целое число
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// getEnvAsBool считывает значение переменной окружения как булево или возвращает значение по умолчанию, если переменная не установлена или не может быть преобразована в булево
func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
