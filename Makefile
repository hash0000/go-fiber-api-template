start_air:
	air
migrate:
	goose postgres "postgres://postgres:root@127.0.0.1:5432/tsFastifyTemplate" GOOSE_MIGRATION_DIR=./app/common/migration/.db up