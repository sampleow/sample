# Quote

Quote Service module

## Getting started

Clone the project in a way that plays well with your `Go` local setup.

You also need a [recent version of Docker](https://docs.docker.com/docker-for-mac/install/) running on your machine.

After having run `make docker_run` you should see the two container running:

```
± |master S:2 U:4 ✗| → docker ps -a
CONTAINER ID        IMAGE                       COMMAND                  CREATED             STATUS                    PORTS                    NAMES
120f84a4d3e1        quote:latest                "/bin/sh -c /srv/quo…"   2 hours ago         Up 2 hours                0.0.0.0:8080->8080/tcp   quote
```
`make docker_stop` stops the container and removes the images. The Makefile uses dep to manage dependencies. [missy](https://github.com/microdevs/missy) is the open source micro-services 
framework being used.

### Run some test requests now?

Use Postman and perform GETs:
```
http://localhost:8080/generateBuyQuote?quantity=1
http://localhost:8080/generateSellQuote?quantity=1
```
### Unit tests
If you want to run tests for whole project, go to the root folder of the project and run:

```
go test ./...
```

## API Specifications
There are primarily two GET REST APIs, one for getting the quote for a BUY request and the other for the SELL request:
```
http://localhost:8080/generateBuyQuote?quantity=1
http://localhost:8080/generateSellQuote?quantity=1
```
Both of them expect `'quantity'` as a _mandatory_ parameter. `mux` (via `missy`) maps these requests to two handler functions:
```
func generateBuyQuoteHandler(w http.ResponseWriter, r *http.Request)
func generateSellQuoteHandler(w http.ResponseWriter, r *http.Request)
```

Internally the logical flow is as follows:

- fetch data from `https://api.mybitx.com/api/1/orderbook?pair=XBTZAR` to populate `OrderBook`
- if BUY, iterate thro `ASK` array from the beginning until the desired quantity is satisfied
- if SELL the iteration is on the 'BID' array

The API returns the quantity requested and total price as a JSON:

```
{
    "quantity": "2736.624252",
    "price": "67134069.72"
}
```

### Possible improvements to the API
- Data fetch is blocking/synchronous when a request comes in. It could have been fetched in an async manner via channels
- It has been observed that ASK is in ascending order of price and SELL in descending order. Maybe a data-structure 
such as a Min/MaxHeap could have been used for ordering instead of assuming order. 
- Better structuring and documentation of error messages in the API Specs
- More test cases 


