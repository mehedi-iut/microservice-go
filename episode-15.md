# Microservice in go
## refactor-1

In previous part, we create enum to add manually the exchange rate. but in this blog, we need fetch the data from bank and use it for conversion.
To do that, we will create *data* folder inside the **currency** module and create file **rates.go**
### **`rates.go`**
```go 
package data

type ExchangeRates struct {
    log hclog.Logger
    rate map[string]float64
}
```
here, we define **struct** to handle the exchange reate data pull from internet.
```go
func NewRates(l hclog.Logger) (*ExchangeRates, error){
    er := &ExchangeRates{log: l, rate: map[string]float64{}}
    err := er.getRates()
    return er, err
}
```
by following good go idiomatic principal, we create **NewRates** for struct **ExchangeRates** where we create an instance of **ExchangeRates**

```go
func (e *ExchangeRates) getRates() error {
    resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
    if err != nil {
        return err
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Expected success code 200 but got %d", resp.StatusCode)
    }
    defer resp.Body.Close()
}
```
Here, we define **getRates** method for **ExchangeRates**. we call the url and store the response in **resp**. if we get any error, we return the error. we also check **resp.StatusCode**, if it is not 200, we show the error. lastly, we close the response bobdy before closing the function.


Now, we need to extract data from xml to golang collection. To do that we need to define **struct**
```go
type Cubes struct {
    CubeData []Cube `xml:"Cube>Cube>wCube"`
}
```
we added *xml* notation **Cube>Cube>Cube** as xml has three **Cube** object

```go
type Cube struct {
    Currency string `xml:"currency,attr"`
    Rate string `xml:"rate,attr"`
}
```

in the above **Cube** struct, we are extracting *currency* and *rate*. we are extracting currency and rate attribute using **attr**

Now we need to parse xml. So to do that, we need to create instance of **Cubes** and decode the xml and save it to the **Cubes** instance. and then we will iterate over the **Cubes** instance and populate the golang collection

```go
func (e *ExchangeRates) getRates() error {
    resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
    if err != nil {
        return err
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Expected success code 200 but got %d", resp.StatusCode)
    }
    defer resp.Body.Close()
    md := &Cubes{}
    xml.NewDecoder(resp.Body).Decode(&md)

    for _, c := range md.CubeData {
        r, err := strconv.ParseFloat(c.Rate, 64)
        if err != nil{
            return err
        }

        e.rate[c.Currency] = r
    }
    return nil
}
```

Now we want to create a simple test in golang. To do that, we need to create a file inside the data folder which must contain **_test.go**. we will create **rates_test.go** and the function name must start with **Test** word and argument will contain ```(t *testing.T)```

```go
package data

import "testing"

func TestNewRates(t *testing.T){

}
```

above code is the structure of golang test codce

Now we will write simple test case following the above structure
```go
package data

import (
    "fmt"
    "testing"
    "github.com/hashicorp/go-hclog"
)

func TestNewRates(t *testing.T){
    tr, err := NewRates(hclog.Default())

    if err != nil{
        t.Fatal(err)
    }

    fmt.Printf("Rates %#v", tr.rates)
}
```
