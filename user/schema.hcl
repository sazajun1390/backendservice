table "bun_migration_locks" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "table_name" {
    null = true
    type = character_varying
  }
  primary_key {
    columns = [column.id]
  }
  unique "bun_migration_locks_table_name_key" {
    columns = [column.table_name]
  }
}
table "bun_migrations" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "name" {
    null = true
    type = character_varying
  }
  column "group_id" {
    null = true
    type = bigint
  }
  column "migrated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
}
table "user_actives" {
  schema = schema.public
  column "user_id" {
    null = false
    type = bigint
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  foreign_key "user_actives_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  unique "user_actives_user_id_key" {
    columns = [column.user_id]
  }
}
table "user_profiles" {
  schema = schema.public
  column "user_id" {
    null = false
    type = bigint
  }
  column "user_multi_id" {
    null = false
    type = character_varying
  }
  column "resource_id" {
    null = false
    type = character_varying
  }
  column "email" {
    null = false
    type = character_varying
  }
  column "password" {
    null = false
    type = character_varying
  }
  column "post_code" {
    null = false
    type = character_varying
  }
  column "address" {
    null = false
    type = bigint
  }
  column "address_kana" {
    null = false
    type = character_varying
  }
  column "nonce" {
    null = false
    type = character_varying
  }
  column "tel" {
    null = true
    type = character_varying
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  column "deleted_at" {
    null = true
    type = timestamptz
  }
  column "purged_expires_at" {
    null = true
    type = timestamptz
  }
  foreign_key "user_profiles_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  unique "user_profiles_email_key" {
    columns = [column.email]
  }
  unique "user_profiles_resource_id_key" {
    columns = [column.resource_id]
  }
  unique "user_profiles_user_id_key" {
    columns = [column.user_id]
  }
  unique "user_profiles_user_multi_id_key" {
    columns = [column.user_multi_id]
  }
}
table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
}
enum "status_enum" {
  schema = schema.public
  values = ["active", "provisioning", "inactive", "deleted", "purged", "unspecified"]
}
schema "public" {
  comment = "standard public schema"
}
