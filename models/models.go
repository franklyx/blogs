// Created by FrankLyx on 2018/1/31
// Desc :
package models

import (
	"time"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/Unknwon/com"
	_ "github.com/mattn/go-sqlite3"

	"github.com/astaxie/beego/orm"
)

const (
	_DB_NAME        = "data/beeblog.db"
	_SQLITE3_DRIVER = "sqlite3"
)

// 分类
type Category struct {
	Id              int64
	Title           string
	Views           int64     `orm:"index"`
	Created         time.Time `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64     `orm:"default(1)"`
	TopicLastUserId int64
}

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category        string
	Labels          string
	Content         string    `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

// 评论
type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB() {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)

}

func AddCategory(name string) error {
	o := orm.NewOrm()
	created := time.Now()
	cate := &Category{Title: name, Created: created, TopicTime: created}
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cid}

	_, err = o.Delete(cate)
	if err != nil {
		return err
	}
	return nil
}

func AddTopic(title, category, label, content, attachment string) error {
	// 处理标签
	// 空格作为多个标签的分割符·
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	o := orm.NewOrm()
	created := time.Now()
	topic := &Topic{
		Title:      title,
		Category:   category,
		Labels:     label,
		Content:    content,
		Created:    created,
		Updated:    created,
		ReplyTime:  created,
		Attachment: attachment,
	}

	_, err := o.Insert(topic)
	if err != nil {
		return err
	}

	// 更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		//	如果不存在，简单忽略操作
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return err

}

func GetAllTopics(cate, label string, isDesc bool) ([]*Topic, error) {
	var err error
	o := orm.NewOrm()

	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")

	if isDesc {
		if len(cate) > 0 {
			qs = qs.Filter("category", cate)
		}
		if len(label) > 0 {
			qs = qs.Filter("labels__contains", "$"+label+"#")
		}
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}

	return topics, err
}

func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")

	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}

	topic.Views++
	_, err = o.Update(topic)

	// 转换标签格式
	topic.Labels = strings.Replace(strings.Replace(
		topic.Labels, "#", " ", -1), "$", "", -1)

	return topic, nil
}

func ModifyTopic(tid, title, category, label, content, attachment string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	// 处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	var oldCate, oldAttath string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}

	if o.Read(topic) == nil {
		oldCate = topic.Category
		oldAttath = topic.Attachment
		topic.Labels = label
		topic.Title = title
		topic.Content = content
		topic.Updated = time.Now()
		topic.Category = category
		topic.Attachment = attachment

		o.Update(topic)
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	// 更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}

	// 删除附件
	if len(oldAttath) > 0 {
		os.Remove(path.Join("attachment", oldAttath))
	}

	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}

	return nil
}

func DelTopic(id string) error {
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Topic{Id: tid}

	_, err = o.Delete(cate)
	if err != nil {
		return err
	}
	return nil
}

func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}

	o := orm.NewOrm()
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}

	// 更新文章
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_, err = o.Update(topic)
	}

	return err
}

func GetAllReplies(tid string) (replies []*Comment, err error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	replies = make([]*Comment, 0)

	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)

	return replies, err
}

func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}

	var tidNum int64
	o := orm.NewOrm()

	// 删除评论
	reply := &Comment{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_, err := o.Delete(reply)
		if err != nil {
			return err
		}
	}

	// 更新评论数量
	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		lenReplies := int64(len(replies))
		// 简单处理删除最后一个回复的问题
		if lenReplies != 0 {
			topic.ReplyTime = replies[0].Created
		}
		topic.ReplyCount = lenReplies
		_, err = o.Update(topic)
	}

	return err
}
