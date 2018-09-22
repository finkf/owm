package owm

// Coord defines the longitude and latitude for a location.
type Coord struct {
	Lon, Lat float64
}

// Sys defines the country and sunset and sunrise timestamps.
type Sys struct {
	Country         string
	Sunrise, Sunset int64
}

// Weather defines basic information about the weather.
type Weather struct {
	ID                      int
	Main, Description, Icon string
}

// Main defines the basic weather information.
type Main struct {
	Temp               Kelvin
	TempMin            Kelvin  `json:"temp_min"`
	TempMax            Kelvin  `json:"temp_max"`
	SeaLevel           float64 `json:"sea_level"`
	Groundlevel        float64 `json:"grnd_level"`
	Humidity, Pressure int
}

// Wind defines the speed and degrees of the wind.
type Wind struct {
	Speed, Deg float64
}

// Volume defines the volumes of rain or snow of
// the last 3 hours.
type Volume struct {
	H3 int `json:"3h"`
}

// Clouds are clouds.
type Clouds struct {
	All int
}

// Current defines all information for the current weather.
type Current struct {
	Coord         Coord
	Sys           Sys
	Weather       []Weather
	Main          Main
	Wind          Wind
	Rain, Snow    Volume
	Clouds        Clouds
	DT            int64
	ID, Cod       int
	Name, Message string
}

// City defines information about a city.
type City struct {
	ID            int
	Name, Country string
	Coord         Coord
}

// Forecast defines information about a 5 day / 3 hours forecast.
type Forecast struct {
	Cod, Cnt int
	Message  string
	City     City
	List     []ForecastItem
}

// ForecastItem defines an item for a forecast.
type ForecastItem struct {
	DT         int64
	DTTXT      int64 `json:"dt_txt"`
	Main       Main
	Weather    []Weather
	Clouds     Clouds
	Wind       Wind
	Rain, Snow Volume
}

// Kelvin defines temperatures in Kelvin.
type Kelvin float64

// Celcius returns the celcius equivalent of the temperature.
func (k Kelvin) Celcius() float64 {
	return float64(k) - 273.15
}

// Fahrenheit returns the fahrenheit equivalent of the temperature.
func (k Kelvin) Fahrenheit() float64 {
	return k.Celcius()*1.8 + 32.0
}
