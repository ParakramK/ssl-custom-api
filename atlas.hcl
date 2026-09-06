# Atlas migrations for the app-owned PostgreSQL database only.
# The legacy MySQL gatepass schema belongs to another app: read-only,
# never diffed or applied here.
#
# Usage (Atlas does not read .env by itself, so export it first):
#   set -a && source .env && set +a
#   atlas migrate diff <name> --env local  # generate migrations/<ts>_<name>.sql
#   atlas migrate apply --env local        # apply to the PG_* database

data "external_schema" "gorm" {
  program = ["go", "run", "./cmd/atlas"]
}

env "local" {
  src = data.external_schema.gorm.url
  # Dev database Atlas uses to normalize the desired schema.
  dev = "docker://postgres/18/dev"
  url = "postgres://${getenv("PG_USER")}:${getenv("PG_PASSWORD")}@${getenv("PG_HOST")}:${getenv("PG_PORT")}/${getenv("PG_DATABASE")}?sslmode=${getenv("PG_SSL_MODE")}"

  migration {
    dir = "file://migrations"
  }
}
