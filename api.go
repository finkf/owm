package owm

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davecheney/errors"
)

const (
	// URL defines the base URL for the open weather map.
	URL = "http://api.openweathermap.org/data/2.5"
)

// Getter defines the default interface for clients
// that can connect to the OWM-API. Note that http.Client
// statifies this interface.
type Getter interface {
	Get(string) (*http.Response, error)
}

// API is the base for all OWM queries.
type API struct {
	// Client is used to send requests to the OWM-API.
	// It will be almost allways an instance of http.Client.
	Client Getter
	// Key defines the OWM appid key.
	Key string
}

// Current queries the API and returns
// all current weather data or an error.
func (api API) Current(q Query) (*Current, error) {
	url := api.url("weather", q.params())
	resp, err := api.Client.Get(url)
	if err != nil {
		return nil, errors.Annotatef(err, "cannot connect to: %s", url)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}
	dec := json.NewDecoder(resp.Body)
	var c Current
	if err := dec.Decode(&c); err != nil {
		return nil, errors.Annotatef(err, "invalid server response")
	}
	if c.Cod != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s [%d]",
			c.Message, c.Cod)
	}
	return &c, nil
}

func (api API) url(what, params string) string {
	return fmt.Sprintf("%s/%s?%s&appid=%s", URL, what, params, api.Key)
}

// Query defines the locations for weather queries.
// You do not need to define all parameters.
// Use combinations that make sense.
// The Lang parameter can be used to change the
// language of weather descriptions.
type Query struct {
	City, Country, Lang string
	Lat, Lon, ID, ZIP   int
}

func (q Query) params() string {
	if q.ID != 0 {
		return q.appendLang(fmt.Sprintf("id=%d", q.ID))
	}
	var params string
	if q.ZIP != 0 {
		params = fmt.Sprintf("zip=%d", q.ZIP)
	}
	if q.City != "" && params == "" {
		params = fmt.Sprintf("q=%s", q.City)
	}
	if q.Country != "" && params != "" {
		params += fmt.Sprintf(",%s", q.Country)
	}
	if params != "" {
		return q.appendLang(params)
	}
	return q.appendLang(fmt.Sprintf("lat=%d&lon=%d", q.Lat, q.Lon))
}

func (q Query) appendLang(params string) string {
	if q.Lang != "" {
		return params + fmt.Sprintf("&lang=%s", q.Lang)
	}
	return params
}
