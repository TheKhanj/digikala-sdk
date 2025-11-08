package config

import (
	"encoding/json"
	"net/http"
)

func (this *CustomerCache) GetCookies() ([]*http.Cookie, error) {
	if this.Cookies == nil {
		return nil, nil
	}

	cookies := make([]*http.Cookie, 0)

	b, err := json.Marshal(this.Cookies)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &cookies)
	if err != nil {
		return nil, err
	}

	return cookies, nil
}

func (this *CustomerCache) SetCookies(cookies []*http.Cookie) error {
	b, err := json.Marshal(cookies)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &this.Cookies)
}
