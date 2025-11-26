# Deployment Guide

This guide will walk you through deploying the Cinema Booking System

## Local Development

### What You'll Need
- Docker Desktop
- Go 1.21+
- Git

### Getting Started
1. Grab the code:
```bash
git clone https://github.com/gcode/cinema-booking.git
cd cinema-booking
```

2. Fire up the services:
```bash
docker-compose up --build
```

3. Check if everything's working:
```bash
curl http://localhost:3000/health
curl http://localhost:3000/api/docs
```

##  Deployment

### Step 1: Build Docker Images

Build production-ready Docker images for all services:

```bash
# Build all services
docker build -t cinema-booking/auth-service:latest ./auth-service
docker build -t cinema-booking/cinema-service:latest ./cinema-service
docker build -t cinema-booking/booking-service:latest ./booking-service
docker build -t cinema-booking/api-gateway:latest ./api-gateway

# Or use docker-compose for batch building
docker-compose -f docker-compose.prod.yml build
```

### Step 2: Push to Container Registry

Push images to your preferred container registry:

```bash
# Tag images for your registry
docker tag cinema-booking/auth-service:latest your-registry.com/cinema-booking/auth-service:latest
docker tag cinema-booking/cinema-service:latest your-registry.com/cinema-booking/cinema-service:latest
docker tag cinema-booking/booking-service:latest your-registry.com/cinema-booking/booking-service:latest
docker tag cinema-booking/api-gateway:latest your-registry.com/cinema-booking/api-gateway:latest

# Push to registry
docker push your-registry.com/cinema-booking/auth-service:latest
docker push your-registry.com/cinema-booking/cinema-service:latest
docker push your-registry.com/cinema-booking/booking-service:latest
docker push your-registry.com/cinema-booking/api-gateway:latest
```

### Step 3: Deploy to Docker-Compatible Cloud Platform

Choose any Docker-compatible cloud service:

#### Popular Options:
- **Railway** - Connect GitHub repo, auto-deploys Dockerfiles
- **Render** - Docker support with render.yaml configuration
- **Fly.io** - Excellent Go support with fly.toml
- **Google Cloud Run** - Serverless containers
- **AWS ECS/Fargate** - Managed container service
- **DigitalOcean App Platform** - Simple Docker deployment
- **Azure Container Instances** - Container hosting
- **Heroku** - Container registry support

#### Generic Deployment Steps:
1. Create a new application/service for each microservice
2. Configure environment variables (see below)
3. Deploy using your pushed Docker images
4. Set up networking between services
5. Configure load balancing and SSL

### Environment Variables

Set these environment variables in production:

```bash
# Database (Go format)
DATABASE_URL=postgres://username:password@host:5432/database?sslmode=disable

# Security
JWT_SECRET=your-super-secret-jwt-key-min-32-chars

# Service URLs (for internal communication)
AUTH_SERVICE_URL=http://auth-service:8080
CINEMA_SERVICE_URL=http://cinema-service:8080
BOOKING_SERVICE_URL=http://booking-service:8080

# Optional
PORT=8080
```

## Database Setup

### Managed PostgreSQL Options

Most cloud platforms offer managed PostgreSQL:
- **Railway** - Built-in PostgreSQL addon
- **Render** - PostgreSQL service
- **Google Cloud** - Cloud SQL
- **AWS** - RDS
- **DigitalOcean** - Managed Databases
- **Azure** - Database for PostgreSQL

Create a database instance and update your `DATABASE_URL` environment variable.

### Database Migration

The application automatically creates tables on startup. For production, consider using a migration tool:

```sql
-- Create database
CREATE DATABASE cinema_booking;

-- Run each service once to create tables
-- Or manually run the CREATE TABLE statements from each service's database.ts
```

## Monitoring and Logging

### Health Checks

Each service exposes health endpoints:
- `/health` - Basic health check
- `/metrics` - Application metrics (if implemented)

### Logging

Services log to stdout. In production, use log aggregation:

```bash
# Docker logs
docker logs <container-id>

```

## Performance Optimization

### Database Optimization

```sql
-- Add indexes for frequently queried columns
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_booking_code ON bookings(booking_code);
CREATE INDEX idx_seats_studio_available ON seats(studio_id, is_available);
```

### Caching Strategy

1. **Redis for sessions**: Store JWT tokens and user sessions
2. **Application caching**: Cache studio/seat data
3. **CDN**: Use CloudFront for static assets

## Troubleshooting

### Common Issues

1. **Service Discovery**: Ensure service URLs are correct
2. **Database Connections**: Check connection limits and pooling
3. **Memory Issues**: Monitor container memory usage
4. **Network Issues**: Verify security groups and network policies

### Debug Commands

```bash
# Check service logs
docker-compose logs auth-service

# Check database connections
docker exec -it postgres psql -U postgres -d cinema_booking

# Test service connectivity
curl http://localhost:3001/health
curl http://localhost:3002/health
curl http://localhost:3003/health
```