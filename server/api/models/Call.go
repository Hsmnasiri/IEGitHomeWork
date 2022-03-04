package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// EndPointCalls - Object for storing endpoints call details
type EndPointCalls struct {
	ID        uint32   `gorm:"primary_key;auto_increment" json:"id"`
	urlID   uint32 `gorm:"index;not null"`
	ResponseCode int
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}



func (epc *EndPointCalls) Prepare() {
	epc.ID = 0
	epc.CreatedAt = time.Now()
	epc.UpdatedAt = time.Now()
}


func (epc *EndPointCalls) Validate()  error{
	 
	 return nil
	
}

func (epc *EndPointCalls) SaveCall(db *gorm.DB) (*EndPointCalls, error) {

	var err error
	err = db.Debug().Create(&epc).Error
	if err != nil {
		return &EndPointCalls{}, err
	}
	return epc, nil
}



func (epc *EndPointCalls) FindCallsByTime(db *gorm.DB,StartTime time.Time,EndTime time.Time) (*[]EndPointCalls, error) {
	var err error
	Calls := []EndPointCalls{}
	err=db.Where("created_at BETWEEN ? AND ?", StartTime, EndTime).Find(&Calls).Error
	if err != nil {
		return &[]EndPointCalls{}, err
	}
	return &Calls, err
}

func (epc *EndPointCalls) FindCallByID(db *gorm.DB, uid uint32) (*EndPointCalls, error) {
	var err error
	err = db.Debug().Model(EndPointCalls{}).Where("id = ?", uid).Take(&epc).Error
	if err != nil {
		return &EndPointCalls{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &EndPointCalls{}, errors.New("Call Not Found")
	}
	return epc, err
}

func (epc *EndPointCalls) FindAllCalls(db *gorm.DB) (*[]EndPointCalls, error) {
	var err error
	Calls := []EndPointCalls{}
	err = db.Debug().Model(&EndPointCalls{}).Limit(100).Find(&Calls).Error
	if err != nil {
		return &[]EndPointCalls{}, err
	}
	return &Calls, err
}

func (epc *EndPointCalls) GetEndPointCallesByUrl(db *gorm.DB,uid uint32) (*[]EndPointCalls, error) {
	Calls := []EndPointCalls{}
	if err := db.Debug().Model(&EndPointCalls{urlID: uid}).Where("urlID = ?", uid).Find(&Calls).Error; err != nil {
		return nil, err
	}
	return &Calls, nil
}
