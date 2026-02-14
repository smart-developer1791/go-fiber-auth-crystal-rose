# 🌹 Crystal Rose Garden Auth

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-v2.52-00ACD7?style=for-the-badge&logo=go&logoColor=white)](https://gofiber.io/)
[![Three.js](https://img.shields.io/badge/Three.js-r160-black?style=for-the-badge&logo=three.js&logoColor=white)](https://threejs.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-3.0-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)](https://tailwindcss.com/)
[![Alpine.js](https://img.shields.io/badge/Alpine.js-3.x-8BC0D0?style=for-the-badge&logo=alpine.js&logoColor=white)](https://alpinejs.dev/)
[![SQLite](https://img.shields.io/badge/SQLite-Pure_Go-003B57?style=for-the-badge&logo=sqlite&logoColor=white)](https://github.com/glebarez/sqlite)
[![Render](https://img.shields.io/badge/Render-Deployed-46E3B7?style=for-the-badge&logo=render&logoColor=white)](https://render.com/deploy)

> 💎 Immersive 3D crystal rose garden authentication with WebGL rendering, realistic glass materials, floating particles, and romantic atmosphere. Built with Go Fiber + Three.js + GORM + Alpine.js + Tailwind CSS.

## ✨ Features

### 3D Visual Experience
- 🌹 **Crystal Roses** — Procedurally generated roses with realistic glass materials
- 💎 **Physical Materials** — MeshPhysicalMaterial with transmission, clearcoat, IOR
- ✨ **Floating Particles** — Dewdrops and sparkles with additive blending
- 🌙 **Dynamic Lighting** — Moving point lights with colored glow
- 🎥 **Ambient Camera** — Subtle automatic camera movement
- 🌫️ **Atmospheric Fog** — Exponential fog for depth

### Authentication
- ✅ **Real-time Validation** — Async field checking
- 📊 **Crystal Strength Meter** — Creative password visualization
- 📱 **Phone Formatting** — Auto US format
- 🍪 **Session Auth** — Secure cookies
- 🔐 **bcrypt** — Password hashing
- 💾 **Pure Go SQLite** — No CGO

## 🎨 Technical Details

### Three.js Materials

```javascript
MeshPhysicalMaterial({
    transmission: 0.9,    // Glass transparency
    thickness: 1.5,       // Refraction depth
    ior: 2.4,            // Index of refraction (diamond-like)
    clearcoat: 1,        // Glossy surface
    metalness: 0.1,      // Subtle metallic
    roughness: 0.05      // Very smooth
})
```

### Rose Generation

Each rose consists of:
- 4 petal layers (5→7→9→11 petals)
- Procedural positioning with rotation
- Extruded geometry with bevel
- Crystal stem and leaves

## 🚀 Quick Start

Clone and run:

```bash
git clone https://github.com/smart-developer1791/go-fiber-auth-crystal-rose
cd go-fiber-auth-crystal-rose
```

Initialize and start:

```bash
go mod tidy
go run .
```

Open [http://localhost:3000](http://localhost:3000) 🌹

## 🔑 Demo Account

| Field | Value |
|-------|-------|
| Email | `rose@crystal.garden` |
| Password | `rose2024` |
| Phone | `+1 (214) 214-2024` |

## 🛠 Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.21+, Fiber v2, GORM |
| Database | SQLite (Pure Go) |
| 3D Engine | Three.js r160 |
| Frontend | Alpine.js 3.x, Tailwind CSS |
| Materials | MeshPhysicalMaterial (PBR) |
| Auth | bcrypt + cookie sessions |

## 📁 Structure

```text
go-fiber-auth-crystal-rose/
├── main.go              # Server & handlers
├── go.mod               # Dependencies
├── render.yaml          # Deploy config
├── .gitignore
├── README.md
└── templates/
    ├── login.html       # 3D login scene
    ├── register.html    # 3D register scene
    └── dashboard.html   # 3D dashboard
```

## 🎭 Color Palette

| Color | Hex | Usage |
|-------|-----|-------|
| Crystal Pink | `#ff6b9d` | Primary roses |
| Crystal Rose | `#c41e3a` | Dark roses, accents |
| Crystal Light | `#ffd6e0` | UI highlights |
| Deep Purple | `#1a0a2e` | Background |

## 🌐 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/login` | Login with 3D roses |
| POST | `/api/login` | Authenticate |
| GET | `/register` | Register scene |
| POST | `/api/register` | Create account |
| POST | `/api/validate` | Real-time validation |
| GET | `/dashboard` | 3D garden view |
| POST | `/logout` | End session |

## 💕 Valentine's Day

Created with love for Valentine's Day 2024. Each crystal rose represents eternal love — beautiful, precious, and timeless.

> *"Love is like a crystal — delicate, beautiful, and reflects light in unexpected ways."*

## ⚡ Performance Notes

- Three.js loaded via ES modules (no build required)
- Optimized particle count for mobile
- Responsive canvas sizing
- RequestAnimationFrame loop
- Efficient geometry reuse

---

## Deploy in 10 seconds

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)
