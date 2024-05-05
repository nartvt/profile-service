package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/nartvt/profile-service/api/profile/v1"
	"github.com/nartvt/profile-service/internal/biz"
)

type ChartService struct {
	pb.UnimplementedChartServiceServer

	log     *log.Helper
	chartUc *biz.ChartUsecase
}

func NewChartService(greeterUc *biz.ChartUsecase) *ChartService {
	return &ChartService{chartUc: greeterUc, log: log.NewHelper(log.DefaultLogger)}
}

func (s *ChartService) ListChart(ctx context.Context, req *pb.ListChartRequest) (*pb.ListChartResponse, error) {
	return s.chartUc.ListChart(ctx, req)
}

func (s *ChartService) SaveChart(ctx context.Context, req *pb.SaveChartRequest) (*pb.SaveChartResponse, error) {
	return s.chartUc.SaveChart(ctx, req)
}

func (s *ChartService) SaveAsChart(ctx context.Context, req *pb.SaveAsChartRequest) (*pb.SaveAsChartResponse, error) {
	return s.chartUc.SaveAsChart(ctx, req)
}

func (s *ChartService) LoadChart(ctx context.Context, req *pb.LoadChartRequest) (*pb.LoadChartResponse, error) {
	return s.chartUc.LoadChart(ctx, req)
}

func (s *ChartService) DeleteChart(ctx context.Context, req *pb.DeleteChartRequest) (*pb.DeleteChartResponse, error) {
	return s.chartUc.DeleteChart(ctx, req)
}

func (s *ChartService) ListStudyTemplate(ctx context.Context, req *pb.ListStudyTemplateRequest) (*pb.ListStudyTemplateResponse, error) {
	return s.chartUc.ListStudyTemplate(ctx, req)
}

func (s *ChartService) SaveStudyTemplate(ctx context.Context, req *pb.SaveStudyTemplateRequest) (*pb.SaveStudyTemplateResponse, error) {
	return s.chartUc.SaveStudyTemplate(ctx, req)
}

func (s *ChartService) LoadStudyTemplate(ctx context.Context, req *pb.LoadStudyTemplateRequest) (*pb.LoadStudyTemplateResponse, error) {
	return s.chartUc.LoadStudyTemplate(ctx, req)
}

func (s *ChartService) DeleteStudyTemplates(ctx context.Context, req *pb.DeleteStudyTemplatesRequest) (*pb.DeleteStudyTemplatesResponse, error) {
	return s.chartUc.DeleteStudyTemplates(ctx, req)
}
