# Saga - Stripe (S), Auth0 (A), Go (G), API Server (A)

Saga is a set of API and redirect endpoints designed for use with mobile and web applications.
The goal of Saga is to provide backend plumbing needed to implement common functions of customer-facing applications:

 - API Middleware including authentication and rate limiting
 - Caching
 - Session Management (via cookies and JWT)
 - Payment Processing (Stripe) including management of payment methods and subscriptions
 - Authentication (Auth0) and user management
 - Database Connection Management (Google Cloud SQL)
 - Sending Mail

Each of the packages in this repository is designed to work with other saga packages, although not all are dependent on other saga packages.

## Google App Engine

Saga is written in Go and currently only supports Google App Engine.


## Using Saga In Your Applications

## Configuration


```
env_variables:
  APP_DOMAIN: "yourlivedomain.com"
  STRIPE_CARD_REDIRECT: '/account'
  STRIPE_TEST_PK: 'YOUR STRIPE TEST PUBLIC KEY'
  STRIPE_TEST_SK: 'YOUR STRIPE TEST SECRET KEY'

  STRIPE_LIVE_PK: 'YOUR STRIPE PUBLIC KEY'
  STRIPE_LIVE_SK: 'YOUR STRIPE SECRET KEY'
  STRIPE_APP_FILTER: 'Optional filter using stipe metadata (tags)'

  CLOUDSQL_CONNECTION_NAME: your_gcp_project_id:region:instance_id
  CLOUDSQL_DB_NAME: 'your_cloudsql_db_name'
  # Replace username and password if you aren't using the root user.
  CLOUDSQL_USER: root
  CLOUDSQL_PASSWORD: ''
  CLOUDSQL_DEV_CONN_STR: 'root@tcp(127.0.0.1:3306)/your_cloudsql_db_name'

  AUTH0_DOMAIN: 'get from auth0'
  AUTH0_CLIENT_ID: 'get from auth0'
  AUTH0_CLIENT_SECRET: 'get from auth0'
  AUTH0_CALLBACK_URI: '/callback'
  AUTH0_SIGNOUT_URI: '/signout'
  AUTH0_CALLBACK_HOST_DEV: 'http://localhost:8080'
  AUTH0_CALLBACK_HOST_LIVE: 'https:/yourlivedomain.com'

  API_RATE_LIMIT_PER_MIN: '10'
```