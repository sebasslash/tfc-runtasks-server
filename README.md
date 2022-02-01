# TFC Run Tasks Server

Run Tasks are a beta feature of Terraform Cloud (TFC) that allow you to perform remote operations during your Terraform Cloud run lifecycle.

This is a simple webserver that exposes a few endpoints that can be used to test different possible Run Task outcomes without integrating your own solution.

`/success` - mocks a successful run task

`/failed` - mocks a failed run task

Both endpoints allow a `timeout` query param (in seconds) that mocks some arbitrary workload the run task performs:

Will return a successful response back to TFC after 30 seconds
```
/success?timeout=30
```
## Get Started

### Using Docker

```sh
docker pull sebasriv/tfc-runtasks-server
docker run -d -p 80:10000 sebasriv/tfc-runtasks-server
```
