# validation
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