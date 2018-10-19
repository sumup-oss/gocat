# gocat

Faster Golang alternative of [socat].

Multi-purpose relay from source to destination.

A relay is a tool for bidirectional data transfer between two independent data
channels.

## Supported source-to-destination relays:

* TCP to Unix,
* Unix to TCP.
* Need something else? Feel free to open an issue to discuss it or shoot a Pull Request.

## Why?

* Significantly faster than [socat] with medium and larger message payloads.
* Static binary, it just works. TM
* Actively health checks the `source` to prevent hanging/zombified `source` connections. (initial reason why [socat] didn't work for us)

## Why not?

* [socat] performs slightly better with small message payloads.

## Where it's used?

At SumUp we use it as a backbone for infrastructure and deployment system(s) that:
* need to relay SSH protocol, 
* proxy to TCP -> Unix or vice-versa where speed
* have reliability as an important concern.

As a now open-source project of SumUp, we hope that we find more use-cases together.

## Benchmarks

### How the benchmarks work

Benchmarks are sending a message from the destination, relaying via `gocat`/`socat` 
 to the source, which is an echo server that relays back to the destination via `gocat`/`socat`.
 
### Reading the benchmarks

`X` axis is the message payload size.

`Y` axis is throughput as per golang `test`'s `-count` argument, which benchmarks only
 the sending and receiving of a message sync or async.

### TCP to UNIX

![tcp-to-unix](/assets/tcp-to-unix.png)

### Unix to TCP

![unix-to-tcp](/assets/unix-to-tcp.png)

### Benchmarking mistakes?

Think we can improve them or got something wrong? Feel free to open an issue to discuss it.

We want the best possible benchmark and opportunity to improve the software!

## Configuration

Check out [config.go](./internal/config/config.go)

## Usage

### Unix Domain Socket to TCP

Example SSH agent forwarding

gocat

```shell
> gocat unix-to-tcp --src /run/ssh-agent.socket --dst 0.0.0.0:56789
```

socat

```shell
# NOTE: `-d -d -d` is to reach at least some level of verbosity
> socat -d -d -d TCP-LISTEN:56789,reuseaddr,fork UNIX-CLIENT:/run/ssh-agent.socket
```

### TCP to Unix Domain Socket

Example TCP to ssh-agent socket forwarding

gocat

```shell
> gocat tcp-to-unix --src 0.0.0.0:56789 --dst /tmp/sshagent.sock
```

socat

```shell
# NOTE: `-d -d -d` is to reach at least some level of verbosity
> socat -t 100000 -v UNIX-LISTEN:/tmp/sshagent.sock,unlink-early,mode=777,fork TCP:0.0.0.0:56789
```

## Contributing

Check out [CONTRIBUTING.md](./CONTRIBUTING.md)

## Code of conduct (CoC)
 
We want to foster an inclusive and friendly community around our Open Source efforts. Like all SumUp Open Source projects, this project follows the Contributor Covenant Code of Conduct. Please, [read it and follow it](CODE_OF_CONDUCT.md).
 
If you feel another member of the community violated our CoC or you are experiencing problems participating in our community because of another individual's behavior, please get in touch with our maintainers. We will enforce the CoC.

## About SumUp
 
![SumUp logo](https://raw.githubusercontent.com/sumup-oss/assets/master/sumup-logo.svg?sanitize=true)
 
It is our mission to make easy and fast card payments a reality across the *entire* world. You can pay with SumUp in more than 30 countries, already. Our engineers work in Berlin, Cologne, Sofia and Sāo Paulo. They write code in JavaScript, Swift, Ruby, Go, Java, Erlang, Elixir and more. Want to come work with us? [Head to our careers page](https://sumup.com/careers) to find out more.

[socat]: https://github.com/craSH/socat
