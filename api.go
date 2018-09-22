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
func (api API) Current(q Queryer) (*Current, error) {
	url := api.url("weather", q)
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

func (api API) url(what string, q Queryer) string {
	return fmt.Sprintf("%s/%s%s&appid=%s", URL, what, q.Query(), api.Key)
}

// Queryer defines the interface for location searches.
type Queryer interface {
	Query() string
}

// ByCity searches for location by city and country.
type ByCity struct {
	City, Country, Lang string
}

// Query implements the Queryer interface.
func (q ByCity) Query() string {
	query := "?q=" + q.City
	if q.Country != "" {
		query += "," + q.Country
	}
	return appendLang(query, q.Lang)
}

// ByID searches for cities by ID.
type ByID struct {
	ID   int
	Lang string
}

// Query implements the Queryer interface.
func (q ByID) Query() string {
	return appendLang(fmt.Sprintf("?id=%d", q.ID), q.Lang)
}

// ByZIP searches location by ZIP code.
type ByZIP struct {
	ZIP           int
	Country, Lang string
}

// Query implements the Queryer interface.
func (q ByZIP) Query() string {
	query := fmt.Sprintf("?zip=%d", q.ZIP)
	if q.Country != "" {
		query += "," + q.Country
	}
	return appendLang(query, q.Lang)
}

// ByCoords searches for locations by longitude an latitute.
type ByCoords struct {
	Lon, Lat int
	Lang     string
}

// Query implements the Queryer interface.
func (q ByCoords) Query() string {
	return appendLang(
		fmt.Sprintf("?lat=%d&lon=%d", q.Lat, q.Lon), q.Lang)
}

func appendLang(params, lang string) string {
	if lang != "" {
		return params + fmt.Sprintf("&lang=%s", lang)
	}
	return params
}

var _ Queryer = ByCity{}
var _ Queryer = ByID{}
var _ Queryer = ByZIP{}
var _ Queryer = ByCoords{}
