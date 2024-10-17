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

type MediaHandler struct {
	Service *service.MediaService
}

// Employer döretmek üçin funksiýa

func (h *MediaHandler) Create(c *fiber.Ctx) error {
	var media domain.Media
	if err := c.BodyParser(&media); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}
	// Täze faýl adyny döretmek
	newCover := "cover_" + time.Now().Format("20060102150405")
	coverPath, err := utils.UploadFile(c, "cover", "uploads/media/cover", newCover)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cover ýüklenip bilinmedi"})
	}
	media.Cover = coverPath

	newFileName := "video_" + time.Now().Format("20060102150405")
	videoPath, err := utils.UploadFile(c, "video", "uploads/media/video", newFileName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Media ýüklenip bilinmedi"})
	}
	media.Video = videoPath

	if err := h.Service.Create(&media); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	api := "api/admin"
	media.Video = fmt.Sprintf("http://localhost:5000/%s/%s", api, media.Video)
	media.Cover = fmt.Sprintf("http://localhost:5000/%s/%s", api, media.Cover)

	return c.Status(fiber.StatusCreated).JSON(media)
}

// Mediai sahypa boýunça getirmek
func (h *MediaHandler) GetPaginated(c *fiber.Ctx) error {
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

	media, total, err := h.Service.GetPaginated(pageInt, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media ýüklenip bilinmedi"})
	}
	api := "api/admin"
	// Her bir isgarin surat URL-ni düzetmek
	for i := range media {
		media[i].Video = fmt.Sprintf("http://localhost:5000/%s/%s", api, media[i].Video)
	}
	for i := range media {
		media[i].Cover = fmt.Sprintf("http://localhost:5000/%s/%s", api, media[i].Cover)
	}

	return c.JSON(fiber.Map{
		"data":  media,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

// ID boýunça isgar getirmek
func (h *MediaHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nädogry ID"},
		)
	}

	media, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media tapylmady"})
	}
	api := "api/admin"
	media.Video = fmt.Sprintf("http://localhost:5000/%s/%s", api, media.Video)
	media.Cover = fmt.Sprintf("http://localhost:5000/%s/%s", api, media.Cover)
	return c.JSON(media)
}

// Isgar pozmak üçin funksiýa
func (h *MediaHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	medai, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media tapylmady"})
	}

	// Media pozmak
	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media pozulyp bilinmedi"})
	}

	// Suraty pozmak
	if medai.Video != "" {
		videoPath := medai.Video

		if err := os.Remove(videoPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Video pozulyp bilinmedi"})
		}
	}
	if medai.Cover != "" {
		coverPath := medai.Cover
		if err := os.Remove(coverPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cover pozulyp bilinmedi"})
		}
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Media üstünlikli pozuldy"})
}

// Isgar üýtgetmek üçin funksiýa
func (h *MediaHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	media, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media tapylmady"})
	}

	var updatedMedia domain.Media
	if err := c.BodyParser(&updatedMedia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"})
	}

	// Täze video faýl ýüklenipmi?
	if file, err := c.FormFile("video"); err == nil && file != nil {
		if media.Video != "" {
			if err := os.Remove(media.Video); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Köne video pozulyp bilinmedi"})
			}
		}
		newFileName := fmt.Sprintf("mediaUpdate_%s", time.Now().Format("20060102150405"))
		videoPath, err := utils.UploadFile(c, "video", "uploads/media/video", newFileName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze video ýüklenip bilinmedi"})
		}
		updatedMedia.Video = videoPath
	} else {
		updatedMedia.Video = media.Video
	}

	// Täze cover faýl ýüklenipmi?
	if file, err := c.FormFile("cover"); err == nil && file != nil {
		if media.Cover != "" {
			if err := os.Remove(media.Cover); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Köne cover pozulyp bilinmedi"})
			}
		}
		newFileName := fmt.Sprintf("mediaCoverUpdate_%s", time.Now().Format("20060102150405"))
		coverPath, err := utils.UploadFile(c, "cover", "uploads/media/cover", newFileName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze cover ýüklenip bilinmedi"})
		}
		updatedMedia.Cover = coverPath
	} else {
		updatedMedia.Cover = media.Cover
	}

	updatedMediaResult, err := h.Service.Update(uint(id), &updatedMedia)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media üýtgedilip bilinmedi"})
	}
	api := "api/admin"
	updatedMediaResult.Video = fmt.Sprintf("http://localhost:5000/%s/%s", api, media.Video)
	return c.JSON(updatedMediaResult)
}
