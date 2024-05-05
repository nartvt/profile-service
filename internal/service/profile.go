package service

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nartvt/go-core/middleware/jwt"
	pb "github.com/nartvt/profile-service/api/profile/v1"
	"github.com/nartvt/profile-service/internal/biz"
	"github.com/nartvt/profile-service/internal/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProfileService struct {
	pb.UnimplementedProfileServiceServer

	log       *log.Helper
	profielUc *biz.ProfileUsecase
}

func NewProfileService(profielUc *biz.ProfileUsecase) *ProfileService {
	return &ProfileService{profielUc: profielUc, log: log.NewHelper(log.DefaultLogger)}
}

func (s *ProfileService) ListenUserAccount(ctx context.Context, req *pb.ListenUserAccountRequest) (*emptypb.Empty, error) {
	userMetadata := req.User.UserMetadata
	mapMetadata := userMetadata.AsMap()
	s.log.Debugf("[METADATA]: %v", mapMetadata)

	name := ""
	if mapMetadata[utils.AttributeName] != nil {
		name = mapMetadata[utils.AttributeName].(string)
	}
	_, err := s.profielUc.CreateProfile(ctx, biz.Profile{UserID: req.User.Id, Email: req.User.Email, EmailConfirmedAt: req.User.EmailConfirmedAt.AsTime(),
		ReferralCode: fmt.Sprintf("%v", mapMetadata[utils.AttributeReferralCode]), Language: fmt.Sprintf("%v", mapMetadata[utils.AttributeLocale]), FullName: name})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *ProfileService) GetUserProfile(ctx context.Context, req *emptypb.Empty) (*pb.GetUserProfileResponse, error) {
	userId, err := jwt.GetUserId(ctx)
	if err != nil {
		log.Error(err)
		return nil, errors.Unauthorized("UNAUTHENTICATED", "")
	}
	s.log.Debugf("USER_ID_FROM_REQUEST: ", userId)
	profile, err := s.profielUc.GetProfile(ctx, userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &pb.GetUserProfileResponse{Data: &pb.UserProfile{
		Id: profile.UserID, Email: profile.Email, EmailConfirmedAt: timestamppb.New(profile.EmailConfirmedAt), Phone: profile.Phone, PhoneConfirmedAt: timestamppb.New(profile.PhoneConfirmedAt), FullName: profile.FullName, Language: profile.Language,
	}}, nil
}

func (s *ProfileService) UpdateLanguage(ctx context.Context, req *pb.UpdateLanguageRequest) (*pb.UpdateLanguageResponse, error) {
	userId, err := jwt.GetUserId(ctx)
	if err != nil {
		log.Error(err)
		return nil, errors.Unauthorized("UNAUTHENTICATED", "")
	}

	profile, err := s.profielUc.UpdateLanguage(ctx, userId, req.Language)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &pb.UpdateLanguageResponse{Data: &pb.UserProfile{
		Id: profile.UserID, Email: profile.Email, EmailConfirmedAt: timestamppb.New(profile.EmailConfirmedAt), Phone: profile.Phone, PhoneConfirmedAt: timestamppb.New(profile.PhoneConfirmedAt), FullName: profile.FullName, Language: profile.Language,
	}}, nil
}

func (s *ProfileService) GetUserProfileInternal(ctx context.Context, req *pb.GetUserProfileInternalRequest) (*pb.GetUserProfileResponse, error) {
	if len(req.UserId) == 0 && len(req.Email) == 0 {
		return nil, errors.BadRequest("INVALID_PARAM", "")
	}
	profile, err := s.profielUc.QueryProfile(ctx, &biz.Profile{
		UserID: req.UserId,
		Email:  req.Email,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &pb.GetUserProfileResponse{Data: s.mapBizToResponse(profile)}, nil
}

func (s *ProfileService) mapBizToResponse(profile *biz.Profile) *pb.UserProfile {
	return &pb.UserProfile{
		Id:               profile.UserID,
		Email:            profile.Email,
		EmailConfirmedAt: timestamppb.New(profile.EmailConfirmedAt),
		Phone:            profile.Phone,
		PhoneConfirmedAt: timestamppb.New(profile.PhoneConfirmedAt),
		FullName:         profile.FullName,
		Language:         profile.Language,
		CreatedAt:        timestamppb.New(profile.CreatedAt),
	}
}

func (s *ProfileService) mapListBizToResponse(profiles []*biz.Profile) (res []*pb.UserProfile) {
	for _, profile := range profiles {
		res = append(res, s.mapBizToResponse(profile))
	}
	return
}

func (s *ProfileService) GetListUserProfileInternal(ctx context.Context, req *pb.GetListUserProfileInternalRequest) (*pb.GetListUserProfileResponse, error) {
	profiles, err := s.profielUc.QueryProfiles(ctx, req.UserIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &pb.GetListUserProfileResponse{Data: s.mapListBizToResponse(profiles)}, nil
}
