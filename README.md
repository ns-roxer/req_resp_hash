## Simple console utility for calculating md5 hash sum of url content

* Utility makes http GET request to provided url(s), calculate md5 has of each response and print "url with hash" to stdout.
* Requests may be preformed concurrently. In that case the order of printed urls with hashes is not determined. 
* Url may be
provided with or without scheme. If the scheme is not provided, "https" is used by default.

### Usage:

#### clone and build:

```bash
$> git clone github.com/ns-roxer/req_resp_hash && \
mv req_resp_hash && \
go build -race .
```

#### run:

```bash
$> ./req_resp_hash [ -parallel ] urls ...
```

* argument `parallel` is limit of parallel requests. Must be above zero. Default value: 10
* `urls` - is list of urls or domains

#### examples:

```bash
$> ./req_resp_hash http://google.com
                                                                                               
http://google.com a5236353a96354fef362cdcb89c1ff52
```

```bash
$> go build -race . && ./req_resp_hash -parallel=5 https://google.com yahoo.com yandex.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
https://google.com ee6ed2efb466a23163daf39be07b37de
reddit.com/r/funny 2c2b4ac4a4ff6fae08cc257566786194
reddit.com/r/notfunny 56528fb545d9f20829a3604ddd77ba4c
yandex.com 213026562149f289d356157f8568a941
baroquemusiclibrary.com 830b4ae0c4d5b9463f64398cc8c35dc2
yahoo.com a4619e77423b7299fb0cb2c6c6c0032c

```
