
# Viroonga Proxy

NGINX is used as the endpoint for all Viroonga traffic. This Deployment is coupled with an SSL-only Ingress for Kubernetes

NGINX is used as a reverse proxy into other services. See `nginx.conf`.

### Ingress

The Ingress controller is defined in `deploy/ingress.yml`. You may want to make changes to the ingress after it is
already running. Deleting and recreating it is problematic, so instead, you can update it on the fly:
```
kubectl edit ing viroonga-ingress
```

### TLS Certs

I used acme.sh to create a multi-domain SSL certificate with DNS validation. See https://github.com/Neilpang/acme.sh
"Use DNS mode" then follow the steps outlined in the readme and the prompt.  Be sure to use the genererated `fullchain.cer`
in order to get an A rating on SSL labs.

```
kubectl create secret tls viroonga-letsencrypt-secret --key acme.sh/viroonga.com.key --cert acme.sh/fullchain.cer
```