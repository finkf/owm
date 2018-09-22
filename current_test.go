package owm

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/davecheney/errors"
)

func TestKelvin(t *testing.T) {
	tests := []struct {
		k            Kelvin
		wantC, wantF float64
	}{
		{0, -273, -460},
		{50, -223, -370},
		{300, 27, 80},
		{310, 37, 98},
		{315, 42, 107},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%f", tc.k), func(t *testing.T) {
			if got := math.Round(tc.k.Celcius()); got != tc.wantC {
				t.Fatalf("expected %f; got %f", tc.wantC, got)
			}
			if got := math.Round(tc.k.Fahrenheit()); got != tc.wantF {
				t.Fatalf("expected %f; got %f", tc.wantF, got)
			}
		})
	}

}

func TestQuery(t *testing.T) {
	tests := []struct {
		query Query
		want  string
	}{
		{Query{}, "lat=0&lon=0"}, /* default query */
		{Query{City: "Berlin"}, "q=Berlin"},
		{Query{City: "Berlin", Country: "de"}, "q=Berlin,de"},
		{Query{City: "Berlin", Lang: "de"}, "q=Berlin&lang=de"},
		{Query{ZIP: 12345}, "zip=12345"},
		{Query{ZIP: 12345, Country: "de"}, "zip=12345,de"},
		{Query{ZIP: 12345, Lang: "de"}, "zip=12345&lang=de"},
		{Query{ID: 12345}, "id=12345"},
		{Query{ID: 12345, Lang: "de"}, "id=12345&lang=de"},
		{Query{Lat: 1, Lon: 1}, "lat=1&lon=1"},
		{Query{Lat: 1, Lon: 1, Lang: "es"}, "lat=1&lon=1&lang=es"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.query.params(); got != tc.want {
				t.Fatalf("expected %q; got %q", tc.want, got)
			}
		})
	}
}

func TestCurrent(t *testing.T) {
	tests := []struct {
		client Getter
		isErr  bool
	}{
		{okGetter{}, false},
		{errorGetter{}, true},
		{badGetter{}, true},
		{badJSONGetter{}, true},
		{notFoundGetter{}, true},
	}
	for _, tc := range tests {
		t.Run(reflect.TypeOf(tc.client).Name(), func(t *testing.T) {
			api := API{Client: tc.client}
			_, err := api.Current(Query{})
			if tc.isErr && err == nil {
				t.Fatalf("expected an error [%v]", err)
			}
			if !tc.isErr && err != nil {
				t.Fatalf("got error: %v", err)
			}
		})
	}
}

type okGetter struct{}

func (okGetter) Get(string) (*http.Response, error) {
	w := httptest.NewRecorder()
	io.WriteString(w, jsonStr)
	return w.Result(), nil
}

type badGetter struct{}

func (badGetter) Get(string) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, jsonStr)
	return w.Result(), nil
}

type errorGetter struct{}

func (errorGetter) Get(string) (*http.Response, error) {
	return nil, errors.New("error")
}

type badJSONGetter struct{}

func (badJSONGetter) Get(string) (*http.Response, error) {
	w := httptest.NewRecorder()
	io.WriteString(w, jsonStr[0:17])
	return w.Result(), nil
}

type notFoundGetter struct{}

func (notFoundGetter) Get(string) (*http.Response, error) {
	w := httptest.NewRecorder()
	io.WriteString(w, jsonStr[0:len(jsonStr)-4]+"401}")
	return w.Result(), nil
}

var jsonStr = `
{"coord":{"lon":139,"lat":35},
"sys":{"country":"JP","sunrise":1369769524,"sunset":1369821049},
"weather":[{"id":804,"main":"clouds","description":"overcast clouds","icon":"04n"}],
"main":{"temp":289.5,"humidity":89,"pressure":1013,"temp_min":287.04,"temp_max":292.04},
"wind":{"speed":7.31,"deg":187.002},
"rain":{"3h":5},
"clouds":{"all":92},
"dt":1369824698,
"id":1851632,
"name":"Shuzenji",
"cod":200}`
