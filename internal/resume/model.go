package resume

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Resume struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	OriginalName   string          `json:"original_name"`
	ExtractedText  string          `gorm:"type:text" json:"-"`
	AnalysisResult *AnalysisResult `gorm:"foreignKey:ResumeID" json:"analysis,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

type Section struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	Feedback string `json:"feedback"`
}

type AnalysisResult struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ResumeID        uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"resume_id"`
	OverallScore    int       `json:"overall_score"`
	SummaryFeedback string    `gorm:"type:text" json:"summary_feedback"`
	Sections        []Section `gorm:"serializer:json" json:"sections"`
	Suggestions     []string  `gorm:"serializer:json" json:"suggestions"`
	Keywords        []string  `gorm:"serializer:json" json:"keywords"`
	CreatedAt       time.Time `json:"created_at"`
}

func (r *Resume) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}

func (a *AnalysisResult) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
