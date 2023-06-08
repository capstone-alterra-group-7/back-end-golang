package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketTravelerDetail struct {
	gorm.Model
	TicketOrderID        uint
	TicketOrder          TicketOrder `gorm:"foreignKey:TicketOrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TravelerDetailID     uint
	TravelerDetail       TravelerDetail `gorm:"foreignKey:TravelerDetailID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TrainID              uint
	Train                Train `gorm:"foreignKey:TrainID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TrainPrice           int
	TrainCarriageID      uint
	TrainCarriage        TrainCarriage `gorm:"foreignKey:TrainCarriageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TrainSeatID          uint
	TrainSeat            TrainSeat `gorm:"foreignKey:TrainSeatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StationOriginID      uint
	StationOrigin        Station `gorm:"foreignKey:StationOriginID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DepartureTime        string
	StationDestinationID uint
	StationDestination   Station `gorm:"foreignKey:StationDestinationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ArrivalTime          string
	DateOfDeparture      time.Time `gorm:"type:DATE"`
	BoardingTicketCode   string
}
