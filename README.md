<h1>Video On Demand Backend</h1>
<h1>About: </h1>
<h4>Project Language: Golang</h4>
<h4>Database: PostgreSQL (saving metadata and location of video)</h4>
<h4>Memory Caching: Redis (saving metadata and location of video)</h4>
<h4>Video Encoding: ffmpeg-go</h4>
<h4>Video Encryption: crypto/aes, crypto/cipher</h4>
<h4>Video Saving: local storage (folder: VideosCollection in the "./vod/modules/video")</h4>
<h4>Testing: testify </h4>
<h4>Mocking: mockery </h4>

&emsp;<h1>Assumptions</h1>
<h4>1. There are only two types of users: ADMIN, CUSTOMER; <br/>
Admin is responsible for uploading the video
Customer is view the videos available on the portal.
</h4>
<h4>2. Video Meta Data is saved using different API(can be done using ffmpeg-go package)</h4>
<h4>3. Encryption of videos are done on encoded videos</h4>
<h4>4. Encrypted Videos are decrypted by different API</h4>
<h4>5. Currently Video is saved in the local storage, but can be saved to a dedicated location using cloud storage, other than the project source code storage. Then the stored video can be distributed to different video streaming dedicated servers at different location for faster streaming and less bandwidth usage by the customer.</h4>

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

<h1>8. Video Description to use the APIs: </h1>
Follow the below link:

```
https://drive.google.com/file/d/1caC4P_sHuBQCTr_T7EIOVCaa3tPGGWkh/view?usp=sharing
```
