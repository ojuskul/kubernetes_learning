brew install helm

create service & deployment templates, add corresponding values in values.yaml
delete ingress and other templates not used

helm install okulkarni-node-app ./node-app-chart --namespace operator

```
$kubectl get pods -n operator
NAME                                                 READY   STATUS    RESTARTS   AGE
http-operator-controller-manager-7777bbbf48-5w2b4    1/1     Running   0          77m
node-app-74677f9cf6-qhxjz                            3/3     Running   0          5m45s
okulkarni-node-app-node-app-chart-788846cb6b-b8qq7   3/3     Running   0          73s

The first two are from previous experiment of creating CRD & Node App


helm uninstall okulkarni-node-app --namespace operator
$kubectl get pods -n operator
NAME                                                 READY   STATUS        RESTARTS   AGE
http-operator-controller-manager-7777bbbf48-5w2b4    1/1     Running       0          80m
node-app-74677f9cf6-qhxjz                            3/3     Running       0          8m59s
okulkarni-node-app-node-app-chart-788846cb6b-b8qq7   3/3     Terminating   0          4m27s

After updating replicaCount to 2
kubectl get pods -n operator
NAME                                                 READY   STATUS    RESTARTS   AGE
http-operator-controller-manager-7777bbbf48-5w2b4    1/1     Running   0          80m
node-app-74677f9cf6-qhxjz                            3/3     Running   0          9m38s
okulkarni-node-app-node-app-chart-788846cb6b-6kp69   3/3     Running   0          5s
okulkarni-node-app-node-app-chart-788846cb6b-bl7gw   3/3     Running   0          5s

$kubectl get deployment okulkarni-node-app-node-app-chart -n operator
NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
okulkarni-node-app-node-app-chart   2/2     2            2           49s
```