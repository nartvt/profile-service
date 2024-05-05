package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nartvt/profile-service/ent"
	"github.com/nartvt/profile-service/ent/predicate"
	"github.com/nartvt/profile-service/ent/profile"
	"github.com/nartvt/profile-service/internal/biz"
)

type profileRepo struct {
	data *Data
	log  *log.Helper
}

func NewProfileRepo(data *Data) biz.ProfileRepo {
	return &profileRepo{
		data: data,
		log:  log.NewHelper(log.DefaultLogger),
	}
}

// CreateProfile implements biz.ProfileRepo.
func (r *profileRepo) CreateProfile(ctx context.Context, bizProfile *biz.Profile) (*biz.Profile, error) {
	id, err := r.data.db.Profile.Create().SetFullName(bizProfile.FullName).SetUserID(bizProfile.UserID).
		SetEmail(bizProfile.Email).SetEmailConfirmedAt(bizProfile.EmailConfirmedAt).
		SetPhone(bizProfile.Phone).SetPhoneConfirmedAt(bizProfile.PhoneConfirmedAt).Save(ctx)
	// SetPhone(bizProfile.Phone).SetPhoneConfirmedAt(bizProfile.PhoneConfirmedAt).OnConflictColumns(profile.FieldUserID).UpdateNewValues().ID(ctx)
	if err != nil {
		return nil, err
	}
	rs, err := r.data.db.Profile.Get(ctx, id.ID)
	if err != nil {
		return nil, err
	}
	return r.mapEntToBiz(rs), nil
}

// GetProfile implements biz.ProfileRepo.
func (r *profileRepo) GetProfile(ctx context.Context, userId string) (*biz.Profile, error) {
	rs, err := r.data.db.Profile.Query().Where(profile.UserID(userId)).First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return &biz.Profile{}, nil
		}
		return nil, err
	}
	return r.mapEntToBiz(rs), nil
}

// GetProfile implements biz.ProfileRepo.
func (r *profileRepo) QueryProfile(ctx context.Context, req *biz.Profile) (*biz.Profile, error) {
	var queryCondition []predicate.Profile
	if req.UserID != "" {
		queryCondition = append(queryCondition, profile.UserIDEQ(req.UserID))
	}
	if req.Email != "" {
		queryCondition = append(queryCondition, profile.EmailEQ(req.Email))
	}

	rs, err := r.data.db.Profile.Query().Where(queryCondition...).First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return &biz.Profile{}, nil
		}
		return nil, err
	}
	return r.mapEntToBiz(rs), nil
}

func (r *profileRepo) UpdateLanguage(ctx context.Context, userId string, language string) (*biz.Profile, error) {
	_, err := r.data.db.Profile.Update().SetLanguage(language).Where(profile.UserID(userId)).Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return &biz.Profile{}, nil
		}
		return nil, err
	}
	rs, err := r.data.db.Profile.Query().Where(profile.UserID(userId)).First(ctx)
	return r.mapEntToBiz(rs), nil
}

func (r *profileRepo) QueryProfilesByListUserIds(ctx context.Context, userIds []string) ([]*biz.Profile, error) {
	var queryCondition []predicate.Profile
	queryCondition = append(queryCondition, profile.UserIDIn(userIds...))

	rs, err := r.data.db.Profile.Query().Where(queryCondition...).All(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return []*biz.Profile{}, nil
		}
		return nil, err
	}
	return r.mapListEntToBiz(rs), nil
}

func (*profileRepo) mapEntToBiz(p *ent.Profile) *biz.Profile {
	return &biz.Profile{ID: p.ID, UserID: p.UserID, FullName: p.FullName, Email: p.Email, Phone: p.Phone, EmailConfirmedAt: p.EmailConfirmedAt, PhoneConfirmedAt: p.PhoneConfirmedAt, IsSSOUser: p.IsSSOUser, Language: p.Language, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func (r *profileRepo) mapListEntToBiz(ps []*ent.Profile) (res []*biz.Profile) {
	for _, p := range ps {
		res = append(res, r.mapEntToBiz(p))
	}
	return
}
