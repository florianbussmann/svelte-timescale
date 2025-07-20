# svelte-timescale-backend

## Developing

Once you've installed dependencies with `uv sync`, perform database migration after running db-container and start a development server:

```bash
docker compose up timescaledb -d
alembic upgrade head
fastapi dev
```
