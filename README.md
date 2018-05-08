# Kubernetes Scheduler

## Introduction

This repository covers the use case of creating a custom Kubernetes scheduler written in Go.

It's just an example of how a Scheduler could be built, but currently there are some things not being handled at the moment in this scheduler:

- Pod labels like `NoSchedule`, `NoExecute`...
- Race conditions in case other scheduler schedules the same Pod.
- Prometheus metrics (we are working on it)
- Pod deployment
- Advanced scheduling (node affinity/anti-affinity, taints and tolerations, pod affinity/anti-affinity, ...)
 
When you write a custom scheduler you have to take all this things into account because you are on your own.

You can find more information about why and how create custom Kubernetes schedulers in our [blog](https://sysdig.com/blog/kubernetes-scheduler/).

## Requirements

- [Golang](https://golang.org/) for compilation
- Operating system:
    - GNU/Linux
    - MacOS
    - Windows
- x86 / x86_64 architecture
- [Git](https://git-scm.com/)
- [g++](https://gcc.gnu.org/) (Linux / MacOS) / [MinGW](http://mingw.org/) (Windows) 

## Build

### Using `go get`

```sh
go get -u -v -t github.com/draios/kubernetes-scheduler
```

The app should be compiled in `$GOPATH/bin/kubernetes-scheduler`

## TODO

- Print useful events in the Kubernetes event log
- Add timeouts & timeout handling functions
- Abstract away the decision & metrics source functions (make this scheduler more generic and vendor neutral)

## Copyright

License: Apache 2.0

Read file [LICENSE](https://github.com/draios/kubernetes-scheduler/blob/master/LICENSE).

## Links

- [Sysdig Webpage](https://sysdig.com/)
- [Sysdig Custom Scheduler Blog](https://sysdig.com/blog/kubernetes-scheduler/)
- [Sysdig Custom Scheduler KubeCon EU 2018 Talk](https://www.youtube.com/watch?v=4TaHQgG9wEg)
