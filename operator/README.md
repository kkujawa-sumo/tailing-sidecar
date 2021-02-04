# Tailing sidecar operator

## Test Tailing sidecar operator

*Notice*: Commands below will use the current kubeconfig context

Install the CRDs into the cluster:

```bash
make install
```

Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```bash
make run
```

Install Instances of Custom Resources:

```bash
kubectl apply -f config/samples/
```

Run It On the Cluster

Build and push your image to the location specified by IMG:

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
```

Deploy the controller to the cluster with image specified by IMG:

```bash
make deploy IMG=<some-registry>/<project-name>:tag
```

To delete your CRDs from the cluster:

```bash
make uninstall
```
