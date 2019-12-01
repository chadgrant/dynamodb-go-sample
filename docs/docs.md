# DYNAMODB SAMPLE

# Product Categories

Lists the defined product categories

> GET /category

```json
["Hats","Shirts","Pants","Shoes","Ties","Belts","Socks","Accessory"]
```

# Products Paged

Enumerates all the products in a specific category.

## Query Parameters
| Name      | Description                               |
|-----------|-------------------------------------------|
| last      | key of the product to start the page      |
| lastPrice | price of the product to start the page    |
| size      | the page size                             |

> GET /product/hats

```json
{
    "results":[{
        "id":"4438a096-09d1-4189-9fdf-a06ad54b8c7a",
        "name":"Test Hat",
        "category":"hats",
        "price":992.67,
        "description":"You should buy this"
    }],
    "next" :"/product/hats/?last=5cbd673f-e6e7-436b-97f0-94cb2e18b8a6&lastPrice=787.85",
    "total" :25
}
```

# Get Product

Returns a single product by id

> GET /product/5cbd673f-e6e7-436b-97f0-94cb2e18b8a6

```json
{
  "id":"4438a096-09d1-4189-9fdf-a06ad54b8c7a",
  "name":"Test Hat",
  "category":"hats",
  "price":992.67,
  "description":"You should buy this"
}
```

# Create Product

Adds a product. Returns a 201 (Created) status code and location header of where to obtain the newly created resource. (i.e. Location: /product/4438a096-09d1-4189-9fdf-a06ad54b8c7a)

> POST /product

```json
{
  "name":"Test Hat",
  "category":"hats",
  "price":992.67,
  "description":"You should buy this"
}
```

# Update Product

Update a product. Returns a 204 (No Content) status code on success

> PUT /product/4438a096-09d1-4189-9fdf-a06ad54b8c7a

```json
{
  "name":"Test Hat",
  "category":"hats",
  "price":992.67,
  "description":"You should buy this"
}
```

# Delete a Product

Delete a product. Returns a 204 (No Content) status code on success

> DELETE /product/4438a096-09d1-4189-9fdf-a06ad54b8c7a


# Metadata

Returns the metadata of the service.

> GET /metadata

| Property           | Description                                              |     Example                                |
| -------------------|----------------------------------------------------------|--------------------------------------------|
| application        | name of the service, normalized with no spaces           | sampleApi                                  |
| applicationFriedly | friendly name of the service                             | Sample API                                 |
| buildNumber        | The build number                                         | 2.5.7                                      |
| buildNumber        | The build number                                         | 2.5.7                                      |
| builtBy            | The user that did the build                              | cgrant                                     |
| builtWhen          | When the build was done                                  | 2015-03-12T19:40:18.877Z                   |
| compilerVersion    | The compiler version                                     | 1.5                                        |
| currentTime        | Time of request                                          | 2015-03-12T19:40:18.877Z                   |
| gitSha1            | The git sha1 hash of the build                           | d567d2650318f704747204815adedd2396a203f5   |
| gitBranch          | The git branch of the build                              | master                                     |
| groupId            | Team responsible for service                             | api                                        |
| machineName        | The name of the machine responding to this request       | server22                                   |
| osName             | Name of the OS of the machine responding to the request  | Linux                                      |
| osNumProcessors    | Number of processors of the machine responding           | 4                                          |
| upSince            | Time the service was started                             | 2015-03-12T19:40:18.877Z                   |
| version            | Current version of the service                           | 2                                          |

# HealthCheck

Returns the healthcheck of the service.

> GET /health

| Property          | Description                                              |     Example    |
| ------------------|----------------------------------------------------------|----------------|
| reportAsOf        | The time at which this report was generated (this may not be the current time) | 2015-03-12T19:40:18.877Z         |
| tests             | array of healthcheck test reports                        |  |
| interval          | How often the health checks are run in seconds                        | 10 |
| tests[].durationMilliseconds | Number of milliseconds taken to run the test  | 100 |
| tests[].name      | name of the healthcheck test                    | sql |
| tests[].result    | The state of the test, may be "notrun", "running", "passed", "failed" | passed |
| tests[].testedAt  | The last time the test was run | passed |

# Liveness / GTG - Good to Go


> GET /health/liveness


The "Good To Go" (GTG) returns a successful response in the case that the service is in an operational state and is able to receive traffic. This resource is used by load balancers and monitoring tools to determine if traffic should be routed to this service or not.

Note that GTG is not used to determine if the service is healthy or not, only if it is able to receive traffic. A healthy instance may not be able to accept traffic due to the failure of critical downstream dependencies.

A successful response is a 200 OK with a content of the text OK and a media type of "plain/text"

A failed response is a 5XX reponse with either a 500 or 503 response preferred. Failure to respond within a predetermined timeout typically 2 seconds is also treated as a failure.

# Readiness / Service Canary

> GET /health/readiness

The "Service Canary" (ASG) returns a successful response in the case that the service is in a healthy state. If a service returns a failure response or fails to respond within a predefined timeout then the service can expect to be terminated and replaced. (Typically this resouce is used in auto-scaling group healthchecks.)

A successful response is a 200 OK with a content of the text "OK" (including quotes) and a media type of "plain/text"

A failed response is a 5XX reponse with either a 500 or 503 response preferred. Failure to respond within a predetermined timeout typically 2 seconds is also treated as a failure.

# Debug Environment Name

Echo's back the configured environment name. Typically PRODUCTION,STAGING,DEVELOPMENT

> GET /debug/environment

# Debug Headers

Echo's back the headers sent to the service. Extremely useful in debugging CDN/Edge/Load balancer headers.

> GET /debug/headers

```json
{
  "Accept": ["text/html,application/xhtml+xml,application/xml"],
  "Accept-Encoding": ["gzip, deflate, br"],
  "Accept-Language": ["en,id;q=0.9,es;q=0.8,en-CA;q=0.7"],
  "Connection": ["keep-alive"]
}
```

# Debug Error

Useful to test error recovery, logging systems, etc.. intentionally throws exception/panics

> GET /debug/error