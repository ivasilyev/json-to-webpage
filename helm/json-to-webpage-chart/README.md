# `json-to-webpage-chart`

Helm chart to deploy `json-to-webpage` web server

## Setup 

```shell script
cd "${HELM_HOME}"
export CHART_NAME="json-to-webpage"
git clone "https://github.com/ivasilyev/json-to-webpage.git"
cd json-to-webpage/helm/json-to-webpage-chart

helm list --all
helm delete "${CHART_NAME}"

git reset --hard
git pull
rm -f *.tgz
export CHART="$(helm package "${CHART_NAME}" | awk '{print $NF}')"
echo "Install '${CHART}'"
helm install "${CHART_NAME}" "${CHART}"

kubectl get pods
kubectl get services

echo "Access IP address: $(
    kubectl describe pod "$(
        kubectl get pods \
        | grep --only-matching "${CHART}[^ \t]*"
    )" \
    | grep 'Node:' \
    | grep --only-matching '[0-9]*\.[0-9]*\.[0-9]*\.[0-9]*'
)"
```
