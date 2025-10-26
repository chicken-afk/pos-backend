package entities

type Outlet struct {
	ID         uint   `gorm:"primaryKey"`
	OutletName string `gorm:"not null"`
}
