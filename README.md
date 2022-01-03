# GO Request

Go Request is a package built to decode `*http.Request` data into a custom struct. Using struct tags we are about to pull HTTP request data from the request's:
- Body
- Header
- Query
- Path

By default the request body will be decoded into the input struct based off of the `Content-Type` header, unless a field with the body tag is specified. 

## Tags
Use struct tags to define where a field should be pulled from. Specify the name to lookup the value by in the tag value. Some tags support options using comma separated value strings, the first of which always being the lookup name.

### `query`
Assigns values by query parameter. Conversion can be controlled with the following options on the tag following a `,` after query parameter name.
- `explode` when set, the request will be decoded expecting multiple query parameters by the same name, following the [OAS Specification](https://swagger.io/docs/specification/serialization/) serialization keyword. Otherwise the request will be decoded expecting the parameter to be delineated with commas.

### `header`
Assigns values by http header.

### `path`
Using [gorilla.mux](github.com/gorilla/mux) router path values, assigns values by path vars.

### `body`
Assigns value from http request body. Useful if the request body is an array, because `request.Decode` only accepts struct inputs.

## Types
Go Request supports the following types as well as slices of these types:
```
string, bool, int, int64, int32, int16, int8, float64, float32, uint, uint64, uint32, uint16, uint8, complex128, complex64, time.Time, time.Duration
```

## Notes
> To avoid potentially overwriting fields not pulled from the request body with values pulled from the request body. use a `body` tag on a sub field or add a tag to ignore the field when decoding, i.e. `json:"-"`.

> To decode a request body that is an array, decode into a field using a `body` tag.

> If a struct tag has multiple Go Request tags the value will be assigned by the following hierarchy `body` > `header` > `query`

---

## Example

```bash 
curl --location --request POST 'www.example.com/users/adam?active=true&friend=bob&friend=steve' \
--header 'X-DELAY: 60' \
--header 'Content-Type: application/json' \
--data-raw '{
    "state": "idle"
}'
```


```go
import ("github.com/jesse0michael/go-request")

type MyRequest struct {
	User    string   `path:"user"`
	Active  bool     `query:"active"`
    Friends []string `query:"friend,explode"`
	State   string   `json:"state"`
	Delay   int      `header:"X-DELAY"`
}


func(w http.ResponseWriter, r *http.Request) {
    var req MyRequest
    err := request.Decode(r, &req)
    if err != nil {
        w.WriteHeader(400)
    }

    fmt.Printf("%+v",req)
}
```

```sh
{User:adam Active:true Friend:[bob steve] State:idle Delay:60}
```

