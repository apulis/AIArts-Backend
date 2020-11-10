package services

import (
	"encoding/json"
	"fmt"
	urllib "net/url"
	"strconv"
	"strings"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
)

func OperateVC(userName string, opType models.VCOperateType, item *models.VCItem) error {

	if len(userName) == 0 || item == nil || item.VCName == nil {
		return fmt.Errorf("invalid paramater")
	}

	url := fmt.Sprintf("%s/%s?userName=%s&vcName=%s", configs.Config.DltsUrl, opType.GetAPIName(), userName, *(item.VCName))
	if opType == models.VC_OPTYPE_ADD || opType == models.VC_OPTYPE_UPDATE {

		if item.Metadata != nil {
			url = url + fmt.Sprintf("&metadata=%s", urllib.PathEscape(*(item.Metadata)))
		}

		if item.Quota != nil {
			url = url + fmt.Sprintf("&quota=%s", urllib.PathEscape(*(item.Quota)))
		}
	} else if opType == models.VC_OPTYPE_DEL {

	} else if opType == models.VC_OPTYPE_GET {

	} else {
		return fmt.Errorf("wrong operate type")
	}

	err := DoRequest(url, "GET", nil, nil, item)
	if err != nil {
		fmt.Printf("operate vc err[%+v]\n", err)
		return err
	}

	return nil
}


func ListVC(userName string, paging models.Paging) (*models.VCRsp, error) {

	url := fmt.Sprintf("%s/ListVCs?userName=%s&page=%d&size=%d&name=%s",
		configs.Config.DltsUrl, userName, paging.PageNum, paging.PageSize, urllib.PathEscape(paging.Keyword))

	vcRsp := &models.VCRsp{
	}

	err := DoRequest(url, "GET", nil, nil, vcRsp)
	if err == nil {
		return vcRsp, nil
	} else {
		fmt.Printf("create training err[%+v]\n", err)
		return nil, err
	}
}

func GetVCStatistic(userName string, req models.VCStatisticReq) (*models.VCStatisticRsp, error) {

	var url string
	var err error
	var jobRsp string

	if req.Type == models.VC_STATISTIC_JOB {

		url = fmt.Sprintf("%s/CountJobByStatus?userName=%s&targetStatus=%s&vcName=%s",
					configs.Config.DltsUrl, userName, urllib.PathEscape(req.TargetStatus),  urllib.PathEscape(req.VCName))

		err, jobRsp = DoGetRequest(url, nil, nil)
		if err == nil {
			count, err2 := strconv.Atoi(strings.Trim(jobRsp, "\n"))
			if err2 == nil {
				return &models.VCStatisticRsp{
					JobCount: count,
				}, nil
			} else {
				err = err2
			}
		}
	} else if req.Type == models.VC_STATISTIC_VC_CONFIG {

		// 返回：
		// 1. 未分配给VC的计算设备数
		// 2. 每个VC下的用户配额数量
		url = fmt.Sprintf("%s/GetAllDevice?userName=%s", configs.Config.DltsUrl, userName)

		devices := make(map[string]models.DeviceItem2)
		unallocated := make(map[string]int)
		alloc := make(map[string]int)
		userQuota := make(map[string]map[string]int)

		err = DoRequest(url, "GET", nil, nil, &devices)
		if err != nil {
			fmt.Printf("get all devices err[%+v]\n", err)
			return nil, err
		}

		// 获取集群设备总数
		total := make(map[string]int)
		for k, v := range(devices) {
			total[k] += v.Capacity
		}

		// 获取集群已分配给VC的设备总数
		page := models.Paging{}
		page.PageNum = 1
		page.PageSize = 9999

		var vcRsp *models.VCRsp
		vcRsp, err = ListVC(userName, page)

		if vcRsp != nil {

			for _, v := range(vcRsp.CurrPage) {

				quota := make(map[string]int)
				if len(*(v.Quota)) > 0 {
					err = json.Unmarshal([]byte(*(v.Quota)), &quota)
					if err == nil {
						userQuota[*(v.VCName)] = make(map[string]int)
						for deviceStr, num := range(quota) {
							alloc[deviceStr] += num
							userQuota[*(v.VCName)][deviceStr] = num
						}
					}
				}
			}
		}

		// 计算剩余数据
		for k, _ := range(devices) {
			if _, ok := alloc[k]; ok {
				if _, totalOk := total[k]; totalOk {
					unallocated[k] += total[k] - alloc[k]
				}
			}
		}

		return &models.VCStatisticRsp{
			UnallocatedDevice: unallocated,
			UserQuota: userQuota,
		}, nil

	} else if req.Type == models.VC_STATISTIC_USER_UNUSED {


	}

	fmt.Printf("get vc err[%+v]\n", err)
	return nil, err
}
