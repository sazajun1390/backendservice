
-- +migrate Up

-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";

-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';

-- Create "users" table
CREATE TABLE "public"."users" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NOT NULL,
    PRIMARY KEY ("id")
);

-- Create "user_actives" table
CREATE TABLE "public"."user_actives" (
    "user_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL,
    CONSTRAINT "user_actives_user_id_key" UNIQUE ("user_id"),
    CONSTRAINT "user_actives_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create "user_deletes" table
CREATE TABLE "public"."user_deletes" (
    "user_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL,
    "purged_expires_at" timestamptz NULL,
    CONSTRAINT "user_deletes_user_id_key" UNIQUE ("user_id"),
    CONSTRAINT "user_deletes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create "user_profiles" table
CREATE TABLE "public"."user_profiles" (
    "user_id" bigint NOT NULL,
    "user_multi_id" varchar(32) NOT NULL,
    "resource_id" char(26) NOT NULL,
    "email" varchar(320) NOT NULL,
    "password" char(64) NOT NULL,
    "post_code" varchar(8) NOT NULL,
    "address" varchar(512) NOT NULL,
    "address_kana" varchar(512) NOT NULL,
    "tel" varchar(16) NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    "deleted_at" timestamptz NULL,
    "purged_expires_at" timestamptz NULL,
    CONSTRAINT "user_profiles_email_key" UNIQUE ("email"),
    CONSTRAINT "user_profiles_resource_id_key" UNIQUE ("resource_id"),
    CONSTRAINT "user_profiles_user_id_key" UNIQUE ("user_id"),
    CONSTRAINT "user_profiles_user_multi_id_key" UNIQUE ("user_multi_id"),
    CONSTRAINT "user_profiles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create "user_provision" table
CREATE TABLE "public"."user_provision" (
    "user_id" bigint NOT NULL,
    "magic_link_token" varchar(64),
    "created_at" timestamptz NOT NULL,
    "expired_at" timestamptz NOT NULL,
    CONSTRAINT "user_provision_user_id_key" UNIQUE ("user_id"),
    CONSTRAINT "user_provision_magic_link_token" UNIQUE ("magic_link_token"),
    CONSTRAINT "user_provision_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);


-- +migrate Down
DROP TABLE "public"."user_provision";
DROP TABLE "public"."user_profiles";
DROP TABLE "public"."user_deletes";
DROP TABLE "public"."user_actives";
DROP TABLE "public"."users";