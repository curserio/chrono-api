-- Table of masters (specialists providing services)
CREATE TABLE masters
(
    id                UUID PRIMARY KEY      DEFAULT uuidv7(), -- unique identifier for the master
    name              VARCHAR(255) NOT NULL,                  -- full name of the master
    email             VARCHAR(255) NOT NULL UNIQUE,           -- contact email
    phone             VARCHAR(50)  NOT NULL,                  -- contact phone number
    telegram_id       BIGINT UNIQUE,                          -- Telegram user ID (nullable)
    telegram_username VARCHAR(255),                           -- Telegram username (nullable)
    description       TEXT,                                   -- profile description of the master
    city              VARCHAR(100),                           -- city where the master operates (nullable)
    timezone          TEXT         NOT NULL DEFAULT 'UTC',    -- time zone of the master
    language          VARCHAR(10)  NOT NULL DEFAULT 'en',     -- ISO 639-1, e.g. 'en', 'ru'
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),    -- record creation timestamp
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT now()     -- last update timestamp
);

COMMENT ON TABLE masters IS 'Specialists providing services';
COMMENT ON COLUMN masters.id IS 'Unique identifier for the master';
COMMENT ON COLUMN masters.name IS 'Full name of the master';
COMMENT ON COLUMN masters.email IS 'Contact email';
COMMENT ON COLUMN masters.phone IS 'Contact phone number';
COMMENT ON COLUMN masters.telegram_id IS 'Telegram user ID (nullable)';
COMMENT ON COLUMN masters.telegram_username IS 'Telegram username (nullable)';
COMMENT ON COLUMN masters.description IS 'Profile description of the master, e.g., expertise and specialties';
COMMENT ON COLUMN masters.city IS 'City where the master operates (optional)';
COMMENT ON COLUMN masters.timezone IS 'Time zone of the master, used to display local schedule, e.g., Asia/Irkutsk';
COMMENT ON COLUMN masters.language IS 'Preferred language of the master (ISO 639-1), e.g., en, ru';
COMMENT ON COLUMN masters.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN masters.updated_at IS 'Last update timestamp';

-- Table of services provided by masters
CREATE TABLE services
(
    id          UUID PRIMARY KEY DEFAULT uuidv7(), -- unique identifier for the service
    master_id   UUID REFERENCES masters (id),      -- master providing this service
    name        VARCHAR(255)   NOT NULL,           -- service name
    description TEXT,                              -- service description (optional)
    duration    INTEGER        NOT NULL,           -- duration in minutes
    price       DECIMAL(10, 2) NOT NULL,           -- service price
    created_at  TIMESTAMPTZ    NOT NULL,           -- record creation timestamp
    updated_at  TIMESTAMPTZ    NOT NULL            -- last update timestamp
);

COMMENT ON TABLE services IS 'Services offered by masters';
COMMENT ON COLUMN services.id IS 'Unique identifier for the service';
COMMENT ON COLUMN services.master_id IS 'Reference to the master providing the service';
COMMENT ON COLUMN services.name IS 'Service name';
COMMENT ON COLUMN services.description IS 'Optional description of the service';
COMMENT ON COLUMN services.duration IS 'Duration of the service in minutes';
COMMENT ON COLUMN services.price IS 'Price of the service';
COMMENT ON COLUMN services.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN services.updated_at IS 'Last update timestamp';

-- Table of clients (customers)
CREATE TABLE clients
(
    id                UUID PRIMARY KEY      DEFAULT uuidv7(), -- unique identifier for the client
    name              VARCHAR(255) NOT NULL,                  -- full name of the client
    email             VARCHAR(255) NOT NULL UNIQUE,           -- contact email
    phone             VARCHAR(50)  NOT NULL,                  -- contact phone number
    telegram_id       BIGINT UNIQUE,                          -- Telegram user ID (nullable)
    telegram_username VARCHAR(255),                           -- Telegram username (nullable)
    city              VARCHAR(100),                           -- city where the client is located (nullable)
    timezone          TEXT         NOT NULL DEFAULT 'UTC',    -- client's time zone, e.g., Europe/Moscow
    language          VARCHAR(10)  NOT NULL DEFAULT 'en',     -- ISO 639-1, e.g. 'en', 'ru'
    created_at        TIMESTAMPTZ  NOT NULL,                  -- record creation timestamp
    updated_at        TIMESTAMPTZ  NOT NULL                   -- last update timestamp
);

