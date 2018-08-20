# oauth2 proxy operator

This is a demo operator for the [bitly oauth2 proxy](https://github.com/bitly/oauth2_proxy).

It is not yet fully featured, but has just enough functionality to create oauth2 authentication back to github.

Once the Operator and CRD is installed you can request an oauth2 proxy with a manifest that looks like:

```yaml
apiVersion: "oauth2proxy.com/v1alpha1"
kind: "Proxy"
metadata:
  name: "example"
spec:
  name: "example"
  replicas: 1
  config:
    provider: "github"
    upstream: "http://example:80"
    emailDomain: "*"
    address: "0.0.0.0:4180"
    cookieSecure: "false"
    cookieSecret: "sdasdsadasdsadsa"
    cookieDomain: "example.35.xxx.131.181.xip.io"
    clientID: "XXXX"
    clientSecret: "XXXX"
```

## Example

* Deploy an example application:

```bash
kubectl run example --image=nginx:1.13.5-alpine
kubectl expose deployment example --port=80
```

* Deploy the oauth2 proxy operator:

```bash
git apply -f deploy/operator.yaml
git apply -f deploy/crd.yaml
git apply -f deploy/operator.yaml
```

* Create a new github authorization at [https://github.com/settings/developers](https://github.com/settings/developers)

* Use a URL from a domain that you control to configure it:
  * https://<domain>/oauth2/callback

* Record the resultant Client ID and Client Secret.

* edit `deploy/zzproxy.yaml` and replace the client id and client url.

* Deploy the oauth2proxy:

```bash
$ kubectl apply -f deploy/zzproxy.yaml
proxy "example" created
$ kubectl get pods
NAME                                    READY     STATUS    RESTARTS   AGE
example-6f59c6cd77-b2k27                1/1       Running   0          1h
oauth2-proxy-7cb45848-b6vnq             1/1       Running   0          5m
oauth2-proxy-example-5d67dd5848-4d8wf   1/1       Running   0          5m
```

* Create a Service of type LoadBalancer for the newly created oauth2-proxy:

```bash
kubectl expose deployment oauth2-proxy-example --port=80 --target-port=4180 --type=LoadBalancer
```

* After a few minutes it should be online and you can assign your DNS to the IP address of the service's external IP:

```
kubectl get svc oauth2-proxy-example
NAME                   TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
oauth2-proxy-example   LoadBalancer   10.100.200.140   35.224.131.181   80:32207/TCP   1h
```

* Once that is done you can point your web browser at the DNS and it should redirect you through the github oauth2 authorization and then back to your application.
