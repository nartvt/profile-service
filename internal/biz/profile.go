package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/nartvt/profile-service/internal/nat"
)

type Profile struct {
	ID               uuid.UUID `json:"id,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
	UserID           string    `json:"user_id,omitempty"`
	FullName         string    `json:"full_name,omitempty"`
	Email            string    `json:"email,omitempty"`
	EmailConfirmedAt time.Time `json:"email_confirmed_at,omitempty"`
	Phone            string    `json:"phone,omitempty"`
	PhoneConfirmedAt time.Time `json:"phone_confirmed_at,omitempty"`
	ReferralCode     string    `json:"referral_code,omitempty"`
	IsSSOUser        bool      `json:"is_sso_user,omitempty"`
	Language         string    `json:"language,omitempty"`
}

type ProfileRepo interface {
	CreateProfile(ctx context.Context, profile *Profile) (*Profile, error)
	GetProfile(ctx context.Context, userId string) (*Profile, error)
	QueryProfile(ctx context.Context, req *Profile) (*Profile, error)
	QueryProfilesByListUserIds(ctx context.Context, userIds []string) ([]*Profile, error)
	UpdateLanguage(ctx context.Context, userId string, language string) (*Profile, error)
}

type ProfileUsecase struct {
	repo ProfileRepo
	log  *log.Helper
}

func NewProfileUsecase(repo ProfileRepo) *ProfileUsecase {
	return &ProfileUsecase{repo: repo, log: log.NewHelper(log.DefaultLogger)}
}

func (uc *ProfileUsecase) CreateProfile(ctx context.Context, profile Profile) (*Profile, error) {
	rsp, err := uc.repo.CreateProfile(ctx, &profile)
	if err != nil {
		return nil, err
	}
	go nat.PublishNewProfile(profile.ReferralCode, profile.UserID)
	return rsp, nil
}

func (uc *ProfileUsecase) GetProfile(ctx context.Context, userId string) (*Profile, error) {
	return uc.repo.GetProfile(ctx, userId)
}

func (uc *ProfileUsecase) QueryProfile(ctx context.Context, req *Profile) (*Profile, error) {
	return uc.repo.QueryProfile(ctx, req)
}

func (uc *ProfileUsecase) UpdateLanguage(ctx context.Context, userId, language string) (*Profile, error) {
	return uc.repo.UpdateLanguage(ctx, userId, language)
}

func (uc *ProfileUsecase) QueryProfiles(ctx context.Context, userIds []string) ([]*Profile, error) {
	return uc.repo.QueryProfilesByListUserIds(ctx, userIds)
}
