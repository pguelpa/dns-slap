# dns-slap

A tool to slap around a DNS server.

## Installation

Simple as it takes to type the following command:

```go
% go get github.com/pguelpa/dns-slap
```

## Usage

dns-slap supports setting a concurrency level and the number of iterations to lookup a DNS entry per concurrent processes.

A per-lookup threshold (in milliseconds) can also be configured. If a single lookup takes longer than the configurable threshold (even if successful) it will be considered a failure for reporting purposes.

```bash
Usage of dns-slap
  -concurrency=10: How many concurrent lookups to try
  -iterations=100: How many times to lookup in each concurrent process
  -threshold=500:  How long to wait (in milliseconds) on a single lookup before considering it a failure
```

This is what happens when you run dns-slap

```
% dns-slap -concurrency 100 -iterations 1000 google.com
Starting 100 workers with 1000 lookups each ...
Workers finished, calculating results

Ran 100000 lookups in an average time of 0.017077 seconds
Found 0 errors
```
