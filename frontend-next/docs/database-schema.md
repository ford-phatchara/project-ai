# Durian Farm Management System — Database Schema

## Entity Relationship Diagram (Text)

```
┌─────────────┐       ┌─────────────┐       ┌──────────────────┐
│   users     │       │    plots    │       │     buyers       │
├─────────────┤       ├─────────────┤       ├──────────────────┤
│ id (PK)     │       │ id (PK)     │       │ id (PK)          │
│ email       │       │ user_id(FK) │───┐   │ user_id (FK)     │
│ password    │       │ name        │   │   │ name             │
│ full_name   │       │ size_sqm    │   │   │ contact_phone    │
│ farm_name   │       │ tree_count  │   │   │ contact_email    │
│ phone       │       │ notes       │   │   │ address          │
│ created_at  │       │ created_at  │   │   │ notes            │
│ updated_at  │       │ updated_at  │   │   │ created_at       │
└──────┬──────┘       └──────┬──────┘   │   └────────┬─────────┘
       │                     │          │            │
       │  1                  │ 1        │            │ 1
       │  │                  │ │        │            │ │
       │  ▼ *                │ ▼ *      │            │ ▼ *
       │         ┌───────────┴──────────┴────────────┴───────────┐
       │         │                    sales                      │
       │         ├───────────────────────────────────────────────┤
       │         │ id (PK)                                       │
       └────────►│ user_id (FK) ─────────────────────────────────┤
                 │ plot_id (FK) ─────────────────────────────────┤
                 │ buyer_id (FK, nullable) ──────────────────────┤
                 │ sale_date                                     │
                 │ grade (A/B)                                   │
                 │ weight_kg                                     │
                 │ price_per_kg                                  │
                 │ total_price (computed)                        │
                 │ notes                                         │
                 │ created_at                                    │
                 │ updated_at                                    │
                 └───────────────────────────────────────────────┘

       │
       │  1
       │  │
       │  ▼ *
       │         ┌───────────────────────────────────────────────┐
       │         │              maintenance_logs                 │
       │         ├───────────────────────────────────────────────┤
       │         │ id (PK)                                       │
       └────────►│ user_id (FK) ─────────────────────────────────┤
                 │ plot_id (FK) ─────────────────────────────────┤
                 │ activity_type                                 │
                 │ log_date                                      │
                 │ duration_minutes (nullable)                   │
                 │ quantity (nullable)                           │
                 │ quantity_unit (nullable)                      │
                 │ notes                                         │
                 │ created_at                                    │
                 │ updated_at                                    │
                 └───────────────────────────────────────────────┘
```

---

## Entities & Relationships

| Entity            | Description                                   |
|-------------------|-----------------------------------------------|
| **users**         | Farm owners/operators with login credentials  |
| **plots**         | Individual farm plots owned by a user         |
| **buyers**        | Recurring customers who purchase durians      |
| **sales**         | Individual sale transactions                  |
| **maintenance_logs** | Farm activities (watering, fertilizing, etc.) |

### Relationships

| Relationship                  | Type        | Description                                      |
|-------------------------------|-------------|--------------------------------------------------|
| users → plots                 | One-to-Many | A user can own multiple plots                    |
| users → sales                 | One-to-Many | A user records multiple sales                    |
| users → maintenance_logs      | One-to-Many | A user records multiple maintenance activities   |
| users → buyers                | One-to-Many | A user manages their own buyer contacts          |
| plots → sales                 | One-to-Many | A plot can have multiple sales                   |
| plots → maintenance_logs      | One-to-Many | A plot can have multiple maintenance logs        |
| buyers → sales                | One-to-Many | A buyer can be associated with multiple sales    |

---

## Table Definitions

### 1. `users`
Authentication and user profile.

| Column       | Type                     | Constraints                    | Description                    |
|--------------|--------------------------|--------------------------------|--------------------------------|
| id           | UUID                     | PRIMARY KEY, DEFAULT uuid()    | Unique user identifier         |
| email        | VARCHAR(255)             | UNIQUE, NOT NULL               | Login email                    |
| password_hash| VARCHAR(255)             | NOT NULL                       | Bcrypt hashed password         |
| full_name    | VARCHAR(100)             | NOT NULL                       | User's full name               |
| farm_name    | VARCHAR(150)             |                                | Name of the farm               |
| phone        | VARCHAR(20)              |                                | Contact phone number           |
| created_at   | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Account creation timestamp     |
| updated_at   | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Last update timestamp          |

---

### 2. `plots`
Farm plot/land parcels.

