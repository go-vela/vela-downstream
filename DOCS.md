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
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
```

Sample of triggering a downstream build for a specific branch:

```diff
steps:
  - name: trigger_hello-world
    image: target/vela-downstream:latest
    pull: always
    parameters:
+     branch: main
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
```

Sample of triggering a downstream build for a specific event:

```diff
steps:
  - name: trigger_hello-world
    image: target/vela-downstream:latest
    pull: always
    parameters:
+     event: tag
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
```

Sample of triggering a downstream build for a specific status:

> **NOTE:**
>
> You can provide a list of statuses to the plugin.
>
> The first build found matching either of the statuses will be triggered.

```diff
steps:
  - name: trigger_hello-world
    image: target/vela-downstream:latest
    pull: always
    parameters:
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
+     status: [ success, failure ]
```

Sample of triggering a downstream build for multiple repos:

```diff
steps:
  - name: trigger_multiple
    image: target/vela-downstream:latest
    pull: always
    parameters:
      repos:
        - octocat/hello-world
+       - go-vela/hello-world
      server: https://vela-server.localhost
```

Sample of triggering a downstream build for multiple repos with different branches:

> **NOTE:**
>
> Use the @ symbol at the end of the org/repo to provide a unique branch per repo.
>
> This will override the value set for the `branch` parameter.

```diff
steps:
  - name: trigger_multiple
    image: target/vela-downstream:latest
    pull: always
    parameters:
      repos:
-       - octocat/hello-world
+       - octocat/hello-world@test
-       - go-vela/hello-world
+       - go-vela/hello-world@stage
      server: https://vela-server.localhost
```

## Secrets

> **NOTE:** Users should refrain from configuring sensitive information in your pipeline in plain text.

A personal access token for the related source control manager (GitHub, GitLab, etc.) is required.

The token must be tied to a user that exists in Vela, that is, the user has logged into Vela at least once.
Additionally, the user must have `write` access to both the origin repo as well as any downstream repos.

### Internal

Users can use [Vela internal secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: trigger_hello-world
    image: target/vela-downstream:latest
    pull: always
+   secrets: [ downstream_token ]
    parameters:
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

| Name                    | Description                                           | Required | Default       | Environment Variables                                                   |
| ----------------------- | ----------------------------------------------------- | -------- | ------------- | ----------------------------------------------------------------------- |
| `branch`                | branch to trigger a build on                          | `false`  | `N/A`         | `PARAMETER_BRANCH`<br>`DOWNSTREAM_BRANCH`                               |
| `event`                 | event to trigger a build on                           | `true`   | `push`        | `PARAMETER_EVENT`<br>`DOWNSTREAM_EVENT`                                 |
| `log_level`             | set the log level for the plugin                      | `true`   | `info`        | `PARAMETER_LOG_LEVEL`<br>`DOWNSTREAM_LOG_LEVEL`                         |
| `repos`                 | list of <org>/<repo> names to trigger a build on      | `true`   | `N/A`         | `PARAMETER_REPOS`<br>`DOWNSTREAM_REPOS`                                 |
| `server`                | Vela server to communicate with                       | `true`   | `N/A`         | `PARAMETER_SERVER`<br>`DOWNSTREAM_SERVER`                               |
| `status`                | list of statuses to trigger a build on                | `true`   | `[ success ]` | `PARAMETER_STATUS`<br>`DOWNSTREAM_STATUS`                               |
| `token`                 | SCM (GitHub, GitLab, etc.) personal access token of an existing Vela user | `true`   | `N/A`         | `PARAMETER_TOKEN`<br>`DOWNSTREAM_TOKEN`                                 |
| `report_back`           | whether or not to track downstream build status       | `false`  | `false`       | `PARAMETER_REPORT_BACK`<br>`DOWNSTREAM_REPORT_BACK`                     |
| `target_status`         | list of statuses to look for from downstream builds   | `false`  | `[ success ]` | `PARAMETER_TARGET_STATUS`<br>`DOWNSTREAM_TARGET_STATUS`                 |
| `timeout`               | how long should the plugin wait for downstream builds | `false`  | `30m`         | `PARAMETER_TIMEOUT`<br>`DOWNSTREAM_TIMEOUT`                             |
| `continue_on_not_found` | continue triggering builds on failure to find one     | `false`  | `false`       | `PARAMETER_CONTINUE_ON_NOT_FOUND`<br>`DOWNSTREAM_CONTINUE_ON_NOT_FOUND` |

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
+     log_level: trace
      repos:
        - octocat/hello-world
      server: https://vela-server.localhost
```

Below are a list of common problems and how to solve them:

### `unable to authenticate: user  not found`

Vela does not have a record of the user that owns the token, even if the user exists in the source control management system and the token may be valid.

Log into Vela as this user once to fix the issue.

### `unable to restart build myorg/myrepo/1234`

Vela does not have permission to restart the build with the given token.

Make sure the user that owns the token has write access to the repo that this .vela.yml is inside as well as any repos you want to trigger builds for.

Admin access on these repos is not required.
