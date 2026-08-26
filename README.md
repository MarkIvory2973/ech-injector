# ECH Injector

Inject ECH configuration while querying HTTPS RR.

## Installation

### Deploy from source

#### Requirements

- Go 1.26+
- Node v24.19.0
- Git

Clone the repository:

```bash
git clone https://github.com/MarkIvory2973/ech-injector.git
cd ech-injector
```

Install dependencies:

```bash
npm install
```

Create KV namspace:

```bash
wrangler kv namespace create ech-injector
```

Deploy worker:

```bash
KV_NAMESPACE_ID=<your-kv-namespace-id> npm run deploy
```

Set build variables (Optional):

|Name|Description|
|:-:|:-:|
|KV_NAMESPACE_ID|Cloudflare KV namespace ID used to store cached data for the Worker.|

## Usage

Set the DoH provider to: 

```
https://example.com/dns-query
```

## References

1. [Total ECH](https://github.com/RememberOurPromise/Total-ECH)