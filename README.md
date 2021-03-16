# time-tracker
A simple time-tracking application in golang.

### start postgres container

```docker pull postgres:latest```

```docker run --name postgres -e POSTGRES_PASSWORD=[your_password] -d postgres```

```docker exec -it postgres psql -U postgres```

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

### start the application 

```go run main.go```

access the swagger in terminal 

http://127.0.0.1:8000/swagger/index.html