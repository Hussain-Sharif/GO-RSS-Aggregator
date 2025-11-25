# 📰 GO Feed Aggregator

A production-ready RSS/Atom feed aggregator built with Go, featuring concurrent feed processing, universal feed parsing (RSS & Atom), automatic feed discovery, and API key authentication.

This Feed aggregator solves the problem of manually checking multiple websites for updates. 
- Users can add RSS/Atom feeds, follow them, and get aggregated posts from all their subscribed(followed) feeds in one place. 
- The system automatically scrapes feeds in the background using Go's powerful concurrency features at most 10 different feeds.

### Key Features
-  Smart Feed Discovery - Paste any website URL, automatically discovers RSS/Atom feeds by matching common Feed URL patterns or by xml tags

- Universal Feed Parser - Supports RSS 2.0, RSS 1.0, Atom 1.0, and various date formats

-  Concurrent Scraping - Utilizes Go goroutines for parallel feed processing at most 10 feeds.

-  API Key Authentication - Secure user authentication with auto-generated API keys

- Feed Following System - Users can subscribe to feeds and get personalized content

- Background Worker - Automatic feed updates with new posts every minute 

- PostgreSQL Database - Type-safe database queries using SQLC and Migration tool Goose

- RESTful API - Clean, well-structured API endpoints.

### Architecture
![alt text](Architecture.png)

### Tech Stack
| Component      | Technology                          |
| -------------- | ----------------------------------- |
| Language       | Go                                  |
| Web Framework  | Chi Router                          |
| Database       | PostgreSQL 16                       |
| DB Tool        | SQLC (Type-safe SQL)                |
| Migration      | Goose                               |
| Authentication | API Keys                            |
| Concurrency    | Goroutines + WaitGroups             |
| Feed Parsing   | encoding/xml, golang.org/x/net/html |
| API Testing   | Postman |

### Quick Start

#### Prerequisites
- Go 1.25.1+
- PostgreSQL 16+
- SQLC (for code generation)
- Goose (for migrations)

1\. Clone the Repository
```bash
git clone https://github.com/Hussain-Sharif/GO-RSS-Aggregator.git
cd GO-RSS-Aggregator
```
2\. Install Dependencies
```bash
go mod download
```
3\. Set Up Environment Variables
```env
PORT=8000
DB_URL=postgres://username:password@localhost:5432/rssagg?sslmode=disable
```
4\. Set Up Database
```bash
# Create database
createdb rssagg

# Run migrations
cd sql/schema
goose postgres "postgres://username:password@localhost:5432/rssagg?sslrequire=disable" up
```

5\. Get the update DB Code (Optional-As Already Included)
```bash
sqlc generate
```

6\. Run Server
```bash
go build && ./GO-RSS-Aggregator
```

##### Tests and API Info can refer Postman [Collection Link](https://ersger.postman.co/workspace/My-Workspace~7ccb8262-1200-4f4a-829c-bb40ee401984/collection/32413422-9db14bbe-d743-4d9a-aa36-b646353aa119?action=share&creator=32413422&active-environment=32413422-bd8b8ce7-ea9e-4156-a176-1075c310b739) 


<div align="center">
⭐ Star this repo if you find it useful!

Built with ❤️ using Go

</div>