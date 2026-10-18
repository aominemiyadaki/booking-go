CREATE SCHEMA IF NOT EXISTS booking_bubuka;

-- Специализации
CREATE TABLE IF NOT EXISTS booking_bubuka.specializations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) UNIQUE NOT NULL,
    version         BIGINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    visibility      BOOLEAN DEFAULT TRUE
);

-- Мастера
CREATE TABLE IF NOT EXISTS booking_bubuka.masters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version         BIGINT DEFAULT 1,
    email           VARCHAR(64) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    full_name       VARCHAR(100) NOT NULL CHECK(char_length(full_name) BETWEEN 3 AND 100),
    specialization_id  UUID NOT NULL,
    avatar_url      TEXT,
    session_token   TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    phone_number    VARCHAR(20) UNIQUE NOT NULL CHECK(
        phone_number ~ '^[0-9]+$' AND
        char_length(phone_number) BETWEEN 10 AND 20
    ),
    visibility      BOOLEAN DEFAULT TRUE
);

-- Валюты
CREATE TABLE IF NOT EXISTS booking_bubuka.currencies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) UNIQUE NOT NULL,
    iso_code        VARCHAR(3) UNIQUE NOT NULL,
    version         BIGINT DEFAULT 1
);

-- Услуги
CREATE TABLE IF NOT EXISTS booking_bubuka.services (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_id        UUID NOT NULL,
    title            VARCHAR(255) NOT NULL,
    price            DECIMAL(10,2) NOT NULL DEFAULT 0,
    duration_minutes INT NOT NULL DEFAULT 0,
    description      TEXT DEFAULT '',
    sort_order       INT DEFAULT 0,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    version          BIGINT DEFAULT 1,
    currency_id      UUID NOT NULL,
    visibility       BOOLEAN DEFAULT TRUE
);

-- Клиенты
CREATE TABLE IF NOT EXISTS booking_bubuka.clients (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_id        UUID NOT NULL,
    name             VARCHAR(255) NOT NULL,
    phone_number     VARCHAR(20) UNIQUE NOT NULL CHECK(
        phone_number ~ '^[0-9]+$' AND
        char_length(phone_number) BETWEEN 10 AND 20
    ),
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    last_visit_at    TIMESTAMPTZ,
    total_visits     INT DEFAULT 1,
    notes            TEXT,
    is_blocked       BOOLEAN DEFAULT FALSE,
    version          BIGINT DEFAULT 1,
    visibility       BOOLEAN DEFAULT TRUE,
    birth_date       DATE,
    telegram_id      BIGINT,
    telegram_username VARCHAR(255)
);

-- Посещения
CREATE TABLE IF NOT EXISTS booking_bubuka.appointments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_id       UUID NOT NULL,
    client_id       UUID NOT NULL,
    service_id      UUID NOT NULL,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    status_id       UUID NOT NULL,
    comment         TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    version         BIGINT DEFAULT 1,
    visibility      BOOLEAN DEFAULT TRUE
);

-- Статусы посещений
CREATE TABLE IF NOT EXISTS booking_bubuka.appointment_statuses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) UNIQUE NOT NULL,
    version         BIGINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    visibility      BOOLEAN DEFAULT TRUE
);
    
-- Черный список
CREATE TABLE IF NOT EXISTS booking_bubuka.blacklist (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_id        UUID NOT NULL,
    client_id        UUID NOT NULL,
    version          BIGINT DEFAULT 1,
    created_at       TIMESTAMPTZ DEFAULT NOW(),
    visibility       BOOLEAN DEFAULT TRUE
);

-- Рабочие дни мастера
CREATE TABLE IF NOT EXISTS booking_bubuka.working_days (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_id       UUID NOT NULL,
    date            DATE NOT NULL,
    version         BIGINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    visibility      BOOLEAN DEFAULT TRUE
);
    
-- Рабочие слоты
CREATE TABLE IF NOT EXISTS booking_bubuka.working_slots (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    working_day_id  UUID NOT NULL,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    version         BIGINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    visibility      BOOLEAN DEFAULT TRUE
);
    
-- Типы перерывов
CREATE TABLE IF NOT EXISTS booking_bubuka.break_types (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) UNIQUE NOT NULL,
    version         BIGINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    visibility      BOOLEAN DEFAULT TRUE
);
    
-- Перерывы мастера
CREATE TABLE IF NOT EXISTS booking_bubuka.master_breaks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_id       UUID NOT NULL,
    break_type_id   UUID NOT NULL,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    version         BIGINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    visibility      BOOLEAN DEFAULT TRUE
);
