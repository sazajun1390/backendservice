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
    "user_multi_id" character varying NOT NULL,
    "resource_id" character varying NOT NULL,
    "email" character varying NOT NULL,
    "password" character varying NOT NULL,
    "post_code" character varying NOT NULL,
    "address" bigint NOT NULL,
    "address_kana" character varying NOT NULL,
    "nonce" character varying NOT NULL,
    "tel" character varying NULL,
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
    "created_at" timestamptz NOT NULL,
    CONSTRAINT "user_provision_user_id_key" UNIQUE ("user_id"),
    CONSTRAINT "user_provision_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

