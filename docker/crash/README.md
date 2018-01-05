# crash

A container that crashes after `CRASH_DELAY` seconds. Defaults to 10 if not set.
Useful for testing the behavior of your automation/monitoring/alerts so you
know how it will behave against production workloads.

<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Build an Image](#build-an-image)
- [Running the Image](#running-the-image)

<!-- /TOC -->

## Build an Image

`docker build -t crash .`

## Running the Image

You can specify how long before the container crashes by supplying the
`CRASH_DELAY` environment variable.

```
# time docker run --rm crash
docker run --rm crash  0.01s user 0.01s system 0% cpu 11.661 total
```

```
# time docker run --rm -e CRASH_DELAY=25 crash
docker run --rm -e CRASH_DELAY=25 crash  0.01s user 0.01s system 0% cpu 26.624 total
```