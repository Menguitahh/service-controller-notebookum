# 🚀 Service Controller NotebookUm (Go API Gateway)

![Go](https://img.shields.io/badge/Go-1.22+-blue.svg?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.9+-darkgreen.svg?style=flat-square&logo=go&logoColor=white)
![Consul](https://img.shields.io/badge/Consul-1.15+-red.svg?style=flat-square&logo=hashicorp&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7.0+-red.svg?style=flat-square&logo=redis&logoColor=white)

Este microservicio, desarrollado en **Go utilizando el framework Gin**, actúa como el **API Gateway y Orquestador de Flujos** central de la plataforma NotebookUm. Es la puerta de entrada única para los clientes y coordina la lógica distribuida hacia los servicios internos de usuarios, persistencia, extracción de PDF e inteligencia artificial.

---

## 📋 Responsabilidades

- **Enrutamiento y Proxificación:** Centralizar los puntos de acceso públicos y redirigir el tráfico de red de forma transparente a los microservicios downstream especializados.
- **Validación de Identidad:** Interceptar las peticiones a endpoints protegidos y validar las firmas de tokens **JWT** (evitando llamadas no autorizadas a los backend).
- **Correlación de Peticiones:** Generar y propagar el encabezado `X-Correlation-ID` en todos los flujos internos para consolidar la trazabilidad de logs (Distributed Tracing).
- **Seguridad CORS:** Centralizar la política de intercambio de recursos de origen cruzado para navegadores y aplicaciones web.

---

## ⚡ Características Clave

- **Alto Rendimiento en Go:** Ejecución rápida y concurrente mediante *goroutines* nativas e hilos virtuales ligeros para la integración con Consul.
- **Resolución por Consul:** Carga dinámica de las URLs de los microservicios aguas abajo desde Consul KV al iniciar, eliminando variables estáticas.
- **Auto-registro del Gateway:** Inscribe automáticamente la instancia del Gateway en Consul junto con etiquetas dinámicas de Traefik para balancear carga.
- **Respuestas RFC 9457:** Formato de error estandarizado de `Problem Details` para peticiones inválidas, rutas inexistentes o fallos de red.

---

## 🌐 Endpoints de la API (Públicos y Protegidos)

| Método | Ruta | Autenticación | Rol / downstream | Descripción |
| :--- | :--- | :---: | :--- | :--- |
| **GET** | `/health` | Ninguna | - | Retorna el estado general de salud del Gateway y sus circuitos. |
| **GET** | `/ready` | Ninguna | - | Indica si el servicio está listo para recibir peticiones de red. |
| **GET** | `/status/circuits` | Ninguna | - | Muestra el estado actual de los disyuntores de red controlados. |
| **POST** | `/api/v1/users` | Ninguna | -> `service-user` | Crea una nueva cuenta de usuario (Registrar). |
| **POST** | `/api/v1/users/login` | Ninguna | -> `service-user` | Inicio de sesión. Retorna el par de tokens JWT. |
| **POST** | `/api/v1/users/refresh` | Ninguna | -> `service-user` | Refresca el token de acceso expirado del usuario. |
| **GET** | `/api/v1/users/:id` | Ninguna | -> `service-user` | Obtiene el perfil público o información del usuario. |
| **POST** | `/api/v1/documento/upload`| **JWT Bearer** | -> `service-ai` | Sube archivo PDF. Inicia ingesta y generación de resumen. |
| **GET** | `/api/v1/documents/:id/status`| **JWT Bearer** | -> `service-extractor` | Polling del estado del trabajo de extracción del PDF. |
| **GET** | `/api/v1/summaries/:id` | **JWT Bearer** | -> `service-persistence` | Recupera el resumen consolidado de un documento. |
| **POST** | `/api/v1/summaries/document` | **JWT Bearer**| -> `service-persistence` | Dispara manualmente la creación de resumen para un documento. |

---

## ⚙️ Configuración centralizada (Consul KV)

El Gateway lee sus variables operativas desde el prefijo raíz `notebookum/controller` en Consul:

| Clave en Consul KV | Variable de Entorno Local | Valor por Defecto | Descripción |
| :--- | :--- | :--- | :--- |
| `port` | `PORT` | `5000` | Puerto HTTP donde escuchará el Gateway. |
| `extractor_url` | `EXTRACTOR_URL` | `http://extractor.universidad.localhost:5000` | URL del microservicio de extracción de texto (PDF). |
| `ai_url` | `AI_URL` | `http://ai.universidad.localhost:5000` | URL del microservicio de procesamiento de IA (LLM). |
| `persistence_url` | `PERSISTENCE_URL` | `http://persistence-java.universidad.localhost:8080` | URL del microservicio de persistencia de base de datos. |
| `user_service_url` | `USER_SERVICE_URL`| `http://users.universidad.localhost:5000` | URL del microservicio de gestión de usuarios. |
| `redis_host` | `REDIS_HOST` | `redis` | Servidor Redis utilizado para enrutamientos u optimizaciones. |
| `redis_port` | `REDIS_PORT` | `6379` | Puerto de conexión Redis. |
| `redis_password` | `REDIS_PASSWORD` | `""` | Contraseña de Redis. |
| `traefik_tags` | - | *(Lista predeterminada)* | Tags inyectados dinámicamente en Consul para configuración de Traefik. |

---

## 🚀 Despliegue y Ejecución

### Requisitos Previos

- Go SDK versión 1.22 o superior instalada.
- Instancia activa de Consul (`http://localhost:8500`).

### Ejecución Local

1. Compilar y arrancar la aplicación de Go directamente desde el punto de entrada:
   ```bash
   export CONSUL_URL=http://localhost:8500
   go run ./cmd/controller
   ```

2. Ejecutar las pruebas automatizadas del Gateway:
   ```bash
   go test ./...
   ```

### Despliegue en Docker

El microservicio se compila y se despliega como contenedor asociándolo a la red interna del clúster:

```bash
docker compose up --build -d
```
Al levantar, cargará la parametrización de Consul y expondrá el puerto `5000` balanceado por Traefik para gestionar todas las solicitudes entrantes del ecosistema NotebookUm.
