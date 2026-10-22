# Lahiri API / Vedic Astrology Calculator

Go-based astronomical and Vedic astrology calculator utilizing the Swiss Ephemeris library (`swephgo`). It calculates precise sidereal planetary positions, Lahiri Ayanamsa, 27 Nakshatras with Padas, and Whole Sign house systems.

## Features

- **Sidereal Zodiac (Lahiri Ayanamsa):** Precise astronomical calculations using Swiss Ephemeris.
- **DMS Formatting:** High-precision degrees, arcminutes (lipda), and arcseconds (phlipda).
- **Vedic Nakshatras:** Complete 27 Nakshatras breakdown with Pada (1-4) and Vimshottari Dasha Lords.
- **House Systems:** Whole Sign house calculation and main chart angles (ASC, MC, DSC, IC).
- **Testable Structure:** Includes unit tests for core utilities and calculations.

## Prerequisites

Before running this project, ensure you have the Swiss Ephemeris C-library installed on your system, as `swephgo` binds to it.

- **macOS (via Homebrew):**
  ```bash
  brew install swisseph
  ```
- **Ubuntu/Debian:**
  ```bash
  sudo apt-get install libswe-dev
  ```

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/lahiri-api.git
   cd lahiri-api
   ```

2. Initialize and tidy Go modules:
   ```bash
   go mod tidy
   ```

## Usage

1. Configure your birth details, coordinates, and settings in `config.go`.
2. Run the program:
   ```bash
   go run .
   ```

## Running Tests

To run the unit tests:
```bash
go test -v
```

## License

This project is licensed under the [MIT License](LICENSE).