COMMENT ON TABLE clients IS 'Customers receiving services';
COMMENT ON COLUMN clients.id IS 'Unique identifier for the client';
COMMENT ON COLUMN clients.name IS 'Full name of the client';
COMMENT ON COLUMN clients.email IS 'Contact email';
COMMENT ON COLUMN clients.phone IS 'Contact phone number';
COMMENT ON COLUMN clients.telegram_id IS 'Telegram user ID (nullable)';
COMMENT ON COLUMN clients.telegram_username IS 'Telegram username (nullable)';
COMMENT ON COLUMN clients.city IS 'City where the client is located (optional)';
COMMENT ON COLUMN clients.timezone IS 'Client''s time zone, used to display times in local timezone, e.g., Europe/Moscow';
COMMENT ON COLUMN clients.language IS 'Preferred language of the client (ISO 639-1), e.g., en, ru';
COMMENT ON COLUMN clients.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN clients.updated_at IS 'Last update timestamp';

-- Enum type for booking status
CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'cancelled', 'completed');
COMMENT ON TYPE booking_status IS 'Status of a booking';

-- Table of bookings (appointments)
CREATE TABLE bookings
(
    id         UUID PRIMARY KEY        DEFAULT uuidv7(),  -- unique booking identifier
    master_id  UUID REFERENCES masters (id),              -- booked master
    client_id  UUID REFERENCES clients (id),              -- booking client
    service_id UUID REFERENCES services (id),             -- booked service
    start_time TIMESTAMPTZ    NOT NULL,                   -- appointment start time in UTC
    end_time   TIMESTAMPTZ    NOT NULL,                   -- appointment end time in UTC
    status     booking_status NOT NULL DEFAULT 'pending', -- current status of the booking
    created_at TIMESTAMPTZ    NOT NULL,                   -- record creation timestamp
    updated_at TIMESTAMPTZ    NOT NULL                    -- last update timestamp
);

COMMENT ON TABLE bookings IS 'Appointments made by clients with masters';
COMMENT ON COLUMN bookings.id IS 'Unique booking identifier';
COMMENT ON COLUMN bookings.master_id IS 'Reference to the booked master';
COMMENT ON COLUMN bookings.client_id IS 'Reference to the booking client';
COMMENT ON COLUMN bookings.service_id IS 'Reference to the booked service';
COMMENT ON COLUMN bookings.start_time IS 'Appointment start time in UTC';
COMMENT ON COLUMN bookings.end_time IS 'Appointment end time in UTC';
COMMENT ON COLUMN bookings.status IS 'Current status of the booking';
COMMENT ON COLUMN bookings.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN bookings.updated_at IS 'Last update timestamp';

-- Enum type for schedule type
CREATE TYPE schedule_type AS ENUM ('weekly', 'cyclic', 'custom');
COMMENT ON TYPE schedule_type IS 'Schedule type';


-- Table of schedules for masters
CREATE TABLE schedules
(
    id         UUID PRIMARY KEY       DEFAULT uuidv7(),        -- unique schedule identifier
    master_id  UUID REFERENCES masters (id) ON DELETE CASCADE, -- master this schedule belongs to
    name       VARCHAR(255)  NOT NULL,                         -- schedule name (e.g., "main", "shift 2")
    type       schedule_type NOT NULL,                         -- weekly / cyclic / custom
    start_date DATE          NOT NULL,                         -- schedule start date
    end_date   DATE,                                           -- schedule end date (NULL = indefinitely)
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),           -- record creation timestamp
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT now()            -- last update timestamp
);

