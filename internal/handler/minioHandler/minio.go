package minioHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"minio-gin-crud/internal/common/dto"
	"minio-gin-crud/internal/common/errors"
	"minio-gin-crud/internal/common/responses"
	"minio-gin-crud/pkg/minio/helpers"
	"net/http"
)

// CreateOne обработчик для создания одного объекта в хранилище MinIO из переданных данных.
func (h *Handler) CreateOne(c *gin.Context) {
	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		// Если файл не получен, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "No file is received",
			Details: err,
		})
		return
	}

	// Открываем файл для чтения
	f, err := file.Open()
	if err != nil {
		// Если файл не удается открыть, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Unable to open the file",
			Details: err,
		})
		return
	}
	defer f.Close() // Закрываем файл после завершения работы с ним

	// Читаем содержимое файла в байтовый срез
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		// Если не удается прочитать содержимое файла, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Unable to read the file",
			Details: err,
		})
		return
	}

	// Создаем структуру FileDataType для хранения данных файла
	fileData := helpers.FileDataType{
		FileName: file.Filename, // Имя файла
		Data:     fileBytes,     // Содержимое файла в виде байтового среза
	}

	// Сохраняем файл в MinIO с помощью метода CreateOne
	link, err := h.minioService.CreateOne(fileData)
	if err != nil {
		// Если не удается сохранить файл, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Unable to save the file",
			Details: err,
		})
		return
	}

	// Возвращаем успешный ответ с URL-адресом сохраненного файла
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  http.StatusOK,
		Message: "File uploaded successfully",
		Data:    link, // URL-адрес загруженного файла
	})
}

// CreateMany обработчик для создания нескольких объектов в хранилище MinIO из переданных данных.
func (h *Handler) CreateMany(c *gin.Context) {
	// Получаем multipart форму из запроса
	form, err := c.MultipartForm()
	if err != nil {
		// Если форма недействительна, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "Invalid form",
			Details: err,
		})
		return
	}

	// Получаем файлы из формы
	files := form.File["files"]
	if files == nil {
		// Если файлы не получены, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "No files are received",
			Details: err,
		})
		return
	}

	// Создаем map для хранения данных файлов
	data := make(map[string]helpers.FileDataType)

	// Проходим по каждому файлу в форме
	for _, file := range files {
		// Открываем файл
		f, err := file.Open()
		if err != nil {
			// Если файл не удается открыть, возвращаем ошибку с соответствующим статусом и сообщением
			c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Error:   "Unable to open the file",
				Details: err,
			})
			return
		}
		defer f.Close() // Закрываем файл после завершения работы с ним

		// Читаем содержимое файла в байтовый срез
		fileBytes, err := io.ReadAll(f)
		if err != nil {
			// Если не удается прочитать содержимое файла, возвращаем ошибку с соответствующим статусом и сообщением
			c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Error:   "Unable to read the file",
				Details: err,
			})
			return
		}

		// Добавляем данные файла в map
		data[file.Filename] = helpers.FileDataType{
			FileName: file.Filename, // Имя файла
			Data:     fileBytes,     // Содержимое файла в виде байтового среза
		}
	}

	// Сохраняем файлы в MinIO с помощью метода CreateMany
	links, err := h.minioService.CreateMany(data)
	if err != nil {
		// Если не удается сохранить файлы, возвращаем ошибку с соответствующим статусом и сообщением
		fmt.Printf("err: %+v\n ", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Unable to save the files",
			Details: err,
		})
		return
	}

	// Возвращаем успешный ответ с URL-адресами сохраненных файлов
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Files uploaded successfully",
		Data:    links, // URL-адреса загруженных файлов
	})
}

// GetOne обработчик для получения одного объекта из бакета Minio по его идентификатору.
func (h *Handler) GetOne(c *gin.Context) {
	// Получаем идентификатор объекта из параметров URL
	objectID := c.Param("objectID")

	// Используем сервис MinIO для получения ссылки на объект
	link, err := h.minioService.GetOne(objectID)
	if err != nil {
		// Если произошла ошибка при получении объекта, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Enable to get the object",
			Details: err,
		})
		return
	}

	// Возвращаем успешный ответ с URL-адресом полученного файла
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  http.StatusOK,
		Message: "File received successfully",
		Data:    link, // URL-адрес полученного файла
	})
}

// GetMany обработчик для получения нескольких объектов из бакета Minio по их идентификаторам.
func (h *Handler) GetMany(c *gin.Context) {
	// Объявление переменной для хранения получаемых идентификаторов объектов
	var objectIDs dto.ObjectIdsDto

	// Привязка JSON данных из запроса к переменной objectIDs
	if err := c.ShouldBindJSON(&objectIDs); err != nil {
		// Если привязка данных не удалась, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "Invalid request body",
			Details: err,
		})
		return
	}

	// Используем сервис MinIO для получения ссылок на объекты по их идентификаторам
	links, err := h.minioService.GetMany(objectIDs.ObjectIDs)
	if err != nil {
		// Если произошла ошибка при получении объектов, возвращаем ошибку с соответствующим статусом и сообщением
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Enable to get many objects",
			Details: err,
		})
		return
	}

	// Возвращаем успешный ответ с URL-адресами полученных файлов
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Files received successfully",
		"data":    links, // URL-адреса полученных файлов
	})
}

// DeleteOne обработчик для удаления одного объекта из бакета Minio по его идентификатору.
func (h *Handler) DeleteOne(c *gin.Context) {
	objectID := c.Param("objectID")

	if err := h.minioService.DeleteOne(objectID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Cannot delete the object",
			Details: err,
		})
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  http.StatusOK,
		Message: "File deleted successfully",
	})
}

// DeleteMany обработчик для удаления нескольких объектов из бакета Minio по их идентификаторам.
func (h *Handler) DeleteMany(c *gin.Context) {
	var objectIDs dto.ObjectIdsDto
	if err := c.BindJSON(&objectIDs); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "Invalid request body",
			Details: err,
		})
		return
	}

	if err := h.minioService.DeleteMany(objectIDs.ObjectIDs); err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "Cannot delete many objects",
			Details: err,
		})
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  http.StatusOK,
		Message: "Files deleted successfully",
	})
}
