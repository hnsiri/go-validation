## Installation

You can install `go-validator` using `go get` command. Open your terminal or command prompt and run the following command:

```bash
go get -u github.com/hnsiri/go-validator
```

This command will download and install the `go-validator` package along with its dependencies into your Go workspace.

## Usage

Once installed, you can import and use `go-validator` in your Go projects. Here's a basic example of how to use it:

```go
package main

import (
	"fmt"

	"github.com/hnsiri/go-validator"
)

type Currency struct {
	ISO  string `json:"iso"`
	Name string `json:"name"`
}

type Payload struct {
	Name          string   `json:"name"`
	Handle        string   `json:"handle"`
	Currency      Currency `json:"currency"`
	ContactNumber string   `json:"contact_number"`
	ContactEmail  string   `json:"contact_email"`
}

func main() {
	req := &Payload{ /* initialize your Payload struct here */ }

	v := validator.New(req, validator.Fields{
		"Name":   validator.Rules(validator.Required, validator.MinLength(3), validator.MaxLength(100)),
		"Handle": validator.Rules(validator.Required, validator.MinLength(3), validator.MaxLength(32)),
		"Currency": validator.Rules(validator.Struct(
			req.Currency, validator.Fields{
				"ISO": validator.Rules(validator.Required, validator.ISO4217),
			}),
		),
		"ContactNumber": validator.Rules(validator.Required),
		"ContactEmail":  validator.Rules(validator.Email),
	})

	if ok := v.Validate(); !ok {
		fmt.Println(v.Errors())
		// return bad request 
	}

	// data is valid implement your logic.....
}

```

## Contributing

If you encounter any issues or would like to contribute to the development of `go-validator`, you can do so by submitting issues or pull requests on the GitHub repository.

## License

Distributed under the [MIT License](https://opensource.org/licenses/MIT).
```