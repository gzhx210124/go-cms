package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/validation"
	"log"
)

type User struct {
	Model
	Id          int    `json:"id"         form:"id"         gorm:"default:''"`
	LoginName   string `json:"login_name" form:"login_name" gorm:"default:''" valid:"Required;MaxSize(20);MinSize(2)"`
	UserName    string `json:"user_name"  form:"user_name"  gorm:"default:''" valid:"Required;MaxSize(20);MinSize(6)"`
	UserType    string `json:"user_type"  form:"user_type"  gorm:"default:'00'"`
	Email       string `json:"email"      form:"email"      gorm:"default:''" valid:"Email"`
	Phone       string `json:"phone"      form:"phone"      gorm:"default:''"`
	Phonenumber string `json:"phonenumber"form:"phonenumber"gorm:"default:''"`
	Sex         string `json:"sex"        form:"sex"        gorm:"default:'0'"`
	Avatar      string `json:"avatar"     form:"avatar"     gorm:"default:''"`
	Password    string `json:"password"   form:"password"   gorm:"default:''" valid:"Required;MaxSize(33);MinSize(6)"`
	Salt        string `json:"salt"       form:"salt"       gorm:"default:''"`
	Status      string `json:"status"     form:"status"     gorm:"default:'0'"`
	DelFlag     string `json:"del_flag"   form:"del_flag"   gorm:"default:'0'"`
	LoginIp     string `json:"login_ip"   form:"login_ip"   gorm:"default:''"`
	LoginDate   int64  `json:"login_date" form:"login_date" gorm:"default:''"`
	CreateBy    string `json:"create_by"  form:"create_by"  gorm:"default:''"`
	CreatedAt   int64  `json:"created_at" form:"created_at" gorm:"default:''"`
	UpdateBy    string `json:"update_by"  form:"update_by"  gorm:"default:''"`
	UpdatedAt   int64  `json:"updated_at" form:"updated_at" gorm:"default:''"`
	DeletedAt   int64  `json:"deleted_at" form:"deleted_at" gorm:"default:''"`
	Remark      string `json:"remark"     form:"remark"     gorm:"default:''"`
}

func NewUser() (user *User) {
	return &User{}
}

func (m *User) Pagination(offset, limit int, key string) (res []User, count int) {
	query := Db
	if key != "" {
		query = query.Where("name like ?", "%"+key+"%")
	}
	query.Offset(offset).Limit(limit).Order("id desc").Find(&res)
	query.Model(User{}).Count(&count)
	return
}

func (m *User) Create() (newAttr User, err error) {
	tx := Db.Begin()

	err = Db.Create(m).Error

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	newAttr = *m
	return
}

func (m *User) Update() (newAttr User, err error) {
	tx := Db.Begin()
	if m.Id > 0 {
		err = tx.Where("id=?", m.Id).Save(m).Error
	} else {
		err = errors.New("id参数错误")
	}
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	newAttr = *m
	return
}

func (m *User) Delete() (err error) {
	tx := Db.Begin()
	if m.Id > 0 {
		err = tx.Delete(m).Error
	} else {
		err = errors.New("id参数错误")
	}
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return
}

func (m *User) DelBatch(ids []int) (err error) {
	tx := Db.Begin()
	if len(ids) > 0 {
		err = tx.Where("id in (?)", ids).Delete(m).Error
	} else {
		err = errors.New("id参数错误")
	}
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return
}

func (m *User) FindById(id int) (user User, err error) {
	err = Db.Where("id=?", id).First(&user).Error
	return
}
func (m *User) FindByMap(page, pageCount int, dataMap map[string]interface{}) (user []User, total int, err error) {
	var where string="1=1"
	if status,ok:=dataMap["status"].(int);ok{
		where+=fmt.Sprintf(" AND status=%v",status)
	}
	if loginName,ok:=dataMap["login_name"].(string);ok{
		where+=fmt.Sprintf(" AND login_name like '%v'","%"+loginName+"%")
	}
	if userName,ok:=dataMap["user_name"].(string);ok{
		where+=fmt.Sprintf(" AND user_name like '%v'","%"+userName+"%")
	}
	if startTime,ok:=dataMap["start_time"].(int64);ok{
		where+=fmt.Sprintf(" AND created_at>=%v",startTime)
	}
	if endTime,ok:=dataMap["end_time"].(int64);ok{
		where+=fmt.Sprintf(" created_at<=%v",endTime)
	}
	if phone,ok:=dataMap["phone"].(string);ok{
		where+=fmt.Sprintf(" AND phone like '%v'","%"+phone+"%")
	}
	err = Db.Offset((page - 1) * pageCount).Limit(pageCount).Where(where).Order("created_at DESC").Find(&user).Error
	err = Db.Model(&User{}).Where(where).Count(&total).Error
	return


}

func (m *User) FindByMaps(page, pageSize int, dataMap map[string]interface{}) (user []User, total int, err error) {
	query := Db
	if status,isExist:=dataMap["status"].(int);isExist == true{
		query = query.Where("status = ?", status)
	}
	if loginName,ok:=dataMap["login_name"].(string);ok{
		query = query.Where("login_name LIKE ?", "%"+loginName+"%")
	}
	
	if userName,ok:=dataMap["user_name"].(string);ok{
		query = query.Where("user_name LIKE ?", "%"+userName+"%")
	}
	if startTime,ok:=dataMap["start_time"].(int64);ok{
		query = query.Where("created_at > ?", startTime)
	}
	if endTime,ok:=dataMap["end_time"].(int64);ok{
		query = query.Where("created_at <= ?", endTime)
	}
	if phone,ok:=dataMap["phone"].(string);ok{
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	
	// 获取取指page，指定pagesize的记录
	err = query.Offset((page-1)*pageSize).Order("created_at desc").Limit(pageSize).Find(&user).Error
	if err == nil{
		err = query.Model(&User{}).Count(&total).Error
	}
	return
}

/*****************************************************************新增加的方法*****************************************************************/

func (m *User) FindByUserName(user_name string) (user User, err error) {
	err = Db.Select("id,user_name,password,salt").Where("user_name=?", user_name).First(&user).Error
	return
}

//验证用户信息
func checkUser(m *User) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&m)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}
