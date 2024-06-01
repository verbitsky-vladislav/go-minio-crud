package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"minio-gin-crud/internal/common/config"
	"minio-gin-crud/pkg/minio/helpers"
)

// Client интерфейс для взаимодействия с Minio
type Client interface {
	InitMinio() error                                                             // Метод для инициализации подключения к Minio
	CreateOne(file helpers.FileDataType) (*helpers.CreatedObject, error)          // Метод для создания одного объекта в бакете Minio
	CreateMany(map[string]helpers.FileDataType) ([]*helpers.CreatedObject, error) // Метод для создания нескольких объектов в бакете Minio
	GetOne(objectID string) (string, error)                                       // Метод для получения одного объекта из бакета Minio
	GetMany(objectIDs []string) ([]string, error)                                 // Метод для получения нескольких объектов из бакета Minio
	DeleteOne(objectID string) error                                              // Метод для удаления одного объекта из бакета Minio
	DeleteMany(objectIDs []string) error                                          // Метод для удаления нескольких объектов из бакета Minio
}

// minioClient реализация интерфейса MinioClient
type minioClient struct {
	mc *minio.Client // Клиент Minio
}

// NewMinioClient создает новый экземпляр Minio Client
func NewMinioClient() Client {
	return &minioClient{} // Возвращает новый экземпляр minioClient с указанным именем бакета
}

// InitMinio подключается к Minio и создает бакет, если не существует
// Бакет - это контейнер для хранения объектов в Minio. Он представляет собой пространство имен, в котором можно хранить и организовывать файлы и папки.
func (m *minioClient) InitMinio() error {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	client, err := minio.New(config.AppConfig.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.MinioRootUser, config.AppConfig.MinioRootPassword, ""),
		Secure: config.AppConfig.MinioUseSSL,
	})
	if err != nil {
		return err
	}

	// Установка подключения Minio
	m.mc = client

	// Проверка наличия бакета и его создание, если не существует
	exists, err := m.mc.BucketExists(ctx, config.AppConfig.BucketName)
	if err != nil {
		return err
	}
	if !exists {
		err := m.mc.MakeBucket(ctx, config.AppConfig.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
