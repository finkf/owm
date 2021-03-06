package owm

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
		q    Queryer
		want string
	}{
		{ByCity{City: "Berlin"}, "?q=Berlin"},
		{ByCity{City: "Berlin", Country: "de"}, "?q=Berlin,de"},
		{ByCity{City: "Berlin", Lang: "de"}, "?q=Berlin&lang=de"},
		{ByZIP{ZIP: 12345}, "?zip=12345"},
		{ByZIP{ZIP: 12345, Country: "de"}, "?zip=12345,de"},
		{ByZIP{ZIP: 12345, Lang: "de"}, "?zip=12345&lang=de"},
		{ByID{ID: 12345}, "?id=12345"},
		{ByID{ID: 12345, Lang: "de"}, "?id=12345&lang=de"},
		{ByCoords{Lat: 1, Lon: 1}, "?lat=1&lon=1"},
		{ByCoords{Lat: 1, Lon: 1, Lang: "es"}, "?lat=1&lon=1&lang=es"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.q.Query(); got != tc.want {
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
		{okGetter{jsonStrCurrent}, false},
		{errorGetter{}, true},
		{badGetter{jsonStrCurrent}, true},
		{badJSONGetter{jsonStrCurrent}, true},
		{notFoundGetter{jsonStrCurrent}, true},
	}
	for _, tc := range tests {
		t.Run(reflect.TypeOf(tc.client).Name(), func(t *testing.T) {
			api := API{Client: tc.client}
			_, err := api.Current(ByCity{})
			if tc.isErr && err == nil {
				t.Fatalf("expected an error [%v]", err)
			}
			if !tc.isErr && err != nil {
				t.Fatalf("got error: %v", err)
			}
		})
	}
}

func TestForecast(t *testing.T) {
	tests := []struct {
		client Getter
		isErr  bool
	}{
		{okGetter{jsonStrForecast}, false},
		{errorGetter{}, true},
		{badGetter{jsonStrForecast}, true},
		{badJSONGetter{jsonStrForecast}, true},
		{notFoundGetter{jsonStrForecast}, true},
	}
	for _, tc := range tests {
		t.Run(reflect.TypeOf(tc.client).Name(), func(t *testing.T) {
			api := API{Client: tc.client}
			_, err := api.Forecast(ByCity{})
			if tc.isErr && err == nil {
				t.Fatalf("expected an error [%v]", err)
			}
			if !tc.isErr && err != nil {
				t.Fatalf("got error: %v", err)
			}
		})
	}
}

type okGetter struct{ str string }

func (g okGetter) Get(string) (*http.Response, error) {
	w := httptest.NewRecorder()
	io.WriteString(w, g.str)
	return w.Result(), nil
}

type badGetter struct{ str string }

func (g badGetter) Get(string) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, g.str)
	return w.Result(), nil
}

type errorGetter struct{}

func (errorGetter) Get(string) (*http.Response, error) {
	return nil, errors.New("error")
}

type badJSONGetter struct{ str string }

func (g badJSONGetter) Get(string) (*http.Response, error) {
	w := httptest.NewRecorder()
	io.WriteString(w, g.str[0:17])
	return w.Result(), nil
}

type notFoundGetter struct{ str string }

func (g notFoundGetter) Get(string) (*http.Response, error) {
	w := httptest.NewRecorder()
	io.WriteString(w, strings.Replace(g.str, "200", "401", 1))
	io.WriteString(w, g.str[0:len(g.str)-4]+"401}")
	return w.Result(), nil
}

var jsonStrCurrent = `
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

var jsonStrForecast = `
{"city":{"id":1851632,"name":"Shuzenji",
"coord":{"lon":138.933334,"lat":34.966671},
"country":"JP"},
"cod":"200",
"message":0.0045,
"cnt":38,
"list":[{
        "dt":1406106000,
        "main":{
            "temp":298.77,
            "temp_min":298.77,
            "temp_max":298.774,
            "pressure":1005.93,
            "sea_level":1018.18,
            "grnd_level":1005.93,
            "humidity":87,
            "temp_kf":0.26},
        "weather":[{"id":804,"main":"Clouds","description":"overcast clouds","icon":"04d"}],
        "clouds":{"all":88},
        "wind":{"speed":5.71,"deg":229.501},
        "sys":{"pod":"d"},
        "dt_txt":"2014-07-23 09:00:00"}
        ]}`
