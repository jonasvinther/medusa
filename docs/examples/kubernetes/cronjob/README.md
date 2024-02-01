# Kubernetes cronjob example
In this example we will deploy `Medusa` in Kubernetes via a `cronjob`. The example is to show how you could create a backup of Vault at specific intervals in Kubernetes.

## Creating a Configmap
First off we will create a Configmap to define three values used by the Cronjob, namely
- vault-url
- vault-path
- vault_skip_verify

These three values will be used as environment variables in the container run by the Job.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: vault
data:
  vault-url: "https://192.168.86.41:8201"
  vault-path: "secret/A"
  vault_skip_verify: "true"
```

Deploy the Configmap with:
```
kubectl apply -f configmap.yaml
```

## Creating a Secret
We need a token to access the target Vault instance. In this case we will store that in a Kubernets Secret like this:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: vault
type: Opaque
stringData:
  token: 00000000-0000-0000-0000-000000000000
```
The data in the key `token` will be used as an environment variable in the container run by the Job. 

Deploy the Secret like this:
```
kubectl apply -f secret.yaml
```

## Creating the Cronjob
> Note: In this example we are going to use a hostPath as a Volume. This is not recommended and only serves to show that it works. For the hostPath to work, a folder called "/backup" needs to exist on the workernodes, and own by uid 1000.

### Deploying the Cronjob
Here we create a Cronjob manifest that defines Jobs to be run every hour. It runs `Medusa` and is configured by the environment variables it gets from our Configmap and Secret. It finds data in Vault at the given path from the Configmap as well, and saves them to the hostPath volume at `/backup`.
```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: medusa
spec:
  schedule: "* * * * *"
  concurrencyPolicy: Forbid
  startingDeadlineSeconds: 600
  jobTemplate:
    spec:
      template:
        spec:
          securityContext:
            runAsUser: 1000
          containers:
          - name: medusa
            image: ghcr.io/jonasvinther/medusa:latest
            imagePullPolicy: IfNotPresent
            command: ["./medusa", "export", "$(VAULT_PATH)", "-o", "/backup/backup.vault"]
            env:
            - name: VAULT_SKIP_VERIFY
              valueFrom:
                configMapKeyRef:
                  name: vault
                  key: vault_skip_verify
            - name: VAULT_ADDR
              valueFrom:
                configMapKeyRef:
                  name: vault
                  key: vault-url
            - name: VAULT_PATH
              valueFrom:
                configMapKeyRef:
                  name: vault
                  key: vault-path
            - name: VAULT_TOKEN
              valueFrom:
                secretKeyRef:
                  name: vault
                  key: token
            volumeMounts:
              - name: backup
                mountPath: /backup
          restartPolicy: OnFailure
          volumes:
          - name: backup
            hostPath:
              path: /backup
              type: DirectoryOrCreate
```
Deploy the Configmap like this :
```
kubectl apply -f configmap.yaml
```
A Cronjob is created and can be verified with
```
kubectl get cronjobs
NAME     SCHEDULE    SUSPEND   ACTIVE   LAST SCHEDULE   AGE
medusa   * * * * *   False     0        32s             4h19m
```

Now once every minute `Medusa` will run via a job.
See the jobs like this :
```
kubectl get jobs
NAME                COMPLETIONS   DURATION   AGE
medusa-1615981980   1/1           2s         2m42s
medusa-1615982040   1/1           3s         102s
medusa-1615982100   1/1           2s         42s
```

Every job spawns a Pod, see those here :
```
kubectl get pods
NAME                      READY   STATUS      RESTARTS   AGE
medusa-1615981980-c9t5f   0/1     Completed   0          3m
medusa-1615982040-bklxx   0/1     Completed   0          2m
medusa-1615982100-5tfl8   0/1     Completed   0          60s
medusa-1615982160-4b527   0/1     Completed   0          9s
```

### Using Kubernetes authentication
If you are using the kubernetes authentication method in Vault, it is also possible to use the kubernetes provided JWT token inside a Pod and auth role in order to authenticate.
If your authentication mount point is different from the default of `kubernetes`, for example if your vault instance is supporting multiple clusters, this can be changed with the
`--kubernetes-auth-path` option.

```yaml
command: ["./medusa", "export", "$(VAULT_PATH)", "--kubernetes", "--role=default", "-o", "/backup/backup.vault"]
```

### Further customization
This only serves as an example as to how you could use `Medusa` to take a backup of Vault from a given location. 

The first this you should change is the Volume. We recommend using a storageClass that makes sense in your environment.

### Cleanup
To cleanup, run this
```
kubectl delete -f cronjob.yaml
kubectl delete -f secret.yaml
kubectl delete -f configmap.yaml
```

