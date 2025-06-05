# URL Shortener Telegram Bot

A simple yet powerful Telegram Bot built with Go for shortening URLs. This project is designed with a clean architecture, making it easy to maintain, test, and extend. It is fully containerized with Docker for straightforward deployment.

## Features

- **Telegram Bot Interface**: Interact with the service directly through Telegram.
- **URL Shortening**: Converts long URLs into a short, unique alias.
- **HTTP Redirect Server**: A lightweight server to handle redirects from shortened links.
- **In-Memory Storage**: Simple and fast storage for local development (can be easily swapped for a persistent DB).
- **Dockerized**: Ready for production deployment with a multi-stage Dockerfile.
- **CI/CD Pipeline**: Includes a GitHub Actions workflow for automated testing and Docker image builds.

## Getting Started

Follow these instructions to get a local copy up and running for development and testing purposes.

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.22 or newer)
- [Docker](https://www.docker.com/get-started) (optional, for running in a container)
- A Telegram Bot Token from [@BotFather](https://t.me/BotFather)

### Running Locally

1.  **Clone the repository:**
    ```sh
    git clone git@github.com:iiixor/URLShortener.git
    cd URLShortener
    ```

2.  **Set up your environment:**
    Edit the `.env` file and add your Telegram Bot token:
    ```env
    TELEGRAM_TOKEN="your_telegram_bot_token_here"
    ```

3.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

4.  **Run the application:**
    ```sh
    go run ./cmd/app
    ```
    The bot and the HTTP server will start. You can now interact with your bot on Telegram!

## Deployment

This application is designed to be deployed as a Docker container.

1.  **Build the Docker image.** The CI/CD pipeline on GitHub automatically builds and pushes the image to Docker Hub. Make sure you have configured the necessary secrets (`DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`, `DOCKERHUB_REPOSITORY_NAME`) in your GitHub repository settings.

2.  **Run on your server.**
    Connect to your server and run the following command. Make sure to replace the placeholders with your actual data.

    ```bash
    docker run -d \
      -p 80:8080 \
      --restart=always \
      -e TELEGRAM_TOKEN="your_real_telegram_token" \
      -e BASE_URL="http://your-domain.com" \
      --name url-shortener \
      your_dockerhub_username/your_dockerhub_repo:latest
    ```

## Configuration

The application is configured via `config/local.yml` and can be overridden by environment variables.

- `ENV`: The application environment (`local`, `dev`, `prod`).
- `TELEGRAM_TOKEN`: **(Required)** Your token for the Telegram Bot API.
- `HTTP_SERVER_ADDRESS`: The address the internal HTTP server listens on (e.g., `localhost:8080`).
- `BASE_URL`: The public-facing base URL for your shortened links (e.g., `http://your-domain.com`).
- `ALIAS_LENGTH`: The length of the generated URL alias (defaults to 4).
