<h1> Clone the repo </h1>

```
git clone https://github.com/RheaSidana/VOD.git
```

<h1> Run the command:</h1>

```
go mod tidy
```

<h1>3. DB operations:</h1>
<h3><<<< INSTALLING >>>><br/>
    &emsp;a.PosgtreSQL <br/>
    &emsp;b.Redis</h3>

<br/>
<h3> <<<< CREATING >>>> <br/>&emsp;&emsp;5. create orm db: </h3>

```
CREATE database vod;
```

<h3>  <<<< CONNECTING >>>> <br/>&emsp;&emsp; 6. connect to db: </h3>

```
\c vod
```

<h3>&emsp;&emsp;7. Edit .env file with postgres details</h3>
<br/>


<h1>4. Migrate Tables: </h1>

```
go run .\migrations\migrate.go
```

<h3>&emsp;&emsp;View DB Table Schemas: </h3>

```
\d "tablename"
```

view all dbs : 

```
\dt
```

<h1>5. Seed Data to the Table </h1>

```
go run .\dataSeeding\seedData.go
```


<h1>6. Run the application </h1>

```
go run .
```


<h1>7. Call APIs </h1>
Postman: 

```
https://documenter.getpostman.com/view/28378586/2s9YJaY4C3
```

<!-- <h1>8. Video Description to use the APIs: </h1> -->