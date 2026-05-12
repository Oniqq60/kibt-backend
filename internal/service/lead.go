package service

import (
	"context"
	"fmt"
	"kibit/internal/mailer"
	"kibit/internal/model"
	"kibit/internal/repository"
)

type LeadService struct {
	repo   *repository.LeadRepo
	mailer *mailer.Mailer
}

func NewLeadService(r *repository.LeadRepo, m *mailer.Mailer) *LeadService {
	return &LeadService{r, m}
}

func (s *LeadService) Create(ctx context.Context, req *model.CreateLeadRequest) error {
	id, err := s.repo.Create(ctx, req)
	if err != nil {
		return err
	}

	// Письмо админу
	adminHTML := fmt.Sprintf(`
		<h2>📩 Новая заявка КИБТ #%d</h2>
		<table style="border-collapse: collapse; width: 100%%; max-width: 600px;">
			<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Имя:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">%s</td></tr>
			<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Email:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">%s</td></tr>
			<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Телефон:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">%s</td></tr>
			<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Компания:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">%s</td></tr>
			<tr><td style="padding: 8px; border-bottom: 1px solid #eee;"><strong>Тип приложения:</strong></td><td style="padding: 8px; border-bottom: 1px solid #eee;">%s</td></tr>
			<tr><td style="padding: 8px;"><strong>Сообщение:</strong></td><td style="padding: 8px;">%s</td></tr>
		</table>
		<p style="margin-top: 20px; color: #666; font-size: 12px;">Заявка создана: %s</p>
	`,
		id,
		req.Name,
		req.Email,
		req.Phone,
		req.Company,
		req.AppType,
		req.Message,
	)
	s.mailer.SendAsync(ctx, "sales@kibt.ru", fmt.Sprintf("🔔 Новая заявка #%d", id), adminHTML)

	// Письмо пользователю
	userHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head><style>
			body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background: #2563EB; color: white; padding: 20px; border-radius: 8px 8px 0 0; }
			.content { background: #f9fafb; padding: 24px; border-radius: 0 0 8px 8px; }
			.summary { background: white; padding: 16px; border-radius: 6px; margin: 16px 0; border-left: 4px solid #2563EB; }
			.footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
		</style></head>
		<body>
			<div class="container">
				<div class="header">
					<h2 style="margin: 0;">✅ Заявка принята</h2>
				</div>
				<div class="content">
					<p>Здравствуйте, <strong>%s</strong>!</p>
					<p>Благодарим за обращение в КИБТ. Мы получили вашу заявку и свяжемся с вами в течение 24 часов.</p>
					
					<div class="summary">
						<h4 style="margin: 0 0 12px 0;">Краткая информация:</h4>
						<p style="margin: 4px 0;"><strong>Тип проекта:</strong> %s</p>
						<p style="margin: 4px 0;"><strong>Компания:</strong> %s</p>
						<p style="margin: 4px 0;"><strong>Сообщение:</strong> %s</p>
					</div>
					
					<p>Если у вас возникнут вопросы, просто ответьте на это письмо.</p>
					<p>С уважением,<br>Команда КИБТ</p>
				</div>
				<div class="footer">
					<p>Это автоматическое письмо, пожалуйста, не отвечайте на него напрямую.<br>
					© 2026 КИБТ. Все права защищены.</p>
				</div>
			</div>
		</body>
		</html>
	`,
		req.Name,
		req.AppType,
		func() string {
			if req.Company != "" {
				return req.Company
			}
			return "Не указана"
		}(),
		func() string {
			if req.Message != "" {
				return req.Message
			}
			return "Без комментария"
		}(),
	)

	// Отправляем подтверждение пользователю (на тот email, который он указал)
	s.mailer.SendAsync(ctx, req.Email, "Ваша заявка в КИБТ принята", userHTML)

	return nil
}
