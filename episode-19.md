# Microservice in go
## bi-directional gRPC-2

In this episode, we will work on grpc bi-directional communication. In the previous epidose, we implement the simple bi-directional communication, we will continue to work on the same code to listen to any changes and updated it.

In this project we call European central bank api to get rates, we want to update our rate when it changes in the European Central Bank.

First, we need to modify codes in **currency.go**

```go
type Currency struct {
    rates *data.ExchangeRates
    log hclog.Logger
    subscriptions map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
    go func(){
        ru :=  r.MonitorRates(5 * time.Second)
        for range ru {
            l.Info("Got Updated rates")
        }
    }()
    return &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
}
```

Here, we need to break up subscription because at the moment we are not listening for any particular messages so we want to keep track of the different clients and what they're interested in. so we need to track that subscription and we can just implement a simple Map and use that as a cache. that's why we added ```subscriptions``` in the ```Currency``` *struct*.

Now in the ```NewCurrency```, we are returning that map ```make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)``` and we added ```go func()``` so that it shouldn't block the execution. and ```ru := r.MonitorRates(5 * time.Second)``` we will implement```MonitorRates``` later in this blog in **data** module


so when it comes to using it we can go back down to our **SubscribeRates** so whenever we get a new subscriber rate message we want to update our map so we can update our map by doing ```c.subscriptions[src]``` and then the key that we're going to use is going to be the actual client and what we want to do is we want to append the the rate to that collection so the first thing we need to do is we need to just get the the collection of rate requests and we can do that just like this ```rrs, ok := c.subscriptions[src]``` now if if this is kind of empty we want to create our our original map so we can just do a little check there we can say ```if !ok``` then what we want to do is create our map which is going to be like this
```go
rrs, ok := c.subscriptions[src]
if !ok{
  rrs = []*protos.RateRequest{}
}
```

so that is all that's doing is that's just updating our object then again we'll just set that subscription again like that so now every time we get a subscription message we're just appending it to our little cash. so now rather than kind of just sending this random message what we're going to do is we're only going to send back the rate messages when things sort of change. So we will get rid of the below *for* loop which randomly send messages every 5 seconds
```go
for {
    err := src.Send(&protos.RateResponse{Rate: 12.1})
    if err != nil{
        return err
    }

    time.Sleep(5 * time.Second)
}
```

we remove the above code and also remove the ```go func()```. and our code look like below
```go
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	// go func() {
		for {
			rr, err := src.Recv()
			// below technically not error
			// it happens when client close the connection
			if err == io.EOF {
				c.log.Info("Client has closed connection")
				break
			}

			if err != nil {
				c.log.Error("Unable to read from client", "error", err)
				return err
			}
			// added inbound message
			c.log.Info("Handle client request", "request", rr.GetBase(), "request_dest", rr.GetDestination())

            rrs, ok := c.subscriptions[src]
            if !ok {
                rrs = []*protos.RateRequest{}
            }

            rrs = append(rrs, rr)
            c.subscriptions[src] = rrs
		}
    return nil
	// }()

	// for {
	// 	err := src.Send(&protos.RateResponse{Rate: 12.1})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }
}
```
so let's have a look at how we can actually implement subscription so in our **MonitorRates** what we're doing is whenever we get an update to say there are new rates available we are just logging something now what we want to do is we want to do something a little bit more sophisticated we want to find out who has subscribed to a particular rate from our subscription collection and we want to then send them a message just kind of say hey there's an update to the rate that you're interested in 

now we need to refactor our ```NewCurrency``` and remove the **MonitorRates** and handle that in seperate method

```go
func NewCurrency(r *data.ExchangeRates,l hclog.Logger) *Currency {
    c := &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
    go c.handleUpdates()

    return c
}

func (c *Currency) handleUpdates() {
    ru := c.rates.MonitorRates(5 * time.Second)
    for range ru {
        c.log.Info("Got Updated rates")

        for k, v := range c.subscriptions{
            for _, rr := range v {
                r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
                if err != nil{
                    c.log.Error("Unable to get update rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
                }

                err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})
                if err != nil{
                    c.log.Error("Unable to send updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
                }
            }
        }
    }
}
```

Now, we can't send the response with base and destination currency along with rate as we didn't define that in the proto file, but we need to return the base and destination currency to get more clear context which currency got changed

