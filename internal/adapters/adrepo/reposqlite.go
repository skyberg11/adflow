package adrepo

import (
	"adflow/internal/ads"
	"adflow/internal/app"
	"sync"
	"time"

	"gorm.io/gorm"
)

type sqliteRepository struct {
	db  *gorm.DB
	cnt int64
	m   sync.Mutex
}

func (r *sqliteRepository) Create(ad *ads.Ad) error {
	r.m.Lock()
	defer r.m.Unlock()

	ad.ID = r.cnt + 1
	ad.CreationTime = time.Now().UTC()
	ad.UpdateTime = time.Now().UTC()

	r.db.Create(ad)
	r.cnt += 1

	return nil
}

func (r *sqliteRepository) Get(id int64) (*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var ad ads.Ad

	r.db.Where("ID = ?", id).Find(&ad)

	if ad == (ads.Ad{}) {
		return nil, ads.ErrBadRequest
	}

	return &ad, nil
}

func (r *sqliteRepository) Update(id int64, title, text string) (*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var ad ads.Ad

	r.db.Where("ID = ?", id).Find(&ad)

	if ad == (ads.Ad{}) {
		return nil, ads.ErrBadRequest
	}

	ad.Title = title
	ad.Text = text
	ad.UpdateTime = time.Now().UTC()

	r.db.Save(&ad)
	return &ad, nil
}

func (r *sqliteRepository) UpdateStatus(id int64, published bool) (*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var ad ads.Ad

	r.db.Where("ID = ?", id).Find(&ad)

	if ad == (ads.Ad{}) {
		return nil, ads.ErrBadRequest
	}

	ad.Published = published
	ad.UpdateTime = time.Now().UTC()

	r.db.Save(&ad)
	return &ad, nil
}

func (r *sqliteRepository) GetAllAds() ([]*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var current_ads []ads.Ad
	var ads []*ads.Ad

	result := r.db.Find(&current_ads)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, v := range current_ads {
		if v.Published {
			ads = append(ads, &v)
		}
	}

	return ads, nil
}

func (r *sqliteRepository) GetAds(filter ads.Filter) ([]*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var current_ads []ads.Ad
	var ads_inds []int64

	result := r.db.Find(&current_ads)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, v := range current_ads {
		if CheckFilter(&v, filter) {
			ads_inds = append(ads_inds, v.ID)
		}
	}

	var ads = make([]*ads.Ad, len(ads_inds))
	for i := range ads_inds {
		r.db.Where("ID = ?", ads_inds[i]).Find(&ads[i])
	}

	return ads, nil
}

func (r *sqliteRepository) DeleteAd(id int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	var ad ads.Ad

	result := r.db.Where("ID = ?", id).First(&ad)

	if result.Error != nil {
		return ads.ErrBadRequest
	}

	r.db.Delete(&ad)

	return nil
}

func NewSQLiteAds(db *gorm.DB) app.Repository {
	err := db.AutoMigrate(&ads.Ad{})
	if err != nil {
		panic(err)
	}

	return &sqliteRepository{
		db:  db,
		cnt: 0,
	}
}
