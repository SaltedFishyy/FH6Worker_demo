# FH6 AI Tuning Assistant

Wails v2 + React TypeScript desktop app for FH6 Data Out telemetry, local rule events, tune profiles, and local Markdown tuning reports.

## Current Scope

- Listen for FH6 Data Out UDP packets on a selected local IPv4 address and port.
- Show compact FH6 Data Out target guidance below the network adapter selector.
- Accept only 324-byte packets and track valid, invalid, and parse-error counts.
- Parse core telemetry fields through a centralized versioned packet spec.
- Normalize speed, RPM load, driver inputs, wheel slip, tire temperature, and suspension travel.
- Display a realtime dashboard with driver inputs, wheel state, status, and 10 Hz recent trends.
- Detect local tuning events and show an event timeline.
- Store tune profiles, telemetry sessions, and detected events in local SQLite.
- Generate local Markdown tuning reports from a saved session and the active tune profile.
- Switch GUI language between English and Chinese.

## Development

```powershell
cd D:\FH6Worker\app
pnpm install --dir frontend
wails dev
```

## Build

```powershell
cd D:\FH6Worker\app
wails build
```

The built executable is written to `build\bin\app.exe`.

## UDP Simulator

Use the simulator when FH6 is not available:

```powershell
cd D:\FH6Worker\app
go run .\cmd\udp-sim -addr 127.0.0.1:5301
```

Start the listener in the app first, then run the simulator.

## LAN Setup

When the game runs on another machine, configure FH6 with the telemetry PC address, not the game PC address.

Example:

- Game PC: `192.168.0.109`
- Telemetry PC: `192.168.0.105`
- FH6 Data Out IP: `192.168.0.105`
- FH6 Data Out Port: `5301`
- App network adapter: `192.168.0.105` or `All interfaces (0.0.0.0)`

If packets still do not arrive, allow the app or UDP port `5301` through Windows Firewall on the telemetry PC.
