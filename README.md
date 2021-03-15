# time-tracker
A simple time-tracking application in golang.

###start postgres container

`docker pull postgres:latest`
<br/>
<br/>
`docker run --name postgres -e POSTGRES_PASSWORD=[your_password] -d postgres` 
<br/>
<br/>
`docker exec -it postgres psql -U postgres`

### create table for users and time enteries
```
CREATE TABLE users (
ID SERIAL PRIMARY KEY,
EMAIL TEXT NOT NULL,
PASSWORD TEXT NOT NULL
);
```
<br/>


```
CREATE TABLE enteries(
ID SERIAL PRIMARY KEY,
USER_ID INT,
START_TIME TIMESTAMP,
END_TIME TIMESTAMP,
CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);
```