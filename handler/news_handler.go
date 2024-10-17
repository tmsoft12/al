package handler

import (
	"fmt"
	"os"
	"rr/domain"
	"rr/service"
	"rr/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NewsHandler struct {
	Service *service.NewsService
}

// Habar döretmek üçin funksiýa
func (h *NewsHandler) Create(c *fiber.Ctx) error {
	var news domain.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}

	// Täze faýl adyny döretmek
	newFileName := "news_" + time.Now().Format("20060102150405")
	imagePath, err := utils.UploadFile(c, "news", "uploads/news", newFileName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Surat ýüklenip bilinmedi"})
	}
	news.Image = imagePath

	// Habar döretmek
	if err := h.Service.Create(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	api := "api/admin"
	news.Image = fmt.Sprintf("http://localhost:5000/%s/%s", api, news.Image)
	return c.Status(fiber.StatusCreated).JSON(news)
}

// Habarlary sahypa boýunça getirmek
func (h *NewsHandler) GetPaginated(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	news, total, err := h.Service.GetPaginated(pageInt, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "newsler ýüklenip bilinmedi"})
	}

	api := "api/admin"
	// Her bir newsiň surat URL-ni düzetmek
	for i := range news {
		news[i].Image = fmt.Sprintf("http://localhost:5000/%s/%s", api, news[i].Image)
	}

	return c.JSON(fiber.Map{
		"data":  news,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

// ID boýunça news getirmek
func (h *NewsHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	news, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "news tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "news tapylmady"})
	}
	api := "api/admin"

	news.Image = fmt.Sprintf("http://localhost:5000/%s/%s", api, news.Image)
	return c.JSON(news)
}

// news pozmak üçin funksiýa
func (h *NewsHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	news, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "news tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "news tapylmady"})
	}

	// newsi pozmak
	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "news pozulyp bilinmedi"})
	}

	// Suraty pozmak
	if news.Image != "" {
		imagePath := news.Image
		if err := os.Remove(imagePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Surat pozulyp bilinmedi"})
		}
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "news üstünlikli pozuldy"})
}

func (h *NewsHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	news, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "news tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "news tapylmady"})
	}

	var updatedNews domain.News
	if err := c.BodyParser(&updatedNews); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	// Täze surat ýüklenipmi?
	if file, err := c.FormFile("news"); err == nil && file != nil {
		// Köne suraty pozmak
		if news.Image != "" {
			if err := os.Remove(news.Image); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Köne surat pozulyp bilinmedi"})
			}
		}
		// Täze surat ýüklemek
		newFileName := fmt.Sprintf("newsUpdate_%s", time.Now().Format("20060102150405"))
		imagePath, err := utils.UploadFile(c, "news", "uploads/news", newFileName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze surat ýüklenip bilinmedi"})
		}
		updatedNews.Image = imagePath
	} else {
		// Täze surat ýüklenmedik bolsa, köne suraty sakla
		updatedNews.Image = news.Image
	}

	updatedNewsResult, err := h.Service.Update(uint(id), &updatedNews)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Habarlar tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Habarlar üýtgedilip bilinmedi"})
	}

	// Surat URL-ni düzetmek
	api := "api/admin"
	updatedNewsResult.Image = fmt.Sprintf("http://localhost:5000/%s/%s", api, updatedNewsResult.Image)

	return c.JSON(updatedNewsResult)
}
