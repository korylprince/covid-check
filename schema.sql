CREATE TABLE location (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL CHECK (1 < char_length(name) AND char_length(name) < 256)
);

CREATE TABLE "user" (
    id VARCHAR(255) PRIMARY KEY,
    name TEXT UNIQUE NOT NULL CHECK (1 < char_length(name) AND char_length(name) < 256)
);

CREATE TABLE checkin (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    time TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
    FOREIGN KEY (location_id) REFERENCES location(id),
    FOREIGN KEY (user_id) REFERENCES "user"(id)
);
