<div align="center">

# 🏙️ Vinhomes West Point · 3D GIS

### Explore the city in 3D. Manage the community in one place.
### Khám phá không gian 3D · Quản lý khu dân cư

[![Website](https://img.shields.io/badge/🌐_Live_website-Visit-0891b2?style=for-the-badge)](https://3d-geographic-information-system.vercel.app)
[![ArcGIS](https://img.shields.io/badge/ArcGIS_JS-4.30-0079c1?style=for-the-badge)](https://developers.arcgis.com/javascript/)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Languages](https://img.shields.io/badge/Languages-VI_%7C_EN-6366f1?style=for-the-badge)

**[🚀 Open website / Mở website](https://3d-geographic-information-system.vercel.app)**
&nbsp; · &nbsp;
**[🗺️ Explore 3D / Xem mô hình](https://3d-geographic-information-system.vercel.app/index.html)**

</div>

---

## ✨ About / Giới thiệu

An academic **3D Geographic Information Systems** project combining an interactive model of **Vinhomes West Point, Hanoi** with a property management portal. The scene represents towers **W1, W2, W3**, the podium, roads, landscaping and rooftop amenities.

Đồ án môn **Hệ thống thông tin địa lý 3 chiều**, kết hợp mô hình tương tác Vinhomes West Point với cổng quản lý tòa nhà, cư dân và vận hành. Kích thước và tọa độ công trình được mô phỏng phục vụ trực quan hóa.

## 🧭 Explore / Khám phá

| | English | Tiếng Việt |
|---|---|---|
| 🏙️ **3D scene** | Explore towers, building layers and construction phases | Khám phá tòa nhà, lớp kiến trúc và các giai đoạn thi công |
| 🗺️ **Map tools** | Switch basemaps, reset the view and inspect coordinates | Chọn bản đồ nền, đặt lại góc nhìn và xem tọa độ |
| 🏠 **Apartment status** | Filter apartments by tower, floor and occupancy status | Lọc căn hộ theo tòa, tầng và trạng thái sử dụng |
| 👥 **Management** | Role-based access for administrators, managers, staff and residents | Phân quyền quản trị viên, ban quản lý, nhân viên và cư dân |
| 📋 **Operations** | Residents, contracts, invoices, maintenance and notifications | Quản lý cư dân, hợp đồng, hóa đơn, bảo trì và thông báo |
| 💳 **VNPay sandbox** | QR checkout with server-side payment confirmation | Thanh toán QR thử nghiệm, xác nhận giao dịch qua máy chủ |
| 🌐 **VI / EN** | Switch language across pages and the 3D control panel | Chọn tiếng Việt hoặc Anh trên các trang và bảng điều khiển 3D |

> **Language / Ngôn ngữ:** Your preference is saved in the browser. Lựa chọn ngôn ngữ được lưu trong trình duyệt và dùng cho các trang tiếp theo.

## 🛠️ Technology / Công nghệ

```text
Frontend     HTML · CSS · JavaScript · ArcGIS Maps SDK 4.30
Backend      Go 1.25 · Gin · GORM
Database     PostgreSQL
Hosting      Vercel (frontend) · Railway (backend)
Integration  VNPay PAY 2.1.0 — Sandbox
```

## 🗂️ Project structure / Cấu trúc

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

## 🚀 Run locally / Chạy trên máy

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

### 3. Checks / Kiểm tra

```sh
node --test tests/*.test.mjs
cd backend
go test ./...
./scripts/test-integration.sh
```

The integration script starts and removes an isolated PostgreSQL test container.

## 💳 Payment setup / Cấu hình thanh toán

VNPay is currently configured for **sandbox testing**. Terminal approval and IPN registration are managed by VNPay. See [VNPay deployment notes](backend/VNPAY_DEPLOYMENT.md) for configuration and callback behavior. Keep credentials in environment variables.

## 📚 Academic context / Phạm vi đồ án

Built for learning 3D GIS visualization and property management workflows. The model is an approximate representation, rather than a survey or construction reference.

Mô hình phục vụ học tập và minh họa nghiệp vụ, không thay thế hồ sơ khảo sát hoặc thiết kế thi công.

---

<div align="center">

**[🏙️ Visit Vinhomes West Point 3D GIS](https://3d-geographic-information-system.vercel.app)**

</div>
