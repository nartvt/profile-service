package constant

import "errors"

var (
	ClientIDEmptyErr       = errors.New("CLIENT_ID_REQUIRED")
	UserIDEmptyErr         = errors.New("USER_ID_REQUIRED")
	ChartNameEmptyErr      = errors.New("CHART_NAME_REQUIRED")
	TemplateNameEmptyErr   = errors.New("TEMPLATE_NAME_REQUIRED")
	TemplateIDEmptyErr     = errors.New("TEMPLATE_ID_REQUIRED")
	ContentEmptyErr        = errors.New("TEMPLATE_CONTENT_REQUIRED")
	SymbolEmptyErr         = errors.New("SYMBOL_REQUIRED")
	ResolutionEmptyErr     = errors.New("RESOLUTION_REQUIRED")
	ChartIDInvalidErr      = errors.New("CHART_ID_INVALID")
	GetMaxChartIDFailedErr = errors.New("GET_MAX_CHART_ID_FAILED")
)
