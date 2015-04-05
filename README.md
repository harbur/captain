# Introduction

Captain - Convert your Git workflow to Docker containers ready for Continuous Delivery

Define your workflow in the `captain.yaml` and use captain to your Continuous Delivery service to create containers for each commit, test them and push them to your registry only when tests passes.

* Use `captain build` to build your Dockerfile(s) of your repository. If your repository has local changes the containers will only be tagged as *latest*, otherwise the containers will be tagged as *latest*, *COMMIT_ID* & *BRANCH_NAME*. Now your Git commit tree is reproduced in your local docker repository.
* Use `captain test` to run your tests
* Use `captain push` to send selected images to the remote repository

From the other side, you can now pull the feature branch you want to test, or create distribution channels (such as 'alpha', 'beta', 'stable') using git tags that are propagated to container tags.

![intro](https://cloud.githubusercontent.com/assets/367397/6997822/c9aeadd8-dbcb-11e4-9901-dd62bcb33e5e.gif)

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
