package entities

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	OutletID uint   `gorm:"not null"`
	Outlet   Outlet `gorm:"foreignKey:outlet_id"`
	RoleID   uint   `gorm:"not null"`
	Role     Role   `gorm:"foreignKey:role_id"`
}
