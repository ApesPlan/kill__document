package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// `gorm:""`
// EMBEDDED 将struct设置为嵌入
// EMBEDDED_PREFIX 设置内嵌结构的前缀名
// auto_increment;
// primary_key;
// primary_key:true;
// (指定列为唯一索引)unique;
// (创建唯一索引)unique_index;
// index;index:idx_name_code;
// not null;default:xxx;
// column:xxx;type:varchar(100);
// (默认值255)size:255;
// (指定列精度)precision:(10,2);
// many2many:user_products;

// AUTO_INCREMENT;
// PRIMARY_KEY;
// UNIQUE;
// UNIQUE_INDEX;
// INDEX;
// INDEX:idx_name_code;
// NOT NULL;
// DEFAULT:xxx;
// Column:xxx;
// Type:varchar(100);
// Size:255;
// PRECISION:(10,2);
// MANY2MANY:user_products;

// (指定联接的表名)MANY2MANY:user_products;
// (指定外键)FOREIGNKEY:xxx;
// (指定关联外键)ASSOCIATION_FOREIGNKEY:xxx;
// (指定多态类型)POLYMORPHIC;
// (指定多态值)POLYMORPHIC_VALUE;
// (指定联接表的外键)JOINTABLE_FOREIGNKEY:xxx;
// (指定联接表的联合外键)ASSOCIATION_JOINTABLE_FOREIGNKEY:xxx;
// (是否自动保存关联)SAVE_ASSOCIATIONS;
// (是否自动更新关联)ASSOCIATION_AUTOUPDATE;
// (是否自动创建关联)ASSOCIATION_AUTOCREATE;
// (是否自动保存关联引用)ASSOCIATION_SAVE_REFERENCE;
// (是否自动预加载关联)PRELOAD;

// Product ...
type Product struct {
  gorm.Model
  Name  string  `gorm:"default:'galeone'"`
	Code  string
	Price uint
}

func main() {
  // 数据库连接
  // 处理time.Time，您需要包括parseTime作为参数
	db, err := gorm.Open("mysql", "root:447728@/gorm?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
  defer db.Close()
  
  db.LogMode(true) // 开启debug

  // 全局禁用表名复数
	db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响

  // 检查模型`User`表是否存在
  db.HasTable(&Product{})
  // // 检查表`users`是否存在
  // db.HasTable("products")

  // 创建表时添加表后缀
  db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Product{})
	// // 自动迁移模式
  // db.AutoMigrate(&Product{})
  
  // // 为模型`Product`创建表
  // db.CreateTable(&Product{})
  // // 创建表`products'时将“ENGINE = InnoDB”附加到SQL语句
  // db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Product{})

  // // 添加主键
  // // 1st param : 外键字段
  // // 2nd param : 外键表(字段)
  // // 3rd param : ONDELETE
  // // 4th param : ONUPDATE
  // db.Model(&Product{}).AddForeignKey("city_id", "cities(id)", "RESTRICT", "RESTRICT")

  // // 为`name`列添加索引`idx_uproducts_name`
  // db.Model(&Product{}).AddIndex("idx_uproducts_name", "name")
  // // 为`name`, `code`列添加索引`idx_products_name_code`
  // db.Model(&Product{}).AddIndex("idx_products_name_code", "name", "code")
  // // 添加唯一索引
  // db.Model(&Product{}).AddUniqueIndex("idx_products_name", "name")
  // // 为多列添加唯一索引
  // db.Model(&Product{}).AddUniqueIndex("idx_products_name_code", "name", "code")
  // // 删除索引
  // db.Model(&Product{}).RemoveIndex("idx_products_name")

  productRecordPointer := &Product{Code: "L1212", Price: 1000}
  db.NewRecord(productRecordPointer) // => 主键为空返回`true`
	// 创建记录
  db.Create(productRecordPointer)
  db.NewRecord(productRecordPointer) // => 创建`user`后返回`false`

  //默认值
  var productRecord = Product{Code: "L1213", Name: ""}
  db.Create(&productRecord)// INSERT INTO product("code") values("L1213"); productRecord.Name => 'galeone'
  

	// 读取
	var product Product
	db.First(&product, 1)                   // 查询id为1的product
	db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
	// 更新 - 更新product的price为2000
  db.Model(&product).Update("Price", 2000)
  
  // 修改模型`Product`的code列的数据类型为`text`
  db.Model(&Product{}).ModifyColumn("code", "text")

  // 删除模型`Product`的code列
  db.Model(&Product{}).DropColumn("code")

	// 删除 - 删除product
  db.Delete(&product)
  // 删除模型`Product`的表
  db.DropTable(&Product{})
  // 删除表`products`
  db.DropTable("products")
  // 删除模型`Product`的表和表`products`
  db.DropTableIfExists(&Product{}, "products")
}



  
