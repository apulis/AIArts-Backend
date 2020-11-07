package models

type VCOperateType int
const (
	VC_OPTYPE_ADD	VCOperateType = iota
	VC_OPTYPE_DEL
	VC_OPTYPE_UPDATE
	VC_OPTYPE_GET

	VC_STATISTIC_JOB     int = 1
	VC_STATISTIC_UNALLOC     	= 2
	VC_STATISTIC_USER_UNUSED    = 3
)

func (opType VCOperateType) GetAPIName() string {

	if opType == VC_OPTYPE_ADD {
		return "AddVC"
	} else if opType == VC_OPTYPE_UPDATE {
		return "UpdateVC"
	} else if opType == VC_OPTYPE_DEL {
		return "DeleteVC"
	} else if opType == VC_OPTYPE_GET {
		return "GetVC"
	}

	return ""
}

type VCRsp struct {
	CurrPage 	        []VCItem			`json:"result"`
	Total				int 				`json:"totalNum"`
}

type VCStatisticReq struct {
	Type 				int 				`form:"type" json:"type"`
	TargetStatus        string 				`form:"targetStatus" json:"targetStatus"`
	VCName 				string 				`form:"vcName" json:"vcName"`
}

type VCStatisticRsp struct {
	JobCount	        int 				`json:"jobCount"`
	UnallocatedDevice   map[string]int 		`json:"unallocatedDevice"`   // 未分配给VC的设备
	UserUnusedDevice    map[string]int 		`json:"userUnusedDevice"`    // 用户配额下的未使用设备
}
