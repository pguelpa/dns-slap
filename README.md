# dns-slap

A tool to slap around a DNS server.

## Installation

Simple as it takes to type the following command:

```go
% go get github.com/pguelpa/dns-slap
```

## Usage

dns-slap supports setting a concurrency level and the number of iterations to lookup a DNS entry per concurrent processes.

```bash
Usage of dns-slap
  -concurrency=10: How many concurrent lookups to try
  -iterations=100: How many times to lookup in each concurrent process
```

This is what happens when you run dns-slap

```
% dns-slap -concurrency 100 -iterations 1000 google.com
Starting 100 workers with 1000 lookups each ...
Workers finished, calculating results

Ran 100000 lookups in an average time of 0.017077 seconds
Found 0 errors
```
