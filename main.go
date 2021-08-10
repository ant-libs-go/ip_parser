/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-08-31 20:06:45
# File Name: main.go
# Description:
####################################################################### */

// see: https://github.com/itbdw/ip-database/blob/master/src/IpLocation.php
// see: https://github.com/yinheli/qqwry

/*
use:
func main() {

	dat, _ := ioutil.ReadFile("./qqwry.dat")
	ip := "112.224.67.58"
	r := NewIpParser(ip, dat).Parse()
	fmt.Println(fmt.Sprintf("%+v", r))
}
*/

package ip_parser

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/axgle/mahonia"
)

const (
	INDEX_LEN       = 7
	REDIRECT_MODE_1 = 0x01
	REDIRECT_MODE_2 = 0x02
)

var (
	Provinces         = []string{"北京", "天津", "重庆", "上海", "河北", "山西", "辽宁", "吉林", "黑龙江", "江苏", "浙江", "安徽", "福建", "江西", "山东", "河南", "湖北", "湖南", "广东", "海南", "四川", "贵州", "云南", "陕西", "甘肃", "青海", "台湾", "内蒙古", "广西", "宁夏", "新疆", "西藏", "香港", "澳门"}
	DirectlyCitys     = []string{"北京", "天津", "重庆", "上海"}
	ProvincialCapital = map[string]string{
		"北京市":  "北京市",
		"天津市":  "天津市",
		"重庆市":  "重庆市",
		"上海市":  "上海市",
		"河北省":  "石家庄市",
		"山西省":  "太原市",
		"辽宁省":  "沈阳市",
		"吉林省":  "长春市",
		"黑龙江省": "哈尔滨市",
		"江苏省":  "南京市",
		"浙江省":  "杭州市",
		"安徽省":  "合肥市",
		"福建省":  "福州市",
		"江西省":  "南昌市",
		"山东省":  "济南市",
		"河南省":  "郑州市",
		"湖北省":  "武汉市",
		"湖南省":  "长沙市",
		"广东省":  "广州市",
		"海南省":  "海口市",
		"四川省":  "成都市",
		"贵州省":  "贵阳市",
		"云南省":  "昆明市",
		"陕西省":  "西安市",
		"甘肃省":  "兰州市",
		"青海省":  "西宁市",
		"台湾省":  "台北市",
		"内蒙古省": "呼和浩特市",
		"广西省":  "南宁市",
		"宁夏省":  "银川市",
		"新疆省":  "乌鲁木齐市",
		"西藏省":  "拉萨市",
		"香港省":  "香港市",
		"澳门省":  "澳门市"}
	Carriers = []string{"联通", "移动", "铁通", "电信", "长城", "聚友"}
)

type Info struct {
	Country  string
	Province string
	City     string
	County   string
	Area     string
	Carrier  string
}

type IpParser struct {
	ip  string
	dat []byte
	ptr uint32
}

func NewIpParser(ip string, dat []byte) *IpParser {
	o := &IpParser{ip: ip, dat: dat, ptr: 0}
	return o
}

func (this *IpParser) Parse() (r *Info) {
	offset := this.search(binary.BigEndian.Uint32(net.ParseIP(this.ip).To4()))
	if offset == 0 {
		return
	}
	place, carrier := this.parseQqwry(offset)
	r = this.parsePlace(place)
	r.Carrier = this.parseCarrier(carrier)
	return
}

func (this *IpParser) parseCarrier(carrier string) string {
	for _, c := range Carriers {
		if idx := strings.Index(carrier, c); idx == -1 {
			continue
		}
		return c
	}
	return ""
}

