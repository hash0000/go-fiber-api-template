start_air:
	air
migrate:
	goose -dir ./app/common/database/migration postgres "postgres://postgres:root@127.0.0.1:5432/tsFastifyTemplate" up
jet:
	jet -dsn=postgresql://postgres:root@127.0.0.1:5432/tsFastifyTemplate?sslmode=disable -path="./app/common/database/"