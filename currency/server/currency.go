package server

import (
	"context"
	protos "currency/currency"
	"currency/data"
	"github.com/hashicorp/go-hclog"
	"io"
	"time"
)

type Currency struct {
	rates                              *data.ExchangeRates
	log                                hclog.Logger
	subscriptions                      map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
	protos.UnimplementedCurrencyServer // Embed the UnimplementedServiceServer struct
}

//func (c *Currency) mustEmbedUnimplementedCurrencyServer() {} // Implement the mustEmbedUnimplementedCurrencyServer method

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{rates: r, log: l, subscriptions: make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go c.handleUpdates()

	return c
	//go func() {
	//	ru := r.MonitorRates(5 * time.Second)
	//	for range ru {
	//		l.Info("Got Updated rates")
	//	}
	//}()
	//return &Currency{r, l, make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got Updated rates")

		for k, v := range c.subscriptions {
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get update rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}

				err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})
				if err != nil {
					c.log.Error("Unable to send updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
	//	go func() {
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
			break
		}
		// added inbound message
		c.log.Info("Handle client request", "request", rr)

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}
	//}()

	//for {
	//	err := src.Send(&protos.RateResponse{Rate: 12.1})
	//	if err != nil {
	//		return err
	//	}
	//	time.Sleep(5 * time.Second)
	//}
	return nil
}
