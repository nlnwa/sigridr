# Sigríðr
Sigríðr is a Twitter API client

### Build
```
go build -o sigridr
```

### Usage

First time with consumer key and consumer secret (after which the access token is stored in config file):
```
./sigridr -k [consumer key] -s [consumer secret] search from:nasjonalbibl
```

Access token provided as environment variable:
```
ACCESS_TOKEN=[access token] ./sigridr search from:nasjonalbibl
```

With filters (no replies and no retweets):
```
./sigridr search from:nasjonalbibl -- -filter:replies -filter:retweets
```
