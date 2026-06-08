Motorcycle ticket payment system
Currently, parking at the Seagames dormitory is still using a day-based system, so we decided to build this system

📁 Requirements
1) A login system that separates the roles of the parking attendant and the owner
2) When the motorcycle enters, a ticket will be created and a ticket will be issued (there will be a barcode on the ticket). After the motorcycle leaves, the barcode will be scanned
3) If it exceeds 1 day, there will be a fine of 2000kip per day (for example, in the ticket on 12/4/2026, but when the barcode is scanned on 15/4/2026, there will be a fine of 6000kip)

4) There is a monthly membership system of 60,000kip
5) View daily, weekly, and monthly income
6) Can search for customer motorcycle tickets (in case the customer loses the ticket)

📁 Tech Stack

Language: Go (Golang) , Dart

Framework: Fiber v2 ,  Flutter

Database: PostgreSQL

ORM: GORM



Configuration: Environment-based configuration management.

📁 Project Structure
.
├── frontend/
│
└── backend/
    ├── cmd/
    │   └── api/
    │       └── main.go
    │
    ├── internal/
    │   ├── models/          # GORM model structs
    │   ├── repository/      # Database access layer
    │   ├── service/         # Business logic
    │   ├── handler/         # HTTP handlers (Fiber)
    │   ├── routes/          # Route registration
    │   ├── middleware/       # Auth, logging middleware
    │   ├── server/          # Fiber app setup
    │   └── utils/           # Helpers (QR, date calc, etc.)
    │
    └── database/            # Migration & DB connection

🗄️ Database ERD


