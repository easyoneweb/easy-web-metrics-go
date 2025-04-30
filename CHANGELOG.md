# Change Log

All notable changes to the "easy-web-metrics-go" application will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0] - - - - - -

### Added

- visitDates fields to visitor in DB to store last 30 unique dates of visits.
- New REST API GET /api/v1/metrics/stats/visitor that returns visitors for last 30 days.
- JSON schema for database models.
- MongoDB Collection Visitors indexes on UserData.UserID, Visitor, IP.

## [0.3.0] - 2025-03-14

### Added

- Visitor createdAt and updatedAt timestamp fields.
- Delete visitor if founded visitor by UserID has different Visitor ID then was sent by request, and requested visitor ID has no UserData.UserID.

### Changed

- Visitor is searched for by UserData.UserID first, then by Visitor.Visitor (visitor ID).

## [0.2.1] - 2025-03-14

### Fixed

- Cannot marshal type bson.D to a BSON Document when filter is nil.

## [0.2.0] - 2025-03-14

### Changed

- REST API: /api/v1/metrics/visitor instead of /api/v1/visitor.
- User data now uses userID instead of bitrixID.
- User data update only of there was no userID previously saved.

## [0.1.0] - 2025-03-14

### Added

- Initial working version.

### Fixed

- Initial working version.

### Changed

- Initial working version.

### Removed

- Initial working version.