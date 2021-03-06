---
title: Security of the Secretless Broker
id: docs
layout: docs
description: Secretless Broker Documentation
permalink: docs/overview/security.html
---

Protection of Secretless Broker credential resolution is based on industry-standard practices of keeping the data for the minimal duration needed and hardening of the execution environment. This is ensured by various strategies listed below.

## Hardened container

### Limited User

Our Docker container not only runs within an isolated Docker environment but is also running within a limited-user context. This ensures a layered system of proven SECCOMP kernel security on top of Linux standard user privilege context limitations.

### Minimal Container Runtime

Our container is also based on [Alpine Linux](https://alpinelinux.org/) which has an extremely limited amount of built-in tools, making it much harder for malware to operate in such an environment. Combined with the limited-user runner, the processes within the container are also unable to add additional packages to the container without circumventing Linux ["ring 3"](https://en.wikipedia.org/wiki/Protection_ring) isolation.

## Minimal Credential Retention

One of the biggest security risks for an application is keeping credentials around longer than they are needed. Secretless Broker centralizes the credential management functionality so that the application no longer has to worry about it, and in doing so is able to ensure that credentials are not stored in memory for longer than they are needed to open a new connection. Since each Listener/Handler combination is responsible for credential lifecycles these are the lifecycles of credentials for each of the built-in listeners:

- `http`: Listener fetches the credentials on each request and they are stored only for the duration of an individual connection authentication, after which they are zeroized.
- `mysql`: Credentials are loaded for each connection and then garbage-collected after connecting to the backend.
- `pg`: Credentials are loaded for each connection and then garbage-collected after connecting to the backend.
- `ssh`: Loaded on each new each ssh connection and only stored for the duration of the indvidual connection.
- `ssh-agent`: Loaded at Listener instantiation time.

## Hardened Networking

You can use both `localhost` listening address and/or socket files to exchange information between the applications and the
Secretless Broker which provide a communication channel that does not leave the host/pod. By having an isolated communication
channel between them, you can limit access to the Secretless Broker in a granular way. Additionally,
more security layers can be added to this system (e.g. [encrypted overlay network](https://docs.docker.com/network/overlay/#create-an-overlay-network),
Kubernetes pod collocation, etc) for improvement. Since these additional improvements are specific to individual
infrastructure deployments, however, they are currently outside the scope of this document.

Regardless of the connection strategy, the operating system provides security between the client and Secretless.
It is _very_ important to configure the OS properly so that unauthorized processes and clients can’t connect to Secretless.
With Unix domain sockets, operating system file permissions protect the socket. With TCP connections, Secretless
should be listening only on localhost.

### Future Work

We continue to investigate ways to make Secretless Broker even more secure. Please check the [changelog](https://github.com/cyberark/secretless-broker) and this page for updates on additional safeguards as we implement them.
