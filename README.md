🚀 Quick Start

Prerequisites
- Go 1.21+

🛠 Installation
# Clone & setup
```
git clone https://github.com/RicYap/zleeper-be.git
cd healthcare-be
```
# Install dependencies
```
go mod init
go mod tidy
```
# Configure env
```
APP_ENV=development
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=zleeper
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

🏃 Running
# Development 
```
go run main.go
```

## API Documentation
[Postman Collection](https://crimson-station-529562.postman.co/workspace/zleeper~24b412e0-0a6a-43e4-b5f5-4b821f6c10bd/collection/32767687-264d52fb-be26-4025-ba6f-04498fc3adc0?action=share&creator=32767687)