func (this *IpParser) parsePlace(place string) (r *Info) {
	r = &Info{}
	if idx := strings.Index(place, "省"); idx > -1 {
		r.Province = fmt.Sprintf("%s省", place[:idx])
		place = strings.TrimLeft(place, r.Province)
	} else {
		for _, province := range Provinces {
			if idx := strings.Index(place, province); idx == -1 {
				continue
			}

			isDirectlyCity := false
			for _, city := range DirectlyCitys {
				if city == province {
					isDirectlyCity = true
				}
			}

			if isDirectlyCity == true {
				r.Province = fmt.Sprintf("%s市", province)
				r.City = r.Province
			} else {
				r.Province = fmt.Sprintf("%s省", province)
			}
			place = strings.TrimLeft(place, r.Province)
		}
	}

	if len(r.Province) > 0 {
		r.Country = "中国"
	} else {
		r.Country = place
		return
	}

	if idx := strings.Index(place, "市"); idx > -1 {
		r.City = fmt.Sprintf("%s市", place[:idx])
		place = strings.TrimLeft(place, r.City)
	}
	if idx := strings.Index(place, "县"); idx > -1 {
		r.County = fmt.Sprintf("%s县", place[:idx])
		place = strings.TrimLeft(place, r.County)
	}
	if idx := strings.Index(place, "区"); idx > -1 {
		r.Area = fmt.Sprintf("%s区", place[:idx])
		place = strings.TrimLeft(place, r.Area)
	}

	// Fix lack of city field, using provincial capital fill
	if len(r.City) == 0 {
		r.City = ProvincialCapital[r.Province]
	}
	return
}

func (this *IpParser) parseQqwry(offset uint32) (place string, carrier string) {
	var placeBytes []byte
	var carrierBytes []byte
	mode := this.readMode(offset + 4)
	switch mode {
	case REDIRECT_MODE_1:
		placeOffset := this.readUInt24()
		mode = this.readMode(placeOffset)
		if mode == REDIRECT_MODE_2 {
			c := this.readUInt24()
			placeBytes = this.readString(c)
			placeOffset += 4
		} else {
			placeBytes = this.readString(placeOffset)
			placeOffset += uint32(len(placeBytes) + 1)
		}
		carrierBytes = this.readCarrier(placeOffset)
	case REDIRECT_MODE_2:
		placeOffset := this.readUInt24()
		placeBytes = this.readString(placeOffset)
		carrierBytes = this.readCarrier(offset + 8)
	default:
		placeBytes = this.readString(offset + 4)
		carrierBytes = this.readCarrier(offset + uint32(5+len(placeBytes)))
	}

	enc := mahonia.NewDecoder("gbk")
	place = enc.ConvertString(string(placeBytes))
	carrier = enc.ConvertString(string(carrierBytes))

	if place == " CZ88.NET" || place == "纯真网络" {
		place = ""
	}
	if carrier == " CZ88.NET" {
		carrier = ""
	}
	return
}

func (this *IpParser) readCarrier(offset uint32) []byte {
	mode := this.readMode(offset)
	if mode == REDIRECT_MODE_1 || mode == REDIRECT_MODE_2 {
		carrierOffset := this.readUInt24()
		if carrierOffset == 0 {
			return []byte("")
		} else {
			return this.readString(carrierOffset)
		}
	} else {
		return this.readString(offset)
	}
	return []byte("")
}

func (this *IpParser) readString(offset uint32) []byte {
	data := make([]byte, 0, 30)
	for {
		this.ptr = offset + 1
		buf := this.dat[this.ptr-1 : this.ptr]
		if buf[0] == 0 {
			break
		}
		offset++
		data = append(data, buf[0])
	}
	return data
}

func (this *IpParser) readMode(offset uint32) byte {
	this.ptr = offset + 1
	return this.dat[this.ptr-1 : this.ptr][0]
}

func (this *IpParser) readUInt24() (r uint32) {
	this.ptr += 3
	return this.byte3ToUInt32(this.dat[this.ptr-3 : this.ptr])
}

func (this *IpParser) search(target uint32) uint32 {
	start := binary.LittleEndian.Uint32(this.dat[:4])
	end := binary.LittleEndian.Uint32(this.dat[4:8])

	for start < end {
		mid := (((end-start)/INDEX_LEN)>>1)*INDEX_LEN + start

		cur := binary.LittleEndian.Uint32(this.dat[mid : mid+4])
		// ??
		if start+INDEX_LEN == end {
			t := this.byte3ToUInt32(this.dat[mid+4 : mid+INDEX_LEN])
			if target < binary.LittleEndian.Uint32(this.dat[mid+INDEX_LEN:mid+INDEX_LEN+4]) {
				return t
			}
			return 0
		}
		if cur == target {
			return this.byte3ToUInt32(this.dat[mid+4 : mid+INDEX_LEN])
		}
		if cur < target {
			start = mid
		}
		if cur > target {
			end = mid
		}
	}
	return 0
}

func (this *IpParser) byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
