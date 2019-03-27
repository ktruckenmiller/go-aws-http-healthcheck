# go-aws-http-healtcheck
A small go app for health checking endpoints


### Rationale

I think what I want to do here is provide an endpoint, and post some metrics to AWS based on that endpoint. Golang has some interesting tracing elements so I'll see waht I can do.

```
# reports metrics to cloudwatch  
docker run -it --rm -e URL="http://my-ip.clustermaestro.com" go-http-healtcheck

```

### Metrics

Created in the AWS AppHealth namespace

| Metric Name | Values | Type |
| ------------| ------ | ---- |
| `isup` | `1` or `0`  | int  |     
| `dns`| `0 - 5000` | milliseconds |
|  | 



  - `AppHealth`
    - `isup: 1 or 0`
    - `dns: 0ms - 5000ms`
    -
