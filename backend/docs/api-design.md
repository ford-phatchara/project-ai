# Durian Farm Manager — REST API Design

## Base URL

```
Production: https://api.durianfarm.app/v1
Development: http://localhost:3000/api/v1
```

---

## Authentication

All endpoints (except `/auth/*`) require a valid JWT in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Auth Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Create new user account |
| POST | `/auth/login` | Authenticate and get tokens |
| POST | `/auth/refresh` | Refresh access token |
| POST | `/auth/logout` | Invalidate refresh token |
| POST | `/auth/forgot-password` | Request password reset |
| POST | `/auth/reset-password` | Reset password with token |

---

### POST /auth/register

**Request:**
```json
{
  "email": "somchai@example.com",
  "password": "SecurePass123!",
  "full_name": "Somchai Durian",
  "phone": "+66812345678"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "usr_01HQ3X...",
      "email": "somchai@example.com",
      "full_name": "Somchai Durian",
      "phone": "+66812345678",
      "created_at": "2024-07-15T08:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJl...",
    "expires_in": 3600
  }
}
```

---

### POST /auth/login

**Request:**
```json
{
  "email": "somchai@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "usr_01HQ3X...",
      "email": "somchai@example.com",
      "full_name": "Somchai Durian"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJl...",
    "expires_in": 3600
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email or password"
  }
}
```

---

## Plots

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/plots` | List all plots |
| GET | `/plots/:id` | Get single plot |
| POST | `/plots` | Create new plot |
| PUT | `/plots/:id` | Update plot |
| DELETE | `/plots/:id` | Delete plot |
| GET | `/plots/:id/sales` | Get sales for a plot |
| GET | `/plots/:id/maintenance` | Get maintenance logs for a plot |
| GET | `/plots/:id/stats` | Get plot statistics |

---

### GET /plots

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 20 | Items per page (max 100) |
| `sort` | string | `-created_at` | Sort field (prefix `-` for desc) |
| `search` | string | — | Search by name |

**Request:**
```
GET /plots?page=1&limit=10&sort=name
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "plt_01HQ4A...",
      "name": "North Hill A",
      "size_rai": 5.5,
      "tree_count": 120,
      "notes": "Musang King variety, planted 2018",
      "created_at": "2024-01-10T08:00:00Z",
      "updated_at": "2024-06-15T10:30:00Z"
    },
    {
      "id": "plt_01HQ4B...",
      "name": "South Valley",
      "size_rai": 3.2,
      "tree_count": 75,
      "notes": "Monthong variety",
      "created_at": "2024-02-20T09:00:00Z",
      "updated_at": "2024-05-10T14:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 4,
    "total_pages": 1
  }
}
```

---

### GET /plots/:id

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "plt_01HQ4A...",
    "name": "North Hill A",
    "size_rai": 5.5,
    "tree_count": 120,
    "notes": "Musang King variety, planted 2018",
    "created_at": "2024-01-10T08:00:00Z",
    "updated_at": "2024-06-15T10:30:00Z"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Plot not found"
  }
}
```

---

### POST /plots

**Request:**
```json
{
  "name": "East Ridge",
  "size_rai": 4.0,
  "tree_count": 95,
  "notes": "New planting area, D24 variety"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "plt_01HQ4C...",
    "name": "East Ridge",
    "size_rai": 4.0,
    "tree_count": 95,
    "notes": "New planting area, D24 variety",
    "created_at": "2024-07-15T08:30:00Z",
    "updated_at": "2024-07-15T08:30:00Z"
  }
}
```

**Validation Error (422 Unprocessable Entity):**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "name",
        "message": "Name is required"
      },
      {
        "field": "size_rai",
        "message": "Size must be a positive number"
      }
    ]
  }
}
```

---

### PUT /plots/:id

**Request:**
```json
{
  "name": "East Ridge (Expanded)",
  "size_rai": 4.5,
  "tree_count": 110,
  "notes": "Expanded in July 2024"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "plt_01HQ4C...",
    "name": "East Ridge (Expanded)",
    "size_rai": 4.5,
    "tree_count": 110,
    "notes": "Expanded in July 2024",
    "created_at": "2024-07-15T08:30:00Z",
    "updated_at": "2024-07-16T09:00:00Z"
  }
}
```

---

### DELETE /plots/:id

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "plt_01HQ4C...",
    "deleted": true
  }
}
```

**Error (409 Conflict - Has Related Records):**
```json
{
  "success": false,
  "error": {
    "code": "HAS_DEPENDENCIES",
    "message": "Cannot delete plot with existing sales or maintenance records"
  }
}
```

---

### GET /plots/:id/stats

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `from` | date | 30 days ago | Start date (ISO 8601) |
| `to` | date | today | End date (ISO 8601) |

