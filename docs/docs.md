## DYNAMODB SAMPLE

## Metadata

Returns the metadata of the service.
~~~
GET /metadata
~~~

| Property          | Description                                              |     Example    |
| ------------------|----------------------------------------------------------|----------------|
| buildNumber       | The build number                                         | 2.5.7          |
| builtBy           | The user that did the build                              | cgrant         |
| builtWhen         | When the build was done                                  | 2015-03-12T19:40:18.877Z        |
| compilerVersion   | The compiler version                                     | 1.5            |
| currentTime       | Time of request                                          | 2015-03-12T19:40:18.877Z        |
| gitSha1           | The git sha1 hash of the build                           | d567d2650318f704747204815adedd2396a203f5         |
| groupId           | Team responsible for service                             | api        |
| machineName       | The name of the machine responding to this request       | server22        |
| osName            | The name of the OS of the machine responding to the request | Linux        |
| osNumProcessors   | The number of processors of the machine responding to the request | 4        |
| upSince           | Time the service was started                             | 2015-03-12T19:40:18.877Z        |
| version           | Current version of the service                           | 2 |



## HealthCheck

Returns the healthcheck of the service.
~~~
GET /health
~~~

| Property          | Description                                              |     Example    |
| ------------------|----------------------------------------------------------|----------------|
| reportAsOf        | The time at which this report was generated (this may not be the current time) | 2015-03-12T19:40:18.877Z         |
| tests             | array of healthcheck test reports                        |  |
| interval          | How often the health checks are run in seconds                        | 10 |
| tests[].durationMilliseconds | Number of milliseconds taken to run the test  | 100 |
| tests[].name      | name of the healthcheck test                    | sql |
| tests[].result    | The state of the test, may be "notrun", "running", "passed", "failed" | passed |
| tests[].testedAt  | The last time the test was run | passed |



## GTG - Good to Go

The "Good To Go" (GTG) returns a successful response in the case that the service is in an operational state and is able to receive traffic. This resource is used by load balancers and monitoring tools to determine if traffic should be routed to this service or not.

Note that GTG is not used to determine if the service is healthy or not, only if it is able to receive traffic. A healthy instance may not be able to accept traffic due to the failure of critical downstream dependencies.

A successful response is a 200 OK with a content of the text OK and a media type of "plain/text"

A failed response is a 5XX reponse with either a 500 or 503 response preferred. Failure to respond within a predetermined timeout typically 2 seconds is also treated as a failure.

~~~
GET /health/gtg
~~~


## ASC - Service Canary

The "Service Canary" (ASG) returns a successful response in the case that the service is in a healthy state. If a service returns a failure response or fails to respond within a predefined timeout then the service can expect to be terminated and replaced. (Typically this resouce is used in auto-scaling group healthchecks.)

A successful response is a 200 OK with a content of the text "OK" (including quotes) and a media type of "plain/text"

A failed response is a 5XX reponse with either a 500 or 503 response preferred. Failure to respond within a predetermined timeout typically 2 seconds is also treated as a failure.


~~~
GET /health/asg
~~~