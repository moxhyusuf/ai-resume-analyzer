package resume

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moxhyusuf/ai-resume-analyzer/pkg/groq"
	"github.com/moxhyusuf/ai-resume-analyzer/pkg/pdf"

	"github.com/google/uuid"
)

type ResumeService interface {
	Analyze(ctx context.Context, userID uuid.UUID, fileData []byte, filename string) (*AnalyzeResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]Resume, error)
	GetByID(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error)
	Delete(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) error
}

type resumeService struct {
	repo       ResumeRepository
	groqClient *groq.Client
}

func NewResumeService(repo ResumeRepository, groqClient *groq.Client) ResumeService {
	return &resumeService{repo, groqClient}
}

func (s *resumeService) Analyze(ctx context.Context, userID uuid.UUID, fileData []byte, filename string) (*AnalyzeResponse, error) {
	// 1. Extract text from PDF
	extractedText, err := pdf.ExtractText(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PDF: %w", err)
	}

	// 2. Save resume record
	resume := &Resume{
		UserID:        userID,
		OriginalName:  filename,
		ExtractedText: extractedText,
	}
	if err := s.repo.Create(resume); err != nil {
		return nil, fmt.Errorf("failed to save resume: %w", err)
	}

	// 3. Call Groq API (openai/gpt-oss-20b)
	rawJSON, err := s.groqClient.AnalyzeResume(ctx, extractedText)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	// 4. Parse AI JSON response
	var aiResult struct {
		OverallScore    int       `json:"overall_score"`
		SummaryFeedback string    `json:"summary_feedback"`
		Sections        []Section `json:"sections"`
		Suggestions     []string  `json:"suggestions"`
		Keywords        []string  `json:"keywords"`
	}
	if err := json.Unmarshal([]byte(rawJSON), &aiResult); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// 5. Save analysis result
	analysis := &AnalysisResult{
		ResumeID:        resume.ID,
		OverallScore:    aiResult.OverallScore,
		SummaryFeedback: aiResult.SummaryFeedback,
		Sections:        aiResult.Sections,
		Suggestions:     aiResult.Suggestions,
		Keywords:        aiResult.Keywords,
	}
	if err := s.repo.SaveAnalysis(analysis); err != nil {
		return nil, fmt.Errorf("failed to save analysis: %w", err)
	}

	resume.AnalysisResult = analysis
	return &AnalyzeResponse{
		ResumeID: resume.ID,
		Analysis: analysis,
	}, nil
}

func (s *resumeService) GetAll(ctx context.Context, userID uuid.UUID) ([]Resume, error) {
	return s.repo.FindByUserID(userID)
}

func (s *resumeService) GetByID(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error) {
	resume, err := s.repo.FindByID(resumeID)
	if err != nil {
		return nil, fmt.Errorf("resume not found")
	}
	if resume.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}
	return resume, nil
}

func (s *resumeService) Delete(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) error {
	resume, err := s.repo.FindByID(resumeID)
	if err != nil {
		return fmt.Errorf("resume not found")
	}
	if resume.UserID != userID {
		return fmt.Errorf("access denied")
	}
	return s.repo.Delete(resumeID)
}
