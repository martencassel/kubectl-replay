### What this is about
`kubectl-replay` is a [krew](https://krew.sigs.k8s.io/) plugin that helps you **see and replay what happened in your cluster**.
It takes Kubernetes **audit logs** and **events**, then translates them into reproducible `kubectl` commands with inline context (who did what, from where, and what the result was).
Instead of digging through JSON, you get clear commands you can run or study to debug issues.

---

### Features
- 🎯 **Audit Log Replay** - Convert API server audit logs into kubectl commands
- 📊 **Event Streaming** - Watch live Kubernetes events as kubectl commands
- 🎨 **Color-coded Output** - Green for success, magenta for 404s, yellow for 403s, red for errors
- ⏱️ **Timestamped Requests** - See exactly when each API call happened
- 🔍 **Low-level Details** - HTTP method, URI, status codes, and error messages
- 📝 **Script-friendly** - Output as shell comments, ready to copy/paste

---

### Setup
1. **Enable audit logging** in your cluster (file backend or webhook).
   Example kube‑apiserver flags:
   ```yaml
   audit-policy-file: /etc/kubernetes/audit-policy.yaml
   audit-log-path: /var/log/kubernetes/audit.log
   audit-log-format: json
   ```
2. **Collect events** with:
   ```bash
   kubectl get events -A -o json
   ```

---

### Install
Install via krew:
```bash
kubectl krew install replay
```

Or download the latest release from our repository:
```bash
# Download the latest dist file
wget https://github.com/martencassel/kubectl-replay/releases/download/latest/kubectl-replay_v0.1.0_linux_amd64.tar.gz
tar -xzf kubectl-replay_v0.1.0_linux_amd64.tar.gz
sudo mv kubectl-replay /usr/local/bin/
```

Or build from source:
```bash
git clone https://github.com/martencassel/kubectl-replay.git
cd kubectl-replay
make build
```

---

### Usage

#### Replay Audit Logs
Stream audit logs and see them as kubectl commands with colored status codes:

```bash
kubectl replay audit -f /var/log/kubernetes/audit.log --replay-speed 100
```

**Example Output:**
```bash
# [10:13:07] GET /api/v1/namespaces/default/configmaps/foobar → 404 configmaps "foobar" not found
kubectl get configmaps foobar -n default

# [10:13:08] POST /api/v1/namespaces/default/configmaps → 201 Created
kubectl create configmaps myconfig -n default

# [10:13:09] GET /api/v1/namespaces/kube-system/secrets/admin-token → 403 Forbidden
kubectl get secrets admin-token -n kube-system
```

**Color coding:**
- 🟢 **Green** (200-299): Successful requests
- 🟣 **Magenta** (404): Not found errors
- 🟡 **Yellow** (403): Permission denied
- 🔴 **Red** (400+): Other errors

#### Replay Events
Watch live Kubernetes events as kubectl commands:

```bash
kubectl replay events --replay-speed 10x
```

**Example Output:**
```bash
kubectl describe node kind-control-plane # reason=NodeReady message=kubelet is posting ready status
kubectl describe pod nginx-7d8b49c8d4-abc -n default # reason=Started message=Started container nginx
```

---

### 🧪 Example: Kind Cluster with Audit Logging

The `examples/` directory includes a complete kind setup with audit logging enabled:

1. **Create a kind cluster with audit logging:**
   ```bash
   cd examples
   kind create cluster --config kind-config.yml --name audit-demo
   ```

2. **Watch audit logs in real-time:**
   ```bash
   # Tail the audit log (mounted from kind container)
   sudo kubectl-replay audit -f ./logs/audit.log --replay-speed 100
   ```

3. **Generate some activity:**
   ```bash
   # Try to get a non-existent configmap (generates 404)
   kubectl get configmap foobar -n default

   # Create a configmap (generates 201)
   kubectl create configmap myconfig --from-literal=key=value -n default

   # Try to access forbidden resource (generates 403)
   kubectl get secrets -n kube-system --as=system:unauthenticated
   ```

4. **See the replayed commands with color-coded status:**
   ```bash
   # [10:15:23] GET /api/v1/namespaces/default/configmaps/foobar → 404 configmaps "foobar" not found
   kubectl get configmaps foobar -n default

   # [10:15:45] POST /api/v1/namespaces/default/configmaps → 201 Created
   kubectl create configmaps myconfig -n default

   # [10:16:02] GET /api/v1/namespaces/kube-system/secrets → 403 Forbidden
   kubectl get secrets -n kube-system
   ```

---

### Options

**Audit command:**
```bash
kubectl replay audit -f <audit-log-file> [--replay-speed <multiplier>]
```
- `-f, --file`: Path to audit log file
- `--replay-speed`: Speed multiplier (e.g., 100 for 100x faster, default: 1)

**Events command:**
```bash
kubectl replay events [--replay-speed <speed>] [--kubeconfig <path>]
```
- `--replay-speed`: Replay speed (e.g., "10x", default: "1x")
- `--kubeconfig`: Path to kubeconfig file (defaults to $KUBECONFIG or ~/.kube/config)
- `--from-event-log`: Replay from event log instead of live cluster

---

### Building a Release

Create a release tarball for distribution:
```bash
make release
```

This creates `dist/kubectl-replay_v0.1.0_linux_amd64.tar.gz` with the SHA256 checksum.

---

### Contributing

Contributions welcome! Please feel free to submit issues or pull requests.
