# Museum scanner server

This is the server for the [museum scanner project]().

## Installation

1. Clone the repository.
2. Run `go get` to install dependencies.
3. Run `go build .` to build the server or `air` to run the server in development mode.
4. Run `./scanner-server` to start the server.

## Endpoints

- **GET `/`** Shows the home page + input for tag ID.
- **GET `/?id={tag_id}`** Shows the tag page for the given tag ID.
- **POST `/scan`** Accepts `tag_id` and `location_id` in the request body. Saves the scan event to the database.
- **GET `/export`** Generates and exports HTML files for all tags in the database. Each file contains generated content unique to the tag.

