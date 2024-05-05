package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/util"
	"github.com/rs/zerolog/log"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (r *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	c context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := r.client.EnqueueContext(c, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("task:send_verify_email enqueued")

	return nil
}

func (r *RedisTaskProcessor) ProcessTaskSendVerifyEmail(c context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("could not unmarshal payload: %w", err)
	}

	user, err := r.store.GetUser(c, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("could not get user: %w", err)
	}

	verifyEmail, err := r.store.CreateVerifyEmail(c, db.CreateVerifyEmailParams{
		Username: user.Username,
		Email:    user.Email,
		Token:    util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("could not create verify email: %w", err)
	}

	subject := "Welcome to Bank App"
	verifyUrl := fmt.Sprintf(
		"http://localhost:8080/v1/verify_email?email_id=%d&token=%s",
		verifyEmail.ID,
		verifyEmail.Token,
	)
	content := fmt.Sprintf(
		"Welcome, %s! Please verify your email by clicking this <a href='%s'>link</a>",
		user.Username,
		verifyUrl,
	)
	to := []string{user.Email}
	err = r.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	log.Info().
		Str("username", user.Username).
		Str("email", user.Email).
		Msg("send verification email")

	return nil
}
