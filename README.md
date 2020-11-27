# pod-mutating-webhook
Are you struggling to migrate your docker images from dockerhub to your local repository, because of dockerhub rate limits?
But there are too many images to migrate, and you need something quick to swap to an internal/proxy registry? Well, you have come to right place.
I was also struggling to migrate dockerhub images to an internal registry because of dockerhub rate limiting.
This webhook makes use of kubernetes admission controller and mutates pod images based on rules.
You can specify which image to change and to what using the chart values file, these are called imageRules.
imageRules are `key:val` pairs that defines which image (key) should be replaced with other image (value).

##### Note: For mutation to work for a pod in a namespace, that namespace must be labeled with `pod-mutating-webhook: enabled` 

Example `imageRules` in chart values file:
```yaml
imageRules: |
    # swap all new pods with mysql:5.7 image to internal.proxy.registry/mysql:5.7
    "mysql:5.7": "internal.proxy.registry/mysql:5.7" 
    # swap all new pods with agill17/test:latest image to a specific tag agill17/test:0.1.0
    "agill17/test:latest": "agill17/test:0.1.0"
```


## Deploy
```shell script
$ make install
```

## Uninstall
```shell script
$ make uninstall
```

## Build
```shell script
$ make build
```
