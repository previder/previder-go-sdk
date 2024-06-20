# Changelog
All notable changes to this project will be documented in this file.

## 2024-06
### Changed
- Updated to Go 1.22
- Added Kubernetes cluster support
- Added extra fields to Virtual Machines and Virtual Networks

## 2023-04-28
### Changed
- Updated to Go 19
- Altered API endpoints for non slash trailing endpoints

## 2019-08-22
### Added
- Improved error handling in baseclient.go
- When an error message occurs, the message itself is returned instead of only the error code
- The baseclient will log errors for better debugging purposes

### Changed
- Token header updated from Authentication with a bearer to the X-Auth-Token header

## 2019-06-26
### Changed
- Updated the paths and models for the V2 endpoint

## 2017-12-20
### Initial release
- Base client
- Virtual machine
- Virtual network
- Task
