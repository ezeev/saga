
# Instructions

Create a `deploy/deployment.yml` with the following:

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: viroongacom
spec:
  replicas: 1
  template:
    metadata:
      labels:
        run: viroongacom
    spec:
      containers:
      - name: example
        image: ezeev/viroongacom:latest
        env:
        - name: DEV
          value: "false"
        - name: APP_NAME
          value: "Viroonga"
        - name: APP_DOMAIN
          value: "viroonga.com"
        - name: TEST_DOMAIN
          value: "viroonga.com"
        - name: METRICS_PROVIDER
          value: "prometheus"
        - name: METRICS_URI
          value: "/metrics"
        - name: TEMPLATE_GLOB_PATTERN
          value: "view/*.html"
        - name: STRIPE_CARD_REDIRECT
          value: "/account"
        - name: STRIPE_TEST_PK
          value: "YOUR STRIPE TEST PK"
        - name: STRIPE_TEST_SK
          value: "YOUR STRIPE TEST SK"
        - name: STRIPE_LIVE_PK
          value: "YOUR STRIPE PK"
        - name: STRIPE_LIVE_SK
          value: "YOUR STRIPE LIVE SK"
        - name: STRIPE_APP_FILTER
          value: ""
        - name: CLOUDSQL_DB_NAME
          value: "viroonga"
        - name: CLOUDSQL_USER
          value: "root"
        - name: CLOUDSQL_PASSWORD
          value: ""
        - name: CLOUDSQL_CONN_STR
          value: "root@tcp(127.0.0.1:3306)/viroonga"
        - name: AUTH0_DOMAIN
          value: "YOUR_DOMAIN.auth0.com"
        - name: AUTH0_CLIENT_ID
          value: "YOUR CLIENT ID"
        - name: AUTH0_CLIENT_SECRET
          value: "YOUR SECRET ID"
        - name: AUTH0_CALLBACK_URI
          value: "/callback"
        - name: AUTH0_SIGNOUT_URI
          value: "/signout"
        - name: AUTH0_SIGNOUT_REDIRECT_URI
          value: "/"
        - name: AUTH0_CALLBACK_HOST_DEV
          value: "https://viroonga.com"
        - name: AUTH0_CALLBACK_HOST_LIVE
          value: "https://viroonga.com"
        - name: OAUTH_SUCCESS_REDIRECT
          value: "/account"
        - name: API_RATE_LIMIT_PER_MIN
          value: "10"
        - name: MAIL_GUN_API_KEY
          value: "YOUR MAILGUN KEY"
        - name: MAIL_DEFAULT_RECIPIENT_EMAIL
          value: "admin@viroonga.com"
        - name: MAIL_GUN_DOMAIN
          value: "mail.viroonga.com"
        - name: MAIL_GUN_PUB_KEY
          value: "YOUR MAIL GUN PUB KEY"
        ports:
        - containerPort: 3001
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "100m"
        # Change [INSTANCE_CONNECTION_NAME] here to include your GCP
        # project, the region of your Cloud SQL instance and the name
        # of your Cloud SQL instance. The format is
        # $PROJECT:$REGION:$INSTANCE
        # Insert the port number used by your database.
        # [START proxy_container]
      - image: gcr.io/cloudsql-docker/gce-proxy:1.09
        name: cloudsql-proxy
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "100m"
        command: ["/cloud_sql_proxy", "--dir=/cloudsql",
                    "-instances=cloud-ninjaio:us-central1:cloud-ninja=tcp:3306",
                    "-credential_file=/secrets/cloudsql/credentials.json"]
        volumeMounts:
            - name: cloudsql-instance-credentials
              mountPath: /secrets/cloudsql
              readOnly: true
            - name: ssl-certs
              mountPath: /etc/ssl/certs
            - name: cloudsql
              mountPath: /cloudsql
        # [END proxy_container]
      # [START volumes]
      volumes:
        - name: cloudsql-instance-credentials
          secret:
            secretName: cloudsql-instance-credentials
        - name: ssl-certs
          hostPath:
            path: /etc/ssl/certs
        - name: cloudsql
          emptyDir:
      # [END volumes]
```