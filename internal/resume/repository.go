package resume

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResumeRepository interface {
	Create(resume *Resume) error
	FindByID(id uuid.UUID) (*Resume, error)
	FindByUserID(userID uuid.UUID) ([]Resume, error)
	Delete(id uuid.UUID) error
	SaveAnalysis(result *AnalysisResult) error
}

type resumeRepository struct {
	db *gorm.DB
}

func NewResumeRepository(db *gorm.DB) ResumeRepository {
	return &resumeRepository{db}
}

func (r *resumeRepository) Create(resume *Resume) error {
	return r.db.Create(resume).Error
}

func (r *resumeRepository) FindByID(id uuid.UUID) (*Resume, error) {
	var resume Resume
	err := r.db.Preload("AnalysisResult").First(&resume, "id = ?", id).Error
	return &resume, err
}

func (r *resumeRepository) FindByUserID(userID uuid.UUID) ([]Resume, error) {
	var resumes []Resume
	err := r.db.Preload("AnalysisResult").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&resumes).Error
	return resumes, err
}

func (r *resumeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Resume{}, "id = ?", id).Error
}

func (r *resumeRepository) SaveAnalysis(result *AnalysisResult) error {
	return r.db.Create(result).Error
}
