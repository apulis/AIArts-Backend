package models

type GetEvaluationsReq struct {
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Status   string `form:"status" json:"status"`
	Search   string `form:"search" json:"search"`
	OrderBy  string `form:"orderBy" json:"orderBy"`
	Order    string `form:"order" json:"order"`
	VCName   string `form:"vcName" json:"vcName"`
}
