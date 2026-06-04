package resume

import (
	"io"

	"github.com/moxhyusuf/ai-resume-analyzer/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ResumeHandler struct {
	resumeService ResumeService
}

func NewResumeHandler(resumeService ResumeService) *ResumeHandler {
	return &ResumeHandler{resumeService}
}

// Analyze godoc
// @Summary      Analyze resume
// @Description  Upload PDF resume untuk dianalisis oleh AI
// @Tags         Resume
// @Accept       multipart/form-data
// @Produce      json
// @Param        resume formData file true "File PDF resume"
// @Success      201 {object} response.APIResponse{data=AnalyzeResponse}
// @Failure      400 {object} response.APIResponse
// @Failure      422 {object} response.APIResponse
// @Security     BearerAuth
// @Router       /resumes/analyze [post]
func (h *ResumeHandler) Analyze(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	file, err := c.FormFile("resume")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Resume file is required (field: resume)")
	}

	if file.Header.Get("Content-Type") != "application/pdf" {
		return response.Error(c, fiber.StatusBadRequest, "Only PDF files are accepted")
	}

	if file.Size > 5*1024*1024 {
		return response.Error(c, fiber.StatusBadRequest, "File size must not exceed 5MB")
	}

	f, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to read uploaded file")
	}

	result, err := h.resumeService.Analyze(c.Context(), userID, fileData, file.Filename)
	if err != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Resume analyzed successfully", result)
}

// GetAll godoc
// @Summary      List semua resume
// @Description  Mendapatkan semua resume milik user yang login
// @Tags         Resume
// @Produce      json
// @Success      200 {object} response.APIResponse{data=[]resume.Resume}
// @Security     BearerAuth
// @Router       /resumes [get]
func (h *ResumeHandler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	resumes, err := h.resumeService.GetAll(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch resumes")
	}

	return response.Success(c, fiber.StatusOK, "Resumes retrieved", resumes)
}

// GetByID godoc
// @Summary      Detail resume
// @Description  Mendapatkan detail resume beserta hasil analisis
// @Tags         Resume
// @Produce      json
// @Param        id path string true "Resume ID"
// @Success      200 {object} response.APIResponse{data=resume.Resume}
// @Failure      404 {object} response.APIResponse
// @Security     BearerAuth
// @Router       /resumes/{id} [get]
func (h *ResumeHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid resume ID")
	}

	resume, err := h.resumeService.GetByID(c.Context(), userID, resumeID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Resume retrieved", resume)
}

// Delete godoc
// @Summary      Hapus resume
// @Tags         Resume
// @Produce      json
// @Param        id path string true "Resume ID"
// @Success      200 {object} response.APIResponse
// @Failure      403 {object} response.APIResponse
// @Security     BearerAuth
// @Router       /resumes/{id} [delete]
func (h *ResumeHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	resumeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid resume ID")
	}

	if err := h.resumeService.Delete(c.Context(), userID, resumeID); err != nil {
		return response.Error(c, fiber.StatusForbidden, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Resume deleted", nil)
}
