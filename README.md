# Crawler ESun JPY and Send Mail with GAE

## Setup
- setup your email on app.yaml
- create GAE project
- install gcloud SDK

## Deploy
```
$ gcloud app deploy app.yaml cron.yaml queue.yaml
```

## Setup the Default Expected
```
$ curl -X POST -H 'Content-Type: application/json' -d '{"expected": 0.28}' "https://<YOUR_GAE_PROJECT>.appspot.com/send"
```
