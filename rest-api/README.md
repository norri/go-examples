# rest-api

This is an example of a RESTful API built without using external dependencies.

## Start

1. Run `go run .` and the applications is now running at `http://localhost:3000`.

   
## Endpoints

- `GET /api/v1/books`: Retrieves a list of all books.
  ```sh
  curl -X GET http://localhost:3000/api/v1/books
  ```

- `POST /api/v1/books`: Adds a new book to the collection.
  ```sh
  curl -X POST http://localhost:3000/api/v1/books \
       -H "Content-Type: application/json" \
       -d '{"title":"Title"}'
  ```
