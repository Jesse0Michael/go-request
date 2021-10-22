#GO Request

Go Request is a package built to decode `*http.Request` data into a custom struct. Using struct tags we are about to pull HTTP request data from the request's:
- Body
- Header
- Query
- Path

Example using the following request:
```bash 
curl --location --request POST 'www.example.com/user/adam?game=go' \
--header 'X-DELAY: 60' \
--header 'Content-Type: application/json' \
--data-raw '{
    "state": "idle"
}'
```


```go
import ("github.com/jesse0michael/go-request")

type MyRequest struct {
    Name  string `path:"name"`
    Game  string `query:"game"`
    State string `json:"state"`
    Delay int64  `header:"X-DELAY"`
}


func(w http.ResponseWriter, r *http.Request) {
    var req MyRequest

    err := request.Decode(r, &req)
    if err != nil {
        w.WriteHeader(400)
    }

    fmt.Println(req)
}
```

```

```