**Request:**
```
GET /plots/plt_01HQ4A.../stats?from=2024-01-01&to=2024-06-30
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "plot_id": "plt_01HQ4A...",
    "plot_name": "North Hill A",
    "period": {
      "from": "2024-01-01",
      "to": "2024-06-30"
    },
    "sales": {
      "total_revenue": 485000.00,
      "total_weight_kg": 2150.5,
      "average_price_per_kg": 225.50,
      "transaction_count": 45,
      "grade_breakdown": {
        "A": { "weight_kg": 1450.0, "revenue": 362500.00 },
        "B": { "weight_kg": 700.5, "revenue": 122500.00 }
      }
    },
    "maintenance": {
      "total_activities": 32,
      "watering_count": 15,
      "fertilizing_count": 8,
      "pruning_count": 5,
      "pest_control_count": 2,
      "harvesting_count": 2
    }
  }
}
```

---

## Sales

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sales` | List all sales |
| GET | `/sales/:id` | Get single sale |
| POST | `/sales` | Create new sale |
| PUT | `/sales/:id` | Update sale |
| DELETE | `/sales/:id` | Delete sale |
| GET | `/sales/summary` | Get sales summary/analytics |
| GET | `/sales/export` | Export sales as CSV |

---

### GET /sales

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 20 | Items per page (max 100) |
| `sort` | string | `-sale_date` | Sort field |
| `plot_id` | uuid | — | Filter by plot |
| `grade` | string | — | Filter by grade (A, B) |
| `buyer_id` | uuid | — | Filter by buyer |
| `from` | date | — | Start date (inclusive) |
| `to` | date | — | End date (inclusive) |
| `min_weight` | number | — | Minimum weight in kg |
| `max_weight` | number | — | Maximum weight in kg |
| `min_price` | number | — | Minimum price per kg |
| `max_price` | number | — | Maximum price per kg |

**Request:**
```
GET /sales?plot_id=plt_01HQ4A...&grade=A&from=2024-06-01&to=2024-06-30&sort=-total_price
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "sal_01HQ5X...",
      "sale_date": "2024-06-28",
      "plot": {
        "id": "plt_01HQ4A...",
        "name": "North Hill A"
      },
      "buyer": {
        "id": "buy_01HQ6Y...",
        "name": "Bangkok Fruits Co.",
        "phone": "+66891234567"
      },
      "grade": "A",
      "weight_kg": 85.5,
      "price_per_kg": 280.00,
      "total_price": 23940.00,
      "notes": "Premium quality batch",
      "created_at": "2024-06-28T14:30:00Z",
      "updated_at": "2024-06-28T14:30:00Z"
    },
    {
      "id": "sal_01HQ5Y...",
      "sale_date": "2024-06-25",
      "plot": {
        "id": "plt_01HQ4A...",
        "name": "North Hill A"
      },
      "buyer": null,
      "grade": "A",
      "weight_kg": 62.0,
      "price_per_kg": 275.00,
      "total_price": 17050.00,
      "notes": null,
      "created_at": "2024-06-25T10:00:00Z",
      "updated_at": "2024-06-25T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total_items": 2,
    "total_pages": 1
  }
}
```

---

### POST /sales

**Request:**
```json
{
  "sale_date": "2024-07-15",
  "plot_id": "plt_01HQ4A...",
  "buyer_id": "buy_01HQ6Y...",
  "grade": "A",
  "weight_kg": 95.0,
  "price_per_kg": 285.00,
  "notes": "Early morning harvest, excellent quality"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "sal_01HQ5Z...",
    "sale_date": "2024-07-15",
    "plot": {
      "id": "plt_01HQ4A...",
      "name": "North Hill A"
    },
    "buyer": {
      "id": "buy_01HQ6Y...",
      "name": "Bangkok Fruits Co."
    },
    "grade": "A",
    "weight_kg": 95.0,
    "price_per_kg": 285.00,
    "total_price": 27075.00,
    "notes": "Early morning harvest, excellent quality",
    "created_at": "2024-07-15T08:00:00Z",
    "updated_at": "2024-07-15T08:00:00Z"
  }
}
```

---

### GET /sales/summary

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `from` | date | Start of year | Start date |
| `to` | date | today | End date |
| `group_by` | string | `month` | Grouping: `day`, `week`, `month`, `year` |
| `plot_id` | uuid | — | Filter by plot |

**Request:**
```
GET /sales/summary?from=2024-01-01&to=2024-06-30&group_by=month
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "period": {
      "from": "2024-01-01",
      "to": "2024-06-30"
    },
    "totals": {
      "revenue": 1250000.00,
      "weight_kg": 5420.5,
      "transaction_count": 156,
      "average_price_per_kg": 230.60
    },
    "by_grade": {
      "A": {
        "revenue": 875000.00,
        "weight_kg": 3150.0,
        "percentage": 70.0
      },
      "B": {
        "revenue": 375000.00,
        "weight_kg": 2270.5,
        "percentage": 30.0
      }
    },
    "timeline": [
      {
        "period": "2024-01",
        "revenue": 85000.00,
        "weight_kg": 380.0,
        "transaction_count": 12
      },
      {
        "period": "2024-02",
        "revenue": 120000.00,
        "weight_kg": 520.0,
        "transaction_count": 18
      },
      {
        "period": "2024-03",
        "revenue": 195000.00,
        "weight_kg": 850.0,
        "transaction_count": 28
      },
      {
        "period": "2024-04",
        "revenue": 280000.00,
        "weight_kg": 1200.0,
        "transaction_count": 35
      },
      {
        "period": "2024-05",
        "revenue": 320000.00,
        "weight_kg": 1380.5,
        "transaction_count": 38
      },
      {
        "period": "2024-06",
        "revenue": 250000.00,
        "weight_kg": 1090.0,
        "transaction_count": 25
      }
    ]
  }
}
```

---

### GET /sales/export

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `format` | string | `csv` | Export format: `csv`, `xlsx` |
| `from` | date | — | Start date |
| `to` | date | — | End date |
| `plot_id` | uuid | — | Filter by plot |
| `grade` | string | — | Filter by grade |

**Request:**
```
GET /sales/export?format=csv&from=2024-06-01&to=2024-06-30
```

**Response (200 OK):**
```
Content-Type: text/csv
Content-Disposition: attachment; filename="sales_2024-06.csv"

