# jqrp -- jq Reverse Proxy

![CI status](https://github.com/bauerd/jqrp/actions/workflows/ci.yaml/badge.svg?branch=master)

HTTP reverse proxy that mutates JSON responses according to client-supplied [jq queries](https://stedolan.github.io/jq).

Assuming a JSON backend that responds with:

```
GET /path
Accept: application/json

[
  {
    "id": 1,
    "name": "alpha",
    "active": true
  },
  {
    "id": 2,
    "name": "beta",
    "active": false
  }
]
```

With jqrp in front, the backend response can be mutated by suppling a query in the `JQ` request header:

```
GET /path
Accept: application/json
JQ: .[] | select(.active) | {results: .name}

{ "results": ["alpha"] }
```

Consult the jq [manual](https://stedolan.github.io/jq/manual/#Basicfilters) for query syntax and available filters.

## Installation

[Pre-built binaries](https://github.com/bauerd/jqrp/releases) and [Docker images](https://hub.docker.com/r/bauerd/jqrp) are available.

## Usage

```
$ jqrp https://example.com
```

jqrp's default port is 9898.

### Docker Image

Alternatively, run a container:

```
$ docker run -p 9898:9898 bauerd/jqrp https://example.com
```

## Behaviour

* The [`X-Forwarded-For`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For) header gets set on every proxied request.

* The [`X-Request-ID`](https://stackoverflow.com/questions/25433258/what-is-the-x-request-id-http-header) header is propagated to the backend. If missing, requests are assigned UUIDs.

* jqrp attempts to mutate upstream responses only if all of the following conditions hold:

  1. Both `Accept: application/json` and `JQ: <QUERY>` headers are set on the request.
  2. The upstream response has a 2xx status code.
  3. The upstream response has the `Content-Type: application/json` header set.

* With the above conditions met, jqrp attempts to transform responses to __all request methods__.

* Otherwise, jqrp proxies transparently, and applies no transformation other than setting the `X-{Forwarded-For, Request-ID}` headers.

### Status Codes

The status code of jqrp indicates the operations performed and their outcomes:

| Status Code                           | Description                                                                                        |
|---------------------------------------|----------------------------------------------------------------------------------------------------|
| __203__ Non-Authoritative Information | The upstream response body was successfully transformed by the query.                              |
| __400__ Bad Request                   | The query provided in the `JQ` header is malformed.                                                |
| __408__ Request Timeout               | Applying the query to the upstream response exceeded the transformation timeout.                   |
| __422__ Unprocessable Entity          | The query evaluates to a primitive type that has no valid JSON representation.                     |
| __500__ Internal Server Error         | An unhandled error occured when proxying the upstream response.                                    |
| __502__ Bad Gateway                   | The upstream response body contains invalid JSON, or its `Content-Type` is not `application/json`. |
| __504__ Gateway Timeout               | The upstream host exceeded the proxy timeout.                                                      |
| Other                                 | The upstream response was not transformed, and its original status code preserved.                 |

## Configuration

jqrp can be configured via environment variables.

* All timeout values use milliseconds
* Setting a timeout to 0 disables it
* Setting the `CACHE_SIZE` to 0 disables query caching

| Environment Variable      | Default | Description                                                                                                                                           | Reference                                                                                           |
|---------------------------|---------|-------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| `PORT`                    | 9898    | Port to bind to                                                                                                                                       |                                                                                                     |
| `LOG_LEVEL`               | debug   | Log level. Either `debug`, `info` or `error`                                                                                                          |                                                                                                     |
| `CACHE_SIZE`              | 512     | Size of the LRU query cache. Setting the size to 0 disables query caching                                                                             |                                                                                                     |
| `EVAL_TIMEOUT`            | 0       | Maximum time spent evaluating jq queries                                                                                                              |                                                                                                     |
| `READ_TIMEOUT`            | 0       | Maxium time from when the client connection is accepted to when the request body is fully read                                                        | [Server.ReadTimeout](https://golang.org/pkg/net/http/#Server.ReadTimeout)                           |
| `WRITE_TIMEOUT`           | 0       | Maximum time from the end of the client request header read to the end of the response write                                                          | [Server.WriteTimeout](https://golang.org/pkg/net/http/#Server.WriteTimeout)                         |
| `DIAL_TIMEOUT`            | 0       | Maximum time spent establishing a backend TCP connection                                                                                              | [Dialer.Timeout](https://golang.org/pkg/net/#Dialer.Timeout)                                        |
| `DIAL_KEEPALIVE`          | 0       | Interval between keep-alive probes for an active backend network connection                                                                           | [Dialer.KeepAlive](https://golang.org/pkg/net/#Dialer.KeepAlive)                                    |
| `TLS_HANDSHAKE_TIMEOUT`   | 0       | Maximum time spent performing backend TLS handshake                                                                                                   | [Transport.TLSHandshakeTimeout](https://golang.org/pkg/net/http/#Transport.TLSHandshakeTimeout)     |
| `RESPONSE_HEADER_TIMEOUT` | 0       | Maxium time spent reading the headers of the backend response                                                                                         | [Transport.ResponseHeaderTimeout](https://golang.org/pkg/net/http/#Transport.ResponseHeaderTimeout) |
| `EXPECT_CONTINUE_TIMEOUT` | 0       | Maximum time to wait between sending the backend request headers when including an `Expect: 100-continue` and receiving the go-ahead to send the body | [Transport.ExpectContinueTimeout](https://golang.org/pkg/net/http/#Transport.ExpectContinueTimeout) |

## Security Considerations

* jqrp uses [gojq](https://github.com/itchyny/gojq), a re-implementation of jq. Because jqrp feeds user input untouched to gojq, its security properties depend mainly on gojq.

* If gojq panics on query evaluation, the jqrp process exits.

* Queries may be prohibitively expensive to evaluate. A malicious user may intentionally craft queries that take a long time to evaluate. Therefore jqrp affords setting an evaluation timeout that defaults to 10ms, configurable with the `EVAL_TIMEOUT` environment variable. Requests with queries exceeding the evaluation timeout get closed with status code 408.

* The [original jq](http://stedolan.github.io/jq) query language is Turing-complete, i.e. evaluation of user-supplied queries may loop indefinitely. jqrp uses gojq with all its [compiler options](https://github.com/itchyny/gojq#usage-as-a-library) disabled. Query evaluation in jqrp is likely not Turing-complete.

* Consider stripping the `JQ` header for unauthenticated/unauthorized requests in front of jqrp.

## Performance Considerations

* jqrp stores compiled queries in an LRU (last-recently-used) cache. Compiled queries retrieved from the cache can be applied immediately to upstream response bodies. The cache has a static size which can be configured with the environment variable `CACHE_SIZE`.

* Consider running jqrp behind a caching reverse proxy, that factors in the `JQ` header when computing cache keys. Note that if clients supply dynamically generated queries, this strategy is not viable.

## Edge Cases/Noteworthy

* If a query results in a single primitive result (i.e. a boolean, number, string or null), the response body is empty and the status code 422. If a query results in multiple primitive results, they are contained in an array.

* If a query's result set is empty, the status code is 203, and the body depends on the backend's JSON response:
  * If the top-level type was an object, the response body is the empty object `{}`.
  * If the top-level type was an array, the response body is the empty array `[]`.

* jqrp does not support HTTP content negotiation and only attempts to transform requests that solely `Accept: application/json`.

* jqrp logs only requests applicable to transformation. Requests proxied transparently are not logged.
