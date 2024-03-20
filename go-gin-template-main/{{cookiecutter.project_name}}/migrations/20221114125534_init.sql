-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "{{cookiecutter.database_schema}}"."order_status"(
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW(),
    "deleted_at" TIMESTAMP,
    "order_status_id" SERIAL PRIMARY KEY,
    "order_status_code" VARCHAR(100) UNIQUE,
    "order_status_name" VARCHAR(100),
    "order_status_color" VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS "{{cookiecutter.database_schema}}"."orders"(
    "created_at" TIMESTAMP DEFAULT NOW(),
    "updated_at" TIMESTAMP DEFAULT NOW(),
    "deleted_at" TIMESTAMP,
    "order_id" SERIAL PRIMARY KEY,
    "ip_address" VARCHAR(100),
    "order_session" VARCHAR(100) UNIQUE,
    "channel_id" INT REFERENCES "sale"."channel",
    "order_status_id" INT REFERENCES "test_schema"."order_status"("order_status_id"),
    "buyer_id" INT REFERENCES "sale"."customer",
    "city_id" INT REFERENCES "integration"."city",
    "order_comment" TEXT,
    "order_status_description" VARCHAR(500)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "{{cookiecutter.database_schema}}"."orders";
DROP TABLE "{{cookiecutter.database_schema}}"."order_status";
-- +goose StatementEnd