sale_date,plot_name,buyer_name,grade,weight_kg,price_per_kg,total_price,notes
2024-06-28,North Hill A,Bangkok Fruits Co.,A,85.5,280.00,23940.00,Premium quality batch
2024-06-25,North Hill A,,A,62.0,275.00,17050.00,
...
```

---

## Maintenance Logs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/maintenance` | List all maintenance logs |
| GET | `/maintenance/:id` | Get single log |
| POST | `/maintenance` | Create new log |
| PUT | `/maintenance/:id` | Update log |
| DELETE | `/maintenance/:id` | Delete log |
| GET | `/maintenance/upcoming` | Get scheduled maintenance |

---

### GET /maintenance

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 20 | Items per page (max 100) |
| `sort` | string | `-activity_date` | Sort field |
| `plot_id` | uuid | — | Filter by plot |
| `activity_type` | string | — | Filter by type |
| `from` | date | — | Start date |
| `to` | date | — | End date |

**Activity Types:** `watering`, `fertilizing`, `pruning`, `pest_control`, `harvesting`

**Request:**
```
GET /maintenance?plot_id=plt_01HQ4A...&activity_type=watering&from=2024-06-01
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "mnt_01HQ7A...",
      "activity_date": "2024-06-28",
      "plot": {
        "id": "plt_01HQ4A...",
        "name": "North Hill A"
      },
      "activity_type": "watering",
      "duration_minutes": 45,
      "notes": "Morning watering, soil was very dry",
      "created_at": "2024-06-28T07:30:00Z",
      "updated_at": "2024-06-28T07:30:00Z"
    },
    {
      "id": "mnt_01HQ7B...",
      "activity_date": "2024-06-25",
      "plot": {
        "id": "plt_01HQ4A...",
        "name": "North Hill A"
      },
      "activity_type": "watering",
      "duration_minutes": 30,
      "notes": null,
      "created_at": "2024-06-25T06:00:00Z",
      "updated_at": "2024-06-25T06:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total_items": 2,
    "total_pages": 1
  }
}
```

---

### POST /maintenance

**Request:**
```json
{
  "activity_date": "2024-07-15",
  "plot_id": "plt_01HQ4A...",
  "activity_type": "fertilizing",
  "duration_minutes": 60,
  "notes": "Applied NPK 15-15-15, 2kg per tree"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "mnt_01HQ7C...",
    "activity_date": "2024-07-15",
    "plot": {
      "id": "plt_01HQ4A...",
      "name": "North Hill A"
    },
    "activity_type": "fertilizing",
    "duration_minutes": 60,
    "notes": "Applied NPK 15-15-15, 2kg per tree",
    "created_at": "2024-07-15T08:00:00Z",
    "updated_at": "2024-07-15T08:00:00Z"
  }
}
```

---

