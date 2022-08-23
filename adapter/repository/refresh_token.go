package repository

import (
	"context"
	"fmt"
	"time"
)

type Token struct {
	ID           string
	Role         string
	JTI          string
	RefreshToken string
	RefreshExp   time.Time
}

func (storage Storage) SaveRefreshToken(ctx context.Context, tk Token) error {

	rtKey := fmt.Sprintf("%s:%s", tk.ID, tk.JTI) // userID:JTI

	err := storage.Redis.Set(ctx, rtKey, tk.RefreshToken, time.Since(tk.RefreshExp)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage Storage) DoesRefreshTokenExist(ctx context.Context, tk Token) (bool, error) {

	rtKey := fmt.Sprintf("%s:%s", tk.ID, tk.JTI) // userID:JTI

	cmd := storage.Redis.Get(ctx, rtKey)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	value := cmd.String()
	if value == "" {
		return false, nil
	}

	return true, nil
}

func (storage Storage) DeleteRefreshToken(ctx context.Context, tk Token) error {

	rtKey := fmt.Sprintf("%s:%s", tk.ID, tk.JTI) // userID:JTI

	err := storage.Redis.Del(ctx, rtKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage Storage) DeleteUserRefreshTokens(ctx context.Context, userID string) error {

	prefix := fmt.Sprintf("%s*", userID)

	iter := storage.Redis.Scan(ctx, 0, prefix, 0).Iterator()

	for iter.Next(ctx) {
		err := storage.Redis.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
