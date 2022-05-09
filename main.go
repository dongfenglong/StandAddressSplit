package main

import (
	"bufio"
	"errors"
	"github.com/dongfenglong/StandAddressSplit/logger"
	"io"
	"os"
	"strings"
)

var err error
var Provinces, Cities []string

func init() {
	logger.Log, err = logger.NewLog()
	if err != nil {
		panic("Initial Log fail: " + err.Error())
	}

	//34个省级行政区
	Provinces = []string{
		"北京市",
		"天津市",
		"河北省",
		"山西省",
		"内蒙古自治区",
		"辽宁省",
		"吉林省",
		"黑龙江省",
		"上海市",
		"江苏省",
		"浙江省",
		"安徽省",
		"福建省",
		"江西省",
		"山东省",
		"河南省",
		"湖北省",
		"湖南省",
		"广东省",
		"广西壮族自治区",
		"海南省",
		"重庆市",
		"四川省",
		"贵州省",
		"云南省",
		"西藏自治区",
		"陕西省",
		"甘肃省",
		"青海省",
		"宁夏回族自治区",
		"新疆维吾尔自治区",
		"台湾省",
		"澳门特别行政区",
		"香港特别行政区",
	}

	//二级行政区：地级市、地区、自治州、盟
	//地区、盟列表
	Cities = []string{
		"大兴安岭地区",
		"阿里地区",
		"阿克苏地区",
		"喀什地区",
		"和田地区",
		"塔城地区",
		"阿勒泰地区",
		"兴安盟",
		"锡林郭勒盟",
		"阿拉善盟",
	}
}

func main() {
	var er error
	f, er := os.OpenFile("list.txt", os.O_RDONLY, 0644)
	if er != nil {
		panic("Open file fail:" + er.Error())
	}

	//创建Reader
	rd := bufio.NewReader(f)
	for {
		var addr, province, city, county, township, street string
		var line string
		addr, er = rd.ReadString('\n')
		switch er {
		case nil:
		case io.EOF:
			logger.Log.Println("Over")
			return
		default:
			logger.Log.Println("Read string fail:", er.Error())
		}

		line = addr

		province, addr, er = ReadProvince(addr)
		if er != nil {
			logger.Log.Printf("[Error]addr:%s, error:%s", line, er.Error())
			continue
		}

		switch strings.HasSuffix(province, "市") {
		case false: //非直辖市
			city, addr, er = ReadCity(addr)
			if er != nil {
				logger.Log.Printf("[Error]addr:%s, error:%s", line, er.Error())
				continue
			}
		case true:
			city = province
		}

		county, addr, er = ReadCounty(addr)
		if er != nil {
			logger.Log.Printf("[Error]addr:%s, error:%s", line, er.Error())
			continue
		}
		township, addr, er = ReadTownship(addr)
		if er != nil {
			logger.Log.Printf("[Error]addr:%s, error:%s", line, er.Error())
			continue
		}
		street, addr, er = ReadStreet(addr)
		if er != nil {
			logger.Log.Printf("[Error]addr:%s, error:%s", line, er.Error())
			//continue
		}

		logger.Log.Printf("province:%s|city:%s|county:%s|township:%s|street:%s|addr:%s|err:%#v", province, city, county, township, street, addr, er)
	}
}

// remaining: 剩余地址
func ReadProvince(addr string) (province, remaining string, e error) {
	for _, province = range Provinces {
		if strings.HasPrefix(addr, province) {
			remaining = strings.TrimLeft(addr, province)
			return
		}
	}

	e = errors.New("province not exists in address")
	remaining = addr
	return
}

//地级行政区：地级市、地区、自治州、盟
//中国共计333个地级行政区，包括293个地级市、7个地区、30个自治州、3个盟。
// remaining: 剩余地址
func ReadCity(addr string) (city, remaining string, e error) {
	//1.地区、盟
	for _, city = range Cities {
		if strings.HasPrefix(addr, city) {
			remaining = strings.TrimLeft(addr, city)
			return
		}
	}
	//2.自治州
	var idx int
	idx = strings.Index(addr, "自治州")
	if idx > 0 {
		city = addr[:idx+9]
		remaining = addr[idx+9:]
		return
	}
	//3.地级市：不规范地址可能会得到县级市
	idx = strings.Index(addr, "市")
	if idx > 0 {
		city = addr[:idx+3]
		remaining = addr[idx+3:]
		return
	}

	e = errors.New("city not exists in address")
	remaining = addr
	return
}

