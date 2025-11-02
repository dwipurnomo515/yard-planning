🚢 Yard Planning System
Backend service untuk mengelola perencanaan dan penempatan kontainer di yard pelabuhan.

🏗️ Arsitektur
Handler → Service → Repository → Database
Handler: Menangani HTTP requests dan responses
Service: Business logic dan validasi
Repository: Data access layer
Database: PostgreSQL
🚀 Setup & Installation
Prerequisites
Go 1.21 atau lebih tinggi
PostgreSQL 14 atau lebih tinggi
Redis 7 atau lebih tinggi (optional, untuk caching)
Git
golangci-lint (optional, untuk linting)

1. Clone Repository
   bash
   git clone https://github.com/dwipurnomo515/yard-planning.git
   cd yard-planning
2. Install Dependencies
   bash
   go mod download
3. Setup Database
   Buat database PostgreSQL:

bash
createdb yard_planning
Jalankan migration:

bash
psql -U postgres -d yard_planning -f migrations/001_init_schema.sql 4. Configuration
Copy .env.example ke .env dan sesuaikan:

bash
cp .env.example .env
Edit .env sesuai konfigurasi database Anda:

env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=yard_planning
SERVER_PORT=8080 5. Run Application
bash
go run cmd/api/main.go
Server akan berjalan di http://localhost:8080

📡 API Endpoints

1. Get Suggestion
   Mendapatkan saran posisi untuk kontainer.

Endpoint: POST /suggestion

Request Body:

json
{
"yard": "YRD1",
"container_number": "ALFI000001",
"container_size": 20,
"container_height": 8.6,
"container_type": "DRY"
}
Response:

json
{
"suggested_position": {
"block": "LC01",
"slot": 1,
"row": 1,
"tier": 1
}
} 2. Place Container
Menempatkan kontainer di yard.

Endpoint: POST /placement

Request Body:

json
{
"yard": "YRD1",
"container_number": "ALFI000001",
"block": "LC01",
"slot": 1,
"row": 1,
"tier": 1
}
Response:

json
{
"message": "Success"
} 3. Pickup Container
Mengambil kontainer dari yard.

Endpoint: POST /pickup

Request Body:

json
{
"yard": "YRD1",
"container_number": "ALFI000001"
}
Response:

json
{
"message": "Success"
} 4. Health Check
Endpoint: GET /health

Response: OK

🧪 Testing dengan cURL
Get Suggestion
bash
curl -X POST http://localhost:8080/suggestion \
 -H "Content-Type: application/json" \
 -d '{
"yard": "YRD1",
"container*number": "ALFI000001",
"container_size": 20,
"container_height": 8.6,
"container_type": "DRY"
}'
Place Container
bash
curl -X POST http://localhost:8080/placement \
 -H "Content-Type: application/json" \
 -d '{
"yard": "YRD1",
"container_number": "ALFI000001",
"block": "LC01",
"slot": 1,
"row": 1,
"tier": 1
}'
Pickup Container
bash
curl -X POST http://localhost:8080/pickup \
 -H "Content-Type: application/json" \
 -d '{
"yard": "YRD1",
"container_number": "ALFI000001"
}'
📊 Database Schema
Tables
yards: Yard information
blocks: Block information dalam yard
yard_plans: Perencanaan area untuk tipe kontainer tertentu
containers: Kontainer yang ada di yard
Relationships
yards (1) → (*) blocks
blocks (1) → (\_) yard_plans
blocks (1) → (\*) containers
🔒 Business Rules
Container Size:
20ft: menggunakan 1 slot
40ft: menggunakan 2 slot berurutan
Stacking Rules:
Tier 1 bisa langsung diisi
Tier > 1 hanya bisa diisi jika tier dibawahnya sudah ada kontainer
Pickup Rules:
Container hanya bisa diambil jika tidak ada container di atasnya
Yard Plan:
Setiap area block bisa memiliki plan untuk container dengan spesifikasi tertentu
Plan memastikan container ditempatkan di area yang sesuai
🛠️ Development
Project Structure
yard-planning/
├── cmd/api/ # Application entry point
├── internal/
│ ├── handler/ # HTTP handlers
│ ├── service/ # Business logic
│ ├── repository/ # Data access
│ ├── model/ # Domain models
│ └── middleware/ # HTTP middleware
├── pkg/
│ ├── database/ # Database connection
│ └── response/ # HTTP response helpers
├── migrations/ # Database migrations
└── config/ # Configuration
Best Practices Used
✅ Clean Architecture (Handler → Service → Repository)
✅ Separation of Concerns
✅ Repository Pattern
✅ Error Handling
✅ Input Validation
✅ Middleware (Logging, Recovery, CORS)
✅ Database Connection Pooling
✅ Prepared Statements (SQL Injection Prevention)

📈 Future Improvements
Unit Testing
Integration Testing
Redis Caching
Concurrent Execution
API Documentation (Swagger)
Docker Support
CI/CD Pipeline
Metrics & Monitoring
📄 License
MIT License
