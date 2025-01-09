## Steps I did:
### Postgres (https://medium.com/@roystatham3003/database-connection-golang-docker-dfff9e958e47)
Pull

```
docker pull postgres
```
Run, open 5432 port
```
docker run --name postgres-docker-database -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
```
Create discogs-database
```
docker exec -ti postgres-docker-database createdb -U postgres discogs-database
```
Verify connection
```
> docker exec -ti postgres-docker-database psql -U postgres
> \c discogs-database
```

docker start postgres-docker-database


## TODO



1) postgresclient - communication with postgres
- db.go - connection, setting up client
- queries.go - stuff to query wand write to db
- writeToDB() - stores them in a PostgreSQL
- fetchOn()..... - provide an API to query the database
2) apiclient - communication with discogs api (https://www.discogs.com/)
- client.go - connection, setting up
- discogs.go - fetching
- fetchOnLabel() - fetches all releases for a selected record label
3) routes - contain all service endpoints
- fetch

## Approach

First I did a poc, with inmemory db based on maps, to try out storing, and have some working backend.
Anothere reason is that my sql database knowledge was last used during the studies

Then I set up react app

