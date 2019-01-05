package models

import ( _"github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/orm"

	"time"
)

type User1 struct {
	Id int
	Name string `orm:"unique"`
	Pwd string
	Articles []*Article `orm:"reverse(many)"`
}
type  Article struct {
	Id          int          `orm:"pk;auto"`
	Title       string       `orm:"size(100)"`
	Time        time.Time    `orm:"type(datetime);auto_now"`
	Count       int          `orm:"default(0)"`
	Img         string       `orm:"null"`
	Content     string
	ArticleType *ArticleType `orm:"rel(fk)"`
	Price       float64      `orm:"digits(10);decimals(2)"`
    Users       []*User1     `orm:"rel(m2m)"`
}
type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Articles  []*Article  `orm:"reverse(many)"`
}
func init()  {
	//注册一个数据库      别名                数据类型
 orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/Web1?charset=utf8")
//注册一个表
 orm.RegisterModel(new(User1),new(Article),new(ArticleType))
 //跑起来              别名             更新        可视
 orm.RunSyncdb("default",false,true)
}