package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Admin ...
type Admin struct {
	ID        int        `gorm:"AUTO_INCREMENT;TYPE:int(11);NOT NULL;PRIMARY_KEY;INDEX"`
	Name      string     `gorm:"TYPE: VARCHAR(255); DEFAULT:'';INDEX"`
	Companies []Company  `gorm:"FOREIGNKEY:AdminID;ASSOCIATION_FOREIGNKEY:ID"`
	CreatedAt time.Time  `gorm:"TYPE:DATETIME"`
	UpdatedAt time.Time  `gorm:"TYPE:DATETIME"`
	DeletedAt *time.Time `gorm:"TYPE:DATETIME;DEFAULT:NULL"`
}

// Company ...
type Company struct {
	gorm.Model
	Industry int    `gorm:"TYPE:INT(11);DEFAULT:0"`
	Name     string `gorm:"TYPE:VARCHAR(255);DEFAULT:'';INDEX"`
	Job      string `gorm:"TYPE:VARCHAR(255);DEFAULT:''"`
	AdminID   int    `gorm:"TYPE:int(11);NOT NULL;INDEX"`
}

func main() {
	// 数据库连接
	// 处理time.Time，您需要包括parseTime作为参数
	db, err := gorm.Open("mysql", "root:447728@/gorm?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// db.LogMode(true) // 开启debug

	// // 全局禁用表名复数
	// db.SingularTable(true) // 如果设置为true,`Admin`的默认表名为`Admin`,使用`TableName`设置的表名不受影响

	// // 更改默认表名
	// // 您可以通过定义DefaultTableNameHandler对默认表名应用任何规则。
	// gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	// 	return "prefix_" + defaultTableName;
	// }

	db.AutoMigrate(&Admin{},&Company{})

	// Related
	// 使用 Related 方法, 需要把 Admin 查询好, 然后根据 Admin 定义中指定的 FOREIGNKEY 去查找 Company, 如果没定义, 则调用时需要指定, 如下:

	var u Admin
	db.First(&u)
	// fmt.Println(&u) // &{0 dfsdf [] 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC <nil>}
	db.Model(&u).Related(&u.Companies).Find(&u.Companies)
	// Admin 列表时遍历列表一一查询 Company
	// fmt.Println(&u.Companies) // &[{{1 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC <nil>} 0 sdfsd sdfsdf 1}]

	// Association
	// 使用 Association 方法, 需要把 Admin 查询好, 然后根据 Admin 定义中指定的 AssociationForeignKey ASSOCIATION_FOREIGNKEY 去查找 Company, 必须定义, 如下:

	db.Model(&u).Association("Companies").Find(&u.Companies)
	// fmt.Println(&u.Companies)  // &[{{1 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC <nil>} 0 sdfsd sdfsdf 1}]

	// Preload
	// 使用 Preload 方法, 在查询 Admin 时先去获取 Company 的记录, 如下:

	// 查询单条 Admin
	db.Debug().Preload("Companies").First(&u)
	// // 对应的 sql 语句
	// // SELECT * FROM Admins LIMIT 1;
	// // SELECT * FROM companies WHERE Admin_id IN (1);
	// fmt.Println(&u) // &{1 dfsdf [{{1 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC <nil>} 0 sdfsd sdfsdf 1}] 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC <nil>}

	// 查询所有 Admin
	var list []Admin
	db.Debug().Preload("Companies").Find(&list)
	// // 对应的 sql 语句
	// // SELECT * FROM Admins;
	// // SELECT * FROM companies WHERE Admin_id IN (1,2,3...);
	// fmt.Println(&list) // &[{1 dfsdf [{{1 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC <nil>} 0 sdfsd sdfsdf 1}] 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC <nil>}]

}
