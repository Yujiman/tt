package model

type Order struct {
	TimestampMixin
	DeleteMixin
	OrderID      uint   `gorm:"primaryKey;column:order_id"`
	OrderSession string `gorm:"column:order_session"`
	ChannelID    uint   `gorm:"column:channel_id"`
	BuyerID      uint   `gorm:"column:buyer_id"`
	OrderComment string `gorm:"column:order_comment"`
}

func (Order) TableName() string {
	return "update_templat.orders" // database_schema.table_name
}
