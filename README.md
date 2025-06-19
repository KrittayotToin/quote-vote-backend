# Quote Vote Backend - ระบบโหวตคำคม (Go API)

ระบบ Backend สำหรับโหวตคำคม พัฒนาด้วยภาษา Go ใช้ Fiber framework และฐานข้อมูล PostgreSQL พร้อมระบบยืนยันตัวตนด้วย JWT

## โครงสร้างโปรเจกต์ (เข้าใจง่าย)

```
quote-vote-backend/
├── cmd/
│   └── main.go          # จุดเริ่มต้น - กำหนด route และตั้งค่า server
├── handler/             # Controller (จัดการ HTTP request)
│   ├── user_handler.go  # สมัครสมาชิกและโปรไฟล์ผู้ใช้
│   ├── login_handler.go # เข้าสู่ระบบ
│   └── middleware.go    # JWT Middleware
├── model/               # โครงสร้างข้อมูล (Model)
│   └── models.go        # User, Quote, Vote model
├── dto/                 # Data Transfer Object
│   └── user_input.go    # โครงสร้าง request/response
├── config/              # การตั้งค่าต่าง ๆ
│   ├── db.go           # การเชื่อมต่อฐานข้อมูล
│   └── jwt.go          # ฟังก์ชัน JWT
└── docker-compose.yml   # ตั้งค่าฐานข้อมูล PostgreSQL
```

## วิธีการทำงาน (Flow แบบง่าย)

1. **รับ Request** → `main.go` ส่งไปยัง controller ที่เกี่ยวข้อง
2. **Controller** → จัดการ request และติดต่อฐานข้อมูลโดยตรง
3. **Database** → เก็บ/ดึงข้อมูล
4. **Response** → Controller ส่งกลับเป็น JSON

ไม่มีชั้นซับซ้อน - แค่ Route → Controller → Database!

## API Endpoint

### Public Route (ไม่ต้องยืนยันตัวตน)
- `POST /api/v1/register` - สมัครสมาชิกใหม่
- `POST /api/v1/login` - เข้าสู่ระบบ (รับ JWT token)

### Protected Route (ต้องยืนยันตัวตน)
- `GET /api/v1/profile` - ดูโปรไฟล์ผู้ใช้ปัจจุบัน

## การยืนยันตัวตนด้วย JWT

API นี้ใช้ JWT สำหรับยืนยันตัวตน:

1. **เข้าสู่ระบบ** → รับ JWT token
2. **เข้าถึง route ที่ต้องยืนยันตัวตน** → ใส่ token ใน Authorization header
3. **รูปแบบ Token** → `Bearer <your-jwt-token>`

## วิธีเริ่มต้นใช้งาน

1. เริ่มต้นฐานข้อมูล:
```bash
docker-compose up -d
```

2. รันเซิร์ฟเวอร์:
```bash
go run cmd/main.go
```

3. ทดสอบด้วย curl:
```bash
# สมัครสมาชิก
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","full_name":"John Doe","password":"password123"}'

# เข้าสู่ระบบ (รับ JWT token)
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# ใช้ token ที่ได้จาก login เพื่อเข้าถึง route ที่ต้องยืนยันตัวตน
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## ตัวอย่าง Response

### Login สำเร็จ:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "email": "test@example.com",
      "full_name": "John Doe",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Response จาก Protected Route:
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "email": "test@example.com",
    "full_name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

## โครงสร้างฐานข้อมูล

```sql
users
------
id            INTEGER PRIMARY KEY
email         TEXT UNIQUE NOT NULL
full_name     TEXT NOT NULL
password_hash TEXT NOT NULL

quotes
-------
id        INTEGER PRIMARY KEY
text      TEXT NOT NULL
author    TEXT
votes     INTEGER DEFAULT 0
created_by INTEGER REFERENCES users(id)

votes
------
id        INTEGER PRIMARY KEY
user_id   INTEGER REFERENCES users(id)
quote_id  INTEGER REFERENCES quotes(id)
created_at DATETIME DEFAULT CURRENT_TIMESTAMP
```

## Test Case

### 1. สมัครสมาชิก (Register)
- **Input:**
  - email: test@example.com
  - full_name: John Doe
  - password: password123
- **Expected:**
  - สถานะ 200 OK
  - JSON response มี success = true, message = "User registered successfully"

### 2. สมัครสมาชิกซ้ำ (Duplicate Register)
- **Input:**
  - email: test@example.com (ซ้ำกับที่สมัครไปแล้ว)
  - full_name: John Doe
  - password: password123
- **Expected:**
  - สถานะ 400 หรือ 409
  - JSON response มี success = false, message แจ้งว่า email ถูกใช้งานแล้ว

### 3. เข้าสู่ระบบ (Login)
- **Input:**
  - email: test@example.com
  - password: password123
- **Expected:**
  - สถานะ 200 OK
  - JSON response มี success = true, มี token ใน data

### 4. เข้าสู่ระบบด้วยรหัสผิด (Login Wrong Password)
- **Input:**
  - email: test@example.com
  - password: wrongpassword
- **Expected:**
  - สถานะ 401 Unauthorized
  - JSON response มี success = false, message แจ้งรหัสผ่านผิด

### 5. ดูโปรไฟล์ (Profile) ด้วย token ที่ถูกต้อง
- **Input:**
  - Authorization: Bearer <token ที่ได้จาก login>
- **Expected:**
  - สถานะ 200 OK
  - JSON response มีข้อมูล user

### 6. ดูโปรไฟล์โดยไม่ใส่ token หรือ token ผิด
- **Input:**
  - ไม่ใส่ Authorization header หรือใส่ token ผิด
- **Expected:**
  - สถานะ 401 Unauthorized
  - JSON response มี success = false, message แจ้งว่าไม่ได้รับอนุญาต
