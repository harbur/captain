# Introduction

Captain, the CLI build tool for Docker made for Continuous Integration / Continuous Delivery

Ditch your Makefiles and use a yaml format to describe your docker build process.

## Documentation

Captain will automatically configure itself with sane values without the need for any pre-configuration, so that it will work in most cases. When it doesn't, the `captain.yml` file can be used to configure it properly. This is a simple YAML file placed on the root directory of your git repository. Captain will look for it and use it to be configured.

Here is a full `captain.yml` example:

```yaml
build:
  images:
    - Dockerfile=harbur/hell-world
    - Dockerfile.test=harbur/hello-world-test
test:
  unit:
    - docker run -e NODE_ENV=TEST harbur/hello-world-test node mochaTest
    - docker run -e NODE_ENV=TEST harbur/hello-world-test node karmaTest
```

## build section

**images**: A list of the Dockerfiles to be compiled accompanied by their docker image namespace.

e.g.

```yaml
build:
  images:
    - Dockerfile=harbur/hello-world
    - Dockerfile.dev=harbur/hello-world-dev
    - test/Dockerfile=harbur/hello-world-test
```

## test section

**unit**: List of commands that run unit tests

e.g.

```yaml
test:
  unit:
    - docker run -e NODE_ENV=TEST harbur/hello-world-test node mochaTest
    - docker run -e NODE_ENV=TEST harbur/hello-world-test node karmaTest
```
