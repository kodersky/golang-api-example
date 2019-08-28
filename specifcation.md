## Problem Statement

1. Must be a RESTful HTTP API listening to port `8080`
2. The API must implement 3 endpoints with path, method, request and response body as specified
    - One endpoint to create an order (see sample)
        - To create an order, the API client must provide an origin and a destination (see sample)
        - The API responds an object containing the distance and the order ID (see sample)

    - One endpoint to take an order (see sample)
        - An order can only be taken once.
        - An error response should be returned if a client tries to take an order that is already taken.

    - One endpoint to list orders (see sample)

3. The request input should be validated before processing. The server should return proper error response in case validation fails.
4. You must use following API to get the distance for the order:
- Google Maps API (https://cloud.google.com/maps-platform/routes/)
- **NOTE:** You don't have to provide actual Google Maps API key to us, just describe in the README how to use a custom key with your solution.
5. A Database must be used (SQL or NoSQL). The DB installation & initialisation must be done in `start.sh`.
6. All responses must be in json format no matter in success or failure situations.


## Api Interface

You are expected to follow the API specification as follows. Your implementation should not have any deviations on the method, URI path, request and response body. Such alterations may cause our automated tests to fail.

#### Place order

  - Method: `POST`
  - URL path: `/orders`
  - Request body:

    ```
    {
        "origin": ["START_LATITUDE", "START_LONGTITUDE"],
        "destination": ["END_LATITUDE", "END_LONGTITUDE"]
    }
    ```

  - Response:

    Header: `HTTP 200`
    Body:
      ```
      {
          "id": <order_id>,
          "distance": <total_distance>,
          "status": "UNASSIGNED"
      }
      ```
    or

    Header: `HTTP <HTTP_CODE>`
    Body:

      ```
      {
          "error": "ERROR_DESCRIPTION"
      }
      ```

  - Requirements:

    - Coordinates in request must be an array of exactly **two** strings. The type shall only be strings, not integers or floats.
    - The latitude and longtitude value of coordinates must be correctly validated.
    - Order id in response should be unique. It can be an auto-incremental integer or uuid string
    - Distance in response should be integer in meters


#### Take order

  - Method: `PATCH`
  - URL path: `/orders/:id`
  - Request body:
    ```
    {
        "status": "TAKEN"
    }
    ```
  - Response:
    Header: `HTTP 200`
    Body:
      ```
      {
          "status": "SUCCESS"
      }
      ```
    or

    Header: `HTTP <HTTP_CODE>`
    Body:
      ```
      {
          "error": "ERROR_DESCRIPTION"
      }
      ```

  - Requirements:

    - Since an order can only be taken once, you must be mindful of race condition.
    - When there are concurrent requests to take a same order, we expect only one can take the order while the other will fail.


#### Order list

  - Method: `GET`
  - Url path: `/orders?page=:page&limit=:limit`
  - Response:
    Header: `HTTP 200`
    Body:
      ```
      [
          {
              "id": <order_id>,
              "distance": <total_distance>,
              "status": <ORDER_STATUS>
          },
          ...
      ]
      ```

    or

    Header: `HTTP <HTTP_CODE>` Body:

    ```
    {
        "error": "ERROR_DESCRIPTION"
    }
    ```

  - Requirements:

    - Page number must starts with 1
    - If page or limit is not a valid integer then you should return error response
    - If there is no result, then you should return an empty array json in response body