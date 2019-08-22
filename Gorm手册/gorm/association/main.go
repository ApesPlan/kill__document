package main

import (
	"database/sql"
	"fmt"
	"time"

	// go get -u -v github.com/google/uuid
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// // TestUser 只需要字段 `ID`, `CreatedAt`
// type TestUser struct {
// 	ID        uint
// 	CreatedAt time.Time
// 	Name      string
// }

// // Model 基本模型的定义
// type Model struct {
// 	ID        uint `gorm:"primary_key"`

// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt *time.Time
// }

// // 只需要字段 `ID`, `CreatedAt`
// type User struct {
// 	ID        uint // 列名为 `id` // 字段`ID`为默认主键
// //  AnimalId int64 `gorm:"primary_key"` // 设置AnimalId为主键
// 	CreatedAt time.Time // 列名为 `created_at`
// 	Name      string
// }

// User ...
type User struct {
	gorm.Model
	Birthday time.Time
	Age      int    `gorm:"column:age"`     // 设置列名为`age`
	Name     string `gorm:"size:255"`       // string默认长度为255, 使用这种tag重设。
	Num      int    `gorm:"AUTO_INCREMENT"` // 自增

	// 属于
	// `User`属于`Profile`, `ProfileID`为外键
	Profile   Profile
	ProfileID int
	// 指定外键
	// Profile      Profile `gorm:"ForeignKey:ProfileRefer"` // 使用ProfileRefer作为外键
	// ProfileRefer int
	// 指定外键和关联外键
	// Profile   Profile `gorm:"ForeignKey:ProfileID;AssociationForeignKey:Refer"`
	// ProfileID int

	// 包含一个
	// User 包含一个 CreditCard, UserID 为外键
	CreditCard CreditCard // One-To-One (拥有一个 - CreditCard表的UserID作外键)
	// 指定外键
	// CreditCard CreditCard `gorm:"ForeignKey:UserRefer"`
	// 指定外键和关联外键
	// CreditCard CreditCard `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"`

	// 包含多个
	Emails []Email // One-To-Many (拥有多个 - Email表的UserID作外键)
	// 指定外键
	// Emails     []Email `gorm:"ForeignKey:UserRefer"`
	// 指定外键和关联外键
	// Emails     []Email `gorm:"ForeignKey:UserID;AssociationForeignKey:Refer"`
	// Refer   string

	BillingAddress   Address // One-To-One (属于 - 本表的BillingAddressID作外键)
	BillingAddressID sql.NullInt64

	ShippingAddress   Address // One-To-One (属于 - 本表的ShippingAddressID作外键)
	ShippingAddressID int

	IgnoreMe  int        `gorm:"-"`                                                                                 // 忽略这个字段
	Languages []Language `gorm:"many2many:user_languages;(User)ForeignKey:ID;(CreditCard)AssociationForeignKey:ID"` // Many-To-Many , 'user_languages'是连接表
}

// TableName 设置User的表名为`users`
func (User) TableName() string {
	return "users"
}

// func (u User) TableName() string {
//     if u.Role == "admin" {
//         return "admin_users"
//     } else {
//         return "users"
//     }
// }

// BeforeCreate 在Callbacks中设置主键
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New())
	return nil
}

// Profile ...
type Profile struct {
	gorm.Model
	// Refer string // 关联外键
	Name string
}

// Email ...
type Email struct {
	ID         int
	UserID     int    `gorm:"index"`                          // 外键 (属于), tag `index`是为该列创建索引
	Email      string `gorm:"type:varchar(100);unique_index"` // `type`设置sql类型, `unique_index` 为该列设置唯一索引
	Subscribed bool
	// 指定外键
	// UserRefer uint
}

// Address ...
type Address struct {
	ID       int
	Address1 string         `gorm:"not null;unique"` // 设置字段为非空并唯一
	Address2 string         `gorm:"type:varchar(100);unique"`
	Post     sql.NullString `gorm:"not null"`
}

// Language ...
type Language struct {
	ID   int    `gorm:"primary_key:true"`
	Name string `gorm:"index:idx_name_code"` // 创建索引并命名，如果找到其他相同名称的索引则创建组合索引
	Code string `gorm:"index:idx_name_code"` // `unique_index` also works
}

// CreditCard ...
type CreditCard struct {
	gorm.Model
	UserID uint
	Number string
	Refer  string
	// UserRefer uint
}

