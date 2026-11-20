<div align="center">

# 🏙️ Vinhomes West Point · 3D GIS

### Explore the city in 3D. Manage the community in one place.

[![Website](https://img.shields.io/badge/🌐_Live_website-Visit-0891b2?style=for-the-badge)](https://3d-geographic-information-system.vercel.app)
[![ArcGIS](https://img.shields.io/badge/ArcGIS_JS-4.30-0079c1?style=for-the-badge)](https://developers.arcgis.com/javascript/)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Languages](https://img.shields.io/badge/Languages-VI_%7C_EN-6366f1?style=for-the-badge)

**[🚀 Open website](https://3d-geographic-information-system.vercel.app)**
&nbsp; · &nbsp;
**[🗺️ Explore 3D](https://3d-geographic-information-system.vercel.app/index.html)**

</div>

---

[![Vinhomes West Point 3D model with English controls](docs/images/3d-gis-preview.jpg)](https://3d-geographic-information-system.vercel.app/index.html)

## ✨ About the project

An academic **3D Geographic Information Systems** project combining an interactive model of **Vinhomes West Point, Hanoi** with a property management portal. The scene represents towers **W1, W2, W3**, the podium, roads, landscaping and rooftop amenities.


## 🔑 Try the resident portal

Explore the [live website](https://3d-geographic-information-system.vercel.app/pages/login.html) with the shared demo account:

| Email | Password |
|---|---|
| `resident@vinhomes.com` | `resident123` |

This account is for exploring the demo. Please avoid entering personal or confidential information in the shared account.

Use **Chat with admin** in the resident portal to ask questions or share feedback. Administrators have a dedicated inbox with user search, unread counts, conversation history and open/resolved filters. Messages are saved and delivered through an authenticated WebSocket connection, with automatic reconnection and polling fallback. Both interfaces support English and Vietnamese.

## 🧭 Features

| Feature | Description |
|---|---|
| 🏙️ **3D scene** | Explore towers, building layers and construction phases |
| 🗺️ **Map tools** | Switch basemaps, reset the view and inspect coordinates |
| 🏠 **Apartment status** | Filter apartments by tower, floor and occupancy status |
| 👥 **Management** | Role-based access for administrators, managers, staff and residents |
| 📋 **Operations** | Residents, contracts, invoices, maintenance and notifications |
| 💳 **VNPay sandbox** | QR checkout with server-side payment confirmation |
| 🌐 **VI / EN** | Switch language across pages and the 3D control panel |

> **Language preference:** English and Vietnamese are available throughout the interface. Your choice is saved in the browser and carried across pages.

## 🛠️ Technology

```text
Frontend     HTML · CSS · JavaScript · ArcGIS Maps SDK 4.30
Backend      Go 1.25 · Gin · GORM
Database     PostgreSQL
Hosting      Vercel (frontend) · Railway (backend)
Integration  VNPay PAY 2.1.0 — Sandbox
```

## 🗂️ Project structure

```text
├── index.html                  # Interactive 3D scene
├── pages/                      # Management and resident portal
├── css/                        # Portal, scene and language controls
├── js/
│   ├── app.js                  # ArcGIS scene and interaction
│   ├── i18n.js                 # English / Vietnamese localization
│   ├── api/                    # Backend API clients
│   ├── gis/                    # Apartment status overlay
│   └── towers/, podium/, ...   # Building geometry and scene modules
├── backend/
│   ├── cmd/                    # Server entry point
│   ├── config/                 # Environment configuration
│   ├── internal/               # Models, handlers and services
│   └── scripts/                # PostgreSQL integration tests
└── tests/                      # Frontend regression tests
```

## 🚀 Run locally

### 1. Backend

Requires **Go 1.25+** and **Docker** for the bundled PostgreSQL setup.

```sh
cd backend
cp config/.env.example config/.env
# Edit config/.env with your local settings.
./run.sh
```

The backend listens on `http://localhost:8080` by default. See [backend setup](backend/README.md) for details.

### 2. Frontend

From the repository root, serve the static files:

```sh
python3 -m http.server 5500
```

Open [local portal](http://localhost:5500/pages/login.html) or [local 3D model](http://localhost:5500/index.html). Localhost pages use the local backend automatically; hosted pages use Railway. To use another backend, configure `window.__API_BASE_URL` in `js/env-config.js`.

### 3. Run checks

```sh
node --test tests/*.test.mjs
cd backend
go test ./...
./scripts/test-integration.sh
```

The integration script starts and removes an isolated PostgreSQL test container.

## 💳 Payment setup

VNPay is currently configured for **sandbox testing**. Terminal approval and IPN registration are managed by VNPay. See [VNPay deployment notes](backend/VNPAY_DEPLOYMENT.md) for configuration and callback behavior. Keep credentials in environment variables.

## 🔎 Google Search Console

The public 3D model is listed in [sitemap.xml](https://3d-geographic-information-system.vercel.app/sitemap.xml), advertised by [robots.txt](https://3d-geographic-information-system.vercel.app/robots.txt). Login and portal pages send `X-Robots-Tag: noindex, follow`; authentication remains responsible for protecting private data.

The URL-prefix property is `https://3d-geographic-information-system.vercel.app/`. Keep the Google verification tags in both `index.html` and `pages/login.html`: the site root redirects to the login page, which Google checks during verification. Search Console verification and sitemap submission do not guarantee indexing or rankings.

## 📚 Academic context

Built for learning 3D GIS visualization and property management workflows. The model is an approximate representation, rather than a survey or construction reference.


---

<div align="center">

**[🏙️ Visit Vinhomes West Point 3D GIS](https://3d-geographic-information-system.vercel.app)**

</div>

### Chat delivery

Chat uses a single WebSocket connection per portal page. Message writes still use the authenticated HTTP API; WebSocket events immediately refresh the conversation and unread counts. Reconnecting clients reload saved history. The event hub currently runs in one backend instance; horizontal scaling requires a shared event bus (such as Redis Pub/Sub) before adding replicas.

### Image attachments on Railway

Residents and the building team can attach photos to maintenance requests and chat messages. JPEG and PNG uploads are limited to five images at 2 MiB each (16 megapixels maximum). Images are re-encoded as JPEG to remove embedded metadata. Completion reports can add five more photos to a maintenance request.

Image bytes are stored on the backend's Railway volume mounted at `/data`, under `/data/images`. PostgreSQL stores image metadata and attachment references. The backend reads `RAILWAY_VOLUME_MOUNT_PATH` automatically; set `UPLOAD_DIR` to an absolute persistent directory for local development or another deployment. Without either setting, uploads return unavailable instead of writing to ephemeral storage.

Images are fetched through authenticated endpoints with conversation/request access checks. Original maintenance photos saved as base64 remain readable. Storage is capped at 50 MiB per uploader and 400 MiB overall for the current small volume. The implementation requires one backend replica with the volume attached. Back up the volume together with PostgreSQL, as both are needed to restore attachments.
