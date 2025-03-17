# Conditional Time Capsule: A URL Shortener by JPia

This is a URL shortener application built with Go and Gin. It provides endpoints to shorten URLs, access shortened URLs, check the status of shortened URLs, and manage URLs through admin routes.

## Assumptions

- Since this is just a proof-of-concept project, the API won’t generate a full short URL (that includes a hostname/base URL) and does not perform any redirects. Instead, the API will just return the shortcode every time there’s a requirement to return the short URL.
- The API will run using the New York timezone as the default.
- The API user will use the following datetime format when sending requests in NYC time: “2025-03-16T04:15:00Z”.
- To conserve API calls, the weather API will only be invoked up to 3 times per day. If the API call errored 3 times in a single day, then mark the day as “API_SICK_DAY” and allow eligible pending URLs to be released for the day. A CRITICAL log message will be generated for this event.
- A pending URL is eligible to be released when today is equal to or greater than the release date.
- An eligible pending URL will be released when today’s NYC weather is either clear or if there’s an “API_SICK_DAY”, otherwise, set the release date of the pending URL to 24 hours from today.
- The API should handle concurrency from both the endpoints and background processes. Requests from the endpoints are assumed to be higher priority than background tasks.

## Prerequisites

- Go 1.16 or higher
- Git

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/jpia/ctc.git
    cd ctc
    ```

2. Install the dependencies:

    ```sh
    go mod tidy
    ```

3. Copy one of the sample env files (i.e. '.env.dev') to create an `.env` file in the root directory and modify the values as necessary.

    - `USER_KEY`: The API key for user access. You can use the default during local dev.
    - `ADMIN_KEY`: The API key for admin access. You can use the default during local dev.
    - `SHORTCODE_LENGTH`: The length of the generated shortcodes.
    - `WEATHER_API_KEY`: The API key for accessing the weather service. Josh will provide the key he's using in the submission email.
    - `DEBUG`: Set to `true` to enable debug mode, `false` to disable. There will be a lot more log entries when this is enabled. It is recommended to turn this off during stress testing.
    - `RELEASE_TICKER_INTERVAL`: The time between release background jobs. The requirement is 1 hour, but I recommend changing this to 5 min during dev and stress testing.

4. (Optional) You can run the unit tests by:

    ```sh
    go test -v ./tests/
    ```

## Running the Application

1. Run the application:

    ```sh
    go run main.go
    ```

2. The application will be running at `http://localhost:8080`.

## Basic Routes

### Shorten URL

The `/shorten` route allows you to shorten a URL.

- **Endpoint**: `POST /shorten`
- **Headers**: `X-API-Key: your_user_key`
- **Request Body**: A JSON object with `long_url` and `release_date`.
- **Response**: A JSON object with the generated shortcode.

Example request:

```sh
curl -X POST -H "Content-Type: application/json" -H "X-API-Key: your_user_key" -d '{"long_url": "https://example.com", "release_date": "2025-03-16T04:15:00Z"}' http://localhost:8080/shorten
```

### Get URL Status

The `/status/:shortcode` route allows you to check the status of a shortened URL.

- **Endpoint**: `GET /status/:shortcode`
- **Headers**: `X-API-Key: your_user_key`
- **Response**: A JSON object with the status and release date of the URL.

Example request:

```sh
curl -H "X-API-Key: your_user_key" http://localhost:8080/status/shortcode
```

### Access URL

The `/access/:shortcode` route allows you to access the original URL using the shortcode.

- **Endpoint**: `GET /access/:shortcode`
- **Headers**: `X-API-Key: your_user_key`
- **Response**: A JSON object with the original long URL if the URL is released, or an error message if it is not yet available.

Example request:

```sh
curl -H "X-API-Key: your_user_key" http://localhost:8080/access/shortcode
```

### Override Shortcode

The `/admin/override/:shortcode` route allows an admin to manually release a URL early.

- **Endpoint**: `POST /admin/override/:shortcode`
- **Headers**: `X-API-Key: your_admin_key`
- **Response**: A JSON object indicating success or failure.

Example request:

```sh
curl -X POST -H "X-API-Key: your_admin_key" http://localhost:8080/admin/override/shortcode
```

## EXTRA: Nifty Admin Routes

### List All URLs

The `admin/list` route allows you to list all URLs stored in the system.

- **Endpoint**: `GET /admin/list`
- **Headers**: `X-API-Key: your_admin_key`
- **Response**: A JSON array of all URLs.

Example request:

```sh
curl -H "X-API-Key: your_admin_key" http://localhost:8080/admin/list
```

### Get URL Statistics

The `admin/stats` route provides statistics about the URLs stored in the system, including the total number of URLs and counts for pending, delayed, and released URLs.

- **Endpoint**: `GET /admin/stats`
- **Headers**: `X-API-Key: your_admin_key`
- **Response**: A JSON object with URL statistics.

Example request:

```sh
curl -H "X-API-Key: your_admin_key" http://localhost:8080/admin/stats
```


## EXTRA: A Basic Stress Tester

I have created a basic stress tester called `ctctester` which can be found at [https://github.com/jpia/ctctester](https://github.com/jpia/ctctester). This tool can be used to perform stress testing on the Conditional Time Capsule API to ensure it can handle high loads and concurrency.

By default, it will send 500 "shorten" requests and 100 overrides request per batch. It will also send a batch every 2 seconds.

Just make sure that this ctc repo is running first before running the ctctester.