so in the **proto** file update the ```RateResponse``` message
```
message RateResponse {
    Currencies Base = 1;
    Currencies Destination = 2;
    double rate = 3;
}
```
and run the ```make protos``` to generate the proto file

As we changed our proto file, we need to update our **GetRate** method to return rate with base and destination currencies

```go
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

```

### quick recap
so now just quick recap of what we've just done so subscribe rates allows the client to subscribe to notifications for updated rates that has a combination of the base and the destination we have our method here called handle updates and our client which is connected to the bank whenever we get updates from the bank we're gonna forward those on to our interested clients 

#### MonitorRates method
this is to simulate the changes of currency and send updated changes to the listener. This is only for demonstration purpose, not a real world case
in our currency service in **data/rates.go** add the below method under *GetRate*

```go
// MonitorRates checks the rates in the ECB API every interval and sends a message to the
// returned channel when there are changes
//
// Note: the ECB API only returns data once a day, this function only simulates the changes
// in rates for demonstration purposes
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				// just add a random difference to the rate and return it
				// this simulates the fluctuations in currency rates
				for k, v := range e.rates {
					// change can be 10% of original value
					change := (rand.Float64() / 10)
					// is this a postive or negative change
					direction := rand.Intn(1)

					if direction == 0 {
						// new value with be min 90% of old
						change = 1 - change
					} else {
						// new value will be 110% of old
						change = 1 + change
					}

					// modify the rate
					e.rates[k] = v * change
				}

				// notify updates, this will block unless there is a listener on the other end
				ret <- struct{}{}
			}
		}
	}()

	return ret
}

```

let's break down above code, first, it defines a method called MonitorRates for the ExchangeRates struct. This method continuously monitors rates in the European Central Bank (ECB) API at a specified interval and sends a message to the returned channel when there are changes.

Here is a breakdown of the code with explanations:

```go
// MonitorRates checks the rates in the ECB API every interval and sends a message to the
// returned channel when there are changes
//
// Note: the ECB API only returns data once a day, this function only simulates the changes
// in rates for demonstration purposes
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
```
The first line is a comment explaining what the ```MonitorRates``` method does. The method takes an argument ```interval``` of type ```time.Duration``` and returns a channel of empty structs. This channel will be used to send notifications to any code that is listening to changes in the rates.

```go
ret := make(chan struct{})
```

This line creates a new channel of empty structs and assigns it to the ```ret``` variable.
```go
go func() {
    ticker := time.NewTicker(interval)
    for {
        select {
        case <-ticker.C:
            // ...
            ret <- struct{}{}
        }
    }
}()
```

This section of the code creates a new goroutine (concurrent execution) to monitor the rates. It uses a ```Ticker``` from the ```time``` package to run a loop that wakes up at the specified ```interval```. When the ticker wakes up, it executes the code inside the ```case``` block of the ```select``` statement. This block randomly changes the rates and sends a notification to the ```ret``` channel to indicate that the rates have been updated.
```go
for k, v := range e.rates {
    // change can be 10% of original value
    change := (rand.Float64() / 10)
    // is this a postive or negative change
    direction := rand.Intn(1)

    if direction == 0 {
        // new value with be min 90% of old
        change = 1 - change
    } else {
        // new value will be 110% of old
        change = 1 + change
    }

    // modify the rate
    e.rates[k] = v * change
}
```
This section of the code randomly changes the rates for each currency by a positive or negative amount that is 10% of the original value. The code selects a random direction and sets the change accordingly. If the direction is 0, the new value will be a minimum of 90% of the old value, and if the direction is 1, the new value will be a maximum of 110% of the old value. Finally, the code modifies the rate in the ```e.rates``` map.

```go
ret <- struct{}{}
```
This line sends an empty struct to the ```ret``` channel to indicate that the rates have been updated.
```go
return ret
```

This line returns the ```ret``` channel to the caller.


#### changes in product-api
In the product-api, **get.go**, we need to modify the code, so that it will listen to rate changes and before returning any value, it should check the cache first

we need to go to product-api service  **data/products.go** and **getRate** method which is called from **ListSingle** and **ListAll** from **handlers/get.go**

we need to tidy this up a little bit because what we're actually doing is we're making that rate request every single time but we've now got subscriber method and we know that we can kind of keep track of things we can cache the data so again in the same way as we we had that set up before we can cache our rates. so let's just create a very simple simple cache here on our ```ProductsDB``` struct.
```go
type ProductsDB struct {
    currency protos.CurrencyClient
    log hclog.Logger
    rates map[string]float64
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
    return &ProductsDB{c, l, make(map[string]float64)}
}
```

