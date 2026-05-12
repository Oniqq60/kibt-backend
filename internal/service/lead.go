package service

import (
	"context"
	"fmt"
	"kibit/internal/mailer"
	"kibit/internal/model"
	"kibit/internal/repository"
	"log/slog"
)

type LeadService struct {
	repo   *repository.LeadRepo
	mailer *mailer.Mailer
}

func NewLeadService(r *repository.LeadRepo, m *mailer.Mailer) *LeadService { return &LeadService{r, m} }

func (s *LeadService) Create(ctx context.Context, req *model.CreateLeadRequest) error {
	id, err := s.repo.Create(ctx, req)
	if err != nil {
		return err
	}
	html := fmt.Sprintf(`<h2>Новый лид КИБТ #%d</h2><p>%s (%s)</p><p>Тип: %s</p><p>Сообщение: %s</p>`, id, req.Name, req.Email, req.AppType, req.Message)
	s.mailer.SendAsync(ctx, "sales@kibt.ru", "Новая заявка с лендинга", html) // вставить почту админа вместо sales@kibt.ru
	slog.Info("Lead saved & email queued", "id", id)
	return nil
}
