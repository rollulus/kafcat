[![Build Status](https://travis-ci.org/rollulus/kafcat.svg?branch=master)](https://travis-ci.org/rollulus/kafcat)
[![Go Report Card](https://goreportcard.com/badge/github.com/rollulus/kafcat)](https://goreportcard.com/report/github.com/rollulus/kafcat)

# Kafcat

Kafcat is supposed to be(come) the swiss army knife for Apache Kafka. It is in an early stage right now.

Kafcat is a single statically linked binary. That means that you can just `curl` it into your containers or onto your VMs and go.

# Installation

Get a pre-built Linux/macOS binary from the [releases](https://github.com/rollulus/kafcat/releases) or DIY:

    make bootstrap build

# Examples

To see what is exactly in a topic since 5 minutes ago:

    $ kafcat cat tweets --since=5m
    key: ""
    value: |
      00000000  00 00 00 00 01 80 80 ab  d9 ff c2 a6 95 18 02 38  |...............8|
      00000010  32 30 31 37 2d 30 36 2d  30 32 54 31 37 3a 30 31  |2017-06-02T17:01|
      00000020  3a 32 38 2e 30 30 30 2b  30 30 30 30 aa 9d df be  |:28.000+0000....|
      00000030  02 02 14 44 61 75 6e 74  6c 65 73 73 2e 02 14 4c  |...Dauntless...L|
      00000040  69 76 5f 52 65 6e 6e 69  65 00 00 aa 02 f4 0a d2  |iv_Rennie.......|
      00000050  e9 01 02 fc 01 52 54 20  40 42 42 43 42 72 65 61  |.....RT @BBCBrea|
      00000060  6b 69 6e 67 3a 20 22 57  65 27 72 65 20 67 65 74  |king: "We're get|
      00000070  74 69 6e 67 20 6f 75 74  22 20 2d 20 50 72 65 73  |ting out" - Pres|
      00000080  69 64 65 6e 74 20 44 6f  6e 61 6c 64 20 54 72 75  |ident Donald Tru|
      00000090  6d 70 20 61 6e 6e 6f 75  6e 63 65 73 20 55 53 20  |mp announces US |
      000000a0  69 73 20 77 69 74 68 64  72 61 77 69 6e 67 20 66  |is withdrawing f|
      000000b0  72 6f 6d 20 74 68 65 20  50 61 72 69 73 20 63 6c  |rom the Paris cl|
      000000c0  69 6d 61 74 65 20 61 67  72 65 65 6d 65 6e 74 e2  |imate agreement.|
      000000d0  80 a6 20 02 04 65 6e 01  02 00 02 00 02 02 02 00  |.. ..en.........|
      000000e0  02 00 02 00 00 02 02 e8  bf 93 05 02 22 42 42 43  |............"BBC|
      000000f0  20 42 72 65 61 6b 69 6e  67 20 4e 65 77 73 02 16  | Breaking News..|
      00000100  42 42 43 42 72 65 61 6b  69 6e 67 00              |BBCBreaking.|
    topic: tweets
    partition: 0
    offset: 101
    timestamp: 2017-06-02T19:01:28.741+02:00

To get a quick overview of available topics:

    $ kafcat topics
    - name: tweets
    - name: _schemas

To get more in-depth topic information:

    $ kafcat topics -a
    - name: tweets
      partitions:
        0:
          leader:
            id: 0
            address: shuttle:9092
          offsets:
            oldest: 0
            newest: 883
          writable: true
          replicas: [0]
          insyncreplicas: [0]

# Commands

    $ kafcat --help
    The pretended Swiss army knife for Apache Kafka

    Usage:
      kafcat [command]

    Available Commands:
      cat
      topics
      version

    Flags:
      -b, --broker-list string   brokers (default "localhost:9092")
          --client-cert string   filename of the client certificate in PEM format
          --client-key string    filename of the client's private key in PEM format
          --root-ca string       filename of the root certificate in PEM format
      -v, --verbose              be verbose

    Use "kafcat [command] --help" for more information about a command.


## kafcat cat

Inspect topic contents.

### Usage

    $ kafcat cat --help
    Usage:
      kafcat cat <topic>[:partition[,partition]*] [flags]

    Flags:
      -f, --follow         don't quit on end of log; keep following the topic
      -s, --since string   only return logs newer than a relative duration like 5s, 2m, or 3h, or shorthands 0 and now

    Global Flags:
      -b, --broker-list string   brokers (default "localhost:9092")
          --client-cert string   filename of the client certificate in PEM format
          --client-key string    filename of the client's private key in PEM format
          --root-ca string       filename of the root certificate in PEM format
      -v, --verbose              be verbose

### Example

    $ kafcat cat tweets --since=5m
    key: ""
    value: |
      00000000  00 00 00 00 01 80 80 ab  d9 ff c2 a6 95 18 02 38  |...............8|
      00000010  32 30 31 37 2d 30 36 2d  30 32 54 31 37 3a 30 31  |2017-06-02T17:01|
      00000020  3a 32 38 2e 30 30 30 2b  30 30 30 30 aa 9d df be  |:28.000+0000....|
      00000030  02 02 14 44 61 75 6e 74  6c 65 73 73 2e 02 14 4c  |...Dauntless...L|
      00000040  69 76 5f 52 65 6e 6e 69  65 00 00 aa 02 f4 0a d2  |iv_Rennie.......|
      00000050  e9 01 02 fc 01 52 54 20  40 42 42 43 42 72 65 61  |.....RT @BBCBrea|
      00000060  6b 69 6e 67 3a 20 22 57  65 27 72 65 20 67 65 74  |king: "We're get|
      00000070  74 69 6e 67 20 6f 75 74  22 20 2d 20 50 72 65 73  |ting out" - Pres|
      00000080  69 64 65 6e 74 20 44 6f  6e 61 6c 64 20 54 72 75  |ident Donald Tru|
      00000090  6d 70 20 61 6e 6e 6f 75  6e 63 65 73 20 55 53 20  |mp announces US |
      000000a0  69 73 20 77 69 74 68 64  72 61 77 69 6e 67 20 66  |is withdrawing f|
      000000b0  72 6f 6d 20 74 68 65 20  50 61 72 69 73 20 63 6c  |rom the Paris cl|
      000000c0  69 6d 61 74 65 20 61 67  72 65 65 6d 65 6e 74 e2  |imate agreement.|
      000000d0  80 a6 20 02 04 65 6e 01  02 00 02 00 02 02 02 00  |.. ..en.........|
      000000e0  02 00 02 00 00 02 02 e8  bf 93 05 02 22 42 42 43  |............"BBC|
      000000f0  20 42 72 65 61 6b 69 6e  67 20 4e 65 77 73 02 16  | Breaking News..|
      00000100  42 42 43 42 72 65 61 6b  69 6e 67 00              |BBCBreaking.|
    topic: tweets
    partition: 0
    offset: 101
    timestamp: 2017-06-02T19:01:28.741+02:00

## kafcat topics

List topics and topic properties.

### Usage

    $ kafcat topics --help
    Usage:
      kafcat topics [flags]

    Flags:
      -a, --all          show all possible info; equals -olwrsp
      -s, --isr          show broker ids of in-sync replicas; implies --partitions
      -l, --leader       show leader information for each partition; implies --partitions
      -o, --offsets      show oldest and newest offsets for each partition; implies --partitions
      -p, --partitions   show partitions
      -r, --replicas     show broker ids of replicas; implies --partitions
      -w, --writable     show writable state (valid leader accepting writes) for each partition; implies --partitions

    Global Flags:
      -b, --broker-list string   brokers (default "localhost:9092")
          --client-cert string   filename of the client certificate in PEM format
          --client-key string    filename of the client's private key in PEM format
          --root-ca string       filename of the root certificate in PEM format
      -v, --verbose              be verbose

### Example

    $ kafcat topics -a
    - name: tweets
      partitions:
        0:
          leader:
            id: 0
            address: shuttle:9092
          offsets:
            oldest: 0
            newest: 883
          writable: true
          replicas: [0]
          insyncreplicas: [0]
