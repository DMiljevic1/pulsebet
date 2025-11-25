# Pulsebet – Event-Driven Microservices Demo (Go + Redpanda + Postgres + Prometheus + Grafana)

Pulsebet je jednostavan, ali realističan mikroservisni sustav za prikaz kako
**event-driven** arhitektura funkcionira u praksi koristeći:

- **Go** (services)
- **Redpanda (Kafka)** za komunikaciju između servisa
- **PostgreSQL** (per-service databases)
- **Prometheus** za metrics scraping
- **Grafana** za vizualizaciju
- **Docker Compose** kao orkestrator

Sustav se sastoji od dva mikroservisa:

1. **GameService**
   - Ima endpoint za kreiranje utakmica (`matches`)
   - Spremi match u svoju bazu
   - Emitira `MatchCreated` event u Redpandu

2. **BetService**
   - Consumer koji sluša `MatchCreated` event
   - Upisuje match u svoju vlastitu bazu kao *available match*
   - Osnovni idempotentni upsert (po `match_id`)

---

## 🚀 Pokretanje

### 1. Kloniraj projekt

git clone <tvoj-repo>
cd pulsebet

### 2. Pokreni projekt 

docker compose up --build
Docker Compose automatski podiže:

Postgres (sa tri baze: pulsebet, gameservice_db, betservice_db)

Redpanda (Kafka)

Prometheus

Grafana

GameService

BetService

### 3. Automatska inicijalizacija projekta 

Pri prvom pokretanju:

Postgres inicijalizira novi data folder

Docker automatski pokreće skripte iz db-init/

01-create-databases.sql → kreira baze

02-gameservice-schema.sql → kreira matches tablicu u gameservice_db

03-betservice-schema.sql → kreira matches + bets tablice u betservice_db

Nije potrebno ručno ništa raditi.

### 4. Testiranje flowa

curl -X POST http://localhost:8081/matches \
  -H "Content-Type: application/json" \
  -d '{"home":"Chelsea","away":"Arsenal"}'

GameService:

spremi match u gameservice_db

publisha MatchCreated event

BetService:

consuma event

upiše match u betservice_db u tablicu matches

idempotentno (ON CONFLICT (id) DO NOTHING / DO UPDATE)
