kubectl port-forward consul-server-0 --namespace $NAMESPACE_NAME 8500:8500
kubectl port-forward --namespace default svc/postgresql-dev 5432:5432