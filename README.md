# Cinema - Example of Microservices in Go with Docker, Kubernetes and MongoDB

## Overview

Cinema is an example project which demonstrates the use of microservices for a fictional movie theater.
The Cinema backend is powered by 4 microservices, all of which happen to be written in Go, using MongoDB for manage the database and Docker to isolate and deploy the ecosystem.

 * Movie Service: Provides information like movie ratings, title, etc.
 * Show Times Service: Provides show times information.
 * Booking Service: Provides booking information.
 * Users Service: Provides movie suggestions for users by communicating with other services.

The project structure is based in the knowledge learned in:

* Golang structure: <https://peter.bourgon.org/go-best-practices-2016/#repository-structure>
* Book Let's Go: <https://lets-go.alexedwards.net/>

Container images used support multi-architectures (amd64, arm/v7 and arm64).

## Index

* [Deployment](#deployment)
* [How To Use Cinema Services](#how-to-use-cinema-services)
* [Significant Revisions](#significant-revisions)
* [The big picture](#screenshots)

## Deployment

The application can be deployed in both environments: **local machine** or in a **kubernetes cluster**. You can find the appropriate documentation for each case in the following links:

* [local machine (docker compose)](./docs/localhost.md)
* [kubernetes (helm)](./docs/kubernetes-helm.md)
* [kubernetes (timoni)](./docs/kubernetes-timoni.md)

## How To Use Cinema Services

* [endpoints](./docs/endpoints.md)

## Significant Revisions

* [Microservices - Martin Fowler](http://martinfowler.com/articles/microservices.html)
* [Traefik Proxy Docs](https://doc.traefik.io/traefik/)
* [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
* [MongoDB Golang Channel](https://www.youtube.com/c/MongoDBofficial/search?query=golang)

## Screenshots

### Architecture

![overview](docs/images/overview.jpg)

### Homepage

![website home page](docs/images/website-home.jpg)

### Users List

![users list page](docs/images/website-users.jpg)
