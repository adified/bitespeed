CREATE TYPE link_precedence AS ENUM ('primary', 'secondary');

CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "phone_number" VARCHAR,
    "email" VARCHAR,
    "linked_id" BIGINT REFERENCES users(id),
    "link_precedence" link_precedence NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at" TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_users_linked_id ON users(linked_id);
