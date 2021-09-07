## Description

This plugin enables you to trigger builds for other repos for [Vela](https://go-vela.github.io/docs/) in a pipeline.

Source Code: https://github.com/go-vela/vela-downstream

Registry: https://hub.docker.com/r/target/vela-downstream

## Usage

> **NOTE:**
>
> Users should refrain from using latest as the tag for the Docker image.
>
> It is recommended to use a semantically versioned tag instead.

Sample of triggering a downstream build:

```yaml
steps:
  - name: trigger_hello-world
    image: target/vela-downstream:latest
    pull: always
    parameters:
      branch: master
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
```

Sample of triggering a downstream build for multiple repos:

```diff
steps:
  - name: trigger_multiple
    image: target/vela-downstream:latest
    pull: always
    parameters:
      branch: master
      repos:
        - octocat/hello-world
+       - go-vela/hello-world
      server: https://vela-server.localhost
```

Sample of triggering a downstream build for multiple repos with different branches:

> **NOTE:** Use the @ symbol at the end of the org/repo to provide a unique branch per repo.

```diff
steps:
  - name: trigger_multiple
    image: target/vela-downstream:latest
    pull: always
    parameters:
-     branch: master
      repos:
-       - octocat/hello-world
+       - octocat/hello-world@test
-       - go-vela/hello-world
+       - go-vela/hello-world@stage
      server: https://vela-server.localhost
```

## Secrets

> **NOTE:** Users should refrain from configuring sensitive information in your pipeline in plain text.

### Internal

Users can use [Vela internal secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: trigger_hello-world
    image: target/vela-downstream:latest
    pull: always
+   secrets: [ downstream_token ]
    parameters:
      branch: master
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
-     token: superSecretVelaToken
```

> This example will add the secret to the `trigger_hello-world` step as environment variables:
>
> * `DOWNSTREAM_TOKEN=<value>`

### External

The plugin accepts the following files for authentication:

| Parameter | Volume Configuration                                                  |
| --------- | --------------------------------------------------------------------- |
| `token`   | `/vela/parameters/downstream/token`, `/vela/secrets/downstream/token` |

Users can use [Vela external secrets](https://go-vela.github.io/docs/concepts/pipeline/secrets/origin/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: trigger_hello-world
    image: target/vela-downstream:latest
    pull: always
    parameters:
      branch: master
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
-     token: superSecretVelaToken
```

> This example will read the secret value in the volume stored at `/vela/secrets/`

## Parameters

> **NOTE:**
>
> The plugin supports reading all parameters via environment variables or files.
>
> Any values set from a file take precedence over values set from the environment.

The following parameters are used to configure the image:

| Name        | Description                                             | Required | Default   | Environment Variables                           |
| ----------- | ------------------------------------------------------- | -------- | --------- | ----------------------------------------------- |
| `branch`    | default branch to trigger a build on                    | `true`   | `master`  | `PARAMETER_BRANCH`<br>`DOWNSTREAM_BRANCH`       |
| `log_level` | set the log level for the plugin                        | `true`   | `info`    | `PARAMETER_LOG_LEVEL`<br>`DOWNSTREAM_LOG_LEVEL` |
| `repos`     | list of <org>/<repo> names to trigger a build on        | `true`   | `N/A`     | `PARAMETER_REPOS`<br>`DOWNSTREAM_REPOS`         |
| `server`    | Vela server to communicate with                         | `true`   | `N/A`     | `PARAMETER_SERVER`<br>`DOWNSTREAM_SERVER`       |
| `status`    | list of acceptable build statuses to trigger a build on | `true`   | `success` | `PARAMETER_STATUS`<br>`DOWNSTREAM_STATUS`       |
| `token`     | token for communication with Vela                       | `true`   | `N/A`     | `PARAMETER_TOKEN`<br>`DOWNSTREAM_TOKEN`         |

## Template

COMING SOON!

## Troubleshooting

You can start troubleshooting this plugin by tuning the level of logs being displayed:

```diff
steps:
  - name: trigger_hello-world
    image: target/vela-downstream:latest
    pull: always
    parameters:
      branch: master
+     log_level: trace
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
```

Below are a list of common problems and how to solve them:
