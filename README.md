# orchestra

orchestra is a user space process orchestrator command line tool that makes it easy to 
coordinate the start up phase of multiple applications and tail their logs to follow-up
on them.

### Capabilities

orchestra has a small and focused set of functionalities. It does:

- Start processes
- Process ordering: start first process A, then B then C
- Check liveness and report start up errors: start up process X, then check it's live before starting dependent processes (if the liveness check if defined)
- Tail processes logs: display logs from each process, prefixed with the name of the process

### Configuration

orchestra supports a simple configuration provided by an `orchestra.yml` file.
The tool will by default look for this file in the current directory (e.g. the result of running `pwd`).
It is possible to specify the path to the configuration file by specifying the `--config` flag:
e.g. `orchestra start --config=/path/to/orchestra.yml`

A minimal configuration file is structured as follows

```
# the processes we want to start
processes:
- name: Process A
  command: ./bin/processa
- name: Process B
  command: ./bin/processb
- name: Website
  command: npm run start
```

#### Full list of fields

```
- name: Process A
  command: ./bin/processa
  liveness:
    kind: exec
    spec:
      command: pg_isready
    initialDelay: 3s
    retries: 3
    retryDelay: 10s
- name: Process B
  command: ./bin/processb
  liveness:
    kind: httpGet
    spec:
      path: /health
      port: 8080
      httpHeaders:
      - name: Custom-Header
        value: Awesome
    initialDelay: 3s
    retries: 3
    retryDelay: 10s
- name: Website
  command: npm run start
  startDelay: 1m
```

- `name`: mandatory, used to identify the process to start
- `command`: mandatory, used to define the program to execute and the arguments to pass
- `startDelay`: optional, used to define a delay after which it is possible to start the defined process. The delay is defined via a duration string, a signed sequence of decimal numbers with optional fraction and unit suffix, like “100ms”, “2.3h” or “4h35m”.
- `liveness`: optional, used to define liveness probe, inspired by [k8s](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
  - `kind`: mandatory, valid values are `httpGet|exec`. It tells orchestra to perform a liveness check using either an HTTP GET call or by executing a shell command
  - `initialDelay`: optional, it tells orchestra that it should wait a certain amount of time before performing the first probe. The delay is defined via a duration string (see `startDelay`)
  - `retries`: optional, defines the number of retries orchestra should perform before giving up and reporting failure
  - `retryFrequency`: optional, defines the delay time between each retry. The delay is defined via a duration string (see `startDelay`)

### Commands

#### start

`start` runs all defined processes, in the desired order, and tails logs

##### Flags

- `--verbose` prints information statement (TODO)
- `--configPath` specifies a custom path to the `orchestra.yml` configuration file
