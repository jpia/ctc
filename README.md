# Conditional Time Capsule: A URL Shortener by JPia

This is a URL shortener application built with Go and Gin. It provides endpoints to shorten URLs, access shortened URLs, check the status of shortened URLs, and manage URLs through admin routes.

## Prerequisites

- Go 1.16 or higher
- Git

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/ctc-url-shortener.git
    cd ctc-url-shortener
    ```

2. Install the dependencies:

    ```sh
    go mod tidy
    ```

3. Create a `.env` file in the root directory and add the following environment variables:

    ```env
    USER_KEY=your_user_key
    ADMIN_KEY=your_admin_key
    SHORTCODE_LENGTH=8
    ```

## Running the Application

1. Run the application:

    ```sh
    go run main.go
    ```

2. The application will be running at `http://localhost:8080`.

## Running the Tests

1. Run the tests with verbose output:

    ```sh
    go test -v ./tests/
    ```
