# GOLBO

Golbo is a simple level 7 load balancer implementing multiple common load balancing algorithm.

Can use the following algorithms:

- Round robin
- Least connections
- Random
- ...more in progress

Also, Golbo performs `active health` checks in the background for the hosts, so requests can be distributed across active hosts.

## Local Setup

- Clone the repo.
```
git clone github.com/Yaxhveer/golbo
```

- Start the test server
```
go run test/main.go
```

- Run the application
```
go run main.go
```

Now, one the access the http://localhost:3330 to test the load balancer.

You can also modify the `config.yaml` file.

- Test using Docker
```
docker compose up
```