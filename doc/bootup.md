# **Bootup Guide**

This document explains how to set up and run the project, which includes Golang, PostgreSQL, Loki, Promtail, and Grafana. All services are handled via Docker Compose.

---

## **Prerequisites**

### **1. Install Docker and Docker Compose**
- [Docker Installation Guide](https://docs.docker.com/get-docker/)
- [Docker Compose Installation Guide](https://docs.docker.com/compose/install/)

### **2. Update Docker Daemon Configuration**
If you encounter issues with the Docker setup, add the following configuration to your Docker daemon file (`/etc/docker/daemon.json` on Linux):

```json
{
  "dns": ["10.202.10.202", "10.202.10.102"],
  "insecure-registries": ["https://docker.arvancloud.ir"],
  "registry-mirrors": ["https://docker.arvancloud.ir"]
}
```

#### **Steps to Update Docker Configuration**:
1. Open the `daemon.json` file:
   ```bash
   sudo nano /etc/docker/daemon.json
   ```
2. Add the configuration above.
3. Restart the Docker service:
   ```bash
   sudo systemctl restart docker
   ```

---

## **Project Setup**

### **1. Clone the Repository**
Clone the project repository to your local system:
```bash
git clone <repository-url>
cd <repository-name>
```

### **2. Start the Services**
To build and start the project, use:
```bash
docker-compose up --build
```

This command:
- Builds Docker images for the project.
- Spins up containers for **Golang**, **PostgreSQL**, **Loki**, **Promtail**, and **Grafana**.

### **3. Verify Services**
Once the command completes, verify that all services are running:

- **PostgreSQL**: Accessible on port `5432`.
- **Grafana**: Accessible at [http://localhost:3000](http://localhost:3000).
- **Loki**: Log aggregation is available at [http://localhost:3100](http://localhost:3100).

---

## **Service Overview**

### **1. Golang**
The main application is written in Go. It connects to PostgreSQL as the database and uses Loki and Promtail for logging.

### **2. PostgreSQL**
PostgreSQL serves as the primary database for the project. Ensure the database credentials in `docker-compose.yml` match the ones your application expects.

### **3. Loki and Promtail**
Loki is used for log aggregation, and Promtail acts as an agent to send logs to Loki.

### **4. Grafana**
Grafana provides a web-based UI for visualizing logs and metrics. You can configure dashboards to monitor the application.

---

## **Troubleshooting**

### **1. Docker Configuration Issues**
If you encounter issues like image pulls failing or logs not appearing in Loki:
- Ensure your `daemon.json` is updated (see the **Update Docker Daemon Configuration** section).
- Restart Docker:
  ```bash
  sudo systemctl restart docker
  ```

### **2. Database Connection Errors**
If the application cannot connect to PostgreSQL:
- Ensure the database service is running:
  ```bash
  docker ps | grep postgres
  ```
- Verify database credentials in the application configuration.

### **3. Grafana Dashboards Not Loading**
- Check that Grafana is running on port `3000`.
- Log in to Grafana (default username: `admin`, password: `admin`) and verify your data sources and dashboards.

---

## **Customizing the Setup**

### **1. Environment Variables**
Modify environment variables in the `.env` file (if provided) or directly in the `docker-compose.yml` file to suit your local setup.

### **2. Logs**
To customize Loki logging, update the `loki-config.yaml` or the `log-driver` options in the `docker-compose.yml` file.

---

## **Stopping the Services**
To stop all running containers:
```bash
docker-compose down
```

This stops all services and removes containers, networks, and volumes created by Docker Compose.

## on boarding
now you can start your on-boarding with reading these docs:

- [architecture](./architecture/architecture.png)
- [ERD](./architecture/ERD.jpeg)