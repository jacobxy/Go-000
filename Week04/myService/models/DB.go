package models


import (
	"log"

	"github.com/astaxie/beego/orm"
)

func ReadNew(newid int) *New {
	o := orm.NewOrm()
	newinfo := &New{}
	err := o.Read(&newinfo)
	if err != nil {
		log.Println(err)
	}
	return newinfo
}

func ReadNews(news []int32) []*New {
	res := make([]*New, 0, 100)
	o := orm.NewOrm()
	o.QueryTable("news").Filter("id__in", news).All(&res)
	return res
}

func ReadNewBeginEnd(begin, end int64) []*New {
	return nil
}

func InsertNews(nn *New) *New {
	o := orm.NewOrm()
	n, err := o.Insert(nn)
	if err != nil || n == 0 {
		return nil
	}
	log.Println(nn)
	return nn
}