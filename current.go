package owm

type Coord struct {
	Lon, Lat float64
}

type Sys struct {
	Country         string
	Sunrise, Sunset uint64
}

type Weather struct {
	ID                      int
	Main, Description, Icon string
}

type Main struct {
	Temp               Kelvin
	TempMin            Kelvin `json:"temp_min"`
	TempMax            Kelvin `json:"temp_max"`
	Humidity, Pressure int
}

type Wind struct {
	Speed, Deg float64
}

type Volume struct {
	H3 int `json:"3h"`
}

type Clouds struct {
	All int
}

type Current struct {
	Coord         Coord
	Sys           Sys
	Weather       []Weather
	Main          Main
	Wind          Wind
	Rain, Snow    Volume
	Clouds        Clouds
	DT            uint64
	ID, Cod       int
	Name, Message string
}

type Kelvin float64

func (k Kelvin) Celcius() float64 {
	return float64(k) - 273.15
}

func (k Kelvin) Fahrenheit() float64 {
	return k.Celcius()*1.8 + 32.0
}