// Product ...
type Product struct {
	gorm.Model
	Name  string `gorm:"default:'galeone'"`
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

	// // 全局禁用表名复数
	// db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响

	// // 更改默认表名
	// // 您可以通过定义DefaultTableNameHandler对默认表名应用任何规则。
	// gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	// 	return "prefix_" + defaultTableName;
	// }

	db.AutoMigrate(&User{}, &Profile{}, &Email{}, &Address{}, &Language{}, &CreditCard{})

	var user User
	db.Create(&user) // 将会设置`CreatedAt`为当前时间
	// 要更改它的值, 你需要使用`Update`
	db.Model(&user).Update("CreatedAt", time.Now())

	db.Save(&user)                           // 将会设置`UpdatedAt`为当前时间
	db.Model(&user).Update("name", "jinzhu") // 将会设置`UpdatedAt`为当前时间

	// 属于
	var profile Profile
	userProfileUser := &User{}
	db.First(userProfileUser)
	db.Model(userProfileUser).Related(&profile, "Profile") // SELECT * FROM profiles WHERE user_id = 1; // 1 是 user 的主键
	// db.Model(userProfileUser).Related(&profile)
	// fmt.Println(&profile) // &{{1 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC <nil>} wang}

	// 包含一个
	var card CreditCard
	userCreditCardUser := &User{}
	db.First(userCreditCardUser)
	// CreditCard是user的字段名称，这意味着获得user的CreditCard关系并将其填充到变量
	// 如果字段名与变量的类型名相同，如上例所示，可以省略，如：
	db.Model(userCreditCardUser).Related(&card, "CreditCard") // SELECT * FROM credit_cards WHERE user_id = 1; // 1 是 user 的主键
	// db.Model(userCreditCardUser).Related(&card)
	// fmt.Println(&card)

	// 包含多个
	var emails Email
	userEmailUser := &User{}
	db.First(userEmailUser)
	db.Model(userEmailUser).Related(&emails, "Emails") // SELECT * FROM emails WHERE user_id = 1; // 1 是 user 的主键
	// db.Model(userEmailUser).Related(&emails)
	// fmt.Println(&emails)

	var languages Language
	userLanguageUser := &Language{}
	db.First(userLanguageUser)
	db.Model(userLanguageUser).Related(&languages, "Languages")
	// SELECT * FROM "languages" INNER JOIN "user_languages" ON "user_languages"."language_id" = "languages"."id" WHERE "user_languages"."user_id" = 1
	// db.Model(userLanguageUser).Related(&languages)
	// fmt.Println(&languages)

	productRecordPointer := &Product{Code: "L1212", Price: 1000}
	db.NewRecord(productRecordPointer) // => 主键为空返回`true`
	// 创建记录
	db.Create(productRecordPointer)
	db.NewRecord(productRecordPointer) // => 创建`user`后返回`false`

	//默认值
	var productRecord = Product{Code: "L1213", Name: ""}
	db.Create(&productRecord) // INSERT INTO product("code") values("L1213"); productRecord.Name => 'galeone'

	// // 扩展创建选项
	// // 为Instert语句添加扩展SQL选项
	// db.Set("gorm:insert_option", "ON CONFLICT").Create(&productRecord) // INSERT INTO products (name, code) VALUES ("name", "code") ON CONFLICT;

	// 查询
	// 获取第一条记录，按主键排序
	db.First(&user)	// SELECT * FROM users ORDER BY id LIMIT 1;
	// 获取最后一条记录，按主键排序
	db.Last(&user)	// SELECT * FROM users ORDER BY id DESC LIMIT 1;
	// 获取所有记录
	db.Find(&users)	// SELECT * FROM users;
	// 使用主键获取记录
	db.First(&user, 10)	// SELECT * FROM users WHERE id = 10;

	// Where查询条件 (简单SQL)
	db.Where("name <> ?","jinzhu").Where("age >= ? and role <> ?",20,"admin").Find(&users)//// SELECT * FROM users WHERE name <> 'jinzhu' AND age >= 20 AND role <> 'admin';
	// 获取第一个匹配记录
	db.Where("name = ?", "jinzhu").First(&user)//// SELECT * FROM users WHERE name = 'jinzhu' limit 1;
	// 获取所有匹配记录
	db.Where("name = ?", "jinzhu").Find(&users)//// SELECT * FROM users WHERE name = 'jinzhu';
	// 不等
	db.Where("name <> ?", "jinzhu").Find(&users)
	// IN
	db.Where("name in (?)", []string{"jinzhu", "jinzhu 2"}).Find(&users)
	// LIKE
	db.Where("name LIKE ?", "%jin%").Find(&users)
	// AND
	db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
	// 大于
	db.Where("updated_at > ?", lastWeek).Find(&users)
	// BETWEEN
	db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)

	// Where查询条件 (Struct & Map)
	// Struct
	db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)	// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 LIMIT 1;
	// Map
	db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)	// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;
	// 主键的Slice
	db.Where([]int64{20, 21, 22}).Find(&users)	// SELECT * FROM users WHERE id IN (20, 21, 22);

	// Not条件查询
	db.Not("name", "jinzhu").First(&user)	// SELECT * FROM users WHERE name <> "jinzhu" LIMIT 1;
	// Not In
	db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)	// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");
	// Not In slice of primary keys
	db.Not([]int64{1,2,3}).First(&user)	// SELECT * FROM users WHERE id NOT IN (1,2,3);
	db.Not([]int64{}).First(&user)	// SELECT * FROM users;
	// Plain SQL
	db.Not("name = ?", "jinzhu").First(&user)	// SELECT * FROM users WHERE NOT(name = "jinzhu");
	// Struct
	db.Not(User{Name: "jinzhu"}).First(&user)	// SELECT * FROM users WHERE name <> "jinzhu";

	// 带内联条件的查询
	// 按主键获取
	db.First(&user, 23)	// SELECT * FROM users WHERE id = 23 LIMIT 1;
	// 简单SQL
	db.Find(&user, "name = ?", "jinzhu")	// SELECT * FROM users WHERE name = "jinzhu";

	db.Find(&users, "name <> ? AND age > ?", "jinzhu", 20)	// SELECT * FROM users WHERE name <> "jinzhu" AND age > 20;
	// Struct
	db.Find(&users, User{Age: 20})	// SELECT * FROM users WHERE age = 20; // 条件User{Age: 20}
	// Map
	db.Find(&users, map[string]interface{}{"age": 20})	// SELECT * FROM users WHERE age = 20;

	// Or条件查询
	db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)	// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
	// Struct
	db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2"}).Find(&users)	// SELECT * FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2';
	// Map
	db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2"}).Find(&users)

	// 查询链
	db.Where("name <> ?","jinzhu").Where("age >= ? and role <> ?",20,"admin").Find(&users)//// SELECT * FROM users WHERE name <> 'jinzhu' AND age >= 20 AND role <> 'admin';
	db.Where("role = ?", "admin").Or("role = ?", "super_admin").Not("name = ?", "jinzhu").Find(&users)

	// 扩展查询选项
	// 为Select语句添加扩展SQL选项
	db.Set("gorm:query_option", "FOR UPDATE").First(&user, 1)//// SELECT * FROM users WHERE id = 1 FOR UPDATE;

	// FirstOrInit
	// 获取第一个匹配的记录，或者使用给定的条件初始化一个新的记录（仅适用于struct，map条件）
	// Unfound
	db.FirstOrInit(&user, User{Name: "non_existing"})	// user -> User{Name: "non_existing"}
	// Found
	db.Where(User{Name: "Jinzhu"}).FirstOrInit(&user)	// user -> User{Id: 111, Name: "Jinzhu", Age: 20}
	db.FirstOrInit(&user, map[string]interface{}{"name": "jinzhu"})	// user -> User{Id: 111, Name: "Jinzhu", Age: 20}

	// Attrs
	// 如果未找到记录，则使用参数初始化结构
	// Unfound
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrInit(&user)	// SELECT * FROM users WHERE name = 'non_existing';//// user -> User{Name: "non_existing", Age: 20}
	db.Where(User{Name: "non_existing"}).Attrs("age", 20).FirstOrInit(&user)	// SELECT * FROM users WHERE name = 'non_existing';//// user -> User{Name: "non_existing", Age: 20}
	// Found
	db.Where(User{Name: "Jinzhu"}).Attrs(User{Age: 30}).FirstOrInit(&user)	// SELECT * FROM users WHERE name = jinzhu';//// user -> User{Id: 111, Name: "Jinzhu", Age: 20}

	// Assign
	// 将参数分配给结果，不管它是否被找到
	// Unfound
	db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrInit(&user)	// user -> User{Name: "non_existing", Age: 20}
	// Found
	db.Where(User{Name: "Jinzhu"}).Assign(User{Age: 30}).FirstOrInit(&user)	// SELECT * FROM users WHERE name = jinzhu';
	// user -> User{Id: 111, Name: "Jinzhu", Age: 30}

	// FirstOrCreate
	// 获取第一个匹配的记录，或创建一个具有给定条件的新记录（仅适用于struct, map条件）
	// Unfound
	db.FirstOrCreate(&user, User{Name: "non_existing"})	// INSERT INTO "users" (name) VALUES ("non_existing");
	// user -> User{Id: 112, Name: "non_existing"}
	// Found
	db.Where(User{Name: "Jinzhu"}).FirstOrCreate(&user)	// user -> User{Id: 111, Name: "Jinzhu"}

	// Attrs
	// 如果未找到记录，则为参数分配结构
	// Unfound
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: 20}).FirstOrCreate(&user)	// SELECT * FROM users WHERE name = 'non_existing';
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{Id: 112, Name: "non_existing", Age: 20}
	// Found
	db.Where(User{Name: "jinzhu"}).Attrs(User{Age: 30}).FirstOrCreate(&user)	// SELECT * FROM users WHERE name = 'jinzhu';
	// user -> User{Id: 111, Name: "jinzhu", Age: 20}

	// Assign
	// 将其分配给记录，而不管它是否被找到，并保存回数据库。
	// Unfound
	db.Where(User{Name: "non_existing"}).Assign(User{Age: 20}).FirstOrCreate(&user)	// SELECT * FROM users WHERE name = 'non_existing';
	// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
	// user -> User{Id: 112, Name: "non_existing", Age: 20}
	// Found
	db.Where(User{Name: "jinzhu"}).Assign(User{Age: 30}).FirstOrCreate(&user)	// SELECT * FROM users WHERE name = 'jinzhu';	
	// UPDATE users SET age=30 WHERE id = 111;
	// user -> User{Id: 111, Name: "jinzhu", Age: 30}

	// Select
	// 指定要从数据库检索的字段，默认情况下，将选择所有字段;
	db.Select("name, age").Find(&users)	// SELECT name, age FROM users;
	db.Select([]string{"name", "age"}).Find(&users)	// SELECT name, age FROM users;
	db.Table("users").Select("COALESCE(age,?)", 42).Rows()	// SELECT COALESCE(age,'42') FROM users;

	// Order
	// 在从数据库检索记录时指定顺序，将重排序设置为true以覆盖定义的条件
	db.Order("age desc, name").Find(&users)	// SELECT * FROM users ORDER BY age desc, name;
	// Multiple orders
	db.Order("age desc").Order("name").Find(&users)	// SELECT * FROM users ORDER BY age desc, name;
	// ReOrder
	db.Order("age desc").Find(&users1).Order("age", true).Find(&users2)	// SELECT * FROM users ORDER BY age desc; (users1)
	// SELECT * FROM users ORDER BY age; (users2)

	// Limit
	// 指定要检索的记录数
	db.Limit(3).Find(&users)	// SELECT * FROM users LIMIT 3;
	// Cancel limit condition with -1
	db.Limit(10).Find(&users1).Limit(-1).Find(&users2)	// SELECT * FROM users LIMIT 10; (users1)
	// SELECT * FROM users; (users2)

	// Offset
	// 指定在开始返回记录之前要跳过的记录数
	db.Offset(3).Find(&users)	// SELECT * FROM users OFFSET 3;
	// Cancel offset condition with -1
	db.Offset(10).Find(&users1).Offset(-1).Find(&users2)	// SELECT * FROM users OFFSET 10; (users1)
	// SELECT * FROM users; (users2)

	// Count
	// 获取模型的记录数
	db.Where("name = ?", "jinzhu").Or("name = ?", "jinzhu 2").Find(&users).Count(&count)	// SELECT * from USERS WHERE name = 'jinzhu' OR name = 'jinzhu 2'; (users)
	// SELECT count(*) FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2'; (count)
	db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)	// SELECT count(*) FROM users WHERE name = 'jinzhu'; (count)
	db.Table("deleted_users").Count(&count)	// SELECT count(*) FROM deleted_users;

	// Group & Having
	// rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()for rows.Next() {
	// 	...
	// }
	// rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()for rows.Next() {
	// 	...
	// }
	type Result struct {
		Date  time.Time
		Total int64
	}
	var results Result
	db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)

	// Join
	// 指定连接条件
	// rows, err := db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()for rows.Next() {
	// 	...
	// }
	db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
	// 多个连接与参数
	db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)

	// Pluck
	// 将模型中的单个列作为地图查询，如果要查询多个列，可以使用Scan
	var ages []int64
	db.Find(&users).Pluck("age", &ages)
	var names []string
	db.Model(&User{}).Pluck("name", &names)
	db.Table("deleted_users").Pluck("name", &names)
	// 要返回多个列，做这样：
	db.Select("name, age").Find(&users)

	// Scan
	// 将结果扫描到另一个结构中。
	type ScanResult struct {
		Name string
		Age  int
	}
	var scanResult ScanResult
	db.Table("users").Select("name, age").Where("name = ?", 3).Scan(&scanResult)
	// Raw SQL
	db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&scanResult)

	// Scopes
	// 将当前数据库连接传递到func(*DB) *DB，可以用于动态添加条件
	func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
		return db.Where("amount > ?", 1000)
	}
	func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
		return db.Where("pay_mode_sign = ?", "C")
	}
	func PaidWithCod(db *gorm.DB) *gorm.DB {
		return db.Where("pay_mode_sign = ?", "C")
	}
	func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
		return func (db *gorm.DB) *gorm.DB {
			return db.Scopes(AmountGreaterThan1000).Where("status in (?)", status)
		}
	}
	db.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find(&orders)	// 查找所有信用卡订单和金额大于1000
	db.Scopes(AmountGreaterThan1000, PaidWithCod).Find(&orders)	// 查找所有COD订单和金额大于1000
	db.Scopes(OrderStatus([]string{"paid", "shipped"})).Find(&orders)	// 查找所有付费，发货订单

	// 指定表名
	// 使用User结构定义创建`deleted_users`表
	db.Table("deleted_users").CreateTable(&User{})
	var deleted_users []User
	db.Table("deleted_users").Find(&deleted_users)	// SELECT * FROM deleted_users;
	db.Table("deleted_users").Where("name = ?", "jinzhu").Delete()	// DELETE FROM deleted_users WHERE name = 'jinzhu';
	
	// 预加载
	db.Preload("Orders").Find(&users)	// SELECT * FROM users;
	// SELECT * FROM orders WHERE user_id IN (1,2,3,4);

	db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)	// SELECT * FROM users;
	// SELECT * FROM orders WHERE user_id IN (1,2,3,4) AND state NOT IN ('cancelled');

	db.Where("state = ?", "active").Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users) // SELECT * FROM users WHERE state = 'active';
	// SELECT * FROM orders WHERE user_id IN (1,2) AND state NOT IN ('cancelled');

	db.Preload("Orders").Preload("Profile").Preload("Role").Find(&users)	// SELECT * FROM users;
	// SELECT * FROM orders WHERE user_id IN (1,2,3,4); 
	// has many
	// SELECT * FROM profiles WHERE user_id IN (1,2,3,4); 
	// has one
	// SELECT * FROM roles WHERE id IN (4,5,6); 
	// belongs to

	// 自定义预加载SQL
	// 您可以通过传递func(db *gorm.DB) *gorm.DB（与Scopes的使用方法相同）来自定义预加载SQL，例如：
	db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
		return db.Order("orders.amount DESC")
	}).Find(&users)
	// SELECT * FROM users;
	// SELECT * FROM orders WHERE user_id IN (1,2,3,4) order by orders.amount DESC;

	// 嵌套预加载
	db.Preload("Orders.OrderItems").Find(&users)
	db.Preload("Orders", "state = ?", "paid").Preload("Orders.OrderItems").Find(&users)

	// 更新全部字段
	// Save将包括执行更新SQL时的所有字段，即使它没有更改
	db.First(&user)
	user.Name = "jinzhu 2"
	user.Age = 100
	db.Save(&user)
	// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=1;

	// 更新更改字段
	// 如果只想更新更改的字段，可以使用Update, Updates
	// 更新单个属性（如果更改）
	db.Model(&user).Update("name", "hello")	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=1;
	// 使用组合条件更新单个属性
	db.Model(&user).Where("active = ?", true).Update("name", "hello")	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=1 AND active=true;
	// 使用`map`更新多个属性，只会更新这些更改的字段
	db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})	// UPDATE users SET name='hello', age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=1;
	// 使用`struct`更新多个属性，只会更新这些更改的和非空白字段
	db.Model(&user).Updates(User{Name: "hello", Age: 18})	// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 1;
	// 警告:当使用struct更新时，FORM将仅更新具有非空值的字段
	// 对于下面的更新，什么都不会更新为""，0，false是其类型的空白值
	db.Model(&user).Updates(User{Name: "", Age: 0, Actived: false})

	// 更新选择的字段
	// 如果您只想在更新时更新或忽略某些字段，可以使用Select, Omit
	db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=1;
	db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})	// UPDATE users SET age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=1;
	
	// 更新更改字段但不进行Callbacks
	// 以上更新操作将执行模型的BeforeUpdate, AfterUpdate方法，更新其UpdatedAt时间戳，在更新时保存它的Associations，如果不想调用它们，可以使用UpdateColumn, UpdateColumns
	// 更新单个属性，类似于`Update`
	db.Model(&user).UpdateColumn("name", "hello")	// UPDATE users SET name='hello' WHERE id = 111;
	// 更新多个属性，与“更新”类似
	db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})	// UPDATE users SET name='hello', age=18 WHERE id = 111;

	// Batch Updates 批量更新
	// Callbacks在批量更新时不会运行
	db.Table("users").Where("id IN (?)", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})	// UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);
	// 使用struct更新仅适用于非零值，或使用map[string]interface{}
	db.Model(User{}).Updates(User{Name: "hello", Age: 18})	// UPDATE users SET name='hello', age=18;
	// 使用`RowsAffected`获取更新记录计数
	db.Model(User{}).Updates(User{Name: "hello", Age: 18}).RowsAffected

	// 使用SQL表达式更新
	DB.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))	// UPDATE "products" SET "price" = price * '2' + '100', "updated_at" = '2013-11-17 21:34:10' WHERE "id" = '2';
	DB.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})	// UPDATE "products" SET "price" = price * '2' + '100', "updated_at" = '2013-11-17 21:34:10' WHERE "id" = '2';
	DB.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))	// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = '2';
	DB.Model(&product).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))	// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = '2' AND quantity > 1;

	// 在Callbacks中更改更新值
	// 如果要使用BeforeUpdate, BeforeSave更改回调中的更新值，可以使用scope.SetColumn，例如
	func (user *User) BeforeSave(scope *gorm.Scope) (err error) {
		if pw, err := bcrypt.GenerateFromPassword(user.Password, 0); err == nil {
			scope.SetColumn("EncryptedPassword", pw)
		}
	}

	// 额外更新选项
	// 为Update语句添加额外的SQL选项
	db.Model(&user).Set("gorm:update_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Update("name", "hello")
	// UPDATE users SET name='hello', updated_at = '2013-11-17 21:34:10' WHERE id=111 OPTION (OPTIMIZE FOR UNKNOWN);

	// 删除/软删除
	// 警告 删除记录时，需要确保其主要字段具有值，GORM将使用主键删除记录，如果主要字段为空，GORM将删除模型的所有记录
	// 删除存在的记录
	db.Delete(&email)	// DELETE from emails where id=10;
	// 为Delete语句添加额外的SQL选项
	db.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&email)	// DELETE from emails where id=10 OPTION (OPTIMIZE FOR UNKNOWN);

	// 批量删除
	// 删除所有匹配记录
	db.Where("email LIKE ?", "%jinzhu%").Delete(Email{})	// DELETE from emails where email LIKE "%jinhu%";
	db.Delete(Email{}, "email LIKE ?", "%jinzhu%")	// DELETE from emails where email LIKE "%jinhu%";

	// 软删除
	// 如果模型有DeletedAt字段，它将自动获得软删除功能！ 那么在调用Delete时不会从数据库中永久删除，而是只将字段DeletedAt的值设置为当前时间。
	db.Delete(&user)	// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;
	// 批量删除
	db.Where("age = ?", 20).Delete(&User{})	// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;
	// 软删除的记录将在查询时被忽略
	db.Where("age = 20").Find(&user)	// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
	// 使用Unscoped查找软删除的记录
	db.Unscoped().Where("age = 20").Find(&users)	// SELECT * FROM users WHERE age = 20;
	// 使用Unscoped永久删除记录
	db.Unscoped().Delete(&order)	// DELETE FROM orders WHERE id=10;

	// 关联
	// 默认情况下，当创建/更新记录时，GORM将保存其关联，如果关联具有主键，GORM将调用Update来保存它，否则将被创建。
	user := User{
		Name:            "jinzhu",
		BillingAddress:  Address{Address1: "Billing Address - Address 1"},
		ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
		Emails:          []Email{
							{Email: "jinzhu@example.com"},
							{Email: "jinzhu-2@example@example.com"},
						 },
		Languages:       []Language{
							{Name: "ZH"},
							{Name: "EN"},
						 },
	}
	db.Create(&user)	// BEGIN TRANSACTION;
	// INSERT INTO "addresses" (address1) VALUES ("Billing Address - Address 1");
	// INSERT INTO "addresses" (address1) VALUES ("Shipping Address - Address 1");
	// INSERT INTO "users" (name,billing_address_id,shipping_address_id) VALUES ("jinzhu", 1, 2);
	// INSERT INTO "emails" (user_id,email) VALUES (111, "jinzhu@example.com");
	// INSERT INTO "emails" (user_id,email) VALUES (111, "jinzhu-2@example.com");
	// INSERT INTO "languages" ("name") VALUES ('ZH');
	// INSERT INTO user_languages ("user_id","language_id") VALUES (111, 1);
	// INSERT INTO "languages" ("name") VALUES ('EN');
	// INSERT INTO user_languages ("user_id","language_id") VALUES (111, 2);
	// COMMIT;
	db.Save(&user)

	// 创建/更新时跳过保存关联
	// 默认情况下保存记录时，GORM也会保存它的关联，你可以通过设置gorm:save_associations为false跳过它。
	db.Set("gorm:save_associations", false).Create(&user)
	db.Set("gorm:save_associations", false).Save(&user)

	// tag设置跳过保存关联
	// 您可以使用tag来配置您的struct，以便在创建/更新时不会保存关联
	type AssociationsUser struct {
		gorm.Model
		Name      string
		CompanyID uint
		Company   Company `gorm:"save_associations:false"`
	}
	type AssociationsCompany struct {
		gorm.Model
		Name string
	}

	// Callbacks
	// 您可以将回调方法定义为模型结构的指针，在创建，更新，查询，删除时将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改。

	// // 创建对象
	// // 创建过程中可用的回调
	// // begin transaction 开始事物
	// BeforeSave
	// BeforeCreate// save before associations 保存前关联
	// // update timestamp `CreatedAt`, `UpdatedAt` 更新`CreatedAt`, `UpdatedAt`时间戳
	// // save self 保存自己
	// // reload fields that have default value and its value is blank 重新加载具有默认值且其值为空的字段
	// // save after associations 保存后关联
	// AfterCreate
	// AfterSave// commit or rollback transaction 提交或回滚事务

	// // 更新对象
	// // 更新过程中可用的回调
	// // begin transaction 开始事物
	// BeforeSave
	// BeforeUpdate// save before associations 保存前关联
	// // update timestamp `UpdatedAt` 更新`UpdatedAt`时间戳
	// // save self 保存自己
	// // save after associations 保存后关联
	// AfterUpdate
	// AfterSave// commit or rollback transaction 提交或回滚事务
	
	// // 删除对象
	// // 删除过程中可用的回调
	// // begin transaction 开始事物
	// BeforeDelete// delete self 删除自己
	// AfterDelete// commit or rollback transaction 提交或回滚事务

	// // 查询对象
	// // 查询过程中可用的回调
	// // load data from database 从数据库加载数据
	// // Preloading (edger loading) 预加载（加载）
	// AfterFind

	// 错误处理
	// 执行任何操作后，如果发生任何错误，GORM将其设置为*DB的Error字段
	// if err := db.Where("name = ?", "jinzhu").First(&user).Error; err != nil {
	// 	// 错误处理...
	// }
	// 如果有多个错误发生，用`GetErrors`获取所有的错误，它返回`[]error`
	db.First(&user).Limit(10).Find(&users).GetErrors()
	// // 检查是否返回RecordNotFound错误
	// db.Where("name = ?", "hello world").First(&user).RecordNotFound()
	// if db.Model(&user).Related(&credit_card).RecordNotFound() {
	// 	// 没有信用卡被发现处理...
	// }

	// // 事务
	// // 要在事务中执行一组操作，一般流程如下。
	// // 开始事务
	// tx := db.Begin()
	// // 在事务中做一些数据库操作（从这一点使用'tx'，而不是'db'）
	// tx.Create(...)
	// // ...
	// // 发生错误时回滚事务
	// tx.Rollback()
	// // 或提交事务
	// tx.Commit()

	// 执行原生SQL
	db.Exec("DROP TABLE users;")
	db.Exec("UPDATE orders SET shipped_at=? WHERE id IN (?)", time.Now, []int64{11,22,33})
	// Scan
	type ScanResult struct {
		Name string
		Age  int
	}
	var scanResult ScanResult
	db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&scanResult)

	// sql.Row & sql.Rows
	// 获取查询结果为*sql.Row或*sql.Rows
	row := db.Table("users").Where("name = ?", "jinzhu").Select("name, age").Row() // (*sql.Row)
	row.Scan(&name, &age)

	rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows() // (*sql.Rows, error)
	defer rows.Close()
	for rows.Next() {
		// ...
		rows.Scan(&name, &age, &email)
		// ...
	}
	// Raw SQL
	rows, err := db.Raw("select name, age, email from users where name = ?", "jinzhu").Rows() // (*sql.Rows, error)
	defer rows.Close()
	for rows.Next() {
		// ...
		rows.Scan(&name, &age, &email)
		// ...
	}

	// 迭代中使用sql.Rows的Scan
	rows, err := db.Model(&User{}).Where("name = ?", "jinzhu").Select("name, age, email").Rows() // (*sql.Rows, error)
	defer rows.Close()
	for rows.Next() {
	  var user User
	  db.ScanRows(rows, &user)
	  // do something
	}
	
	// 通用数据库接口sql.DB
	// 从*gorm.DB连接获取通用数据库接口*sql.DB
	// 获取通用数据库对象`*sql.DB`以使用其函数
	db.DB()
	// Ping
	db.DB().Ping()

	// 连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// 复合主键
	// 将多个字段设置为主键以启用复合主键
	type primaryProduct struct {
		ID           string `gorm:"primary_key"`
		LanguageCode string `gorm:"primary_key"`
	}

	// 日志
	// Gorm有内置的日志记录器支持，默认情况下，它会打印发生的错误
	// 启用Logger，显示详细日志
	db.LogMode(true)
	// 禁用日志记录器，不显示任何日志
	db.LogMode(false)
	// 调试单个操作，显示此操作的详细日志
	db.Debug().Where("name = ?", "jinzhu").First(&User{})
	
	// 自定义日志
	// 参考GORM的默认记录器如何自定义它https://github.com/jinzhu/gorm/blob/master/logger.go
	db.SetLogger(gorm.Logger{revel.TRACE})
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	

	// // 架构
	// // Gorm使用可链接的API，*gorm.DB是链的桥梁，对于每个链API，它将创建一个新的关系。
	// db, err := gorm.Open("postgres", "user=gorm dbname=gorm sslmode=disable")
	// // 创建新关系
	// db = db.Where("name = ?", "jinzhu")
	// // 过滤更多
	// if SomeCondition {
	// 	db = db.Where("age = ?", 20)
	// } else {
	// 	db = db.Where("age = ?", 30)
	// }
	// if YetAnotherCondition {
	// 	db = db.Where("active = ?", 1)
	// }
	// // 当我们开始执行任何操作时，GORM将基于当前的*gorm.DB创建一个新的*gorm.Scope实例
	// // 执行查询操作
	// db.First(&user)
	// // 并且基于当前操作的类型，它将调用注册的creating, updating, querying, deleting或row_querying回调来运行操作。
	// // 对于上面的例子，将调用querying，参考查询回调

	// 写插件
	// GORM本身由Callbacks提供支持，因此您可以根据需要完全自定义GORM
	db.Callback().Create().Register("update_created_at", updateCreated)// 注册Create进程的回调

	// 删除现有的callback
	db.Callback().Create().Remove("gorm:create")// 从Create回调中删除`gorm:create`回调

	// 替换现有的callback
	db.Callback().Create().Replace("gorm:create", newCreateFunction)// 使用新函数`newCreateFunction`替换回调`gorm:create`用于创建过程

	// 注册callback顺序
	db.Callback().Create().Before("gorm:create").Register("update_created_at", updateCreated)
	db.Callback().Create().After("gorm:create").Register("update_created_at", updateCreated)
	db.Callback().Query().After("gorm:query").Register("my_plugin:after_query", afterQuery)
	db.Callback().Delete().After("gorm:delete").Register("my_plugin:after_delete", afterDelete)
	db.Callback().Update().Before("gorm:update").Register("my_plugin:before_update", beforeUpdate)
	db.Callback().Create().Before("gorm:create").After("gorm:before_create").Register("my_plugin:before_create", beforeCreate)

	// 预定义回调
	// GORM定义了回调以执行其CRUD操作，在开始编写插件之前检查它们。
	// •	Create callbacks
	// •	Update callbacks
	// •	Query callbacks
	// •	Delete callbacks
	// •	Row Query callbacks Row Query callbacks将在运行Row或Rows时被调用，默认情况下没有注册的回调，你可以注册一个新的回调：
	db.Callback().RowQuery().Register("publish:update_table_name", updateTableName)

}

