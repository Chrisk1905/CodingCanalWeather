# CodingCanal Weather App

This is a technical challenge for CodingCanal. 
The Weather App will call the Open Weather API for 5 cities, display the response periodically, and store the results into a postgres database.  

## Setup instructions

### Prequisites
Go version 1.23+

you can download it here
https://go.dev/doc/install


postgres (PostgreSQL) 16.4+

Download it here: https://www.postgresql.org/download/

OR you can install postgres with homebrew on mac
```
brew install postgresql

```
on Linux/WSL(debian) 
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

### Clone the repository

zsh
```
git clone https://github.com/Chrisk1905/CodingCanalWeather
cd CodingCanalWeather

```

Add dependancies
```
go mod download
```

### Set up the Postgres Database

Create the dattabase on Postgres. I will call my database weather.
```
psql postgres 
CREATE DATABASE weather;
SELECT version();
exit
```

### Create Tables from the Schema with Goose
From here I will be using psql and goose to create our tables.
If you use another SQL manager such as pgAdmin, feel free to use that instead. 
The schemas are in sql/schema. 

Add goose to create our schemas
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Get your postgres connection string
Mac OS
```
protocol://username@host:port/database
```
Linux/WSL
```
protocol://username:password@host:port/database
```
The port is 5432 by default.

for example on macOS:

postgres://username:@localhost:5432/weather

Test the connection string by connectiong to the weather database
```
psql "postgres://username:@localhost:5432/weather"
exit
```

cd into the sql/schema directory and run the up migration.
```
cd sql/schema
goose postgres <connection_string> up
```

Now the tables should be created! 
Check if the tables have been crreated.
```
psql weather
\dt
exit
```

### Run the app and tests

Add enviroment variables
```
export WEATHER_API_KEY='<your open weather API KEY>'
export DB_URL='<your postgres connection string with ssl mode disabled>'
```
Add your postgres connections string with ?sslmode=disable at the end.

`postgres://username:password@host:port/dbname?sslmode=disable`

for example:

`postgres://username:@localhost:5432/weather?sslmode=disable`


return to the root of the project

Run tests
```
go test ./internal/services
```

Start the app
```
go run .
```
You should see the weather JSON data for 5 cities, every 5 seconds. 



## Architecture overview
```
/internal
-/database
-/services

/sql
-/queries
-/schemas

main.go 
```

main.go - contains the dependancy injection, and calls services to periodically call the open weather API. 

**/sql** - contains the raw SQL. 

- sql/queries - contains the queries. I used SQLC to compile it into boilerplate Go code inside internal/database.
- sql/schemas - contains the schemas. I used Goose to create the tables and migrations.  

**/database** - Contains all the sqlc generated code. Models for data transfer, and Repositories for data acess 

- /database/models.go - The database models, or go structs that represent sql rows and columns.
- /database/database_table.sql.go - Contains the queries struct and all its methods to make queries to the database. 

**/services** - Contains the bussiness logic

- /services/weather.go - contains the WeatherService struct, which has the methods StoreWeather() and GetWeather() for storing
    the weather data, and making a GET call to the API
- /services/services_test.go - contains the mock API call to GetWeather()



## Design decisions explanation

I choose to use the tools I knew how to use from previous projects. 
I chose Go for my language a popular one for backend development.
I chose Postgres for my SQL database, along with Goose for database migration. 

### Database design

The database contains 4 tables:
- coordinates
- weather_conditions
- weather_data
- weather_data_conditions

Normalization

I tried to avoid data duplications and normalize to first normal form. 

- coordinates was used to store the coordinates of a city.
- weather_conditions was used to store conditions, it has imformation to discribe sunny clear skys, cloudy with rain, or windy with thunder.
- weather_data was used to store the rest of the data, and form relationships with coordinates and weather_conditions.

Relationships

- coordinates has a one to many relationship with weather_data
- weather_conditions has a many to many relationship with weather_data
- weather_data_conditions is the linking table to the many-to-many relationship between weather_conditions and weather_data


### The flaws
I had to make 3 big concessions due to a lack of time and experience on my part. 

**Synchrous over asynchronous**
    
For the sake of time I choose to call the API synchronously, as an minimum viable product. This made returning any errors easy, and was quicker to write and debug. If I was a more experience developer, or had more time, I would have refactored my code the weather API asynchronously using go routines and channel, one of the killer features of the language. 

**No Docker**

I don't have enough experience with docker yet, and couldn't containizerize the project.

**No database integration tests**

Adding database integration would either require more intrsuctions in the inital setup or using Docker. I choose to not do either. 
    
