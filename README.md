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

### GET /ping

Response example:

```json
{
  "message": "ping"
}
```

### POST /metrics/visitor

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
  "ip": "",
  "userAgent": "",
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

### POST /metrics/stats/visitor

Body example:

```json
{
  "limit": 5000,
  "skip": 0
}
```

Response example:

```json
[
  {
    "createdAt": "2025-03-20T12:15:52.924Z",
    "updatedAt": "2025-03-26T16:46:01.622Z",
    "visitDates": ["2025-03-26T16:46:01.622Z"],
    "visitor": "15e04756-1daf-4b5b-a30b-80aad90ab050",
    "url": [
      {
        "url": "",
        "utm": {
          "utmSource": "",
          "utmMedium": "",
          "utmCampaign": ""
        },
        "referrer": "https://google.com"
      }
    ],
    "ip": "127.0.0.5",
    "userAgent": "test",
    "userData": {
      "userID": "2",
      "login": "",
      "email": "",
      "firstName": "",
      "secondName": "",
      "lastName": "",
      "phone": ""
    }
  }
]
```

## How visitors are being referenced

If not data was provided, new visitor data and ID will be generated. Server is looking for data in a few steps:

- it looks for userData.userID
- if there's visitor data found by userData.userID, it updates it with provided json data and returns visitor ID and status (note: userData will be updated only if no previous userData was present)
- if there's no visitor found by userData.userID, app will look next for visitor by visitor ID and update data if it was found
- if there's no visitor found by visitor ID, app will look next for userAgent && IP address and update data if it was found
- if there's no visitor found by userAgent && IP address, new visitor will be generated

## Additional information

Easy-Web-Metrics-Go is written in Go language (Go 1.24.1), uses: mongodb as DB, mongo-db-driver, uuid, chi, godotenv. Please, before proceed be sure to check official documentation on corresponding technology. Tests are uses test db name: easywebmetricstest.

## Known Issues

There are currently no known issues.

## Release Notes

### 0.4.0

- Added visitDates fields to visitor in DB to store last 30 unique dates of visits.
- Added REST API GET /api/v1/metrics/stats/visitor that returns visitors for last 30 days.
- Added JSON schema for database models.
- Added MongoDB Collection Visitors indexes on UserData.UserID, Visitor, IP.
- Added Processed Visitor createdAt field, includes field in REST API response.

### 0.3.0

- Added createdAt and updatedAt fields to visitor in DB.
- Search visitor by UserID first.
- Delete visitor if founded visitor by UserID has different Visitor ID then was sent by request, and requested visitor ID has no UserData.UserID.

### 0.2.1

- Fixed issue when bson.D filter for mongodb query is nil.

### 0.2.0

- REST API /api/v1/visitor changed to /api/v1/metrics/visitor.
- User data now accepts userID instead of bitrixID.
- User data update only updates user info if there was no userID previously saved.

### 0.1.0

Initial working version.

### About testing

Tests are using easywebmetricstest DB. Please, delete existing DB before running tests for pure results.

---

## For more information

- [GitHub](https://github.com/ikirja/easy-web-metrics-go)
- [EasyOneWeb LLC](https://easyoneweb.ru)

# Copyright

EasyOneWeb LLC 2020 - 2025. All rights reserved. Code author: Kirill Makeev. See LICENSE.md for licensing and usage information.

**Enjoy!**
