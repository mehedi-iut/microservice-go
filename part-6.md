# Micro-service in go
## Part-6

In this section, we will discuss about **JSON Validation** using **go-validator**[link](https://github.com/go-playground/validator)


In our **Product** struct in **data** handler, we need add **validate** tag.
```go
type Product struct {
    ID          int     `json:"id"`
    Name        string  `json:"name" validate:"required"`
    Description string  `json:"description"`
    Price       float32 `json:"price" validate:"gt=0"`
    SKU         string  `json:"sku" validate:"required,sku"`
    CreatedOn   string  `json:"-"`
    UpdatedOn   string  `json:"-"`
    DeletedOn   string  `json:"-"`
}
```

so in the tag we added **validate** tage with **required**. so if we don't pass **Name** field, it will failed. **Price** field, it need to be greater than 0. so we add **gt=0**. and for **SKU**, we need custom validation so we add **sku** and we need to define custom function for validation.

we can write very basic unit test to test the json validation. In the **data** handler, add one go file **products_test.go**
```go
package data

import "testing"

func TestChecksValidation(t *testing.T){
    p := &Product{}

    err := p.Validate()

    if err != nil {
        t.Fatal(err)
    }
}
```

if we run the above test it will fail because ```p := &Product{}``` is empty but we must pass **name**, **price** and **sku**
but before run the test we need to define **Validate** function in **data** handler **products.go**

```go
func (p *Product) Validate() error {
    validate := validator.New()
    return validate.Struct(p)
}
```
So first we create a instance of **validator** and then pass our struct in ```validate.Struct(p)```. if it can't validate, it will return error, so we just return ```validate.Struct(p)```

For Custom Valdation Functions we need to use ```validate.RegisterValidation("custom tag name", customFunc)```

we will create **validateSKU** function accroding to **validator** documentation [link](https://pkg.go.dev/github.com/go-playground/validator#hdr-Custom_Validation_Functions)

```go
func validateSKU(fl validator.FieldLevel) bool {
    // sku is of format abc-absd-dfsdf
    re := regexp.MustCompiled(`[a-z]+-[a-z]+-[a-z]+`)
    matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}

	return true
}
```

Now we need to update the **Validate** function
```go
func (p *Product) Validate() error {
    validate := validator.New()
    validate.RegisterValidation("sku", validateSKU)
    return validate.Struct(p)
}
```
Now we can run test with sku value in struct and it will pass, if we don't provide any sku or wrong sku format, it will show error

Now we need to add the validation in **MiddlewareValidateProduct**

```go
func(p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

        // validate the product
        err = prod.Validate()
        if err != nil{
            p.l.Println("[ERROR] validating product", err)
            http.Error(rw,
            fmt.Sprintf("Error validating product: %s", err),
            http.StatusBadRequest,
            )
            return
        }

		// add the product to the context
		ctx := context.WithValue(r.Context, KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
```



