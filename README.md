# GoUrlShortener
Url Shortener high availability server written in go and deployed with Kubernetes.

## Database server configuration
The server requires by default access to a MongoDB database [^1]. The information required to connect to the database is defined in the server's [ConfigMap](https://github.com/NicolasCamachoP/GoUrlShortener/blob/master/deploy/server/shortener-configmap.yml) and [Secret](https://github.com/NicolasCamachoP/GoUrlShortener/blob/master/deploy/server/shortener-secrets.yml).

There is a Charmed MongoDB Operator available wich is very useful for deploying a single node, replica set or sharded cluster. You will need to [install first Juju](https://snapcraft.io/juju) and then follow the [usage guide for the charm](https://charmhub.io/mongodb).

[^1]: The IRepository interface was created to define all data access functions. The current implementation of the interface is built to use MongoDB but if you want to change the database you would only have to implement a different one and modify the dependency injection to use the new service.