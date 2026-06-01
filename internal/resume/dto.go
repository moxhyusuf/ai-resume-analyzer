package resume

import "github.com/google/uuid"

type AnalyzeResponse struct {
	ResumeID uuid.UUID       `json:"resume_id"`
	Analysis *AnalysisResult `json:"analysis"`
}
