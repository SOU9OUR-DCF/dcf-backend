## Deployment

### Quick Start

1. Clone the repository
2. Copy the example environment file: `cp .env.example .env`
3. Modify the `.env` file with your desired configuration
4. Start the application: `docker-compose -f docker-compose.prod.yaml up -d`
5. Access the API at http://localhost:80
6. Access Swagger documentation at http://localhost:80/swagger/index.html

### Configuration

The application can be configured through environment variables in the `.env` file:

- `DB_USER`: PostgreSQL database username
- `DB_PASSWORD`: PostgreSQL database password
- `DB_NAME`: PostgreSQL database name
- `REDIS_PASSWORD`: Redis password
- `JWT_SECRET`: Secret key for JWT token generation

### Production Deployment

For production deployment, make sure to:

1. Change all default passwords in the `.env` file
2. Set a strong `JWT_SECRET` value
3. Consider using a reverse proxy like Nginx for SSL termination
4. Set up proper monitoring and logging