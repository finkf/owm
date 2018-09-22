# OWM

Stupid [open weather map](http://openweathermap.org) API
implementation in [go](https://golang.org).

## Examples

Basic usage:
```golang
import (
	   "fmt"
	   "net/http"

	   "github.com/finkf/owm"
)

func main() {
	api := owm.API{
	 	 Client: &http.Client{},
		 Key: "secret-key"
	}
	c, err := api.Current(owm.Query{City: "London"})
	if err != nil {
	   panic(err)
	}
	fmt.Printf("%v\n", *c)
}
```