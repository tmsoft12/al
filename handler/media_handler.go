package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"rr/domain"
	"rr/service"
	"rr/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	videoPath = "uploads/media/video/"
	coverPath = "uploads/media/cover/"
	apiBase   = "api/admin"
)

type MediaHandler struct {
	Service *service.MediaService
}

func (h *MediaHandler) Create(c *fiber.Ctx) error {
	var media domain.Media
	if err := c.BodyParser(&media); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Serwerde ýalňyşlyk: Maglumatlar işlenip bilinmedi"})
	}

	// Täze faýl adyny döretmek
	newCover := "cover_" + time.Now().Format("20060102150405")
	coverPath, err := utils.UploadFile(c, "cover", coverPath, newCover)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cover ýüklenip bilinmedi"})
	}
	media.Cover = coverPath

	newFileName := "video_" + time.Now().Format("20060102150405")
	videoPath, err := utils.UploadFile(c, "video", videoPath, newFileName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Media ýüklenip bilinmedi"})
	}

	videoExt := filepath.Ext(videoPath)
	media.Video = newFileName + videoExt
	coverExt := filepath.Ext(coverPath)
	media.Cover = newCover + coverExt

	if err := h.Service.Create(&media); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	api := "api/admin"
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	media.Cover = fmt.Sprintf("http://%s:%s/%s/uploads/media/cover/%s", ip, port, api, media.Cover)
	media.Video = fmt.Sprintf("http://%s:%s/%s/media/video/%s", ip, port, api, media.Video)

	return c.Status(fiber.StatusCreated).JSON(media)
}

// Media sahypa boýunça getirmek
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
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	for i := range media {
		media[i].Cover = fmt.Sprintf("http://%s:%s/%s/%s", ip, port, api, media[i].Cover)
	}
	for i := range media {
		media[i].Video = fmt.Sprintf("http://%s:%s/%s/media/video/%s", ip, port, api, media[i].Video)
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"})
	}

	media, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media tapylmady"})
	}

	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	if ip == "" || port == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sazlamalar ýalňyş"})
	}

	media.Video = fmt.Sprintf("http://%s:%s/%s/media/video/%s", ip, port, apiBase, media.Video)
	media.Cover = fmt.Sprintf("http://%s:%s/%s/uploads/media/cover/%s", ip, port, apiBase, media.Cover)

	return c.JSON(media)
}

// Isgar pozmak üçin funksiýa
func (h *MediaHandler) Delete(c *fiber.Ctx) error {
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

	// Media pozmak
	if err := h.Service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media pozulyp bilinmedi"})
	}

	// Suraty pozmak
	if media.Video != "" {
		videoFilePath := filepath.Join(videoPath, media.Video)
		if err := os.Remove(videoFilePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Video pozulyp bilinmedi"})
		}
	}
	if media.Cover != "" {
		coverFilePath := filepath.Join(coverPath, media.Cover)
		if err := os.Remove(coverFilePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cover pozulyp bilinmedi"})
		}
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "Media üstünlikli pozuldy"})
}

// Isgar üýtgetmek üçin funksiýa
// Isgar üýtgetmek üçin funksiýa
func (h *MediaHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry ID"}) // Invalid ID error
	}

	media, err := h.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"}) // Media not found
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Mediani almakda ýalňyşlyk ýüze çykdy"}) // Error retrieving media
	}

	var updatedMedia domain.Media
	if err := c.BodyParser(&updatedMedia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nädogry maglumatlar"}) // Invalid data error
	}

	// Check if new video file is being uploaded
	if media.Video != "" {
		// Delete the old video file
		videoFilePath := filepath.Join(videoPath, media.Video)
		if err := utils.DeleteFileWithRetry(videoFilePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Wideo pozulyp bilinmedi: %v", err)}) // Error deleting video
		}

		// Create new video file name and upload the new video
		newVideoName := fmt.Sprintf("mediaUpdate_%s", time.Now().Format("20060102150405"))
		newVideoPath, err := utils.UploadFile(c, "video", videoPath, newVideoName) // Uploading new video
		if err != nil {
			fmt.Printf("Video upload error: %v\n", err)                                                       // Log specific upload error
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze video ýüklenip bilinmedi"}) // New video upload failed
		}
		updatedMedia.Video = newVideoName + filepath.Ext(newVideoPath) // Store only the filename with extension
	} else {
		updatedMedia.Video = media.Video // Keep the old video if not updated
	}

	// Check if new cover file is being uploaded
	if media.Cover != "" {
		// Delete the old cover file
		coverFilePath := filepath.Join(coverPath, media.Cover)
		if err := utils.DeleteFileWithRetry(coverFilePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Surat pozulyp bilinmedi: %v", err)}) // Error deleting cover
		}

		// Create new cover file name and upload the new cover
		newCoverName := fmt.Sprintf("mediaCoverUpdate_%s", time.Now().Format("20060102150405"))
		newCoverPath, err := utils.UploadFile(c, "cover", coverPath, newCoverName) // Uploading new cover
		if err != nil {
			fmt.Printf("Cover upload error: %v\n", err)                                                       // Log specific upload error
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Täze cover ýüklenip bilinmedi"}) // New cover upload failed
		}
		updatedMedia.Cover = newCoverName + filepath.Ext(newCoverPath) // Store only the filename with extension
	} else {
		updatedMedia.Cover = media.Cover // Keep the old cover if not updated
	}

	// Update the media record in the database
	updatedMediaResult, err := h.Service.Update(uint(id), &updatedMedia)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media tapylmady"}) // Media not found
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Media üýtgedilip bilinmedi"}) // Media update failed
	}

	// Construct the full URLs for the updated media
	ip := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	if ip == "" || port == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Sazlamalar ýalňyş"}) // Configuration error
	}

	updatedMediaResult.Cover = fmt.Sprintf("http://%s:%s/%s/uploads/media/cover/%s", ip, port, apiBase, updatedMediaResult.Cover)
	updatedMediaResult.Video = fmt.Sprintf("http://%s:%s/%s/media/video/%s", ip, port, apiBase, updatedMediaResult.Video)

	return c.JSON(updatedMediaResult) // Return the updated media
}
