# Project Context: lahiri-api

## 1. Overview
`lahiri-api` เป็น Microservice ที่พัฒนาด้วยภาษา Go สำหรับคำนวณดาราศาสตร์และโหราศาสตร์พระเวท (Vedic Astrology / Sidereal Zodiac - Lahiri Ayanamsa) โดยเชื่อมต่อกับ C-library ของ Swiss Ephemeris (`swephgo`) เพื่อใช้เป็น Engine หลักส่งข้อมูลผลลัพธ์ให้กับ Web Application (`ssr-web`), Postman และ AI Prompt

- **Language & Runtime:** Go 1.26+ (ต้องการ CGO ในการ Build ร่วมกับ Swiss Ephemeris C Library `libswe-dev`)
- **Default Port:** `8083` (กำหนดผ่าน Environment Variable `PORT`)
- **Environment Modes:** `APP_ENV=development` (Default) หรือ `APP_ENV=production`

---

## 2. Architecture & Directory Structure

```text
lahiri-api/
├── internal/
│   ├── api/
│   │   ├── api.go         # HTTP Handlers (Calculate, Dasha, Prompt), Routing, Validation
│   │   ├── api_test.go    # Integration & Unit Tests สำหรับ API Endpoints
│   │   ├── dto.go         # Request / Response Structs พร้อม Validate() method
│   │   └── prompt.go      # FormatRomanCompactPrompt() สร้าง Prompt Roman Ultra-Compact
│   └── calculator/
│       ├── config.go      # ข้อมูลราศี, นักษัตร, ดาวเจ้าเรือน, คำแปลไทย, รหัสย่อ
│       ├── dasha.go       # คำนวณวิมโสตรีทศา 120 ปี (Vimshottari Dasha), Mahadasha/Antardasha
│       ├── dignity.go     # คำนวณมาตรฐานดาว D1/D9 (อุจจ์/เกษตร/มิตร/นิจ), ดาวดับ (Combust)
│       ├── helpers.go     # Concurrency Mutex (sweMu), Julian Day, Ayanamsa wrapper
│       ├── house.go       # คำนวณลัคนาและมุมภพ (ASC, MC, DSC, IC)
│       ├── nakshatra.go   # คำนวณ 27 นักษัตรฤกษ์ และ บาทฤกษ์ (Pada 1-4)
│       ├── navamsa.go     # คำนวณนวางศ์จักร (D9 Navamsa Chart)
│       ├── planet.go      # คำนวณตำแหน่งดาวเคราะห์, ราหู, เกตุ
│       ├── types.go       # Data types และ Data models ภายใน calculator
│       └── utils.go       # ฟังก์ชันคณิตศาสตร์มุม, DMS, ระยะเชิงมุม
├── bin/
│   └── lahiri-api         # Compiled binary (Ignored by Git)
├── response/
│   └── response.json      # ตัวอย่าง Response Data
├── lahiri-api.postman_collection.json  # Postman Collection V2.1
├── logger.go              # Zap Logger configuration
├── main.go                # Entrypoint, Graceful Shutdown, Server Timeout config
├── run.sh                 # Bash script สำหรับ Build/Run binary
├── test.sh                # Automated test & coverage report script
├── go.mod / go.sum
└── GEMINI.md              # Project Context & Agent Rules
```

---

## 3. Core Capabilities & Astrological Rules

1. **Sidereal Lahiri Ayanamsa:**
   - ใช้วงรอบจักรราศีแบบนิรายนะ (Sidereal Zodiac) อิงค่าอายนางศะแบบลาหิรี (Lahiri) ผ่าน Swiss Ephemeris
2. **Whole Sign House System (ภพละ 1 ราศีเต็ม):**
   - ภพที่ 1 (ตนุ) เริ่มต้นที่ราศีที่ลัคนาสถิต (0° - 30°)
   - ไล่ครบทั้ง 12 ภพ (ตนุ, กดุมพะ, สหัชชะ, พันธุ, ปุตตะ, อริ, ปัตนิ, มรณะ, ศุภะ, กัมมะ, ลาภะ, วินาศ)
   - สรุปความสัมพันธ์ `lord_placed_in_house` (ดาวเจ้าเรือนเกษตรไปสถิตอยู่ในภพใด)
