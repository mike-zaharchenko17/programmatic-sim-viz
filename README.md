# Programmatic Sim Viz

A real-time visualization of a simulated programmatic ad auction pipeline. The
backend generates synthetic OpenRTB-style bid requests, runs second-price
auctions, and streams results over WebSocket. The frontend renders the live
auction flow as a Sankey diagram alongside per-scope performance charts.

CS-GY 6313 - Information Visualization, Spring 2026.

## Repository layout

```
programmatic-sim-viz/
├── psv-generator/   # Go backend: auction simulation + WebSocket server
└── psv-svelte/      # SvelteKit frontend: D3 Sankey + exploratory charts
```

## Prerequisites

- Go 1.25+
- Node.js 20+ and npm
- (Optional) Docker, if you want to run the backend in a container

## Running locally

The frontend and backend run as two independent processes. Start the backend
first.

### 1. Backend (`psv-generator`)

```bash
cd psv-generator
go run ./cmd
```

The server listens on `:1323` by default. Override with the `PORT` env var:

```bash
PORT=8080 go run ./cmd
```

Endpoints:
- `GET /health` - liveness probe, returns `ok`
- `GET /ws` - WebSocket stream of `AuctionResult` JSON messages

The auction pipeline is lazy: it only runs while at least one WebSocket client
is connected, and shuts itself down 45 seconds after the last client
disconnects.

### 2. Frontend (`psv-svelte`)

In a second terminal:

```bash
cd psv-svelte
npm install
npm run dev
```

Vite serves the app at `http://localhost:5173`.

The frontend reads its WebSocket URL from `PUBLIC_WS_URL`. A default `.env`
pointing at `ws://localhost:1323/ws` is committed locally; copy
`.env.example` if you need to recreate it:

```bash
cp .env.example .env
```

If you change the backend port, update `PUBLIC_WS_URL` to match and restart
`npm run dev`.

### Optional: running the backend in Docker

```bash
cd psv-generator
docker build -t psv-generator .
docker run --rm -p 1323:1323 psv-generator
```

## Production builds

Backend:

```bash
cd psv-generator
go build -o server ./cmd
./server
```

Frontend (static SPA, output in `psv-svelte/build/`):

```bash
cd psv-svelte
npm run build
npm run preview   # to serve the build locally on :4173
```

## Deployment (Render)

The two services are deployed independently:

- `psv-generator` - Render Web Service, Docker runtime, root directory
  `psv-generator`. Render injects `PORT`; no other env vars required.
  Health check path: `/health`.
- `psv-svelte` - Render Static Site, root directory `psv-svelte`.
  Build command: `npm install && npm run build`. Publish directory: `build`.
  Required env var: `PUBLIC_WS_URL=wss://<psv-generator-host>/ws`.

`PUBLIC_WS_URL` is baked into the JS bundle at build time, so changing it on
Render requires a rebuild of the static site, not just a redeploy.
