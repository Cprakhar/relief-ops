<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h1 align="center">🚨 Relief Ops</h1>

  <p align="center">
    A microservices-based disaster management and relief coordination platform
    <br />
    <strong>Real-time disaster reporting • Resource tracking • Role-based access</strong>
    <br />
    <br />
    <a href="#features"><strong>Explore Features »</strong></a>
    <br />
    <br />
    <a href="https://github.com/Cprakhar/relief-ops/issues">Report Bug</a>
    ·
    <a href="https://github.com/Cprakhar/relief-ops/issues">Request Feature</a>
  </p>
</div>

<!-- BADGES -->
<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.25.1-00ADD8?style=flat&logo=go)
![Kubernetes](https://img.shields.io/badge/Kubernetes-1.34-326CE5?style=flat&logo=kubernetes)
![MongoDB](https://img.shields.io/badge/MongoDB-8.0-47A248?style=flat&logo=mongodb)
![Apache Kafka](https://img.shields.io/badge/Kafka-4.1.0-231F20?style=flat&logo=apache-kafka)
![Redis](https://img.shields.io/badge/Redis-8.2-DC382D?style=flat&logo=redis)

</div>

---

## 📋 Table of Contents

- [About](#about)
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Deployment](#deployment)
- [Contributing](#contributing)

---

## 🌟 About

**Relief Ops** is a cloud-native disaster management platform built with Go microservices. It enables communities and organizations to coordinate disaster response through real-time reporting, resource discovery, and admin oversight.

### Key Capabilities

- **Three User Roles**: Admin, Contributor, and Public users with distinct permissions
- **Disaster Workflow**: Report → Admin Review → Public Visibility
- **Resource Discovery**: Find nearby hospitals, shelters, fire stations, pharmacies, and police stations
- **Event-Driven**: Kafka-based async notifications to admins
- **Geospatial Queries**: MongoDB 2dsphere indexes for location-based searches
- **Production-Ready**: JWT auth, retry logic, horizontal scaling, Kubernetes deployment

---

## ✨ Features

### 👥 User Roles & Permissions

**Public Users**
- View approved disaster reports
- Search disasters by location
- View nearby emergency resources on map

**Contributors**
- All public user capabilities
- Report disasters with location and details
- Upload disaster information
- Wait for admin approval before public visibility

**Admins**
- All contributor capabilities
- Review and approve/reject disaster reports
- Receive email notifications for new disaster reports
- Manage disaster lifecycle

### 🚨 Disaster Management
- Geolocation-based disaster reporting
- Admin approval workflow
- Status tracking (pending, approved, rejected)
- Email notifications to admins via SendGrid
- Event-driven architecture with Kafka

### 🏥 Resource Discovery
- Find nearby emergency resources:
  - Hospitals
  - Police stations
  - Fire stations
  - Shelters
  - Pharmacies
- Geospatial radius search (e.g., "resources within 5km")
- Automatic data sync from OpenStreetMap via Overpass API
- Smart duplicate prevention by name + amenity type

### 🔐 Authentication & Security
- JWT-based stateless authentication
- Role-based access control (RBAC)
- Password hashing with bcrypt
- Secure cookie-based sessions
- Redis-backed token blacklist for logout

---

### Event Flow

```
Contributor Reports Disaster
         ↓
Disaster Service (gRPC)
         ↓
MongoDB (Store)
         ↓
Kafka Producer (Publish Event)
         ↓
User Service Consumer (Receive Event)
         ↓
SendGrid API (Send Email to Admins)
```

---

## 🛠️ Tech Stack

**Backend**
- Go 1.25.1 - Primary language
- gRPC - Inter-service communication
- Gin - HTTP framework
- Protocol Buffers - Service definitions

**Databases**
- MongoDB 6.0 - Document database with geospatial support
- Redis 7.0 - Session management and caching

**Message Queue**
- Apache Kafka 4.0.0 - Event streaming (KRaft mode, 3 brokers)

**Infrastructure**
- Kubernetes - Container orchestration
- Docker - Containerization
- Minikube - Local development

**External APIs**
- SendGrid - Email notifications
- Overpass API - OpenStreetMap resource data

---

## 🚀 Getting Started

### Prerequisites

```bash
# Required tools
go version          # Go 1.25.1+
docker --version    # Docker
kubectl version     # Kubernetes CLI
minikube version    # Minikube
protoc --version    # Protocol Buffers compiler
```

### Quick Start

1. **Clone and setup**
   ```bash
   git clone https://github.com/Cprakhar/relief-ops.git
   cd relief-ops
   go mod download
   make generate-proto
   ```

2. **Start local Kubernetes**
   ```bash
   minikube start --memory=8192 --cpus=4
   kubectl create namespace relief-ops
   ```

3. **Create secrets**
   ```bash
   # JWT Secret
   kubectl create secret generic jwt-secret \
     --from-literal=jwt-secret=$(openssl rand -base64 48) \
     -n relief-ops

   # MongoDB
   kubectl create secret generic mongo-secret \
     --from-literal=mongo-uri="mongodb://mongo-db:27017" \
     -n relief-ops

   # Redis
   kubectl create secret generic redis-secret \
     --from-literal=redis-password=$(openssl rand -base64 32) \
     -n relief-ops

   # SendGrid (replace with your API key)
   kubectl create secret generic sendgrid-secret \
     --from-literal=api-key="YOUR_SENDGRID_API_KEY" \
     -n relief-ops

   # Kafka brokers
   kubectl create secret generic kafka-config \
     --from-literal=kafka-brokers="apache-kafka-0.apache-kafka.relief-ops.svc.cluster.local:9092,apache-kafka-1.apache-kafka.relief-ops.svc.cluster.local:9092,apache-kafka-2.apache-kafka.relief-ops.svc.cluster.local:9092" \
     -n relief-ops
   ```

4. **Deploy infrastructure**
   ```bash
   kubectl apply -f k8s/base/mongo-db.yaml
   kubectl apply -f k8s/base/redis-db.yaml
   kubectl apply -f k8s/base/apache-kafka.yaml
   
   # Wait for readiness
   kubectl wait --for=condition=ready pod -l app=mongo-db -n relief-ops --timeout=300s
   kubectl wait --for=condition=ready pod -l app=redis-db -n relief-ops --timeout=300s
   kubectl wait --for=condition=ready pod -l app=apache-kafka -n relief-ops --timeout=300s
   ```

5. **Build and deploy services**
   ```bash
   make build-all
   kubectl apply -f k8s/base/user-service.yaml
   kubectl apply -f k8s/base/disaster-service.yaml
   kubectl apply -f k8s/base/resource-service.yaml
   kubectl apply -f k8s/base/api-gateway.yaml
   ```

6. **Verify deployment**
   ```bash
   kubectl get pods -n relief-ops
   kubectl get svc -n relief-ops
   ```

---

## 💻 API Endpoints

### Authentication

**Register User**
```bash
POST /auth/register
{
  "email": "user@example.com",
  "password": "SecurePass123",
  "name": "John Doe"
}
```

**Login**
```bash
POST /auth/login
{
  "email": "user@example.com",
  "password": "SecurePass123"
}
# Returns JWT token in cookie
```

**Logout**
```bash
POST /auth/logout
# Requires authenticated session
```

### Disasters

**Report Disaster** (Contributors only)
```bash
POST /disasters
{
  "title": "Earthquake in San Francisco",
  "description": "6.5 magnitude earthquake",
  "location": {
    "latitude": 37.7749,
    "longitude": -122.4194
  },
  "severity": 8,
  "tags": ["earthquake", "urgent"]
}
```

**Get All Disasters** (Public)
```bash
GET /disasters
```

**Get Disaster by ID** (Public)
```bash
GET /disasters/{id}
```

**Review Disaster** (Admins only)
```bash
POST /admin/review/{id}
{
  "approved": true
}
```

### Resources

**Get Nearby Resources** (Public)
```bash
GET /resources/nearby?lat=37.7749&lon=-122.4194&radius=5000&type=hospital
```

**Sync Resources from OpenStreetMap**
```bash
POST /resources/sync?lat=37.7749&lon=-122.4194&radius=10000
```

---

## ☁️ Deployment

### Microservices

| Service | Port | Purpose |
|---------|------|---------|
| API Gateway | 8080 | HTTP REST entry point |
| User Service | 9001 | Authentication & notifications |
| Disaster Service | 9002 | Disaster reporting & management |
| Resource Service | 9003 | Resource discovery & sync |

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `JWT_SECRET` | Secret for JWT signing | Yes |
| `MONGO_URI` | MongoDB connection string | Yes |
| `KAFKA_BROKERS` | Kafka broker addresses | Yes |
| `SENDGRID_API_KEY` | SendGrid API key | Yes |
| `REDIS_PASSWORD` | Redis password | Yes |

### Production Considerations

- Enable horizontal pod autoscaling
- Configure resource limits (CPU, memory)
- Set up persistent volumes for databases
- Use ingress for external access
- Enable TLS/SSL certificates
- Configure health checks and probes

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

<div align="center">

[⬆ Back to Top](#-relief-ops)

</div>