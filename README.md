# Kustomize Concepts
- [Kustomize Tutorial](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/)
- [The Kustomization File](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/)
- [Transformer Configuration](https://github.com/kubernetes-sigs/kustomize/blob/master/examples/transformerconfigs/README.md)
- [Containerized KRM Functions](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/containerized_krm_functions/)
- Containerized KRM Functions framework:
  * https://pkg.go.dev/sigs.k8s.io/kustomize/kyaml/yaml
  * https://pkg.go.dev/sigs.k8s.io/kustomize/kyaml/fn/framework


# Kubernetes API Validation

An example of how to reuse the Kpt kubeval function.

```yaml
transformers:
  - |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: validate-k8s-objects
      annotations:
        config.kubernetes.io/function: |
          container:
            image: gcr.io/kpt-fn/kubeval:v0.2.0
    data:
      strict: true
```


# StatefulSet Headless Service Lookup

All consumers should be able to construct stateful endpoints/URIs from a headless service DNS record.

DNS records example:
```bash
~ $ nslookup kuard
Server:    172.30.0.10
Address 1: 172.30.0.10 dns-default.openshift-dns.svc.cluster.local

Name:      kuard
Address 1: 10.129.7.107 kuard-0.kuard.dizaharov-dev.svc.cluster.local
Address 2: 10.128.7.90 kuard-1.kuard.dizaharov-dev.svc.cluster.local
Address 3: 10.128.4.90 kuard-2.kuard.dizaharov-dev.svc.cluster.local
```


# Java App Generator

A custom generator configuration example to be used in a KRM function.

```yaml
apiVersion: java.example.com/v1alpha1
kind: JavaAppGenerator
metadata:
  # a common base name for a Deployment/StatefulSet, a Service and a client ConfigMap
  name: koffee
spec:
  # Deployment or StatefulSet
  stateful: true
  image: zadenis/koffee
  ports:
    - https=8443
    - https-jmx=12000
  replicas: 1
  jvm:
    # -Xms1024m -Xmx1024m
    heap: 1024
    # heap / container memory limit ratio
    ratio: 0.7
    extra: -Dmy.sys.prop=my-value
  # an app container env spec
  env:
    - SERVER_PORT=8080
    - ANOTHER=value
  # configMapRefs to be used in a container envFrom spec
  # the referenced ConfigMaps must be generated in overlay environments
  upstreams:
    - database
    - ext-service
```

An optional configMapGenerator that produces an API properties exposed
by this process for consumers.

```yaml
configMapGenerator:
  # the name must be equal to the JavaApp name
  - name: koffee
    literals:
      # convention: <app-name | upper>_SVC_NAME=<app-name>
      - KOFFEE_SVC_NAME=koffee
      - KOFFEE_PROTO=https
      - KOFFEE_PORT=8443
```

To enforce the service name update in overlays add a name referece configuraion:

```yaml
configuration:
  - svc-name-ref.yaml

# svc-name-ref.yaml content
nameReference:
  - kind: Service # get a metadata.name from Service objects
    fieldSpecs:
      - kind: ConfigMap
        path: data/KOFFEE_SVC_NAME
```


# Execution

Embedded (non-standalone) mode as run by kustomize.

```sh
cat resource-list-example.yaml | go run main.go
```

Standalone mode for local test runs.
```sh
go run main.go debug FunctionConfig.yaml rl-item1.yaml rl-item2.yaml
```