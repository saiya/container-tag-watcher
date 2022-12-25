# container-tag-watcher

A simple tool that watches container registry's tag update and executes any command 👀

## Use cases

- Do `kubectl rollout restart deployments/something` (k8s) when your container's `:latest` tag updated
- Do `systemctl restart something.service` when your container's `:latest` tag updated

## Synopsis

```
container-tag-watcher [--debug][--aws-ecr] path-to-config-file.yml
```

- `--debug` : Show DEBUG logs (opt-in option)
- `--aws-ecr` : Enable AWS ECR credentials handling (opt-in option)
  - You need to provide AWS credentials with AWS SDK's standard way (such as setting `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`)

## Config file

Example:

```yaml
targets:
  # Minimal example
  "saiya/container-tag-watcher:latest":
    commands:
      - "echo container image update detected!"
  
  # Full example
  "12345678890.dkr.ecr.ap-northeast-1.amazonaws.com/my-application:latest":
    # Platform of the container image to check (default: linux/amd64)
    platform: "linux/amd64"

    # Interval of polling (default: 3m)
    # Syntax of the duration complies to https://pkg.go.dev/time#ParseDuration
    polling-interval: "3m"

    # If false and if any command failed, won't execute succeeding `commands`  (default: false).
    # If true, continue `commands` execution even on error.
    continue-on-error: false

    # Queue size of `commands` execution queue (default: 1)
    # If event fired during `commands` run, it will be queued until the queue is not full.
    # Events that cannot be queued will be discarded; this behavior is not only to prevent system overloading but also useful to "aggregate" events
    backlog-limit: 1

    # List of commands to execute when the tag update detected
    commands:
      # Commands will be executed sequentially.
      - "echo command #1"

      # Can use array instead of string.
      # Note: command will be executed *without* shell, and won't use $PATH (= need to write full path of the command)
      -
        - "/bin/echo"
        - "this is an single argument, no need to escape space characters"
```
