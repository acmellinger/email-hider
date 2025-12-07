# Email Hiding Service

Backend service that returns an email (or any other payload) upon succesful verification by Cloudflare Turnstile.

Return values and turnstile secret are stored as env vars on the worker:

```
site = return
```

## Deploying to gcp

```
./deploy.sh
```

## Calling

```
curl --request POST \
  --url <URL> \
  --header 'Content-Type: application/json' \
  --data '{"token":"<recaptcha_token>","site":"<site_name>"}'
```
