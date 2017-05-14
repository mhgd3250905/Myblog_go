package models

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"github.com/mattn/go-sqlite3"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	_DB_NAME        = "data/beeblog.db" //数据库名字
	_SQLITE3_DRIVER = "sqlite3"         //数据区引擎
)

type Category struct {
	Id            int64                                    //索引ID
	Title         string                                   //索引标题
	Created       time.Time `orm:"index" orm:"default(0)"` //索引创建时间，可索引排序
	Views         int64     `orm:"index" `                 //浏览次数，可索引排序
	TopicTime     time.Time `orm:"index" orm:"default(0)"` //文章时间，可索引排序
	TopicCount    int64                                    //文章数量
	TopicLastUser int64                                    //最后文章编辑者
}

type Topic struct {
	Id              int64   `form:"tid"`                      //文章ID
	Uid             int64   `form:"uid"`                      //用户ID
	Title           string  `form:"title"`                    //标题
	Category        string  `form:"category"`                 //分类
	Label           string  `form:"label"`                    //标签
	Content         string  `form:"content" orm:"size(5000)"` //内容
	Attachment      string                                    //浏览次数
	Created         time.Time `orm:"index"`                   //创建时间，可索引排序
	Updated         time.Time `orm:"index"`                   //更新时间，可索引排序
	Views           int64                                     //浏览次数
	Author          string                                    //作者
	ReplyTime       time.Time                                 //回复时间
	ReplyCount      int64                                     //回复数量
	ReplyLastUserId int64                                     //最后回复用户
}

type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string `orm:"size(1000)"`
	Created time.Time  `orm:"index"`
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

	cate := &Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	}

	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}

	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return err
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
	return err
}

func DelTopic(category, tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}
	if len(oldCate) > 0 {
		//更新分类统计
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", category).One(cate)
		if err == nil {
			cate.TopicTime = time.Now()
			cate.TopicCount--
			_, err = o.Update(cate)
			return err
		}
	}

	return err
}

//添加文章
func AddTopic(topic *Topic) error {
	//处理标签
	topic.Label = "$" + strings.Join(strings.Split(topic.Label, " "), "#$") + "#"
	//空格作为多个标签的分隔
	//"beego orm"-->[beego,orm]-->"$beego#$orm#"


	o := orm.NewOrm()
	topic.Created=time.Now()
	topic.Updated=time.Now()
	topic.ReplyTime=time.Now()

	_, err := o.Insert(topic)
	if err != nil {
		return err
	}
	//更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", topic.Category).One(cate)
	if err == nil {
		cate.TopicTime = time.Now()
		cate.TopicCount++
		_, err = o.Update(cate)
		return err
	}
	return err
}

////添加文章
//func AddTopic(title, category, label, content, attachment string) error {
//	//处理标签
//	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
//	//空格作为多个标签的分隔
//	//"beego orm"-->[beego,orm]-->"$beego#$orm#"
//
//	o := orm.NewOrm()
//	topic := &Topic{
//		Title:      title,
//		Content:    content,
//		Attachment: attachment,
//		Created:    time.Now(),
//		Updated:    time.Now(),
//		ReplyTime:  time.Now(),
//		Category:   category,
//		Label:      label,
//	}
//	_, err := o.Insert(topic)
//	if err != nil {
//		return err
//	}
//	//更新分类统计
//	cate := new(Category)
//	qs := o.QueryTable("category")
//	err = qs.Filter("title", category).One(cate)
//	if err == nil {
//		cate.TopicTime = time.Now()
//		cate.TopicCount++
//		_, err = o.Update(cate)
//		return err
//	}
//	return err
//}

/**
	获取所有文章
 */
func GetAllTopics(cate, label string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")

	var err error
	if isDesc { //这里是首页，设置时间倒序排序，以及是否有category过滤
		if len(cate) > 0 {
			qs = qs.Filter("category", cate)
		}
		if len(label) > 0 {
			qs = qs.Filter("label__contains", "$"+label+"#")
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

	topic.Label = strings.Replace(strings.Replace(topic.Label,
		"#", " ", -1), "$", "", -1)

	return topic, err
}



func ModifyTopic(tid, title, category, label, content, attachment string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	//处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	var oldCate, oldAttach string //旧的分类
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}

	if o.Read(topic) == nil {
		oldCate = topic.Category //获取旧的分类
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Content = content
		topic.Updated = time.Now()
		topic.Category = category
		topic.Label = label
		topic.Attachment = attachment
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}

	//删除旧的附件
	if len(oldAttach) >= 0 && !strings.EqualFold(oldAttach, attachment) {
		os.Remove(path.Join("attachment", oldAttach))
	}

	//更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
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

//取出回复
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

//增加回复
func AddReply(tid, nickName, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	reply := &Comment{
		Tid:     tidNum,
		Name:    nickName,
		Content: content,
		Created: time.Now(),
	}

	o := orm.NewOrm()
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}

	//修改文章的评论次数，最后评论时间
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		beego.Warn("\n", "topic.ReplyTime", topic.ReplyTime, " topic.ReplyCount", topic.ReplyCount)
		_, err = o.Update(topic)
	}
	return err
}

//删除回复
func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()

	var tidNum int64 //评论对应文章的Tid
	reply := &Comment{Id: ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
	}
	_, err = o.Delete(reply)
	if err != nil {
		return err
	}
	//获取所有的回复：1.计算个数 2.获取删除后的最后回复时间
	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}

	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		if len(replies) > 0 {
			topic.ReplyTime = replies[0].Created
			topic.ReplyCount = int64(len(replies))
		} else {
			topic.ReplyTime = time.Now()
			topic.ReplyCount = 0
		}
		_, err = o.Update(topic)
	}

	return err
}
