1. Install Docker Desktop
2. Install minikune `brew install minikube`
3. Install kubectl `brew install kubectl`, I didn't install it. Maybe it came with minikube.
4. Run minikube `minikube start --driver=docker --container-runtime=containerd --cpus=4 --memory=6g`
5. Run `kubectl get nodes` to check if the local k8s cluster is setup. Output should be as follows
```
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   46s   v1.32.0
```
6. Now you need to develop the node app.
```
brew install node
node -v //should return version number
npm -v  //should return version number
brew install nvm

cd operator
mkdir node-app && cd node-app
npm init -y
npm install express

create a dummy app
to test, run node index.js

verify if it works. (Hint: It should, this is a dummy app)

Now create the docker image and push to minikube

docker build -t node-app:latest .
minikube image load node-app:latest

check if the image is present by running
minikube image ls

create dedicated namespace
kubectl create namespace operator

create configmap for nginx => 
kubectl create configmap nginx-config --from-file=nginx.conf --namespace=operator

check
kubectl describe configmap --namespace=operator
you should see a ca-cert & nginx namespace

add k8s manifests, yaml files
then run

kubectl apply -f deployment.yaml
service/node-app created
deployment.apps/node-app created

I had error running the deployment, because k8s was trying to pull image from global registry
After setting imagePullPolicy: Never, the pod ran succesfully

To test, enable port forwarding
kubectl --namespace=operator port-forward service/node-app 8080:80

Now use curl/postman to test the app
Your app is running!
```

7. Now the operator part

```
brew install operator-sdk

mkdir http-operator && cd http-operator
operator-sdk init --domain=mydomain.com --repo=github.com/you/http-operator
operator-sdk create api --group=monitor --version=v1alpha1 --kind=HTTPMonitor --resource --controller

make code changes. Right now this just reads logs from nginx

after that push
make docker-build docker-push IMG=httpmonitor:latest
minikube image load httpmonitor:latest
make deploy IMG=httpmonitor:latest #this deploys the new CRD to minikube cluster

take2
make docker-build IMG=httpmonitor:latest
minikube image load httpmonitor:latest
make deploy IMG=httpmonitor:latest

take3
had to update the apis module update in controller due to pod error
go build ./...


now create httpmonitor.yaml in node-app dir
apply this crd
kubectl apply -f httpmonitor.yaml

now to test the logs
kubectl logs -n operator -l control-plane=controller-manager -f

other helpful commands:
minikube ssh
check if mnt is being correctlt written to

force http-monitor to pickup new config
kubectl delete httpmonitor monitor-sample -n operator 
kubectl delete pod <http-monitor-pod>
kubectl apply -f httpmonitor.yaml

```