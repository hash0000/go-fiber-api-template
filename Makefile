start_air:
	air
migrate:
	goose -dir ./app/common/database/migration postgres "postgres://postgres:root@127.0.0.1:5432/tsFastifyTemplate" up