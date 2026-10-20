# 🏔️ APU — Control parental transparente

> *Apu*: en la cosmovisión andina, el espíritu protector de las montañas que cuida a la comunidad.

**APU** es una aplicación de control parental para Android. Se instala en el celular de la hija o del hijo y permite a madres, padres o tutores ver, desde un panel web, información de uso, ubicación y (en fases posteriores) actividad de mensajería, con límites configurables.

---

## 📑 Tabla de contenidos

1. [Principios del proyecto](#-principios-del-proyecto)
2. [Funcionalidades y roadmap](#-funcionalidades-y-roadmap)
3. [Arquitectura](#-arquitectura)
4. [Stack tecnológico](#-stack-tecnológico)
5. [Estructura del repositorio](#-estructura-del-repositorio)
6. [Modelo de datos](#-modelo-de-datos)
7. [API (resumen)](#-api-resumen)
8. [Flujo de emparejamiento](#-flujo-de-emparejamiento)
9. [Monitoreo de mensajes](#-monitoreo-de-mensajes-fase-4)
10. [Seguridad y privacidad](#-seguridad-y-privacidad)
11. [Instalación y desarrollo](#-instalación-y-desarrollo)
12. [Distribución (APK directo)](#-distribución-apk-directo)
13. [Permisos de Android](#-permisos-de-android)
14. [Consideraciones legales y éticas](#-consideraciones-legales-y-éticas)
15. [Contribuir](#-contribuir)
16. [Licencia](#-licencia)

---

## 🧭 Principios del proyecto

| Principio | Qué significa en la práctica |
|---|---|
| **Transparencia** | Icono que se puede cambiar, no hay sospechas de que el sistema esta funcionando. **Existe modo oculto y posble camuflaje.** |
| **Conocimiento del menor** | El primer arranque muestra una pantalla clara con lo que se monitorea. El menor puede ver en su propia app qué datos se comparten. |
| **Minimización de datos** | Se recolecta solo lo necesario para cada función activada. Nada se recolecta "por si acaso". |
| **Control del tutor, con límites** | El tutor decide qué módulos activar; el menor ve cuáles están activos. |
| **Seguridad primero** | Cifrado en tránsito y en reposo, tokens ligados al dispositivo, retención limitada. Son datos de una menor. |
| **Uso exclusivo familiar** | APU es para tutores legales y sus hijos. No está diseñada ni se debe usar para vigilar a otros adultos. |

---

## 🚀 Funcionalidades y roadmap

### Fase 1 — MVP: emparejamiento y ubicación
- [ ] Registro y login de tutores
- [ ] Emparejamiento por código/QR entre panel y celular
- [ ] Ubicación en tiempo real (foreground service)
- [ ] Historial de ubicaciones con retención configurable
- [ ] Zonas seguras (geofencing): casa, colegio, universidad
- [ ] Botón **SOS** del menor → notificación push al tutor con última ubicación

### Fase 2 — Tiempo de uso
- [ ] Estadísticas de uso por aplicación (`UsageStatsManager`)
- [ ] Tiempo total de pantalla diario y semanal
- [ ] Reportes en el panel (gráficos por día/semana)
- [ ] Estado del dispositivo: batería, conectividad

### Fase 3 — Límites y filtrado
- [ ] Límites de tiempo por app y por categoría
- [ ] Horarios (ej. "sin celular de 22:00 a 06:00")
- [ ] Bloqueo remoto de apps (comando vía FCM)
- [ ] Filtro web por DNS mediante `VpnService` local
- [ ] Solicitudes del menor: "pedir más tiempo" → aprobación del tutor

### Fase 4 — Mensajes y chats *(opt-in, transparente)*
- [ ] Lectura de **contenido de notificaciones** de apps de mensajería
- [ ] Alertas por palabras o patrones de riesgo configurables
- [ ] Indicador visible en el celular del menor cuando este módulo está activo
- [ ] Ver [Monitoreo de mensajes](#-monitoreo-de-mensajes-fase-4)

### Fase 5 — Madurez
- [ ] Cumplimiento de políticas de Google Play (si se publica)
- [ ] Soporte de múltiples hijos y múltiples tutores por familia
- [ ] Auditoría de seguridad externa
- [ ] Internacionalización (es / en / qu / ay)

---

## 🏗️ Arquitectura

```
┌─────────────────────┐        HTTPS / TLS         ┌──────────────────────────┐
│  APU Hija (Android) │ ─────────────────────────► │      API Backend (Go)    │
│  Kotlin             │ ◄───────────────────────── │  REST + WebSocket        │
│  - Ubicación        │      comandos (FCM)        │  - Auth (JWT + refresh)  │
│  - Uso de apps      │                            │  - Reglas y eventos      │
│  - Filtro DNS (VPN) │                            │  - Notificaciones        │
│  - Notificaciones   │                            └────────────┬─────────────┘
└─────────────────────┘                                         │
                                                                ▼
┌─────────────────────┐        HTTPS / TLS         ┌──────────────────────────┐
│  Panel de tutores   │ ─────────────────────────► │      PostgreSQL          │
│  Next.js + React    │                            │  (datos sensibles        │
│  shadcn/ui + Zustand│                            │   cifrados a nivel campo)│
└─────────────────────┘                            └──────────────────────────┘
```

**Flujo general:**
1. La app hija recolecta datos y los envía por lotes al backend (WorkManager para ahorrar batería).
2. El backend valida, almacena y evalúa reglas (zonas seguras, límites, alertas).
3. El panel consulta la API y recibe actualizaciones en tiempo real por WebSocket.
4. El tutor envía comandos (bloquear, pedir ubicación) → backend → **FCM** → app hija.

---

## 🧰 Stack tecnológico

| Capa | Tecnología |
|---|---|
| **App hija** | Kotlin, Jetpack Compose, WorkManager, Room (cola local), Retrofit/OkHttp |
| **Backend** | Go, Gin o Chi, GORM, JWT, WebSocket (`gorilla/websocket`) |
| **Base de datos** | PostgreSQL |
| **Panel web** | Next.js, React, TypeScript, Tailwind CSS, shadcn/ui, Zustand |
| **Push / comandos** | Firebase Cloud Messaging (FCM) |
| **Mapas** | Leaflet + OpenStreetMap (sin dependencia de Google Maps) |
| **Infra / tooling** | Docker, Docker Compose, GitHub Actions; scripts de apoyo en Python |

---

## 📁 Estructura del repositorio

```
apu/
├── README.md
├── LICENSE
├── docker-compose.yml
├── docs/
│   ├── arquitectura.md
│   ├── modelo-amenazas.md        # Threat model (STRIDE)
│   ├── privacidad.md             # Política de datos
│   └── api/                      # OpenAPI / Swagger
├── backend/                      # Go
│   ├── cmd/apu-server/
│   ├── internal/
│   │   ├── auth/
│   │   ├── device/
│   │   ├── location/
│   │   ├── usage/
│   │   ├── rules/
│   │   ├── messages/
│   │   ├── push/
│   │   └── storage/
│   ├── migrations/
│   └── go.mod
├── dashboard/                    # Next.js
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   ├── lib/
│   │   └── stores/               # Zustand
│   └── package.json
├── android/                      # Kotlin
│   ├── app/
│   │   └── src/main/java/bo/apu/child/
│   └── build.gradle.kts
└── tools/                        # Python: scripts auxiliares
```

---

## 🗄️ Modelo de datos

Esquema inicial (simplificado):

```sql
-- Tutores (padres/madres/tutores legales)
CREATE TABLE guardians (
    id            UUID PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,           -- argon2id
    created_at    TIMESTAMPTZ DEFAULT now()
);

-- Familias (permite varios tutores y varios hijos)
CREATE TABLE families (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE family_guardians (
    family_id   UUID REFERENCES families(id),
    guardian_id UUID REFERENCES guardians(id),
    PRIMARY KEY (family_id, guardian_id)
);

-- Dispositivos de los menores
CREATE TABLE devices (
    id            UUID PRIMARY KEY,
    family_id     UUID REFERENCES families(id),
    child_alias   TEXT NOT NULL,
    fcm_token     TEXT,
    public_key    BYTEA NOT NULL,          -- clave del dispositivo
    paired_at     TIMESTAMPTZ DEFAULT now(),
    last_seen_at  TIMESTAMPTZ
);

-- Módulos activos (lo que el menor puede ver que está activo)
CREATE TABLE device_modules (
    device_id  UUID REFERENCES devices(id),
    module     TEXT NOT NULL,              -- location | usage | web_filter | messages
    enabled    BOOLEAN NOT NULL DEFAULT false,
    enabled_at TIMESTAMPTZ,
    PRIMARY KEY (device_id, module)
);

CREATE TABLE locations (
    id          BIGSERIAL PRIMARY KEY,
    device_id   UUID REFERENCES devices(id),
    lat         DOUBLE PRECISION NOT NULL,
    lng         DOUBLE PRECISION NOT NULL,
    accuracy_m  REAL,
    recorded_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE app_usage (
    id          BIGSERIAL PRIMARY KEY,
    device_id   UUID REFERENCES devices(id),
    package     TEXT NOT NULL,
    seconds     INTEGER NOT NULL,
    day         DATE NOT NULL,
    UNIQUE (device_id, package, day)
);

CREATE TABLE rules (
    id         UUID PRIMARY KEY,
    device_id  UUID REFERENCES devices(id),
    type       TEXT NOT NULL,              -- app_limit | schedule | geofence | web_block
    payload    JSONB NOT NULL,
    active     BOOLEAN DEFAULT true
);

CREATE TABLE audit_log (                   -- quién consultó qué y cuándo
    id          BIGSERIAL PRIMARY KEY,
    guardian_id UUID REFERENCES guardians(id),
    device_id   UUID REFERENCES devices(id),
    action      TEXT NOT NULL,
    at          TIMESTAMPTZ DEFAULT now()
);
```

---

## 🔌 API (resumen)

Base: `https://<host>/api/v1`

| Método | Ruta | Descripción |
|---|---|---|
| `POST` | `/auth/register` | Registro de tutor |
| `POST` | `/auth/login` | Login → access + refresh token |
| `POST` | `/pairing/code` | Genera código de emparejamiento (TTL corto) |
| `POST` | `/pairing/confirm` | La app hija confirma el código y registra su clave pública |
| `GET`  | `/devices` | Lista dispositivos de la familia |
| `PUT`  | `/devices/:id/modules` | Activa/desactiva módulos |
| `POST` | `/devices/:id/locations` | La app hija envía ubicaciones (lote) |
| `GET`  | `/devices/:id/locations` | Historial de ubicaciones |
| `POST` | `/devices/:id/usage` | La app hija envía uso por app |
| `GET`  | `/devices/:id/usage` | Reporte de uso |
| `GET/POST/DELETE` | `/devices/:id/rules` | CRUD de reglas |
| `POST` | `/devices/:id/commands` | Envía comando (bloquear, pedir ubicación) |
| `POST` | `/devices/:id/sos` | Evento SOS desde la app hija |
| `WS`   | `/ws` | Actualizaciones en tiempo real para el panel |

---

## 🔗 Flujo de emparejamiento

1. El tutor inicia sesión en el panel y pulsa **"Añadir dispositivo"**.
2. El backend genera un código de un solo uso (ej. 8 caracteres) con **TTL de 5 minutos**, mostrado como texto y QR.
3. En el celular de la hija se abre APU y se escanea el QR o se escribe el código.
4. La app muestra **la pantalla de consentimiento**: qué módulos se activarán y qué datos se comparten. El menor (con su tutor presente) la acepta.
5. La app genera un par de claves en el **Android Keystore**, envía la clave pública y recibe un token ligado al dispositivo.
6. Aparece la notificación permanente de APU y el dispositivo queda vinculado.

---

## 💬 Monitoreo de mensajes (Fase 4)

Esta es la función más sensible del proyecto y se trata con reglas propias.

**Cómo funciona técnicamente:**
- Se usa `NotificationListenerService`, que lee el contenido de las **notificaciones** que muestran las apps de mensajería (WhatsApp, Telegram, Messenger, SMS, etc.).
- Apps como WhatsApp usan cifrado de extremo a extremo; APU **no rompe ese cifrado**, solo accede a lo que el propio sistema muestra en pantalla como notificación.
- Limitaciones: no se ven mensajes que el menor lee con la app abierta sin notificación, ni los borrados antes de llegar a la bandeja.

**Reglas de diseño obligatorias:**
1. **Opt-in explícito**: el módulo `messages` está desactivado por defecto y se activa por separado, con su propia pantalla de consentimiento.
2. **Indicador visible**: mientras está activo, la notificación permanente de APU lo indica ("Monitoreo de mensajes activo").
3. **Sin modo sigiloso**: no se permite ocultar el icono, disfrazar la app ni desactivar la notificación.
4. **Alertas antes que lectura completa** (modo recomendado): el tutor configura categorías de riesgo (acoso, grooming, autolesión, drogas, extorsión) y recibe **alertas con contexto mínimo**, no un espejo de toda la conversación.
5. **Cifrado extremo a tutor**: el contenido se cifra en el dispositivo con la clave pública del tutor; el servidor solo almacena texto cifrado.
6. **Retención corta**: por defecto 7 días, configurable, con borrado automático.
7. **Registro de auditoría**: cada acceso del tutor a mensajes queda registrado.

---

## 🔐 Seguridad y privacidad

**Comunicación**
- TLS 1.2+ obligatorio, con *certificate pinning* en la app Android.
- Tokens de acceso de vida corta + refresh tokens rotativos; revocación por dispositivo.

**Autenticación**
- Contraseñas con **argon2id**. MFA (TOTP) opcional para tutores.
- Claves de dispositivo en **Android Keystore**; las peticiones de la app hija se firman.

**Datos**
- Cifrado en reposo (disco) y **cifrado a nivel de campo** para ubicaciones y mensajes.
- Retención limitada y borrado automático configurable.
- Minimización: solo se recolecta lo que el módulo activo requiere.

**Backend**
- Validación estricta de entrada, *rate limiting*, consultas parametrizadas (GORM).
- Aislamiento por familia: toda consulta se filtra por `family_id` (previene IDOR).
- Secretos fuera del repositorio (variables de entorno / gestor de secretos).

**Proceso**
- [ ] Modelo de amenazas (STRIDE) en `docs/modelo-amenazas.md`
- [ ] Análisis estático (`gosec`, `staticcheck`, Detekt) en CI
- [ ] Escaneo de dependencias (`govulncheck`, `npm audit`, Dependabot)
- [ ] Pentest antes de cualquier uso fuera de pruebas internas
- [ ] Política de divulgación responsable (`SECURITY.md`)

---

## 🛠️ Instalación y desarrollo

### Requisitos
- Go 1.22+
- Node.js 20+ y pnpm
- PostgreSQL 15+
- Android Studio (Koala o superior), JDK 17
- Docker y Docker Compose (opcional)
- Proyecto de Firebase (para FCM)

### 1. Clonar

```bash
git clone https://github.com/VillaAlexTor/apu.git
cd apu
```

### 2. Backend (Go)

```bash
cd backend
cp .env.example .env        # editar credenciales y claves
go mod download
go run ./cmd/apu-server
```

Variables principales (`.env`):

```env
APU_ENV=development
APU_PORT=8080
DATABASE_URL=postgres://apu:apu@localhost:5432/apu?sslmode=disable
JWT_SECRET=cambia-esto
FCM_CREDENTIALS_FILE=./secrets/fcm.json
FIELD_ENCRYPTION_KEY=base64-32-bytes
```

### 3. Panel (Next.js)

```bash
cd dashboard
pnpm install
cp .env.example .env.local   # NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
pnpm dev
```

### 4. App Android

1. Abrir la carpeta `android/` en Android Studio.
2. Colocar `google-services.json` (de Firebase) en `android/app/`.
3. Configurar `API_BASE_URL` en `local.properties`.
4. Ejecutar en un dispositivo o emulador.

### 5. Todo con Docker

```bash
docker compose up --build
```

---

## 📦 Distribución (APK directo)

En esta etapa APU se distribuye **solo como APK para pruebas**, instalado en dispositivos propios de la familia.

```bash
cd android
./gradlew assembleRelease
# → app/build/outputs/apk/release/app-release.apk
```

**Instalación en el celular de la hija:**
1. Transferir el APK al dispositivo.
2. Permitir "instalar apps de orígenes desconocidos" solo para el gestor de archivos usado.
3. Instalar y completar el flujo de consentimiento y permisos.
4. Verificar la firma del APK (SHA-256 publicado en cada *release*).

> ⚠️ **Nota sobre Google Play:** distribuir por Play implicará cumplir sus políticas sobre apps de monitoreo, `AccessibilityService`, `NotificationListenerService`, ubicación en segundo plano y la declaración de uso de permisos sensibles. Esto está planificado para la Fase 5.

---

## 📱 Permisos de Android

| Permiso / API | Para qué se usa | Módulo |
|---|---|---|
| `ACCESS_FINE_LOCATION` / `ACCESS_BACKGROUND_LOCATION` | Ubicación en tiempo real y geofencing | Ubicación |
| `FOREGROUND_SERVICE` (+ tipo `location`) | Servicio persistente con notificación visible | Todos |
| `PACKAGE_USAGE_STATS` | Tiempo de uso por app | Uso |
| `BIND_VPN_SERVICE` | Filtro DNS local | Filtro web |
| `BIND_ACCESSIBILITY_SERVICE` | Detectar app en primer plano para aplicar bloqueos | Límites |
| `BIND_NOTIFICATION_LISTENER_SERVICE` | Leer notificaciones de mensajería | Mensajes |
| `POST_NOTIFICATIONS` | Notificación permanente y avisos | Todos |
| `INTERNET`, `ACCESS_NETWORK_STATE` | Sincronización | Todos |

Cada permiso se solicita **solo cuando se activa su módulo**, con una explicación en lenguaje claro.

---

## ⚖️ Consideraciones legales y éticas

- **Uso previsto:** madres, padres o tutores legales de **menores de edad**, con conocimiento del menor.
- **Prohibido:** instalar APU en dispositivos de otros adultos, parejas, exparejas o terceros. Eso constituye vigilancia no consentida y es delito en muchas jurisdicciones.
- **Privacidad del menor:** en Bolivia, el Código Niña, Niño y Adolescente (Ley 548) reconoce derechos de privacidad e intimidad de los menores. Conviene revisar con un abogado cómo encaja el monitoreo de mensajes antes de cualquier uso más allá de pruebas.
- **Proporcionalidad:** el monitoreo debe adecuarse a la edad. Un menor de 8 años y uno de 16 no necesitan el mismo nivel de supervisión.
- **Acompañar, no solo vigilar:** APU busca complementar el diálogo familiar, no sustituirlo.
- **Datos personales:** al tratar datos de menores se deben aplicar principios de licitud, finalidad, minimización y seguridad.

Este proyecto se ofrece sin garantías; quien lo use es responsable del cumplimiento de las leyes aplicables.

---

## 🤝 Contribuir

1. Haz un *fork* del repositorio.
2. Crea una rama: `git checkout -b feat/mi-mejora`.
3. Haz *commits* siguiendo [Conventional Commits](https://www.conventionalcommits.org/).
4. Abre un *Pull Request* describiendo el cambio.

**No se aceptarán contribuciones** que añadan modo sigiloso, ocultamiento del icono, evasión de la notificación permanente o cualquier función diseñada para que el menor no sepa que está siendo monitoreado.

Para reportar vulnerabilidades, ver `SECURITY.md`.

---

## 📄 Licencia

Por definir (se sugiere AGPL-3.0 o MIT). Ver `LICENSE`.

---

<p align="center">
  Hecho en La Paz, Bolivia 🇧🇴 — <strong>APU</strong>: cuidar sin esconderse.
</p>