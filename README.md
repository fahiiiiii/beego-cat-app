
# Beego Cat App

This is a web application built using the [Beego](https://beego.vip/) framework. It integrates with The Cat API to display cat images, manage favorites, voting on images, and provide information about different cat breeds. The app is designed to showcase various features like fetching cat images, saving favorites, voting, and displaying breed details.

## Features

- **Cat Image Fetching**: Get random cat images or a list of images.
- **Favorites Management**: Save and delete favorite images.
- **Voting**: Vote for images and track votes.
- **Breed Information**: Get a list of cat breeds along with breed-specific images.
- **Concurrent Data Loading**: Fetch all necessary data concurrently for better performance.

## Setup and Run

### 1. Clone the Repository

Clone the project to your local machine:

```bash
git clone https://github.com/fahiiiiii/beego-cat-app.git
cd beego-cat-app
cd beego-cat-app

```

### 2. Install Dependencies

Make sure you have Go (v1.18 or above) installed. Install the necessary Go modules:

```bash
go mod tidy
```

### 3. Configure API Key

In the `conf/app.conf` file, set up your Cat API key:

```ini
cat_api_key = "your-cat-api-key"
```

### 4. Run the Application

Start the Beego server:

```bash
bee run
```

Visit the app at `http://localhost:8080`.

### 5. Run Unit Tests

To run unit tests for the controllers:

```bash
go test ./controllers -v
```

To generate a comprehensive coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

This will generate an HTML report of the test coverage.

## Project Structure

The project is organized as follows:

```
beego-cat-app/
├── conf/
│   └── app.conf
├── controllers/
│   ├── breed_controller.go
│   ├── cat_controller.go
│   ├── cat_controller_test.go
│   ├── concurrent_controller.go
│   ├── favourite_controller.go
│   ├── main_controller.go
│   └── unified_controller.go
├── models/
│   └── cat_image.go
├── routers/
│   └── router.go
├── views/
│   ├── breeds.html
│   ├── debug.html
│   ├── favorites.html
│   └── voting.html
├── .gitignore
├── beego-cat-app
├── coverage.html
├── coverage.out
├── go.mod
├── go.sum
└── main.go
```

## API Endpoints

### CatController Endpoints

- **GET** `/api/getImages`: Fetch images of cats.
- **GET** `/api/cat/random-image`: Fetch a random cat image.
- **POST** `/api/cat/save-favorite/:subId`: Save a favorite cat image.
- **GET** `/api/cat/favorites/:subId`: Get the user's favorite cat images.
- **DELETE** `/api/cat/delete-favorite/:favoriteId`: Delete a favorite cat image.
- **POST** `/api/vote/:subId`: Vote on a cat image.

### BreedController Endpoints

- **GET** `/api/breeds`: Fetch the list of cat breeds.
- **GET** `/api/breed-images`: Fetch images for a specific breed.

### Concurrent Data Loading Endpoints

- **GET** `/api/concurrent/voting-data`: Fetch concurrent voting data.
- **GET** `/api/concurrent/favorites-data`: Fetch concurrent favorites data.
- **GET** `/api/concurrent/breeds-data`: Fetch concurrent breed data.

## Repository

For more details, visit the [GitHub repository](https://github.com/fahiiiiii/beego-cat-app).