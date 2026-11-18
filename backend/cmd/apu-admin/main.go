package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/VillaAlexTor/Apu/backend/internal/database"
	"github.com/VillaAlexTor/Apu/backend/internal/identity"
	"github.com/VillaAlexTor/Apu/backend/internal/password"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("could not create platform administrator", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	databaseURL := strings.TrimSpace(os.Getenv("APU_DATABASE_URL"))
	email := strings.ToLower(strings.TrimSpace(os.Getenv("APU_BOOTSTRAP_EMAIL")))
	displayName := strings.TrimSpace(os.Getenv("APU_BOOTSTRAP_NAME"))
	plainPassword := os.Getenv("APU_BOOTSTRAP_PASSWORD")
	if databaseURL == "" || email == "" || displayName == "" || plainPassword == "" {
		return errors.New("APU_DATABASE_URL, APU_BOOTSTRAP_EMAIL, APU_BOOTSTRAP_NAME and APU_BOOTSTRAP_PASSWORD are required")
	}
	parsedAddress, err := mail.ParseAddress(email)
	if err != nil || strings.ToLower(parsedAddress.Address) != email {
		return errors.New("APU_BOOTSTRAP_EMAIL must be a valid email address")
	}
	if len(displayName) > 160 {
		return errors.New("APU_BOOTSTRAP_NAME must contain at most 160 characters")
	}
	passwordHash, err := password.Hash(plainPassword)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	user, err := identity.NewStore(pool).CreatePlatformAdmin(ctx, email, displayName, passwordHash)
	if err != nil {
		if errors.Is(err, identity.ErrEmailConflict) {
			return fmt.Errorf("a user with email %s already exists", email)
		}
		return err
	}
	logger.Info("platform administrator created", "user_id", user.ID, "email", user.Email)
	return nil
}
