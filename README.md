# traefik-hydrate-headers
Hydrate headers from remote endpoint

## Configuration

`remote.url` *(string)* - Url to be fetched

`remote.method` *(string)* [optional, default: `GET`] - Request method

`fetchOn.cookies` *([]string)* [optional] - Only fetch remote when any of cookies from this list present in original request

`appendOn.statusCodes` *([]int)* [optional] - Remote status codes at which headers will be appended

`nextOn.statusCodes` *([]int)* [optional] - Remote status codes at which next middleware will be called

`forwardHeaders` *([]string)* [optional] - List of headers from original request to be passed to remote

`headers` *(map[string]string)* - Map of headers in which to set remote response (*key* is header name, *value* is Go template)

## Headers template examples

### Set remote body to request header

```yaml
X-Example: '{{ .RemoteBody }}'
```

### Unmarshal JSON body and set field to header

```yaml
X-User-Id: |
  {{ $user := unmarshalJson .RemoteBody }}
  {{ $user.id }}
```

### Set remote header value to request header

```yaml
X-Example: '{{ .RemoteResponse.Header.Get "X-Remote-Header" }}'
```

### Set original request header value to another request header

```yaml
X-Real-Ip: '{{ .Request.Header.Get "Cf-Connecting-Ip" }}'
```

