# GoGrab

GoGrab is a web scraping tool written in Go that performs web crawling and data extraction. The project allows users to start a web crawl process, register and log in users, and manage scraped data.

## Features

- **User Authentication**: Register and log in users with JWT-based authentication.
- **Web Crawling**: Initiate and manage web scraping processes to collect data from websites.
- **Data Management**: Download, delete, and manage scraped data.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker (optional, for MySQL setup)

### Installation

1. **Clone the repository:**
2. **Install dependencies**
-- go mod tidy
3. **Set up the MySQL database with Docker**
-- docker-compose up -d
4. **Run the application**
-- go run main.go
5. **Swagger Documentation**
-- Swagger UI is available at http://localhost:8080/swagger/ to view and interact with the API documentation.

**This was my intern project as back-end developer**
