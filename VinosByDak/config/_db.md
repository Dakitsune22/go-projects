# SQL

# DB
```
CREATE DATABASE vinos;
```

# Tables
```
CREATE TABLE users (
   Id          CHAR(50) PRIMARY KEY NOT NULL,
   Password    CHAR(64)             NOT NULL,
   Nickname    CHAR(30)             NOT NULL,
   Name        CHAR(30),
   Fname       CHAR(30),
   Lname       CHAR(30),
   Bdate       DATE                 DEFAULT '1900-01-01'
);
```