# Email Hiding Service

Backend service that returns an email (or any other payload) upon succesful completion of a ReCaptcha v2.

Return values are stored as env vars on the lambda:
```
site = return
```

The ReCaptcha secret key should be stored in an env var `RECAPTCHA_SECRET`

Environment variables are set in a `secrets.tfvars` file (not in git). 

## Deploying to aws
```
./deploy.sh
```

## Calling
```
curl --request POST \
  --url <API GW URL> \
  --header 'Content-Type: application/json' \
  --data '{"token":"<recaptcha_token>","site":"<site_name>"}'