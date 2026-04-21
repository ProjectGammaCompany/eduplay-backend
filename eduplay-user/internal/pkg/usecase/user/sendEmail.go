package user

import (
	"bytes"
	"context"
	"eduplay-user/internal/model"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"time"
)

const verificationCodeSize = 6
const sendEmailTimeout = 10 * time.Second

func (a *UseCase) SendVerificationCodeEmail(email string, code string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), sendEmailTimeout)
		defer cancel()

		message, err := generateVerificationEmail(code)
		if err != nil {
			a.log.Error("failed to generate verification email", slog.String("error", err.Error()))
			return
		}

		emailDate, err := json.Marshal(model.EmailRequest{
			To:      email,
			Subject: "Сброс пароля",
			Body:    message,
		})
		if err != nil {
			a.log.Error("failed to marshal email data", slog.String("error", err.Error()))
			return
		}

		err = a.emailClient.SendEmail(emailDate)
		if err != nil {
			a.log.Error("failed to send email", slog.String("error", err.Error()))
			return
		}

		err = a.storage.PutVerificationCode(ctx, email, code)
		if err != nil {
			a.log.Error("failed to put verification code", slog.String("error", err.Error()))
			return
		}

	}()
}

func generateVerificationEmail(code string) (string, error) {
	data := struct {
		Code string
	}{
		Code: code,
	}

	tmpl, err := template.New("verification_code").Parse(sendCodeTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
