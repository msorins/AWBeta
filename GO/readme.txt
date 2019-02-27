proxy: ssh -R 80:localhost:3000 serveo.net

Setup Google Cloud:

gcloud config set project awbeta-212719
gcloud container clusters create awbeta-kube --zone europe-west2-a
docker tag smircea/awbeta:1.0.0 eu.gcr.io/awbeta-212719/awbeta-go:1.0.0

gcloud docker push nameOfImage
kubectl create -f deployments/deployment.yaml

gcloud container clusters list


# Crete static global ip
gcloud compute addresses create awbeta-ip --global

# Starting sequence
kubectl apply -f ingress.yaml
kubectl create -f deployment.yaml
kubectl create -f service.yaml

# Getting SSL for domain awbeta.cf (TUTORIAL: https://estl.tech/configuring-https-to-a-web-service-on-google-kubernetes-engine-2d71849520d )
https://cloud.google.com/load-balancing/docs/ssl-certificates
certbot -d awbeta.cf --manual --logs-dir certbot --config-dir certbot --work-dir certbot --preferred-challenges dns certonly
