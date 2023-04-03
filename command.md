## kubernetes

### instances

```bash
curl localhost:11111/api/kubernetes/instances
```

### deployments

```bash
curl localhost:11111/api/kubernetes/instances/default/namespaces/default/deployments
```

### statefulsets

```bash
curl localhost:11111/api/kubernetes/instances/default/namespaces/default/stateful-sets
```

## kafka

### topics
```
curl localhost:11111/api/kafka/instances/default/topics
```

### consumer groups
```
curl localhost:11111/api/kafka/instances/default/consumer-groups
```
