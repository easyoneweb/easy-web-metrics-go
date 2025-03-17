# Easy-Web-Metrics-Go

Server application to track web metrics for your visitors, version written in Go. It provides REST API which you can use from your own website to save visitor data: user id, last 10 pages visited, ip address, user-agent, additional user data if provided. Currently it provides visitor id and status ("new" or "updated") as visitor data for website.

## Setup

Make sure to install the dependencies:

```bash
go mod tidy
```

Environment variables should be put in .env file.

## Production

Build the application for production:

```bash
go build
```

## Environment variables

Application is using environment variables. You have to define:

- PORT (on which the server will run locally)
- API_KEY (to access application's REST API)
- DB_NAME (name of mongodb to use, it will be appended to mongodb URI as follows: mongodb://127.0.0.1:27017/name-of-database).

You can define all needed variables in .env file in root folder of this application.

## How to use Easy-Web-Metrics-Go REST API

After deployment of this app on some server, you should have access to it's REST API. For example we will use: your-domain.com.

REST API host: your-domain.com/api/v1
Api-Key to access routes should be provided in headers["api-key"]

GET /ping

Response example:
```json
{
  "message": "ping"
}
```

POST /visitor

Body example:
```json
{
    "visitor": "",
    "url": "some-domain.com/page/id",
    "utm": {
        "utmSource": "",
        "utmMedium": "",
        "utmCampaign": ""
    },
    "referrer": "",
    "userData": {
        "userID": "",
        "login": "",
        "email": "",
        "firstName": "",
        "secondName": "",
        "lastName": "",
        "phone": ""
    }
}
```
There are no required fields. But you should at least send empty json {}.

Response example:
```json
{
    "visitor": "e69a379a-a6e8-4686-b686-3b94e59545d3",
    "status": "new"
}
```
Visitor will have generated UUID and a status new, if no credentials were found or empty json {} has been sent. You should keep track of visitor ID to update data.

## How visitors are being referenced

If not data was provided, new visitor data and ID will be generated. Server is looking for data in a few steps:
- it looks for visitor ID
- if there's visitor data by visitor ID, it updates it with provided json data and returns ID and status
- if there's no visitor found by ID, app will look next for userData.bitrixID and update data is visitor was found
- if there's no visitor found by userData.bitrixID, app will look next for userAgent && IP address and update data is visitor was found
- if there's no visitor found by userAgent && IP address, new visitor will be generated

## Additional information

Easy-Web-Metrics-Go is written in Go language (Go 1.24.1), uses: mongodb as DB, mongo-db-driver, uuid, chi, godotenv. Please, before proceed be sure to check official documentation on corresponding technology. Tests are uses test db name: easywebmetricstest.

## Known Issues

There are currently no known issues.

## Release Notes

### 0.1.0

Initial working version.

---

## For more information

* [GitHub](https://github.com/ikirja/easy-ollama)
* [EasyOneWeb LLC](https://easyoneweb.ru)

# Copyright

EasyOneWeb LLC 2020 - 2025. All rights reserved. Code author: Kirill Makeev. See LICENSE.md for licensing and usage information.

**Enjoy!**

# TODO Roadmap:

- [X] Refactor to packages
- [X] Move db logic from visitor pkg to database pkg
- [X] Refactor models of visitor and visitorDB
- [X] Write tests
- [X] /api/v1/metrics/visitor instead of /api/v1/visitor
- [X] User UserData.UserID instead of UserData.BitrixID
- [X] database.VisitorUpdate UPDATE ONLY WHAT IS NEEDED TO BE UPDATED!
- [ ] use context in database methods instead of Context.TODO