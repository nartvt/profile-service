package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nartvt/profile-service/adapter"
	"github.com/nartvt/profile-service/ent"
	"github.com/nartvt/profile-service/ent/chart"
	"github.com/nartvt/profile-service/ent/predicate"
	"github.com/nartvt/profile-service/internal/biz"
	"github.com/nartvt/profile-service/internal/constant"
)

type chartRepo struct {
	data *Data
	log  *log.Helper
}

func NewChartRepo(data *Data) biz.ChartRepo {
	return &chartRepo{
		data: data,
		log:  log.NewHelper(log.DefaultLogger),
	}
}

func (r *chartRepo) ListChart(ctx context.Context, client, userId, typeChart string) (results []*adapter.ChartData, err error) {
	var queryCondition []predicate.Chart
	queryCondition = append(queryCondition, chart.ClientIDEQ(client))
	queryCondition = append(queryCondition, chart.UserIDEQ(userId))
	queryCondition = append(queryCondition, chart.TypeEQ(typeChart))

	charts, err := r.data.db.Chart.Query().Where(queryCondition...).Order(ent.Asc("chart_id")).All(ctx)
	if err != nil {
		return nil, err
	}

	return adapter.BuildListChartDataFromEnt(charts), nil
}

func (r *chartRepo) GetChart(ctx context.Context, client, userId string, chartId uint32) (result *adapter.ChartData, err error) {
	var queryCondition []predicate.Chart
	queryCondition = append(queryCondition, chart.ClientIDEQ(client))
	queryCondition = append(queryCondition, chart.UserIDEQ(userId))
	queryCondition = append(queryCondition, chart.TypeEQ(constant.ChartType))
	queryCondition = append(queryCondition, chart.ChartIDEQ(chartId))

	chart, err := r.data.db.Chart.Query().Where(queryCondition...).First(ctx)
	if err != nil {
		return nil, err
	}
	return adapter.BuildChartDataFromEnt(chart), nil
}

func (r *chartRepo) InsertChart(ctx context.Context, c *adapter.ChartData) (err error) {
	now := time.Now()
	_, err = r.data.db.Chart.Create().
		SetUserID(c.UserId).
		SetTemplateID(c.TemplateId).
		SetClientID(c.Client).
		SetChartID(c.ChartId).
		SetType(c.Type).
		SetName(c.Name).
		SetContent(c.Content).
		SetSymbol(c.Symbol).
		SetResolution(c.Resolution).
		SetCreatedAt(now.UTC()).
		SetUpdatedAt(now.UTC()).
		Save(ctx)

	return err
}

func (r *chartRepo) UpdateChart(ctx context.Context, c *adapter.ChartData) (err error) {
	var queryCondition []predicate.Chart
	queryCondition = append(queryCondition, chart.ClientIDEQ(c.Client))
	queryCondition = append(queryCondition, chart.UserIDEQ(c.UserId))
	queryCondition = append(queryCondition, chart.TypeEQ(constant.ChartType))
	queryCondition = append(queryCondition, chart.ChartIDEQ(c.ChartId))

	chart, err := r.data.db.Chart.Query().Where(queryCondition...).Only(ctx)
	if err != nil {
		return err
	}

	_, err = r.data.db.Chart.UpdateOneID(chart.ID).
		SetName(c.Name).
		SetContent(c.Content).
		SetResolution(c.Resolution).
		SetUpdatedAt(time.Now().UTC()).
		Save(ctx)
	return
}

func (r *chartRepo) MaxChartId(ctx context.Context, client, userId string) (id uint32, err error) {
	var queryCondition []predicate.Chart
	queryCondition = append(queryCondition, chart.ClientIDEQ(client))
	queryCondition = append(queryCondition, chart.UserIDEQ(userId))
	queryCondition = append(queryCondition, chart.TypeEQ(constant.ChartType))

	chart, err := r.data.db.Chart.Query().Where(queryCondition...).Order(ent.Desc("chart_id")).First(ctx)
	if ent.IsNotFound(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	id = chart.ChartID
	return
}

func (r *chartRepo) DeleteChart(ctx context.Context, client, userId string, chartId uint32) (err error) {
	var queryCondition []predicate.Chart
	queryCondition = append(queryCondition, chart.ClientIDEQ(client))
	queryCondition = append(queryCondition, chart.UserIDEQ(userId))
	queryCondition = append(queryCondition, chart.TypeEQ(constant.ChartType))
	queryCondition = append(queryCondition, chart.ChartIDEQ(chartId))

	chartEnt, err := r.data.db.Chart.Query().Where(queryCondition...).First(ctx)
	if chartEnt != nil {
		err = r.data.db.Chart.DeleteOneID(chartEnt.ID).Exec(ctx)
	}
	return
}

func (r *chartRepo) UpsertTemplate(ctx context.Context, c *adapter.ChartData) (err error) {
	var queryCondition []predicate.Chart
	queryCondition = append(queryCondition, chart.ClientIDEQ(c.Client))
	queryCondition = append(queryCondition, chart.UserIDEQ(c.UserId))
	queryCondition = append(queryCondition, chart.TypeEQ(constant.TemplateType))
	queryCondition = append(queryCondition, chart.NameEQ(c.Name))
	chartEnt, err := r.data.db.Chart.Query().Where(queryCondition...).First(ctx)
	if chartEnt == nil {
		now := time.Now()
		_, err = r.data.db.Chart.Create().
			SetUserID(c.UserId).
			SetClientID(c.Client).
			SetType(constant.TemplateType).
			SetTemplateID(c.Name).
			SetName(c.Name).
			SetContent(c.Content).
			SetSymbol(c.Symbol).
			SetResolution(c.Resolution).
			SetCreatedAt(now.UTC()).
			SetUpdatedAt(now.UTC()).
			Save(ctx)
	} else {
		now := time.Now()
		_, err = r.data.db.Chart.UpdateOneID(chartEnt.ID).
			SetUserID(c.UserId).
			SetClientID(c.Client).
			SetType(constant.TemplateType).
			SetTemplateID(c.Name).
			SetName(c.Name).
			SetContent(c.Content).
			SetUpdatedAt(now.UTC()).
			Save(ctx)
	}
	return
}

func (r *chartRepo) GetTemplate(ctx context.Context, client, userId, templateId string) (result *adapter.ChartData, err error) {
	var queryCondition []predicate.Chart
	queryCondition = append(queryCondition, chart.ClientIDEQ(client))
	queryCondition = append(queryCondition, chart.UserIDEQ(userId))
	queryCondition = append(queryCondition, chart.TypeEQ(constant.TemplateType))
	queryCondition = append(queryCondition, chart.TemplateIDEQ(templateId))
	chartEnt, err := r.data.db.Chart.Query().Where(queryCondition...).First(ctx)
	if err == nil {
		return adapter.BuildChartDataFromEnt(chartEnt), nil
	}
	return
}

func (r *chartRepo) DeleteTemplate(ctx context.Context, client, userId, templateId string) (err error) {
	var queryCondition []predicate.Chart
	queryCondition = append(queryCondition, chart.ClientIDEQ(client))
	queryCondition = append(queryCondition, chart.UserIDEQ(userId))
	queryCondition = append(queryCondition, chart.TypeEQ(constant.TemplateType))
	queryCondition = append(queryCondition, chart.TemplateIDEQ(templateId))
	chartEnt, err := r.data.db.Chart.Query().Where(queryCondition...).First(ctx)
	if chartEnt != nil {
		err = r.data.db.Chart.DeleteOneID(chartEnt.ID).Exec(ctx)
	}
	return
}
