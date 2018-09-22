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
		Key: "personal-owm-api-key"
	}
	c, err := api.Current(owm.ByCity{City: "London"})
	if err != nil {
	   panic(err)
	}
	fmt.Printf("%v\n", *c)
}
```