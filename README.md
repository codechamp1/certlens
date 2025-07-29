# certlens
*certlens* is terminal UI for inspecting Kubernetes TLS Secrets Certificates.

## Features
- Inspect Kubernetes TLS Secrets interactively in the terminal
- View both raw/formatted PEM data with additional computed certificate details (expiry status, time until expiry, validity used, self-signed...)
- Navigate certificate chains in a single TLS secret
- Paginated and filterable secrets list for easy navigation
- Copy certificate or private key data to clipboard 

## Demo


## Installation
### Using `go install`
```bash
go install github.com/abapcp/certlens@latest
```
### Build from source
```bash
git clone https://github.com/codechamp1/certlens
cd certlens
make install
```
This will build the binary and place it in your `$GOPATH/bin` directory.

## Usage
```bash
certlens [flags]
```

## Flags
```bash
certlens --help
Usage of certlens:
  -context string
        context to use from kubeconfig, if not set, the current context will be used
  -kubeconfig string
        path to a kubeconfig (default "~/.kube/config")
        name of the secret to lens, if not set, all secrets will be listed
  -namespace string
        namespace to lens, if not set, all namespaces will be used
```

### Example
```bash
certlens -kubeconfig ~/.kube/config -namespace my-namespace
```


