# traefik-hydrate-headers
Hydrate headers from remote endpoint

## Configuration

| Key | Type | Required | Default | Description |
| :-- | :-- | :-: | :-- | :-- |
| `appendOn.statusCodes` | []int | | | Remote status codes at which headers will be appended |
| `fetchOn.cookies` | []string | | | Only fetch remote when any of cookies from this list present in original request |
| `forwardHeaders` | []string | | | List of headers from original request to be passed to remote |
| `headers` | map[string]string | :white_check_mark: | | Map of headers in which to set remote response (*key* is header name, *value* is Go template) |
| `nextOn.statusCodes` | []int | | | Remote status codes at which next middleware will be called |
| `remote.method` | string | | `GET` | Request method |
| `remote.url` | string | :white_check_mark: | | Url to be fetched |

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
