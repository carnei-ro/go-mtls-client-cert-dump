# Client Certificate Dumper

Simple app to help debug mTLS connection printing out the client certificate

## Usage

```
.
├── main.go
├── root-ca.crt
├── server.crt
└── server.key
```

```bash
go run main.go
```

```bash
$ curl http://localhost:8081/health
:8443/client-cert to dump client certificate
```

```bash
$ curl --cacert root-ca.crt --cert bundle-leandro.carneiro.crt --key leandro.carneiro.key https://localhost:8443/client-cert

Length of PeerCertificates: 3
----------------------------------------------------------------

-----BEGIN CERTIFICATE-----
MIIEmDCCA4CgAwIBAgIUKd5SdXBUSsXGG1xBvnd9ENe0VSUwDQYJKoZIhvcNAQEL
...
```

## Istio mTLS example

```bash
kubectl apply -f k8s-manifests.yaml
kubectl exec -it $(kubectl get pod -l app=debug -o jsonpath='{.items[0].metadata.name}') -c debug -- curl http://mtls.carnei.ro/client-cert
```

The server side will be running with let's encrypt certificate (so, no need to pass `caCertificates` at ServiceEntry), and root_ca configured will be able to validate the istio-proxy certificate.

K8s Manifests will:

- Mount istio-proxy client cert at /var/run/secrets/istio/certs using the annotations
- Create a ServiceEntry pointing to mtls.carnei.ro
- Create a VirtualService to map the port 18443 to 80 to be easy to call internally
- Create a DestinationRule to indicate the connection to mtls.carnei.ro must be mTLS and pass the istio-proxy client certificates
- Make the call as a simple "http" from inside the container app