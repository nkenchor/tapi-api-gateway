# TAPI API Gateway Service

A high-performance, modular **API Gateway** built in Go and powered by **KrakenD**, designed to orchestrate and unify multiple microservices within the TAPI ecosystem.

---

## 🚀 Overview

The **TAPI API Gateway Service** acts as the central entry point to the TAPI microservices platform. It provides:

- Centralized routing and service aggregation
- Declarative configuration using `krakend.json`
- Secure, scalable API gateway operations
- Support for dynamic service discovery, request filtering, and transformation

This gateway facilitates smooth interaction between external clients (e.g., web, mobile apps) and backend services including:
- Identity & Authentication
- Cart & Checkout
- Payment Processing
- User & Organization Management
- Product Catalog
- Dashboard Analytics
- OTP & Notification services
- PolPred tidal prediction service

---

## 🧠 Key Features

| Feature                     | Description |
|----------------------------|-------------|
| **Built with Go**          | Performance and concurrency optimized |
| **Powered by KrakenD**     | Extensible and declarative API Gateway framework |
| **Microservice Integration** | Routes requests to over 10+ backend services |
| **Environment-based Config** | `.env` powered service URLs and secrets |
| **GitLab CI Ready**        | Includes `.gitlab-ci.yml` for automated builds/deployments |
| **Template Architecture**  | Customizable `api_gateway_template.json` for service-specific routing |

---

## 🧱 Architecture

```text
            ┌────────────┐
            │   Client   │
            └─────┬──────┘
                  │
         ┌────────▼────────┐
         │ TAPI API Gateway│
         │ (KrakenD + Go)  │
         └────────┬────────┘
                  │
      ┌───────────┼──────────────┬─────────────┬────────────┐
      ▼           ▼              ▼             ▼            ▼
 Identity   Cart Service   Payment Service   User Svc   Checkout Svc
 Service      (JSON)         (REST API)       ...         ...

```

 This architecture allows centralized control over:

- **Rate limiting**  
- **Caching**  
- **Header transformation**  
- **Request merging**  
- **Service composition**  

---

## ⚙️ Getting Started

### 1. Clone the Repo

```bash
git clone https://github.com/your-username/tapi-api-gateway-service.git
cd tapi-api-gateway-service
```

### 2. Set Up Environment

Create a .env file with service URLs and required secrets (sample provided):

```bash
TAPI_IDENTITY_SERVICE_URL=http://localhost:30301
TAPI_CART_SERVICE_URL=http://localhost:30302
```

### 3. Run the Gateway

Ensure Go is installed.

```bash
go run main.go
```
This launches the KrakenD gateway using the krakend.json config.

### 4. Access the Gateway

```bash
http://localhost:8080/

```

### Folder Structure

```text

.
├── main.go
├── go.mod / go.sum
├── .env
├── .gitlab-ci.yml
├── configuration/
│   ├── krakend.json                # Main gateway config
│   ├── template/                   # Template for auto-generating routes
│   └── services/                   # Per-service route config files
└── README.md

```

## 📦 Dependencies

- **Go** v1.21+
- **KrakenD** v2.x
- **GitLab CI** (optional)

---

## 🔐 Security Considerations

- API access should be protected using a layer of authentication middleware or KrakenD plugins.
- Services are modular and can be containerized for zero-trust environments.

---

## 🧪 Testing & Deployment

- CI/CD supported via GitLab pipelines (`.gitlab-ci.yml`)
- Can be dockerized or deployed to Kubernetes (support not included yet in repo)

---

## 🧠 Author & Vision

This gateway was developed as part of the **TAPI Microservices Suite** — a scalable, cloud-native platform supporting modular service domains, built with clean architecture principles, domain-driven design, and polyglot services (Go, Python, Julia, etc.).

---

## 📎 Related Repositories

| Service                | Language | Description                         |
|------------------------|----------|-------------------------------------|
| `tapi-cart-service`    | Go       | Handles shopping cart logic         |
| `tapi-payment-service` | Python   | Integrates with Stripe, GoCardless  |
| `tapi-polpred-service` | Julia    | Tidal prediction engine             |
| `tapi-identity-service`| Go       | Auth and identity provider          |
| `tapi-checkout-service`| Go       | Checkout and session manager        |
