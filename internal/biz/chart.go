package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nartvt/profile-service/adapter"
	pb "github.com/nartvt/profile-service/api/profile/v1" // TODO: need to remove it. we should not use proto in biz
	"github.com/nartvt/profile-service/internal/constant"
)

type ChartRepo interface {
	ListChart(ctx context.Context, client, userId, typeChart string) (results []*adapter.ChartData, err error)
	GetChart(ctx context.Context, client, userId string, chartId uint32) (result *adapter.ChartData, err error)
	InsertChart(ctx context.Context, c *adapter.ChartData) (err error)
	UpdateChart(ctx context.Context, c *adapter.ChartData) (err error)
	MaxChartId(ctx context.Context, client, userId string) (id uint32, err error)
	DeleteChart(ctx context.Context, client, userId string, chartId uint32) (err error)
	UpsertTemplate(ctx context.Context, c *adapter.ChartData) (err error)
	GetTemplate(ctx context.Context, client, userId, name string) (result *adapter.ChartData, err error)
	DeleteTemplate(ctx context.Context, client, userId, name string) (err error)
}

type ChartUsecase struct {
	repo ChartRepo
	log  *log.Helper
}

func NewChartUsecase(repo ChartRepo) *ChartUsecase {
	return &ChartUsecase{repo: repo, log: log.NewHelper(log.DefaultLogger)}
}

func (uc *ChartUsecase) ListChart(ctx context.Context, req *pb.ListChartRequest) (*pb.ListChartResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}

	charts, err := uc.repo.ListChart(ctx, req.ClientId, req.UserId, constant.ChartType)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListChartResponse{
		Status: constant.ResponseStatusOK,
	}
	for _, chart := range charts {
		resp.Data = append(resp.Data, &pb.ListChartResponse_Data{
			Timestamp:  uint32(chart.UpdatedAt.Unix()),
			Symbol:     chart.Symbol,
			Resolution: chart.Resolution,
			Id:         chart.ChartId,
			Name:       chart.Name,
		})
	}
	return resp, nil
}

func (uc *ChartUsecase) SaveChart(ctx context.Context, req *pb.SaveChartRequest) (*pb.SaveChartResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}
	if req.Symbol == "" {
		return nil, constant.SymbolEmptyErr
	}
	if req.Resolution == "" {
		return nil, constant.ResolutionEmptyErr
	}
	if req.Name == "" {
		return nil, constant.ChartNameEmptyErr
	}
	if req.Content == "" {
		return nil, constant.ContentEmptyErr
	}

	chartId, err := uc.repo.MaxChartId(ctx, req.ClientId, req.UserId)
	if err != nil {
		return nil, err
	}
	chartId += 1
	newChart := &adapter.ChartData{
		UserId:     req.UserId,
		ChartId:    chartId,
		Client:     req.ClientId,
		Type:       constant.ChartType,
		Name:       req.Name,
		Content:    req.Content,
		Symbol:     req.Symbol,
		Resolution: req.Resolution,
	}
	err = uc.repo.InsertChart(ctx, newChart)
	if err != nil {
		return nil, err
	}
	resp := &pb.SaveChartResponse{
		Status: constant.ResponseStatusOK,
		Id:     chartId,
	}
	return resp, nil
}

func (uc *ChartUsecase) SaveAsChart(ctx context.Context, req *pb.SaveAsChartRequest) (*pb.SaveAsChartResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}
	if req.Symbol == "" {
		return nil, constant.SymbolEmptyErr
	}
	if req.Resolution == "" {
		return nil, constant.ResolutionEmptyErr
	}
	if req.Name == "" {
		return nil, constant.ChartNameEmptyErr
	}
	if req.Content == "" {
		return nil, constant.ContentEmptyErr
	}
	if req.ChartId <= 0 {
		return nil, constant.ChartIDInvalidErr
	}

	updateChart := &adapter.ChartData{
		UserId:     req.UserId,
		ChartId:    req.ChartId,
		Client:     req.ClientId,
		Type:       constant.ChartType,
		Name:       req.Name,
		Content:    req.Content,
		Symbol:     req.Symbol,
		Resolution: req.Resolution,
	}
	err := uc.repo.UpdateChart(ctx, updateChart)
	if err != nil {
		return nil, err
	}
	resp := &pb.SaveAsChartResponse{
		Status: constant.ResponseStatusOK,
	}
	return resp, nil
}

func (uc *ChartUsecase) LoadChart(ctx context.Context, req *pb.LoadChartRequest) (*pb.LoadChartResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}
	if req.ChartId <= 0 {
		return nil, constant.ChartIDInvalidErr
	}

	chart, err := uc.repo.GetChart(ctx, req.ClientId, req.UserId, req.ChartId)
	if err != nil {
		return nil, err
	}
	data := &pb.LoadChartResponse_Data{
		Content:   chart.Content,
		Timestamp: uint32(chart.UpdatedAt.Unix()),
		Id:        chart.ChartId,
		Name:      chart.Name,
	}
	resp := &pb.LoadChartResponse{
		Status: constant.ResponseStatusOK,
		Data:   data,
	}
	return resp, nil
}

func (uc *ChartUsecase) DeleteChart(ctx context.Context, req *pb.DeleteChartRequest) (*pb.DeleteChartResponse, error) {
	if req.ClientId == "" {
		return nil, constant.ClientIDEmptyErr
	}
	if req.UserId == "" {
		return nil, constant.UserIDEmptyErr
	}
	if req.ChartId <= 0 {
		return nil, constant.ChartIDInvalidErr
	}

	err := uc.repo.DeleteChart(ctx, req.ClientId, req.UserId, req.ChartId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteChartResponse{
		Status: constant.ResponseStatusOK,
	}, nil
}
