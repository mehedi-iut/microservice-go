# Microservice in go
## refactor-2

In our previous session, we implemented the getRates service to retrieve exchange rates from the European Central Bank, which returns data in XML format. We then parse the XML and create our own data structure.

One important consideration is that the European Central Bank returns rates based on the euro as the reference currency. For example, the value you get for the US dollar is the number of US dollars equivalent to one euro. This means we need to account for the fact that the exchange rate for the euro itself is always 1, which is crucial for calculating different rates as a ratio.

To add the euro as a rate into our exchange rate currency, we simply added ```e.rates["EUR"] = 1``` before returning.

Now, we can query this getRates service to get the exchange rate by creating a new function called GetRate. This function calculates the exchange rate ratio from one currency to another currency. For example, if we provide the base currency as USD and the destination currency as GBP, the function will return the rate by dividing the destination by the base. If the base currency is EURO, the function will divide by 1 as we set earlier, and we will get the same rate provided by the European Central Bank, as they took the EURO as the base currency.

```go
func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected success code but got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		e.rates[c.Currency] = r
	}

    e.rates["EUR"] = 1

	return nil
}
```

```go
func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
    br, ok := e.rates[base]
    if !ok {
        return 0, fmt.Errorf("Rate not found for currency %s", base)
    }

    dr, ok := e.rates[dest]
    if !ok {
        return 0, fmt.Errorf("Rates not found for currency %s", dest)
    }

    return dr/br, nil
}
```

in **server/currency.go**, we need to get rates and return that rate

```go
package server

import (
	"context"
	protos "currency/currency"
    "currency/data"
	"github.com/hashicorp/go-hclog"
)

type Currency struct {
    rates *data.ExchangeRates
	log hclog.Logger
	protos.UnimplementedCurrencyServer
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{rates: r, log: l}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

    rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
    if err != nil{
        return nil, err
    }

	return &protos.RateResponse{Rate: rate}, nil
}
```
here, we add ```rates *data.ExchangeRates``` in **Currency** struct and then added in the **NewCurrency** function. and in the **GetRate** method we added this code ```rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())``` to get rate and then return it.

Here we have a small problem, **protos** *RateResponse* we are returning *float32* but our *rate* is *float64*. so we need to change protos RateResponse definition.

so in **protos/currency.proto** we need to change
```go
message RateResponse {
  float rate = 1;
}
```
to *double*

```
message RateResponse {
  double rate = 1;
}
```

then run the makefile to generate the code
```make protos```

now in the **main.go** we need to create *NewRates* instance and pass to the *NewCurrency*

```go
log := hclog.Default()

rates, err := data.NewRates(log)
if err != nil{
    log.Error("Unable to generate rates", "error", err)
    os.Exit(1)
}

gs := grpc.NewServer()

c := server.NewCurrency(rates, log)
```

When we created our product API, we started implementing the currency service logic inline with the product calls, which was not ideal. This is because we were mutating the price inside of the get handler, even though the price is a property of the data object that contains our products. To improve this, we will remove the currency service calls from the get handler and move them into the data class.

In the **data/products.go** file, we initially created a simple set of functions, but we now need to refactor the code to make it more scalable. Specifically, we want to be able to pass in an object representing our rates currency service, which requires a constructor. Although we could use the **init** function, this is not recommended. Instead, we will follow the **constructor injection** approach, which will make testing easier.

Now we need to refactor the **product-api**
in **data/products.go**, we need to add below ```type Products []*Product```
```go
type ProductDB struct {
    currency protos.CurrencyClient
    log hclog.Logger
}

func NewProductDB(c protos.CurrencyClient, l hclog.Logger) *ProductDB {
    return &ProductDB{c, l}
}
```

Now we need to modify the *func* call to method call

```go
func GetProducts() Products {
	return productList
}
```

we have this previously, now we will convert it to the method
```go
func (p *ProductsDB) GetProducts(currency string)(Products, error){
    if currency == ""{
        return productList, nil
    }

    // copied from get.go handler
    rr := &protos.RateRequest{
        Base:        protos.Currencies(protos.Currencies_value["EUR"]),
        Destination: protos.Currencies(protos.Currencies_value["GBP"]),
    }

    resp, err := p.currency.GetRate(context.Background(), rr)
    if err != nil{
        p.log.Error("Unable to get rate", "currency", currency, "error", err)
        return nil, err
    }
    // above portion copied from get.go handler

    pr := Products{}
    for _, p := range productList{
        np := *p
        np.Price = np.Price * resp.Rate
        pr = append(pr, &np)
    }

    return pr, nil
}
```

Now in the above code 
```go 
if currency == ""{
    return productList, nil
}
```
so here, we don't supply the currency, we will take Euro as base currency and return the productList as we don't need to manipulate the price

in the below code we manipulate the price
```go
pr := Products{}
for _, p := range productList{
    np := *p
    np.Price = np.Price * resp.Rate
    pr = append(pr, &np)
}
```

we already define **Products** in the *data/products.go* using this code ```type Products []*Product``` 
we initialize the *Products* using ```pr := Products{}``` then we iterate over the *productList* but we can't use the **p** to manipulate the price, if we do that, it will change the original price that's why we create the *np* variable ```np := *p```, we deference the *p* to get the copy of *p* which is flat but not deep and will not manipulate the original value. then we manipulate the price ```np.Price = np.Price * resp.Rate``` then we append our new price in **pr** ```pr = append(pr, &np)```

Now, the below code will be added in other function as well, so we will keep that in another helper function to reduce the code duplication
```go
rr := &protos.RateRequest{
    Base:        protos.Currencies(protos.Currencies_value["EUR"]),
    Destination: protos.Currencies(protos.Currencies_value["GBP"]),
}

resp, err := p.currency.GetRate(context.Background(), rr)

```

we also hardcoded the destination currency, we need to handle that as well
so we will create helper function **getRate**

```go
func (p *ProductDB) getRate(destination string) (float64, error){
    rr := &protos.RateRequest{
    Base:        protos.Currencies(protos.Currencies_value["EUR"]),
    Destination: protos.Currencies(protos.Currencies_value[destination]),
    }

    resp, err := p.currency.GetRate(context.Background(), rr)
    return resp.Rate, err
}
```

now we need to modify our **GetProducts** method that we define in above
```go
func (p *ProductsDB) GetProducts(currency string)(Products, error){
    if currency == ""{
        return productList, nil
    }

    rate, err := p.getRate(currency)
    if err != nil{
        p.log.Error("Unable to get rate", "currency", currency, "error", err)
        return nil, err
    }
    // above portion copied from get.go handler

    pr := Products{}
    for _, p := range productList{
        np := *p
        np.Price = np.Price * rate
        pr = append(pr, &np)
    }

    return pr, nil
}

```

now we will modify the **GetProductByID** which is like below
```go
func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}
```
now we will convert it to the method
```go
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error){
    i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

    if currency == ""{
        return productList[i], nil
    }

    rate, err := p.getRate(currency)
    if err != nil{
        p.log.Error("Unable to get rate", "currency", currency, "error", err)
        return nil, err
    }

    // we are doing things below to avoid changing the original price
    np := *productList[i]
    np.Price = np.Price*rate

    return &np, nil
}
```

Now **UpdateProduct**, we don't change the currency while updating the product, we will assume, it will always be in base currency
but we will change it to the method

```go
func UpdateProduct(p Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &p

	return nil
}
```

now we will change it to the method by just adding ```(p *ProductsDB)```
```go
func (p *ProductsDB) UpdateProduct(pr Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &pr

	return nil
}
```