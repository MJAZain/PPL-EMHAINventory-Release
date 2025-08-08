package model

type SystemConfig struct {
	ID              uint `gorm:"primaryKey"`
	MaxFailedLogin  int  `gorm:"default:5"`  // Batas maksimal login gagal
	LockoutDuration int  `gorm:"default:30"` // Durasi lockout dalam menit
}
