### 🧩 Docker Compose Services

| Service Name      | Type      | Source / Image            | Ports              | Depends On                        | Language | Purpose / Notes                     |
| ----------------- | --------- | ------------------------- | ------------------ | --------------------------------- | -------- | ----------------------------------- |
| `consul`          | Container | `hashicorp/consul:1.17`   | `8500`, `8600/udp` | –                                 | –        | Service discovery and configuration |
| `gateway`         | Built     | `./backend/gateway`       | `8080`             | `consul` (healthcheck)            | Java 17  | API Gateway (Spring Cloud Gateway)  |
| `cms-main-system` | Built     | `./backend/cms-sys`       | `8081`             | `postgres-cms`, `redis`, `consul` | Go 1.24  | Main CMS backend                    |
| `lms-main-system` | Built     | `./backend/lms-sys`       | `8084`, `9090`     | `postgres-lms`, `consul`          | Go 1.24  | LMS backend service                 |
| `redis`           | Container | `redis:7-alpine`          | `6379`             | –                                 | –        | Redis for caching                   |
| `postgres-cms`    | Container | `postgres:15-alpine`      | `5432`             | –                                 | –        | PostgreSQL DB for CMS               |
| `postgres-lms`    | Container | `postgres:15-alpine`      | `5433` (host)      | –                                 | –        | PostgreSQL DB for LMS (optional)    |
| `email-service`   | Built     | `./backend/email-service` | `8082`             | `nats`, `consul`                  | Java 17  | Email service with SMTP + NATS      |
| `nats`            | Container | `nats:latest`             | `4222`, `8222`     | –                                 | –        | Message broker (JetStream enabled)  |
| `dozzle`          | Container | `amir20/dozzle:latest`    | `8083`             | –                                 | –        | Real-time Docker logs viewer        |

> ✅ `lms-main-system` was commented out in Compose file, but listed here assuming future use.

---

### 🌐 Networks

| Network Name           | Driver | Description               |
| ---------------------- | ------ | ------------------------- |
| `whole-system-network` | bridge | Shared network for system |

---

### 💾 Volumes

| Volume Name     | Purpose                                |
| --------------- | -------------------------------------- |
| `postgres-cms`  | Data storage for CMS PostgreSQL        |
| `postgres-lms`  | Data storage for LMS PostgreSQL        |
| `redis-data`    | Data storage for Redis                 |
| `consul-data`   | Persistent Consul data                 |
| `consul-config` | Consul configuration files             |
| `vault-data`    | Persistent data for Vault (future use) |

---
