package model

type Customer struct {
	CustomerID uint   `gorm:"primaryKey;column:customer_id"`
	FullName   string `gorm:"column:full_name"`
	Phone      string `gorm:"column:phone"`
	Email      string `gorm:"column:email"`
}

func (Customer) TableName() string {
	return "update_templat.customer" // database_schema.table_name
}

