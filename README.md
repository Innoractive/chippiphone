# Chippi Phone

This is a web application for retrieving business phone numbers. Specify an area
in Malaysia (e.g. "taman pelangi, johor bahru"). It search for F&B businesses
within that area and displays information (name, phone number, address).

## Getting Started

### Dependencies

* Google Cloud Project
  * The application uses Google Places API (New) and Geocoding API.
  * [Setup](https://bit.ly/3FaNs63) a Google Cloud project and enable the following APIs
    * Places API (New)
    * Geocoding API
  * [Setup](https://bit.ly/4ikbEBu) an API key
* Redis Server is required for caching query results.

### Executing Program

Configure Redis credentials and Google API key in `.env` file.

```shell
# Download the source code
git clone git@github.com:Innoractive/chippiphone.git
# Instantlly run the application
cd chippiphone
go run cmd/main.go
```

Note: make sure the _running_ directory contains `.env` configuration and `view/` folder.

### Development

[air-verse/air](https://github.com/air-verse/air) is used to automatically reload the application during development. Install `air` and `make`. A default `.air.toml` air configuration files has been provided. Just execute `air` command and it will automatically reload the web service when code changes.