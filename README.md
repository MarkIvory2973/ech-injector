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

Build binaries:

```bash
npm run build
```

Deploy worker:

```bash
npm run deploy
```

## Usage

Set the DoH provider to: 

```
https://example.com/dns-query
```

## References

1. [Total ECH](https://github.com/RememberOurPromise/Total-ECH)