| Column       | Type                     | Constraints                    | Description                    |
|--------------|--------------------------|--------------------------------|--------------------------------|
| id           | UUID                     | PRIMARY KEY, DEFAULT uuid()    | Unique plot identifier         |
| user_id      | UUID                     | NOT NULL, FK → users(id)       | Owner of the plot              |
| name         | VARCHAR(100)             | NOT NULL                       | Plot name (e.g., "Hillside A") |
| size_sqm     | DECIMAL(10,2)            | NOT NULL, CHECK (> 0)          | Size in square meters          |
| tree_count   | INTEGER                  | DEFAULT 0, CHECK (>= 0)        | Number of durian trees         |
| notes        | TEXT                     |                                | Additional notes               |
| created_at   | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Creation timestamp             |
| updated_at   | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Last update timestamp          |

**Index:** `idx_plots_user_id` on `user_id`

---

### 3. `buyers`
Customer/buyer contacts.

| Column        | Type                     | Constraints                    | Description                    |
|---------------|--------------------------|--------------------------------|--------------------------------|
| id            | UUID                     | PRIMARY KEY, DEFAULT uuid()    | Unique buyer identifier        |
| user_id       | UUID                     | NOT NULL, FK → users(id)       | Owner of this contact          |
| name          | VARCHAR(150)             | NOT NULL                       | Buyer/company name             |
| contact_phone | VARCHAR(20)              |                                | Phone number                   |
| contact_email | VARCHAR(255)             |                                | Email address                  |
| address       | TEXT                     |                                | Physical address               |
| notes         | TEXT                     |                                | Additional notes               |
| created_at    | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Creation timestamp             |
| updated_at    | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Last update timestamp          |

**Index:** `idx_buyers_user_id` on `user_id`

---

### 4. `sales`
Durian sale transactions.

| Column        | Type                     | Constraints                    | Description                       |
|---------------|--------------------------|--------------------------------|-----------------------------------|
| id            | UUID                     | PRIMARY KEY, DEFAULT uuid()    | Unique sale identifier            |
| user_id       | UUID                     | NOT NULL, FK → users(id)       | User who recorded the sale        |
| plot_id       | UUID                     | NOT NULL, FK → plots(id)       | Plot the durians came from        |
| buyer_id      | UUID                     | FK → buyers(id), NULLABLE      | Buyer (optional)                  |
| sale_date     | DATE                     | NOT NULL                       | Date of sale                      |
| grade         | VARCHAR(1)               | NOT NULL, CHECK IN ('A', 'B')  | Durian grade (A=premium, B=standard) |
| weight_kg     | DECIMAL(10,2)            | NOT NULL, CHECK (> 0)          | Total weight sold in kg           |
| price_per_kg  | DECIMAL(10,2)            | NOT NULL, CHECK (> 0)          | Price per kilogram                |
| total_price   | DECIMAL(12,2)            | GENERATED (weight_kg * price_per_kg) | Computed total price         |
| notes         | TEXT                     |                                | Additional notes                  |
| created_at    | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Creation timestamp                |
| updated_at    | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Last update timestamp             |

**Indexes:**
- `idx_sales_user_id` on `user_id`
- `idx_sales_plot_id` on `plot_id`
- `idx_sales_sale_date` on `sale_date`
- `idx_sales_grade` on `grade`

---

### 5. `maintenance_logs`
Farm maintenance activities.

| Column           | Type                     | Constraints                    | Description                       |
|------------------|--------------------------|--------------------------------|-----------------------------------|
| id               | UUID                     | PRIMARY KEY, DEFAULT uuid()    | Unique log identifier             |
| user_id          | UUID                     | NOT NULL, FK → users(id)       | User who recorded the activity    |
| plot_id          | UUID                     | NOT NULL, FK → plots(id)       | Plot where activity occurred      |
| activity_type    | VARCHAR(20)              | NOT NULL, CHECK IN (...)       | Type of maintenance               |
| log_date         | DATE                     | NOT NULL                       | Date of activity                  |
| duration_minutes | INTEGER                  | CHECK (>= 0), NULLABLE         | Duration in minutes (for watering, etc.) |
| quantity         | DECIMAL(10,2)            | NULLABLE                       | Amount used (fertilizer kg, etc.) |
| quantity_unit    | VARCHAR(20)              | NULLABLE                       | Unit of quantity (kg, liters)     |
| notes            | TEXT                     |                                | Additional notes                  |
| created_at       | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Creation timestamp                |
| updated_at       | TIMESTAMP WITH TIME ZONE | DEFAULT NOW()                  | Last update timestamp             |

**Activity Types:** `'watering'`, `'fertilizing'`, `'pruning'`, `'pest_control'`, `'harvesting'`

**Indexes:**
- `idx_maintenance_user_id` on `user_id`
- `idx_maintenance_plot_id` on `plot_id`
- `idx_maintenance_log_date` on `log_date`
- `idx_maintenance_activity_type` on `activity_type`

---

## SQL Schema (PostgreSQL)

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- USERS TABLE
-- ============================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    farm_name VARCHAR(150),
    phone VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================
