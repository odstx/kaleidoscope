# FrontendApp (Swift)

SwiftUI frontend application for Kaleidoscope.

## Features

- **SwiftUI** - Declarative UI framework
- **Cross-platform** - macOS and iOS support
- **i18n** - Internationalization (English/Chinese)
- **MVVM Architecture** - Clean separation of concerns
- **API Integration** - Full backend integration

## Requirements

- Swift 5.9+
- macOS 14.0+ / iOS 17.0+

## Project Structure

```
swift/
├── Package.swift              # Swift Package Manager config
├── Sources/
│   ├── FrontendApp.swift      # App entry point
│   ├── Models/
│   │   └── User.swift         # Data models
│   ├── Services/
│   │   ├── APIService.swift   # API client
│   │   └── AuthState.swift    # Authentication state
│   ├── Utils/
│   │   └── Localization.swift # i18n support
│   └── Views/
│       ├── ContentView.swift  # Main view with routing
│       ├── LoginView.swift    # Login page
│       ├── RegisterView.swift # Registration page
│       ├── DashboardView.swift# Dashboard page
│       ├── ProfileView.swift  # User profile page
│       ├── NavbarView.swift   # Navigation bar
│       └── FooterView.swift   # Footer with system info
└── Tests/
    └── FrontendAppTests/
        ├── LocalizationTests.swift
        └── ModelTests.swift
```

## Development

### Build

```bash
swift build
```

### Run

```bash
swift run
```

### Test

```bash
swift test
```

## Views

| View | Description |
|------|-------------|
| LoginView | User login with email/password |
| RegisterView | User registration |
| DashboardView | Main dashboard after login |
| ProfileView | User profile management |
| NavbarView | Navigation bar with language switch |
| FooterView | Footer with system info |

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/users/login` | POST | User login |
| `/api/v1/users/register` | POST | User registration |
| `/api/v1/users/info` | GET | Get user info |
| `/api/v1/system/info` | GET | Get system info |

## Localization

Supported languages:
- English (en)
- Chinese (zh)

Language files are defined in `Sources/Utils/Localization.swift`.

## Architecture

```
┌─────────────────────────────────────────┐
│                  Views                  │
│  (SwiftUI Views - Login, Dashboard...)  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│               Services                  │
│     (AuthState, APIService)             │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│                Models                   │
│   (User, LoginRequest, SystemInfo...)   │
└─────────────────────────────────────────┘
```

## License

MIT
