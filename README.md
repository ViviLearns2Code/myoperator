# About

This repo was created using kubebuilder v3.0.0-beta.0 using the tutorial: https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html

It was tested on minikube v1.18.1 on Ubuntu 20.04
```bash
$ minikube start --cpus 2 --memory 3500 --disk-size=10g --driver=docker --insecure-registry=192.168.49.1:5000
```
The cluster uses Kubernetes v1.20.2 on Docker 20.10.3.

# Learning notes
## Certs
* certs are necessary so that all services in a k8s cluster can authenticate each other, e.g. the k8s apiserver needs to trust the webhook server so it can invoke it. [(More details here)](https://www.youtube.com/watch/gXz4cq3PKdg)
* self-signed certificate: certicicates not signed by any CA but by the server itself ([_"A self-signed certificate is like making a gold-colored badge looking thing in your home and then going around showing it to people saying you're a police officer."_](https://security.stackexchange.com/questions/112768/why-are-self-signed-certificates-not-trusted-and-is-there-a-way-to-make-them-tru#comment202398_112768))
    * special case of in-house PKI: self-sign your own internal CA and add its public key servers and clients within your org so that they can authenticate each other based on this trusted CA
* kubebuilder tutorial uses `cert-manager`'s self-signed issuer: [_"The SelfSigned issuer doesn’t represent a certificate authority as such, but instead denotes that certificates will be signed through “self signing” using the given private key. [...] Clients consuming these certificates have no way to trust this certificate since there is no CA signer apart from itself, and as such, would be forced to trust the certificate as is."_ ](https://cert-manager.io/docs/configuration/selfsigned/)

### Sources
* https://speakerdeck.com/govargo/inside-of-kubernetes-controller
* https://engineering.bitnami.com/articles/a-deep-dive-into-kubernetes-controllers.html