package projects

import (
	"context"
	"errors"
	"time"

	"qr-saas/internal/qr"

	"github.com/google/uuid"
)

type Service interface {
	CreateProject(ctx context.Context, userID, name, color string) (*Project, error)
	ListProjects(ctx context.Context, userID string) ([]Project, error)
	GetProject(ctx context.Context, userID, id string) (*Project, error)
	UpdateProject(ctx context.Context, userID, id string, req UpdateProjectRequest) (*Project, error)
	DeleteProject(ctx context.Context, userID, id string) error

	// NEW FEATURES
	ListProjectQRs(ctx context.Context, userID, projectID string) ([]qr.QRCode, error)
	AssignQR(ctx context.Context, userID, qrID, projectID string) error
	RemoveQR(ctx context.Context, userID, qrID string) error
}

type service struct {
	repo   Repository
	qrRepo qr.Repository
}

func NewService(repo Repository, qrRepo qr.Repository) Service {
	return &service{
		repo:   repo,
		qrRepo: qrRepo,
	}
}

// ----------------------------
// CREATE
// ----------------------------
func (s *service) CreateProject(ctx context.Context, userID, name, color string) (*Project, error) {
	if color == "" {
		color = "#3B82F6" // default Tailwind blue
	}

	p := &Project{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		Color:     color,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

// ----------------------------
// LIST PROJECTS
// ----------------------------
func (s *service) ListProjects(ctx context.Context, userID string) ([]Project, error) {
	return s.repo.ListByUser(ctx, userID)
}

// ----------------------------
// GET PROJECT BY ID
// ----------------------------
func (s *service) GetProject(ctx context.Context, userID, id string) (*Project, error) {
	return s.repo.GetByID(ctx, userID, id)
}

// ----------------------------
// UPDATE PROJECT
// ----------------------------
func (s *service) UpdateProject(ctx context.Context, userID, id string, req UpdateProjectRequest) (*Project, error) {
	p, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Color != nil {
		p.Color = *req.Color
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

// ----------------------------
// DELETE PROJECT
// ----------------------------
func (s *service) DeleteProject(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

// =====================================================================
// NEW: LIST QRs INSIDE A PROJECT
// =====================================================================
func (s *service) ListProjectQRs(ctx context.Context, userID, projectID string) ([]qr.QRCode, error) {

	// ensure project exists & belongs to user
	p, err := s.repo.GetByID(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("project not found")
	}

	return s.repo.ListProjectQRs(ctx, userID, projectID)
}

// =====================================================================
// NEW: ASSIGN QR TO PROJECT
// =====================================================================
func (s *service) AssignQR(ctx context.Context, userID, qrID, projectID string) error {
	// make sure QR belongs to the user
	q, err := s.qrRepo.GetByID(ctx, qrID, userID) // ✅ FIXED ORDER
	if err != nil {
		return err
	}
	if q == nil {
		return errors.New("qr not found")
	}

	// If projectID is not empty, ensure project belongs to the user
	if projectID != "" {
		p, err := s.repo.GetByID(ctx, userID, projectID)
		if err != nil {
			return err
		}
		if p == nil {
			return errors.New("project not found")
		}
	}

	return s.repo.AssignQR(ctx, userID, qrID, projectID)
}

// =====================================================================
// NEW: REMOVE QR FROM ANY PROJECT
// =====================================================================
func (s *service) RemoveQR(ctx context.Context, userID, qrID string) error {
	return s.repo.AssignQR(ctx, userID, qrID, "")
}
