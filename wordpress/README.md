# Wordpress Deployment Using Cloud SQL

Steps for using this deployment

helm install --name viroonga-wp-4 --set wordpressUsername=admin,wordpressPassword=frog7Tom! stable/wordpress
kubectl expose deployment viroonga-wp-4-wordpress --target-port=80 --type=NodePort --name=viroonga-wp-np