//县级行政区：地级市的市辖区、县级市、县、自治县、旗、自治旗、林区、特区等8种。
// remaining: 剩余地址
func ReadCounty(addr string) (county, remaining string, e error) {
	var idx int

	//1.林区
	idx = strings.Index(addr, "林区")
	if idx > 0 {
		county = addr[:idx+6]
		remaining = addr[idx+6:]
		logger.Log.Println("林区：", county)
		return
	}
	//2.特区
	idx = strings.Index(addr, "特区")
	if idx > 0 {
		county = addr[:idx+6]
		remaining = addr[idx+6:]
		logger.Log.Println("特区：", county)
		return
	}
	//3.旗、自治旗
	idx = strings.Index(addr, "旗")
	if idx > 0 {
		county = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("旗：", county)
		return
	}

	//4.县、自治县
	idx = strings.Index(addr, "县")
	if idx > 0 {
		county = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("县、自治县：", county)
		return
	}

	//5.矿区
	idx = strings.Index(addr, "矿区")
	if idx > 0 {
		county = addr[:idx+6]
		remaining = addr[idx+6:]
		logger.Log.Println("矿区：", county)
		return
	}
	//6.市辖区
	idx = strings.Index(addr, "区")
	if idx > 3 && (strings.Index(addr, "小区") != idx-3) { //排除“小区”
		county = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("市辖区：", county)
		return
	}
	//7.县级市
	idx = strings.Index(addr, "市")
	if idx > 0 {
		county = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("县级市：", county)
		return
	}

	e = errors.New("county not exists in address")
	remaining = addr
	return
}

//乡级行政区：即行政地位与乡相同的行政区，包括街道、镇、乡、民族乡、苏木、民族苏木、县辖区，为四级行政区。
// remaining: 剩余地址
func ReadTownship(addr string) (township, remaining string, e error) {
	var idx int

	//1.街道
	idx = strings.Index(addr, "街道")
	if idx > 0 {
		township = addr[:idx+6]
		remaining = addr[idx+6:]
		logger.Log.Println("街道：", township)
		return
	}
	//2.镇
	idx = strings.Index(addr, "镇")
	if idx > 0 {
		township = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("镇：", township)
		return
	}
	//3.旗、自治旗
	idx = strings.Index(addr, "旗")
	if idx > 0 {
		township = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("旗：", township)
		return
	}

	//4.乡、民族乡
	idx = strings.Index(addr, "乡")
	if idx > 0 {
		township = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("乡、民族乡：", township)
		return
	}

	//5.苏木、民族苏木
	idx = strings.Index(addr, "苏木")
	if idx > 0 {
		township = addr[:idx+6]
		remaining = addr[idx+6:]
		logger.Log.Println("苏木、民族苏木：", township)
		return
	}
	//6.县辖区：区公所，标准地址忽略区公所
	//idx = strings.Index(addr, "区公所")
	//if idx > 0 {
	//	township = addr[:idx+3]
	//	remaining = addr[idx+3:]
	//	logger.Log.Println("区公所：", township)
	//	return
	//}

	e = errors.New("township not exists in address")
	remaining = addr
	return
}

//道路：街、路、巷、弄、胡同、桥、道、国道、高速
func ReadStreet(addr string) (street, remaining string, e error) {
	var idx int
	//1.1辅路、
	idx = strings.Index(addr, "辅路")
	if idx > 0 {
		street = addr[:idx+6]
		remaining = addr[idx+6:]
		logger.Log.Println("辅路：", street)
		return
	}
	//1.2路、
	idx = strings.Index(addr, "路")
	if idx > 0 {
		street = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("路：", street)
		return
	}
	//2.街、
	idx = strings.Index(addr, "街")
	if idx > 0 {
		street = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("街：", street)
		return
	}
	//3.巷、
	idx = strings.Index(addr, "巷")
	if idx > 0 {
		street = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("巷：", street)
		return
	}
	//4.弄、
	idx = strings.Index(addr, "弄")
	if idx > 0 {
		street = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("弄：", street)
		return
	}
	//5.胡同、
	idx = strings.Index(addr, "胡同")
	if idx > 0 {
		street = addr[:idx+6]
		remaining = addr[idx+6:]
		logger.Log.Println("胡同：", street)
		return
	}
	//6.桥、
	idx = strings.Index(addr, "桥")
	if idx > 0 {
		street = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("桥：", street)
		return
	}
	//7.道、国道、
	idx = strings.Index(addr, "道")
	if idx > 0 {
		street = addr[:idx+3]
		remaining = addr[idx+3:]
		logger.Log.Println("道：", street)
		return
	}
	//8.高速
	idx = strings.Index(addr, "高速")
	if idx > 0 {
		street = addr[:idx+6]
		remaining = addr[idx+6:]
		logger.Log.Println("高速：", street)
		return
	}
	e = errors.New("street not exists in address")
	remaining = addr
	return
}