// 回调示例
func (u *User) BeforeUpdate() (err error) {
	if u.readonly() {
		err = errors.New("read only user")
	}
	return
}
// 如果用户ID大于1000，则回滚插入func (u *User) AfterCreate() (err error) {
	if (u.Id > 1000) {
		err = errors.New("user id is already greater than 1000")
	}
	return
}
// gorm中的保存/删除操作正在事务中运行，因此在该事务中所做的更改不可见，除非提交。 如果要在回调中使用这些更改，则需要在同一事务中运行SQL。 所以你需要传递当前事务到回调，像这样：
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(u).Update("role", "admin")
	return
}
func (u *User) AfterCreate(scope *gorm.Scope) (err error) {
	scope.DB().Model(u).Update("role", "admin")
	return
}

// // CreateAnimals 事务
// func CreateAnimals(db *gorm.DB) err {
//   tx := db.Begin()
//   // 注意，一旦你在一个事务中，使用tx作为数据库句柄

//   if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
//      tx.Rollback()
//      return err
//   }

//   if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
//      tx.Rollback()
//      return err
//   }

//   tx.Commit()
//   return nil
// }

// 注册新的callback
func updateCreated(scope *Scope) {
	if scope.HasColumn("Created") {
		scope.SetColumn("Created", NowFunc())
	}
}

func updateTableName(scope *gorm.Scope) {
	scope.Search.Table(scope.TableName() + "_draft") // append `_draft` to table name
}