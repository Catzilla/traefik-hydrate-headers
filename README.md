# traefik-hydrate-headers
Hydrate headers from remote endpoint

## Configuration

`remote.url` *(string)* - Url to be fetched

`remote.method` *(string)* - Request method

`appendOn.statusCodes` *([]int)* - Remote status codes at which headers will be appended

`nextOn.statusCodes` *([]int)* - Remote status codes at which next middleware will be called

`forwardHeaders` *([]string)* - List of headers from original request to be passed to remote

`headers` *(map[string]string)* - Map of headers in which to set remote response (*key* is header name, *value* does not used yet and is reserved for next releases)
