-- -------------------------------------------------------------
-- TablePlus 6.4.8(608)
--
-- https://tableplus.com/
--
-- Database: bitespeed
-- Generation Time: 2025-06-10 10:59:48.3950
-- -------------------------------------------------------------


DROP TABLE IF EXISTS "public"."users";
-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
DROP TYPE IF EXISTS "public"."link_precedence";
CREATE TYPE "public"."link_precedence" AS ENUM ('primary', 'secondary');

-- Table Definition
CREATE TABLE "public"."users" (
    "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "phone_number" varchar NOT NULL DEFAULT ''::character varying,
    "email" varchar NOT NULL DEFAULT ''::character varying,
    "linked_id" int8,
    "link_precedence" "public"."link_precedence" NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz,
    PRIMARY KEY ("id")
);

ALTER TABLE "public"."users" ADD FOREIGN KEY ("linked_id") REFERENCES "public"."users"("id");


-- Indices
CREATE INDEX idx_users_email ON public.users USING btree (email);
CREATE INDEX idx_users_phone_number ON public.users USING btree (phone_number);
CREATE INDEX idx_users_linked_id ON public.users USING btree (linked_id);
