# $3 per month Common Crawl Web Graph in AWS

[![Coverage Status](https://coveralls.io/repos/github/dharnitski/cc-hosts/badge.svg)](https://coveralls.io/github/dharnitski/cc-hosts)

try it - https://api.cc.dharnitski.com/domain/badssl.com


## Why

Common Crawl regularly releases host- and domain-level graphs, for visualising the crawl data.




```mermaid
graph LR;
    Gateway>API Gateway] --> Lambda;
    Lambda -->S3[(AWS S3)];
```
