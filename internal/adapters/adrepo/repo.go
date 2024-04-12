package adrepo

import (
	"adflow/internal/ads"
	"adflow/internal/app"
	"strings"
	"sync"
	"time"
)

type localRepository struct {
	ads map[int64]*ads.Ad
	cnt int64
	m   sync.Mutex
}

func (r *localRepository) Create(ad *ads.Ad) error {
	r.m.Lock()
	defer r.m.Unlock()

	ad.ID = r.cnt
	ad.CreationTime = time.Now().UTC()
	ad.UpdateTime = time.Now().UTC()

	r.ads[r.cnt] = ad
	r.cnt += 1

	return nil
}

func (r *localRepository) Get(id int64) (*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	_, ok := r.ads[id]

	if !ok {
		return nil, ads.ErrBadRequest
	}
	return r.ads[id], nil
}

func (r *localRepository) Update(id int64, title, text string) (*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.ads[id]; !ok {
		return nil, ads.ErrBadRequest
	}

	r.ads[id].Title = title
	r.ads[id].Text = text
	r.ads[id].UpdateTime = time.Now().UTC()

	return r.ads[id], nil
}

func (r *localRepository) UpdateStatus(id int64, published bool) (*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.ads[id]; !ok {
		return nil, ads.ErrBadRequest
	}

	r.ads[id].Published = published
	r.ads[id].UpdateTime = time.Now().UTC()

	return r.ads[id], nil
}

func (r *localRepository) GetAllAds() ([]*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var ads []*ads.Ad

	for _, v := range r.ads {
		if v.Published {
			ads = append(ads, v)
		}
	}

	return ads, nil
}

func CheckFilter(ad *ads.Ad, filter ads.Filter) bool {
	if filter.Published != nil {
		if filter.Published != ad.Published {
			return false
		}
	}
	if filter.AuthorID != nil {
		if filter.AuthorID != ad.AuthorID {
			return false
		}
	}
	if filter.TitlePrefix != nil {
		if !strings.HasPrefix(ad.Title, filter.TitlePrefix.(string)) {
			return false
		}
	}
	if filter.CreationTime != nil {
		if (filter.CreationTime.(time.Time)).After(ad.CreationTime) {
			return false
		}
	}
	return true
}

func (r *localRepository) GetAds(filter ads.Filter) ([]*ads.Ad, error) {
	r.m.Lock()
	defer r.m.Unlock()

	var ads []*ads.Ad

	for _, v := range r.ads {
		if CheckFilter(v, filter) {
			ads = append(ads, v)
		}
	}

	return ads, nil
}

func (r *localRepository) DeleteAd(id int64) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.ads[id]; !ok {
		return ads.ErrBadRequest
	}

	delete(r.ads, id)

	return nil
}

func New() app.Repository {
	return &localRepository{
		ads: make(map[int64]*ads.Ad),
		cnt: 0,
	}
}
