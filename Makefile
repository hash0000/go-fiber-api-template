start_rebuild_mode:
	air
migrate_local:
	goose -dir ./app/common/database/migrations postgres "postgres://postgres:root@127.0.0.1:5432/neurotales" up
jet:
	jet -dsn=postgresql://postgres:root@127.0.0.1:5432/neurotales?sslmode=disable -path="./app/common/types/jet"