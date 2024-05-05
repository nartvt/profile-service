package biz

import (
	"context"

	"github.com/nartvt/profile-service/adapter"
	pb "github.com/nartvt/profile-service/api/profile/v1"
	"github.com/nartvt/profile-service/internal/constant"
)

func (uc *ChartUsecase) ListStudyTemplate(ctx context.Context, req *pb.ListStudyTemplateRequest) (*pb.ListStudyTemplateResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}

	charts, err := uc.repo.ListChart(ctx, req.ClientId, req.UserId, constant.TemplateType)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListStudyTemplateResponse{
		Status: constant.ResponseStatusOK,
	}
	for _, chart := range charts {
		resp.Data = append(resp.Data, &pb.ListStudyTemplateResponse_Data{
			Name: chart.Name,
		})
	}
	return resp, nil
}

func (uc *ChartUsecase) SaveStudyTemplate(ctx context.Context, req *pb.SaveStudyTemplateRequest) (*pb.SaveStudyTemplateResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}
	if req.Name == "" {
		return nil, constant.TemplateNameEmptyErr
	}
	if req.Content == "" {
		return nil, constant.ContentEmptyErr
	}

	upsertTemplate := &adapter.ChartData{
		UserId:     req.UserId,
		TemplateId: req.Name,
		Client:     req.ClientId,
		Type:       constant.TemplateType,
		Name:       req.Name,
		Content:    req.Content,
	}
	err := uc.repo.UpsertTemplate(ctx, upsertTemplate)
	if err != nil {
		return nil, err
	}
	return &pb.SaveStudyTemplateResponse{
		Status: constant.ResponseStatusOK,
	}, nil
}

func (uc *ChartUsecase) LoadStudyTemplate(ctx context.Context, req *pb.LoadStudyTemplateRequest) (*pb.LoadStudyTemplateResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}

	template, err := uc.repo.GetTemplate(ctx, req.ClientId, req.UserId, req.Name)
	if err != nil {
		return nil, err
	}
	resp := &pb.LoadStudyTemplateResponse{
		Status: constant.ResponseStatusOK,
		Data: &pb.LoadStudyTemplateResponse_Data{
			Name:    template.Name,
			Content: template.Content,
		},
	}
	return resp, nil
}

func (uc *ChartUsecase) DeleteStudyTemplates(ctx context.Context, req *pb.DeleteStudyTemplatesRequest) (*pb.DeleteStudyTemplatesResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}
	if req.Name == "" {
		return nil, constant.TemplateNameEmptyErr
	}

	err := uc.repo.DeleteTemplate(ctx, req.ClientId, req.UserId, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteStudyTemplatesResponse{
		Status: constant.ResponseStatusOK,
	}, nil
}
