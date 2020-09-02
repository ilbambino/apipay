你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# APIPAY

[![Build Status](https://cloud.drone.io/api/badges/ilbambino/apipay/status.svg)](https://cloud.drone.io/ilbambino/apipay)

This is a small project, creating a webservice for a payment service. It only implements basic CRUD operations. It is not fully finished and it can be used only as a starting point of what a production service would be.

## Running it

It needs to have Go with Modules support (it has been tested only with 1.12+). Taking into account that you have Go installed, all you need to do is:

```
make
```

and you will get the binary `apipay` that you can run.

You need to have Mongo 3.6 running also. If you don't have one, you can easily have one with Docker. You can run `docker-compose -f needed-services.yml -d up` to run one locally.

In case you want to run `apipay` against a different Mongo (not running in localhost), you can change it using environment variables:

- `APIPAY_MONGOHOST` for the mongo host name or IP
- `APIPAY_MONGOPORT` for the port of the mongo server.
- `APIPAY_MONGOUSER` to provide a username.
- `APIPAY_MONGOPASSWORD` to provide a password for Mongo.

## Tests

The tests run against a real mongo, not mocked version. You need to have mongo running. They will use different DBs for the tests, so they won't mess up your data. To run all of them:

```
make tests
```

Tests are also run in CI (using [drone](https://drone.io)) with real Mongo.

## Documentation

With the project there is also a proof of concept to create the [API documentation](api.pdf) automatically from the code. It uses `swagger` but it would require more effort to really work well.

To do so, if you have installed [swag](https://github.com/swaggo/swag) you can run `swag init` and it will generate the Swagger _json_ definitions.

Once you have the _json_ you can use [swagger2pdf](https://www.npmjs.com/package/swagger-spec-to-pdf) to generate a _PDF_. But there are some errors and formatting issues that would need to be fixed. `swagger2pdf -j -s docs/swagger.json -o .`

## TODOs

This project it is just a starting point, and it could have lots of improvements. Not done because the lack of time.

- **Improve the documentation generation**. So it is really useful and looks decent. Also automating the generation of it in the `Makefile`
- **Improve Logging**. Add a request ID when a request is done to be able to trace requests. `Gin` also should use the same logger (and json formatting).
- **Tracing** With [Jaeger](https://www.jaegertracing.io/) it should be easy to do.
- **Health** It would be nice if the service would expose some health APIs so external parties can know if the service is working properly, eg. `/health/status` API.
- **TLS** Depending on how this would be deployed, it might need to do the TLS termination.
- **Authentication** Currently the API does not perform any Auth on the request.
- **Data Model improvements** There should be proper validation on the models. Also when serializing to Mongo and _json_ `omitempty` could be added if needed.
- **Pagination** The get list of payments is hardcoded to a max 100 elements. Pagination could be implemented.
- **Tests** Some basic tests have been added. But there should be more, testing the errors, etc.
- **CI** Project currently uses CI, but the binary generated is not saved anywhere. And easy one would be to build a Docker image and push it to Docker Hub. But it depends on how this would be run in practice.
- **Git** Instead of using _master_ use proper PRs and reviews.
