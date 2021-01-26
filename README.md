# Overview

The present repository contains an exemplarily implementation of a *Kafka*
producer. The implementation is based on the *Confluent Kafka Go* library and
mostly a slight extension of *Confluent's* example implementation, which can be
found
[here](https://github.com/confluentinc/examples/tree/5.5.0-post/clients/cloud/go).

## Usage

To build the *Kafka* producer just execute:

```
$ make
```

Make sure that the
[Confluent Kafka Go Build System](https://github.com/devsecurity-io/confluent-kafka-go-build-system)
was installed before.

The example application can be run with the command:

```
kafka-producer -c config.yaml process
```
