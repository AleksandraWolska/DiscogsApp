
### Postgres

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
later: 
```
docker start postgres-docker-database
```

### Backend
```
cd backend
go mod tidy
go run main go
```

### Frontend
There are some errors, but works
```
cd frontend
npm i
npm start
```

## About docker and docker compose files
My try to run all three parts at once. But I run out of time while debugging, so they don't work


