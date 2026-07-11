# FH6Worker

FH6Worker is a telemetry analysis and vehicle tuning toolkit designed for Forza Horizon 6.

It receives real-time vehicle telemetry through the game's Data Out UDP interface, analyzes vehicle behavior, detects common handling issues, and helps users create, test, and improve tuning setups.

The repository includes a Windows desktop application, a WeChat Mini Program, tuning calculation models, documentation, vehicle data, and supporting development tools.

## Features

### Desktop Application

The desktop application is built with Wails, Go, React, and TypeScript. It provides:

- Real-time FH6 Data Out telemetry monitoring
- Speed, RPM, throttle, brake, steering, and gear visualization
- Tire temperature and wheel-slip analysis
- Suspension travel monitoring
- Understeer, oversteer, traction, and braking issue detection
- Quick, professional, and tire-focused diagnostic modes
- Tune profile management
- Telemetry session recording and replay
- Track baseline capture and comparison
- Before-and-after tuning evaluation
- Local SQLite storage
- Markdown tuning report generation
- English and Chinese interface support

### Static Tuning Calculator

FH6Worker can generate a starting tune from vehicle information such as:

- Performance Index
- Vehicle weight and front weight distribution
- Power and torque
- Drivetrain
- Tire compound and tire width
- Gear count
- Intended use, including Road and Drift

Generated recommendations may include:

- Tire pressure
- Camber, toe, and caster
- Springs and anti-roll bars
- Rebound and bump damping
- Ride-height and aero adjustment tiers
- Brake balance and pressure
- Differential settings
- Final drive and individual gear ratios

These values are intended as a baseline. Track testing and driver feedback are still required for final tuning.

### WeChat Mini Program

The `weChatApp` directory contains a lightweight tuning assistant for WeChat with:

- Quick vehicle tuning calculations
- Saved tuning profiles
- Tuning sharing
- Recommended vehicles
- Upgrade guidance
- CloudBase cloud functions
- Local fallback data

## Project Structure

```text
FH6Worker/
|-- app/          Windows desktop application
|-- weChatApp/    WeChat Mini Program and cloud functions
|-- docs/         Tuning formulas and development documentation
|-- img/          FH6 vehicle images
|-- tools/        Data and asset preparation tools
`-- demo.ahk      AutoHotkey test automation script
```

## Running the Windows Application

A prebuilt Windows executable is available as `app.exe` in the repository.

Run it from PowerShell:

```powershell
.\app.exe
```

## FH6 Data Out Configuration

To receive live telemetry:

1. Open the FH6 settings.
2. Enable the Data Out option.
3. Set the Data Out IP address to the computer running FH6Worker.
4. Set the Data Out port to `5301`.
5. Open FH6Worker.
6. Select the matching local network interface.
7. Set the listener port to `5301`.
8. Start telemetry capture and begin driving.

When the game and FH6Worker run on the same computer, try:

```text
IP address: 127.0.0.1
Port: 5301
```

When they run on different computers, use the local network IP address of the computer running FH6Worker. You may also need to allow UDP port `5301` through Windows Firewall.

## Desktop Development

### Requirements

- Windows
- Go
- Node.js
- pnpm
- Wails v2
- Microsoft Edge WebView2 Runtime

### Install Frontend Dependencies

```powershell
cd app\frontend
pnpm install
```

### Start Development Mode

```powershell
cd app
wails dev
```

### Build

```powershell
cd app
wails build
```

The compiled executable will be written to:

```text
app/build/bin/app.exe
```

### UDP Simulator

When the game is unavailable, the included simulator can send test telemetry packets:

```powershell
cd app
go run .\cmd\udp-sim -addr 127.0.0.1:5301
```

Start the telemetry listener in FH6Worker before running the simulator.

## WeChat Mini Program Development

1. Install WeChat Developer Tools.
2. Import the `weChatApp` directory.
3. Configure the CloudBase environment ID in `weChatApp/miniprogram/config.js`.
4. Deploy the `calculateTune` and `shareTune` cloud functions.

Private WeChat project configuration is intentionally excluded from this repository.

## Local Data and Privacy

The desktop application stores tune profiles, telemetry sessions, events, and analysis results locally using SQLite.

The following local or generated files are excluded from version control:

- SQLite databases and telemetry recordings
- Build outputs and dependency folders
- Compiler caches and temporary files
- Private WeChat configuration
- Generated executables
- Downloaded third-party tool archives

Do not commit personal telemetry data, credentials, API keys, or private CloudBase configuration.

## Important Notice

FH6Worker is an independent community project and is not affiliated with, endorsed by, or sponsored by Microsoft, Xbox Game Studios, Playground Games, or the Forza franchise.

Forza, Forza Horizon, and related names and assets are trademarks of their respective owners.

Vehicle images and game-related data are provided for development and informational purposes. Review the applicable licenses and platform policies before redistributing assets.

## Disclaimer

Generated tuning values are starting points rather than guaranteed optimal setups.

Vehicle behavior can vary depending on upgrades, tire compound, driving style, controller settings, track conditions, and game updates. Always validate recommendations through in-game testing.

## Contributing

Contributions are welcome. When submitting changes:

1. Keep generated files and local data out of commits.
2. Add or update tests when changing tuning calculations.
3. Document changes to formulas or parameter units.
4. Verify both English and Chinese interface text where applicable.
5. Describe the reasoning behind tuning-model changes.

## License

No open-source license has been selected yet.

Until a license is added, the source code remains under the copyright of its respective author and may not automatically be copied, modified, or redistributed.
