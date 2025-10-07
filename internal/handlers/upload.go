package handlers

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Загрузить файл
// @Description Загружает файл на сервер и возвращает URL
// @Tags upload
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Router /upload [post]

func (h *Handler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		respondErr(c, 400, "file required")
		return
	}
	exts := filepath.Ext(file.Filename)
	name := fmt.Sprintf("%d_%d%s", uid(c), time.Now().UnixNano(), exts)
	path := filepath.Join(h.uploadDir, name)
	if err := c.SaveUploadedFile(file, path); err != nil {
		respondErr(c, 500, "save failed")
		return
	}
	url := fmt.Sprintf("%s/uploads/%s", h.staticBase, name)
	c.JSON(201, gin.H{"url": url})
}