3. **Dignity & Combustion (มาตรฐานดาว และ ดาวดับ):**
   - **มาตรฐานดาว (Dignity):** ตรวจสอบทั้งกระดาน D1 และ D9 ได้แก่ `Exalted (อุจจ์)`, `Moolatrikona (มูลไตรโกณ)`, `Own Sign (เกษตร)`, `Friend (มิตร)`, `Neutral (ปกติ/สมพล)`, `Enemy (ศัตรู)`, `Debilitated (นิจ)`
   - **วรโคตมนวางศ์ (Vargottama):** ตรวจสอบกรณีดาวสถิตราศีเดียวกันทั้งใน D1 และ D9 (`[Varg]`)
   - **ดาวดับ/อัสดงคต (Combust):** คำนวณระยะเชิงมุมกับดวงอาทิตย์ตามเกณฑ์องศาพระเวท (แยกกรณีดาวพักร์องศาแคบลง)
4. **Vimshottari Dasha (วิมโสตรีทศา 120 ปี):**
   - คำนวณจากองศาของดวงจันทร์ (Moon's Nakshatra & Pada)
   - บอกทศาคงเหลือแรกเกิด (`balance_at_birth`)
   - คำนวณมหาทศา (Mahadasha) และ อันตรทศา (Antardasha) ในช่วงเวลาปัจจุบันอย่างแม่นยำ
5. **ชื่อมุมหลัก 4 จุด (Cardinals):**
   - Ascendant (ASC) = **ลัคนา**
   - Midheaven (MC) = **มัธยลัคนา**
   - Descendant (DSC) = **อัสตลัคนา**
   - Imum Coeli (IC) = **ปาตาลลัคนา**

---

## 4. API Endpoints

### 4.1 `POST /api/v1/hora/calculate`
คำนวณข้อมูลดวงชะตาฉบับสมบูรณ์ (Metadata, Planets, HousesAndAngles, WholeSignHouses, VimshottariDasha)

- **Request Body:**
  ```json
  {
    "calculateReq": {
      "year": 1988,          // รองรับปี พ.ศ. (> 2400) โดยจะแปลงเป็น ค.ศ. ให้อัตโนมัติ
      "month": 12,
      "day": 23,
      "hour_thai": 10.35,     // เวลาไทย (เช่น 10:21 น. = 10.35)
      "time_zone": 7.0,
      "latitude": 13.7563,
      "longitude": 100.5018,
      "house_system": "W"     // Default "W" (Whole Sign)
    }
  }
  ```

### 4.2 `POST /api/v1/hora/dasha`
คำนวณเฉพาะข้อมูลวิมโสตรีทศา 120 ปี และทศาปัจจุบัน (ใช้ Request Body แบบเดียวกับ `/calculate`)

### 4.3 `POST /api/v1/hora/prompt`
สร้างข้อความสรุปดวงชะตาแบบ **Roman Ultra-Compact** สำหรับนำไปป้อนเข้า System Prompt ของ LLM/AI เพื่อประหยัด Token สูงสุด (~75-100 Tokens)
- ค่าเริ่มต้นคืนเป็น JSON: `{"status": "success", "prompt": "..."}`
- หากส่ง Query Parameter `?format=raw` หรือ Header `Accept: text/plain` จะส่งกลับเป็นข้อความ Plain Text ทันที

---

## 5. Build & Execution Rules

1. **WSL / Linux Environment:**
   - เนื่องจาก `swephgo` มี CGO Dependency (`libswe-dev`), การ Build และ Run จะต้องทำผ่านสภาพแวดล้อมที่มี GCC/Clang เช่น **WSL (Ubuntu/Debian)**
2. **การรันสคริปต์:**
   - Build Binary: `./run.sh --build-only`
   - รันเซิร์ฟเวอร์จาก Binary: `./run.sh --no-build`
   - Build ใหม่พร้อมรัน: `./run.sh`
   - รัน Unit Test ทั้งหมด: `go test -v ./...` หรือ `./test.sh`
3. **Concurrency Safety:**
   - การเรียกใช้ฟังก์ชัน Swiss Ephemeris (`swephgo`) ต้องถูกล็อกผ่าน `sweMu` (Mutex ใน `helpers.go`) เสมอ เนื่องจาก C Library มี Shared Global State