## Buyers

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/buyers` | List all buyers |
| GET | `/buyers/:id` | Get single buyer |
| POST | `/buyers` | Create new buyer |
| PUT | `/buyers/:id` | Update buyer |
| DELETE | `/buyers/:id` | Delete buyer |
| GET | `/buyers/:id/sales` | Get sales to a buyer |

---

### GET /buyers

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 20 | Items per page |
| `search` | string | — | Search by name |
| `sort` | string | `name` | Sort field |

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "buy_01HQ6Y...",
      "name": "Bangkok Fruits Co.",
      "contact_name": "Khun Somying",
      "phone": "+66891234567",
      "email": "somying@bangkokfruits.co.th",
      "address": "123 Ratchada Rd, Bangkok 10400",
      "notes": "Preferred buyer, always pays on time",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-03-20T14:30:00Z"
    },
    {
      "id": "buy_01HQ6Z...",
      "name": "Chanthaburi Export",
      "contact_name": "Khun Prasert",
      "phone": "+66897654321",
      "email": null,
      "address": "Chanthaburi Province",
      "notes": null,
      "created_at": "2024-02-10T09:00:00Z",
      "updated_at": "2024-02-10T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total_items": 2,
    "total_pages": 1
  }
}
```

---

### POST /buyers

**Request:**
```json
{
  "name": "Southern Markets Ltd.",
  "contact_name": "Khun Wichai",
  "phone": "+66823456789",
  "email": "wichai@southernmarkets.com",
  "address": "456 Songkhla Rd, Hat Yai 90110",
  "notes": "New buyer, interested in Grade A only"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "buy_01HQ6A...",
    "name": "Southern Markets Ltd.",
    "contact_name": "Khun Wichai",
    "phone": "+66823456789",
    "email": "wichai@southernmarkets.com",
    "address": "456 Songkhla Rd, Hat Yai 90110",
    "notes": "New buyer, interested in Grade A only",
    "created_at": "2024-07-15T08:00:00Z",
    "updated_at": "2024-07-15T08:00:00Z"
  }
}
```

---

## Dashboard / Analytics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/dashboard` | Get dashboard overview |
| GET | `/dashboard/revenue` | Get revenue analytics |
| GET | `/dashboard/activity` | Get recent activity feed |

---

### GET /dashboard

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `from` | date | 30 days ago | Start date |
| `to` | date | today | End date |

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "period": {
      "from": "2024-06-15",
      "to": "2024-07-15"
    },
    "stats": {
      "total_revenue": 485000.00,
      "total_weight_kg": 2150.5,
      "average_price_per_kg": 225.50,
      "transaction_count": 45
    },
    "comparison": {
      "revenue_change_percent": 12.5,
      "weight_change_percent": 8.2
    },
    "grade_distribution": {
      "A": { "percentage": 68, "weight_kg": 1462.34 },
      "B": { "percentage": 32, "weight_kg": 688.16 }
    },
    "top_plots": [
      { "id": "plt_01HQ4A...", "name": "North Hill A", "revenue": 185000.00 },
      { "id": "plt_01HQ4B...", "name": "South Valley", "revenue": 142000.00 }
    ],
    "recent_sales": [
      {
        "id": "sal_01HQ5X...",
        "sale_date": "2024-07-14",
        "plot_name": "North Hill A",
        "grade": "A",
        "weight_kg": 85.5,
        "total_price": 23940.00
      }
    ],
    "upcoming_maintenance": [
      {
        "id": "mnt_01HQ7D...",
        "activity_date": "2024-07-16",
        "plot_name": "South Valley",
        "activity_type": "fertilizing"
      }
    ]
  }
}
```

---

## Error Codes

| HTTP Status | Code | Description |
|-------------|------|-------------|
| 400 | `BAD_REQUEST` | Malformed request syntax |
| 401 | `UNAUTHORIZED` | Missing or invalid authentication |
| 403 | `FORBIDDEN` | Authenticated but not authorized |
| 404 | `RESOURCE_NOT_FOUND` | Resource does not exist |
| 409 | `CONFLICT` | Resource conflict (duplicate, dependencies) |
| 422 | `VALIDATION_ERROR` | Request validation failed |
| 429 | `RATE_LIMITED` | Too many requests |
| 500 | `INTERNAL_ERROR` | Server error |

**Standard Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable error message",
    "details": [
      {
        "field": "weight_kg",
        "message": "Weight must be greater than 0"
      }
    ],
    "request_id": "req_01HQ8X..."
  }
}
```

---

## Rate Limiting

- **Standard:** 100 requests per minute
- **Export endpoints:** 10 requests per minute
- **Auth endpoints:** 20 requests per minute

**Rate Limit Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1720000000
```

---

## Pagination

All list endpoints support cursor-based pagination for large datasets:

**Query Parameters:**
```
?page=1&limit=20
```

**Response Meta:**
```json
{
  "meta": {
    "page": 1,
    "limit": 20,
    "total_items": 156,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

---

## Versioning

API versioning is included in the URL path (`/v1/`). Breaking changes will be introduced in new versions while maintaining backward compatibility for at least 12 months.

---

## CORS

The API supports CORS for browser-based applications:

```
Access-Control-Allow-Origin: https://app.durianfarm.app
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Authorization, Content-Type
Access-Control-Max-Age: 86400
```
