package adapter

import (
	"time"

	"github.com/nartvt/profile-service/ent"
)

type ChartData struct {
	ID         string    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     string    `json:"user_id" bson:"user_id,omitempty"`
	ChartId    uint32    `json:"chart_id" bson:"chart_id,omitempty"`
	TemplateId string    `json:"template_id" bson:"template_id,omitempty"`
	Client     string    `json:"client" bson:"client,omitempty"`
	Type       string    `json:"type" bson:"type,omitempty"`
	Name       string    `json:"name" bson:"name,omitempty"`
	Content    string    `json:"content" bson:"content,omitempty"`
	Symbol     string    `json:"symbol" bson:"symbol,omitempty"`
	Resolution string    `json:"resolution" bson:"resolution,omitempty"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

func BuildChartDataFromEnt(chart *ent.Chart) *ChartData {
	return &ChartData{
		ID:         chart.ID.String(),
		UserId:     chart.UserID,
		ChartId:    chart.ChartID,
		TemplateId: chart.TemplateID,
		Client:     chart.ClientID,
		Type:       chart.Type,
		Name:       chart.Name,
		Content:    chart.Content,
		Symbol:     chart.Symbol,
		Resolution: chart.Resolution,
		CreatedAt:  chart.CreatedAt,
		UpdatedAt:  chart.UpdatedAt,
	}
}

func BuildListChartDataFromEnt(charts []*ent.Chart) []*ChartData {
	res := []*ChartData{}
	for _, chart := range charts {
		res = append(res, BuildChartDataFromEnt(chart))
	}
	return res
}