Now, in the **getRate** method, we add if condition to check the cache
```go
func (p *ProductsDB) getRate(destination string) (float64, error){
    // if cached return
    if r, ok := p.rates[destination]; ok {
        return r, nil
    }
    // ......
}
```

Now we need to subscribe to the future update
but we can't add this code ```err = p.currency.SubscribeRates(context.Background(), rr)``` because, we need to create SubscribeRates Client. so we need to create client for bi-directinal stream messaging. as we are calling gRPC service from products-api, we are now client.

we're handling it from the client side now so I need to subscribe for rates and I need to get that subscription client and it's gonna be like this ```sub, err := p.currency.SubscribeRates(context.Background())``` this is going to allow me to receive messages.
so we add receive method ```sub.Recv()``` now this is exactly the same as what was going on in my currency service 
```go
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
    // handle client messages
    for {
        rr, err := src.Recv()
        // ....
    }
} 
``` 
it's pretty much just the reverse we're implementing it on the client side because it's bi-directional so in both sides we have sanding on both sides we have receiving

now ```sub.Recv()``` return rates and error and we need to add it in the loop to get continuous message and ```sub.Recv()``` will block the execution. once we get the updated rate, we will update it in out cache(map)
so our updated code in **data/products.go** added after **func NewProductDB**
```go
func(p *ProductsDB) handleUpdates(){
    sub, err := p.currency.SubscribeRates(context.Background())
    if err != nil{
        p.log.Error("Unable to subscribe for rates", "error", err)
    }
    for {
        rr, err := sub.Recv()
        if err != nil{
            p.log.Error("Error receiving message", "error", err)
            return
        }
        p.rates[rr.Destination.String()] = rr.Rate
    }
}
```

Now we don't want to call currency service everytime rather use our above **handleUpdates** method
```go
func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB{
    pb := &ProductsDB{c, l, make(map[string]float64)}
    go pb.handleUpdates()
    return pb
}
```

now, we're creating the client in **handleRates** we need to use this client  that's returned from **SubscribeRates** in order to send messages and we don't have a reference to it in ```type ProductsDB struct``` so let's let's add the reference

```go
type ProductsDB struct {
    currency protos.CurrencyClient
    log hclog.Logger
    rates map[string]float64
    client protos.Currency_SubscribeRatesClient
}
```
and in **handleUpdates** we need to add this line ```p.Client = sub``` to have the reference

```go
func(p *ProductsDB) handleUpdates(){
    sub, err := p.currency.SubscribeRates(context.Background())
    if err != nil{
        p.log.Error("Unable to subscribe for rates", "error", err)
    }
    // added the client reference
    p.Client = sub

    for {
        rr, err := sub.Recv()
        // added log to see the output
        p.log.Info("Received updated rate from server", "dest", rr.GetDestination().String())
        if err != nil{
            p.log.Error("Error receiving message", "error", err)
            return
        }
        p.rates[rr.Destination.String()] = rr.Rate
    }
}
```

Finally, we can now send messages and subscribe for getting update in **getRate** method

```go
func (p *ProductsDB) getRate(destination string) (float64, error){
    // if cached return
    if r, ok := p.rates[destination]; ok {
        return r, nil
    }

    rr := &protos.RateRequest{
        Base: protos.Currencies(protos.Currencies_value["EUR"]),
        Destination: protos.Currencies(protos.Currencies_value[destination]),
    }

    // get initial rate
    resp, err := p.currency.GetRate(context.Background(), rr)
    // then update the cache
    p.rates[destination] = resp.Rate

    // subscribe for updates
    p.client.Send(rr)

    return resp.Rate, err
}
```

now whenever we call this get handler in products-api **get.go**,  we're gonna go and call **GetProducts** which is going to get the currency rate, if needed it's gonna mutate it to get the correct conversion rate. then we call **getRate** which is going to try and use the cash if it has a cached item and it's gonna subscribe for updates if it doesn't. Also when we create a new instance of ```type ProductsDB```,  we're also creating a bi-directional stream to the currency server which is awaiting updates that we're subscribing to

Now we can run our currency service and products-api service by running ```go run main.go``` and then call our products api with currency USD
```bash
curl "localhost:9090/products?currency=USD" | jq
```

we need to give "" around url because *?* will be interpreted by bash and will show error