-- PLOTS TABLE
-- ============================================
CREATE TABLE plots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    size_sqm DECIMAL(10,2) NOT NULL CHECK (size_sqm > 0),
    tree_count INTEGER DEFAULT 0 CHECK (tree_count >= 0),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_plots_user_id ON plots(user_id);

-- ============================================
-- BUYERS TABLE
-- ============================================
CREATE TABLE buyers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    contact_phone VARCHAR(20),
    contact_email VARCHAR(255),
    address TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_buyers_user_id ON buyers(user_id);

-- ============================================
-- SALES TABLE
-- ============================================
CREATE TABLE sales (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plot_id UUID NOT NULL REFERENCES plots(id) ON DELETE RESTRICT,
    buyer_id UUID REFERENCES buyers(id) ON DELETE SET NULL,
    sale_date DATE NOT NULL,
    grade VARCHAR(1) NOT NULL CHECK (grade IN ('A', 'B')),
    weight_kg DECIMAL(10,2) NOT NULL CHECK (weight_kg > 0),
    price_per_kg DECIMAL(10,2) NOT NULL CHECK (price_per_kg > 0),
    total_price DECIMAL(12,2) GENERATED ALWAYS AS (weight_kg * price_per_kg) STORED,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sales_user_id ON sales(user_id);
CREATE INDEX idx_sales_plot_id ON sales(plot_id);
CREATE INDEX idx_sales_sale_date ON sales(sale_date);
CREATE INDEX idx_sales_grade ON sales(grade);

-- ============================================
-- MAINTENANCE LOGS TABLE
-- ============================================
CREATE TABLE maintenance_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plot_id UUID NOT NULL REFERENCES plots(id) ON DELETE RESTRICT,
    activity_type VARCHAR(20) NOT NULL CHECK (activity_type IN ('watering', 'fertilizing', 'pruning', 'pest_control', 'harvesting')),
    log_date DATE NOT NULL,
    duration_minutes INTEGER CHECK (duration_minutes >= 0),
    quantity DECIMAL(10,2),
    quantity_unit VARCHAR(20),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_maintenance_user_id ON maintenance_logs(user_id);
CREATE INDEX idx_maintenance_plot_id ON maintenance_logs(plot_id);
CREATE INDEX idx_maintenance_log_date ON maintenance_logs(log_date);
CREATE INDEX idx_maintenance_activity_type ON maintenance_logs(activity_type);

-- ============================================
-- UPDATE TRIGGER FUNCTION
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to all tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_plots_updated_at BEFORE UPDATE ON plots FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_buyers_updated_at BEFORE UPDATE ON buyers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sales_updated_at BEFORE UPDATE ON sales FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_maintenance_updated_at BEFORE UPDATE ON maintenance_logs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

---

## Example Records

### Users
| id | email | full_name | farm_name | phone |
|----|-------|-----------|-----------|-------|
| `a1b2c3d4-...` | `ahmad@example.com` | Ahmad bin Hassan | Ladang Durian Raub | +60-12-345-6789 |

### Plots
| id | user_id | name | size_sqm | tree_count | notes |
|----|---------|------|----------|------------|-------|
| `p001-...` | `a1b2c3d4-...` | Plot A — Hillside | 2400.00 | 48 | Older trees, best for Grade A |
| `p002-...` | `a1b2c3d4-...` | Plot B — Valley | 1800.00 | 36 | Younger trees planted 2021 |
| `p003-...` | `a1b2c3d4-...` | Plot C — Eastern Ridge | 3200.00 | 64 | Mixed grades |
| `p004-...` | `a1b2c3d4-...` | Plot D — North Field | 1200.00 | 24 | Newly planted 2023 |

### Buyers
| id | user_id | name | contact_phone | contact_email |
|----|---------|------|---------------|---------------|
| `b001-...` | `a1b2c3d4-...` | Ahmad Traders | +60-13-111-2222 | ahmad.traders@mail.com |
| `b002-...` | `a1b2c3d4-...` | KL Export Co. | +60-3-8888-9999 | sales@klexport.com |
| `b003-...` | `a1b2c3d4-...` | Local Market | +60-12-555-6666 | NULL |

### Sales
| id | plot_id | buyer_id | sale_date | grade | weight_kg | price_per_kg | total_price |
|----|---------|----------|-----------|-------|-----------|--------------|-------------|
| `s001-...` | `p001-...` | `b001-...` | 2024-07-15 | A | 320.00 | 28.00 | 8960.00 |
| `s002-...` | `p003-...` | `b003-...` | 2024-07-18 | B | 210.00 | 14.00 | 2940.00 |
| `s003-...` | `p001-...` | `b002-...` | 2024-07-22 | A | 415.00 | 30.00 | 12450.00 |
| `s004-...` | `p002-...` | `b001-...` | 2024-07-28 | A | 180.00 | 27.00 | 4860.00 |

### Maintenance Logs
| id | plot_id | activity_type | log_date | duration_minutes | quantity | quantity_unit | notes |
|----|---------|---------------|----------|------------------|----------|---------------|-------|
| `m001-...` | `p001-...` | watering | 2024-07-10 | 45 | NULL | NULL | Soil dry, extra 20 min |
| `m002-...` | `p002-...` | fertilizing | 2024-07-11 | NULL | 25.00 | kg | NPK 15-15-15 blend |
| `m003-...` | `p003-...` | pest_control | 2024-07-13 | 60 | 2.50 | liters | Spotted fruit borers |
| `m004-...` | `p002-...` | pruning | 2024-07-18 | 90 | NULL | NULL | Removed dead branches |

---

## Row Level Security (RLS) Policies (Supabase)

```sql
-- Enable RLS on all tables
ALTER TABLE plots ENABLE ROW LEVEL SECURITY;
ALTER TABLE buyers ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales ENABLE ROW LEVEL SECURITY;
ALTER TABLE maintenance_logs ENABLE ROW LEVEL SECURITY;

-- Users can only access their own data
CREATE POLICY "Users can view own plots" ON plots FOR SELECT USING (auth.uid() = user_id);
CREATE POLICY "Users can insert own plots" ON plots FOR INSERT WITH CHECK (auth.uid() = user_id);
CREATE POLICY "Users can update own plots" ON plots FOR UPDATE USING (auth.uid() = user_id);
CREATE POLICY "Users can delete own plots" ON plots FOR DELETE USING (auth.uid() = user_id);

CREATE POLICY "Users can view own buyers" ON buyers FOR SELECT USING (auth.uid() = user_id);
CREATE POLICY "Users can insert own buyers" ON buyers FOR INSERT WITH CHECK (auth.uid() = user_id);
CREATE POLICY "Users can update own buyers" ON buyers FOR UPDATE USING (auth.uid() = user_id);
CREATE POLICY "Users can delete own buyers" ON buyers FOR DELETE USING (auth.uid() = user_id);

CREATE POLICY "Users can view own sales" ON sales FOR SELECT USING (auth.uid() = user_id);
CREATE POLICY "Users can insert own sales" ON sales FOR INSERT WITH CHECK (auth.uid() = user_id);
CREATE POLICY "Users can update own sales" ON sales FOR UPDATE USING (auth.uid() = user_id);
CREATE POLICY "Users can delete own sales" ON sales FOR DELETE USING (auth.uid() = user_id);

CREATE POLICY "Users can view own maintenance logs" ON maintenance_logs FOR SELECT USING (auth.uid() = user_id);
CREATE POLICY "Users can insert own maintenance logs" ON maintenance_logs FOR INSERT WITH CHECK (auth.uid() = user_id);
CREATE POLICY "Users can update own maintenance logs" ON maintenance_logs FOR UPDATE USING (auth.uid() = user_id);
CREATE POLICY "Users can delete own maintenance logs" ON maintenance_logs FOR DELETE USING (auth.uid() = user_id);
```

---

## Useful Queries

### Monthly Revenue Report
```sql
SELECT 
    DATE_TRUNC('month', sale_date) AS month,
    SUM(total_price) AS revenue,
    SUM(weight_kg) AS total_weight,
    COUNT(*) AS sale_count
FROM sales
WHERE user_id = $1 AND EXTRACT(YEAR FROM sale_date) = $2
GROUP BY DATE_TRUNC('month', sale_date)
ORDER BY month;
```

### Grade Distribution
```sql
SELECT 
    grade,
    SUM(weight_kg) AS total_weight,
    ROUND(SUM(weight_kg) * 100.0 / SUM(SUM(weight_kg)) OVER (), 2) AS percentage
FROM sales
WHERE user_id = $1
GROUP BY grade;
```

### Plot Performance
```sql
SELECT 
    p.name AS plot_name,
    COUNT(s.id) AS sale_count,
    SUM(s.weight_kg) AS total_kg,
    SUM(s.total_price) AS total_revenue,
    ROUND(AVG(s.price_per_kg), 2) AS avg_price_per_kg
FROM plots p
LEFT JOIN sales s ON s.plot_id = p.id
WHERE p.user_id = $1
GROUP BY p.id, p.name
ORDER BY total_revenue DESC;
```

### Recent Maintenance by Plot
```sql
SELECT 
    p.name AS plot_name,
    m.activity_type,
    m.log_date,
    m.duration_minutes,
    m.notes
FROM maintenance_logs m
JOIN plots p ON p.id = m.plot_id
WHERE m.user_id = $1
ORDER BY m.log_date DESC
LIMIT 20;
```
