package models

import "gorm.io/gorm"

type TrainStation struct {
	gorm.Model
	TrainID    uint
	StationID  uint
	Train      Train   `gorm:"foreignKey:TrainID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Station    Station `gorm:"foreignKey:StationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ArriveTime string
}