COMMENT ON TABLE schedules IS 'Schedules defining working periods for masters';
COMMENT ON COLUMN schedules.id IS 'Unique schedule identifier';
COMMENT ON COLUMN schedules.master_id IS 'Reference to the master';
COMMENT ON COLUMN schedules.name IS 'Schedule name (e.g., main, shift 2)';
COMMENT ON COLUMN schedules.type IS 'Schedule type: weekly (by weekday), cyclic (e.g 2 days on/2 days off), custom (specific slots)';
COMMENT ON COLUMN schedules.start_date IS 'Start date of the schedule';
COMMENT ON COLUMN schedules.end_date IS 'End date of the schedule (NULL = indefinite)';
COMMENT ON COLUMN schedules.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN schedules.updated_at IS 'Last update timestamp';

-- Table for weekly schedule details (hours per weekday) / cyclic schedule details (hours per day index)
CREATE TABLE schedule_days
(
    id          UUID PRIMARY KEY     DEFAULT uuidv7(),            -- unique identifier
    schedule_id UUID REFERENCES schedules (id) ON DELETE CASCADE, -- parent schedule
    weekday     INT,                                              -- optional, for weekly schedules, 1 = Monday … 7 = Sunday
    day_index   INT,                                              -- optional, for cyclic schedules
    start_time  TIME,                                             -- work start time
    end_time    TIME,                                             -- work end time
    is_day_off  BOOLEAN     NOT NULL DEFAULT FALSE,               -- true if master has a day off
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),               -- record creation timestamp
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),               -- last update timestamp
    UNIQUE (schedule_id, weekday)
);

COMMENT ON TABLE schedule_days IS 'Daily working hours for schedules';
COMMENT ON COLUMN schedule_days.id IS 'Unique identifier';
COMMENT ON COLUMN schedule_days.schedule_id IS 'Reference to parent schedule';
COMMENT ON COLUMN schedule_days.weekday IS 'Day of week for weekly schedules (1 = Monday … 7 = Sunday), NULL for cyclic';
COMMENT ON COLUMN schedule_days.day_index IS 'Index of day in cyclic schedule (1..cycle_length), NULL for weekly';
COMMENT ON COLUMN schedule_days.start_time IS 'Start time of working day';
COMMENT ON COLUMN schedule_days.end_time IS 'End time of working day';
COMMENT ON COLUMN schedule_days.is_day_off IS 'Indicates if the day is off';
COMMENT ON COLUMN schedule_days.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN schedule_days.updated_at IS 'Last update timestamp';

-- Table for specific date slots / day offs
CREATE TABLE schedule_slots
(
    id          UUID PRIMARY KEY     DEFAULT uuidv7(),            -- unique slot identifier
    schedule_id UUID REFERENCES schedules (id) ON DELETE CASCADE, -- parent schedule
    date        DATE        NOT NULL,                             -- specific date
    start_time  TIME,                                             -- slot start time (optional)
    end_time    TIME,                                             -- slot end time (optional)
    is_day_off  BOOLEAN     NOT NULL DEFAULT FALSE,               -- true if master has a day off
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),               -- record creation timestamp
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),               -- last update timestamp
    UNIQUE (schedule_id, date)
);

COMMENT ON TABLE schedule_slots IS 'Specific working slots or day-offs for schedules';
COMMENT ON COLUMN schedule_slots.id IS 'Unique slot identifier';
COMMENT ON COLUMN schedule_slots.schedule_id IS 'Reference to parent schedule';
COMMENT ON COLUMN schedule_slots.date IS 'Specific date';
COMMENT ON COLUMN schedule_slots.start_time IS 'Slot start time';
COMMENT ON COLUMN schedule_slots.end_time IS 'Slot end time';
COMMENT ON COLUMN schedule_slots.is_day_off IS 'Indicates if the day is off';
COMMENT ON COLUMN schedule_slots.created_at IS 'Record creation timestamp';
COMMENT ON COLUMN schedule_slots.updated_at IS 'Last update timestamp';

-- Create indexes
CREATE INDEX idx_bookings_master_id ON bookings (master_id);
CREATE INDEX idx_bookings_client_id ON bookings (client_id);
CREATE INDEX idx_bookings_start_time ON bookings (start_time);
CREATE INDEX idx_services_master_id ON services (master_id);
