package models

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Urls struct {
	ID        uint32   `gorm:"primary_key;auto_increment" json:"id"`
	Name, URL string
	Type      string          `gorm:"DEFAULT:'GET'"`
	Calls     []EndPointCalls `gorm:"ForeignKey:EndPointID"`
	Owner    User      `json:"owner"`
	OwnerID  uint32   `sql:"type:int REFERENCES users(id)" json:"owner_id"`
	Threshold    int
	FailedTimes  int
	SuccessTimes int
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Urls) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.URL = html.EscapeString(strings.TrimSpace(p.URL))
	p.Owner = User{}
	p.FailedTimes  =0
	p.SuccessTimes =0
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Urls) Validate() error {

	if p.Name == "" {
		return errors.New("Required Name")
	}

	if p.URL == "" {
		return errors.New("Required URL")
	}

	if p.Type == "" {
		return errors.New("Required Type")
	}

	if p.OwnerID < 1 {
		return errors.New("Required owner")
	}
	if p.Threshold < 1 {
		return errors.New("Required Threshold")
	}
	return nil
}

func (url *Urls) SaveUrl(db *gorm.DB) (*Urls, error) {
	var err error
	err = db.Debug().Model(&Urls{}).Create(&url).Error
	if err != nil {
		fmt.Print("errrorrr")
		return &Urls{}, err
	}
	if url.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", url.OwnerID).Take(&url.Owner).Error
		if err != nil {
			fmt.Print("errrorrr")
			return &Urls{}, err
		}
	}
	return url, nil
}

func (p *Urls) FindAllUrlses(db *gorm.DB) (*[]Urls, error) {
	var err error
	Urlses := []Urls{}
	err = db.Debug().Model(&Urls{}).Limit(100).Find(&Urlses).Error
	if err != nil {
		return &[]Urls{}, err
	}
	if len(Urlses) > 0 {
		for i := range Urlses {
			err := db.Debug().Model(&User{}).Where("id = ?", Urlses[i].OwnerID).Take(&Urlses[i].Owner).Error
			if err != nil {
				return &[]Urls{}, err
			}
		}
	}
	return &Urlses, nil
}

func (url *Urls) FindUrlByID(db *gorm.DB, pid uint64) (*Urls, error) {
	var err error
	err = db.Debug().Model(&Urls{}).Where("id = ?", pid).Take(&url).Error
	if err != nil {
		return &Urls{}, err
	}
	if url.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", url.OwnerID).Take(&url.Owner).Error
		if err != nil {
			return &Urls{}, err
		}
	}
	return url, nil
}

func (p *Urls) UpdateAUrl(db *gorm.DB) (*Urls, error) {

	var err error

	err = db.Debug().Model(&Urls{}).Where("id = ?", p.ID).Updates(Urls{Threshold: p.Threshold,Name: p.Name, URL: p.URL, 
		UpdatedAt: time.Now(),FailedTimes: p.FailedTimes,SuccessTimes: p.SuccessTimes}).Error
	if err != nil {
		return &Urls{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.OwnerID).Take(&p.Owner).Error
		if err != nil {
			return &Urls{}, err
		}
	}
	return p, nil
}

func (p *Urls) DeleteAUrl(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Urls{}).Where("id = ? and owner_id = ?", pid, uid).Take(&Urls{}).Delete(&Urls{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("URL not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (url *Urls) ShouldTriggerAlarm() bool {
	return url.FailedTimes >= url.Threshold
}

func (url *Urls) SendRequest() (*EndPointCalls, error) {
	resp, err := http.Get(url.URL)
	req := new(EndPointCalls)
	req.urlID = url.ID
	if err != nil {
		return req, err
	}
	req.ResponseCode = resp.StatusCode
	return req, nil
}