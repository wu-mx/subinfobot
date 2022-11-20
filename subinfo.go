package main

import (
	"errors"
	"fmt"
	_ "fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"subinfobot/utils"
	"time"
)

type Subinfo struct {
	Link       string
	ExpireTime string
	TimeRemain string
	Upload     string
	Download   string
	Used       string
	Total      string
	Expired    int //0:not Expired,1:Expired,2:unknown
	Available  int //0:Available,1:unavailable,2:unknown
	DataRemain string
}

func getSinf(link string) (error, Subinfo) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", link, nil)
	req.Header.Add("User-Agent", "ClashforWindows/0.19.21")
	if err != nil {
		return err, Subinfo{}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err, Subinfo{}
	}
	if res.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("获取失败，服务器返回了代码%s", strconv.Itoa(res.StatusCode))), Subinfo{}
	}
	if sinfo := res.Header["Subscription-Userinfo"]; sinfo == nil {
		return errors.New("未获取到订阅详细信息"), Subinfo{}
	} else {
		sinf := Subinfo{Link: link}
		sinfmap := make(map[string]int64)
		parseExp := regexp.MustCompile("[A-Za-z]+=[0-9]+")
		sslice := parseExp.FindAllStringSubmatch(sinfo[0], -1)
		for _, val := range sslice {
			kvslice := strings.Split(val[0], "=")
			if len(kvslice) == 2 {
				i, err := strconv.ParseInt(kvslice[1], 10, 64)
				if err == nil {
					sinfmap[kvslice[0]] = i
				}
			}
		}
		if upload, oku := sinfmap["upload"]; oku {
			sinf.Upload = utils.FormatFileSize(upload)
		} else {
			sinf.Upload = "没有说明捏"
		}
		if download, okd := sinfmap["download"]; okd {
			sinf.Download = utils.FormatFileSize(download)
		} else {
			sinf.Download = "没有说明捏"
		}
		if total, okt := sinfmap["total"]; okt {
			sinf.Total = utils.FormatFileSize(total)
			down, oka := sinfmap["download"]
			up, okb := sinfmap["upload"]
			if (oka) && (okb) {
				sinf.Used = utils.FormatFileSize(up + down)
				remain := total - (up + down)
				if remain >= 0 {
					if remain > 0 {
						sinf.Available = 0
						sinf.DataRemain = utils.FormatFileSize(remain)
					} else {
						sinf.Available = 1
						sinf.DataRemain = utils.FormatFileSize(remain)
					}
				} else {
					sinf.Available = 1
					sinf.DataRemain = "逾量" + utils.FormatFileSize(int64(math.Abs(float64(remain))))
				}
			} else {
				sinf.Available = 2
				sinf.DataRemain = "未知"
			}
		} else {
			sinf.Available = 2
			sinf.Total = "没有说明捏"
		}
		if exp, oke := sinfmap["expire"]; oke {
			//get expire time and remain time
			timeStamp := time.Now().Unix()
			timeExp := time.Unix(exp, 0)
			sinf.ExpireTime = timeExp.String()
			if timeStamp >= exp {
				sinf.Expired = 1
				sinf.Available = 1
				remain := timeExp.Sub(time.Now())
				if remain.Hours() > 24 {
					sinf.TimeRemain = "逾期<code>" + strconv.Itoa(int(math.Floor(remain.Hours()/24))) + "天" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Hours()))%24)))) + "小时" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Minutes()))%60)))) + "分" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Seconds()))%60)))) + "秒" + "</code>"
				} else if remain.Minutes() > 60 {
					sinf.TimeRemain = "逾期<code>" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Hours()))%24)))) + "小时" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Minutes()))%60)))) + "分" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Seconds()))%60)))) + "秒" + "</code>"
				} else if remain.Seconds() > 60 {
					sinf.TimeRemain = "逾期<code>" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Minutes()))%60)))) + "分" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Seconds()))%60)))) + "秒" + "</code>"
				} else {
					sinf.TimeRemain = "逾期<code>" + strconv.Itoa(int(math.Floor(remain.Seconds()))) + "秒" + "</code>"
				}
			} else {
				sinf.Expired = 0
				remain := timeExp.Sub(time.Now())
				if remain.Hours() > 24 {
					sinf.TimeRemain = "距离到期还有<code>" + strconv.Itoa(int(math.Floor(remain.Hours()/24))) + "天" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Hours()))%24)))) + "小时" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Minutes()))%60)))) + "分" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Seconds()))%60)))) + "秒" + "</code>"
				} else if remain.Minutes() > 60 {
					sinf.TimeRemain = "距离到期还有<code>" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Hours()))%24)))) + "小时" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Minutes()))%60)))) + "分" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Seconds()))%60)))) + "秒" + "</code>"
				} else if remain.Seconds() > 60 {
					sinf.TimeRemain = "距离到期还有<code>" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Minutes()))%60)))) + "分" + strconv.Itoa(int(math.Floor(float64(int(math.Floor(remain.Seconds()))%60)))) + "秒" + "</code>"
				} else {
					sinf.TimeRemain = "距离到期还有<code>" + strconv.Itoa(int(math.Floor(remain.Seconds()))) + "秒" + "</code>"
				}
			}
		} else {
			sinf.ExpireTime = "没有说明捏"
			sinf.TimeRemain = "未知"
		}
		return nil, sinf
	}
}
