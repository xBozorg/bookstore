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

	if err := storage.Redis.Set(
		ctx,
		fmt.Sprintf("%s:%s:rt:%s", tk.Role, tk.ID, tk.JTI), // {role}:{id}:rt:{jti}
		tk.RefreshToken,
		time.Since(tk.RefreshExp),
	).Err(); err != nil {
		return err
	}

	return nil
}

func (storage Storage) DoesRefreshTokenExist(ctx context.Context, tk Token) (bool, error) {

	cmd := storage.Redis.Get(
		ctx,
		fmt.Sprintf("%s:%s:rt:%s", tk.Role, tk.ID, tk.JTI), // {role}:{id}:rt:{jti}
	)
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

	if err := storage.Redis.Del(
		ctx,
		fmt.Sprintf("%s:%s:rt:%s", tk.Role, tk.ID, tk.JTI), // {role}:{id}:rt:{jti}
	).Err(); err != nil {
		return err
	}

	return nil
}

func (storage Storage) DeleteRefreshTokens(ctx context.Context, role, id string) error {

	iter := storage.Redis.Scan(
		ctx,
		0,
		fmt.Sprintf("%s:%s*", role, id), // {role}:{id}
		0,
	).Iterator()

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
