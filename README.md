# i-am-spartakus

![I am spartakus](media/i-am-spartacus.jpg)

This is 'copy' of the archived project  

https://github.com/kubernetes-retired/spartakus

The last release from that closed project was many years ago.

The main purpose of this project is to allow maintenance builds.

## Motivation

Fixing spartakus issues from out date client libraries.

If you are using the BigQuery datastore and experience issues such the following:
```
Error: failed to store record: googleapi:
```
Then try replacing the existing 'gcr.io/google_containers/spartakus-amd64:v1.1.0' image with
```
docker.io/seldonio/spartakus-amd64:v1.1.0-custom-0.3
```
