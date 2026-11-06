# Project Context: lahiri-api

## 1. Overview
`lahiri-api` เป็น Microservice ประสิทธิภาพสูงที่พัฒนาด้วยภาษา Go สำหรับคำนวณดาราศาสตร์และโหราศาสตร์พระเวท (Vedic Astrology / Sidereal Zodiac - Lahiri Ayanamsa) โดยเชื่อมต่อกับ C-library ของ Swiss Ephemeris (`swephgo`) ใช้เป็น Engine หลักในการส่งข้อมูลผลลัพธ์ดวงชะตาให้กับ Web Application (`ssr-web`), Postman และระบบแปลง Prompt เพื่อป้อนให้กับ LLM/AI สำหรับพยากรณ์

- **Language & Runtime:** Go 1.26+ (ต้องการ CGO ในการ Build ร่วมกับ Swiss Ephemeris C Library `libswe-dev`)
- **Default Port:** `8083` (กำหนดผ่าน Environment Variable `PORT`)
- **Environment Modes:** `APP_ENV=development` (Default) หรือ `APP_ENV=production`
- **Logging:** ใช้ Uber Zap ร่วมกับ `file-rotatelogs` จัดเก็บแบบ Daily Rotation ใน `logs/` เก็บย้อนหลัง 30 วัน

---

## 2. Architecture & Directory Structure

```text
lahiri-api/
├── internal/
│   ├── api/
│   │   ├── api.go         # HTTP Handlers (Calculate, Dasha, Prompt), Routing, Validation
│   │   ├── api_test.go    # Integration & Unit Tests สำหรับ API Endpoints
│   │   ├── dto.go         # Request / Response Structs พร้อม Validate() method
│   │   ├── prompt.go      # FormatRomanCompactPrompt() สร้าง Prompt Roman Ultra-Compact (Drishti, Orb, PD)
│   │   └── prompt_test.go # Unit Test สำหรับตรวจรูปแบบและฟังก์ชันสร้าง AI Prompt
│   └── calculator/
│       ├── config.go      # ข้อมูลราศี, นักษัตร, ดาวเจ้าเรือน, คำแปลไทย, รหัสย่อสากล
│       ├── dasha.go       # คำนวณวิมโสตรีทศา 120 ปี (Mahadasha, Antardasha, Pratyantardasha)
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
├── logs/                  # Daily Rotated Logs (e.g. app-YYYY-MM-DD.log, app.log symlink)
├── response/
│   └── response.json      # ตัวอย่าง Response Data
├── lahiri-api.postman_collection.json  # Postman Collection V2.1 (ครอบคลุมทั้ง 3 Endpoints)
├── logger.go              # Zap Logger + file-rotatelogs Daily Rotation config
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
3. **Dignity, Vargottama & Combustion:**
   - **มาตรฐานดาว (Dignity):** ตรวจสอบทั้งกระดาน D1 และ D9 ได้แก่ `Exalted (อุจจ์)`, `Moolatrikona (มูลไตรโกณ)`, `Own Sign (เกษตร)`, `Friend (มิตร)`, `Neutral (ปกติ/สมพล)`, `Enemy (ศัตรู)`, `Debilitated (นิจ)`
   - **วรโคตมนวางศ์ (Vargottama):** ตรวจสอบกรณีดาวสถิตราศีเดียวกันทั้งใน D1 และ D9 (`[Varg]`)
   - **ดาวดับ/อัสดงคต (Combust):** คำนวณระยะเชิงมุมกับดวงอาทิตย์ตามเกณฑ์องศาพระเวท (แยกกรณีดาวพักร์องศาแคบลง) และตรวจจับ **Deep Combust** (กุมอาทิตย์ไม่เกิน 5°)
4. **เกณฑ์สายตาของดาว (Drishti):**
   - **เกณฑ์พิเศษ (Special Aspects):** เสาร์ (`Sa`: 3, 7, 10), อังคาร (`Ma`: 4, 7, 8), พฤหัสบดี (`Ju`: 5, 7, 9)
   - **เกณฑ์เล็งข้ามภพ (Mutual Aspects):** ตรวจจับคู่ภพตรงข้าม 180° ที่มีดาวสถิตทั้งสองฝั่ง
5. **Vimshottari Dasha (วิมโสตรีทศา 120 ปี):**
   - คำนวณจากองศาของดวงจันทร์ (Moon's Nakshatra & Pada)
   - บอกทศาคงเหลือแรกเกิด (`balance_at_birth`)
   - คำนวณครบ 3 ลำดับชั้น: **มหาทศา (Mahadasha)**, **อันตรทศา (Antardasha)** และ **ปรัตยันตรทศา (Pratyantardasha / PD)** ณ ช่วงเวลาปัจจุบันอย่างแม่นยำ
6. **ชื่อมุมหลัก 4 จุด (Cardinals):**
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
คำนวณเฉพาะข้อมูลวิมโสตรีทศา 120 ปี มหาทศา (Mahadasha), อันตรทศา/อนุทศา (Antardasha), และปรัตยันตรทศา (Pratyantardasha)
- รองรับทั้ง Request Body แบบห่อใน `calculateReq` และแบบ Flat JSON
- สามารถส่ง `transit_year` และ `transit_month` เพิ่มเติมเพื่อระบุปีและเดือนจรที่ต้องการหาทศาได้ (ยึดตามตัวเลขวันเกิดเดิม และปรับเป็นวันสิ้นเดือนอัตโนมัติหากเดือนนั้นมีจำนวนวันน้อยกว่า)

### 4.3 `POST /api/v1/hora/transit`
คำนวณดาวจรประจำปีรายเดือนทั้ง 12 เดือน (`monthly_transits`) ณ วันเกิดและเวลาเดิมในปีจร (`transit_year`)
- กำหนดราศีลัคนาดวงกำเนิดเป็นเรือนที่ 1 (Whole Sign)
- แสดงดาวนวเคราะห์ 9 ดวง (`1-Sun` ถึง `9-Ketu`) พร้อมเรือนภพเทียบลัคนากำเนิด

### 4.4 `POST /api/v1/hora/transitmonth`
คำนวณดาวจรเฉพาะเดือนที่ระบุ (`transit_year` และ `transit_month`) คืนค่า `monthly_transit` ของเดือนนั้นๆ

### 4.5 `POST /api/v1/hora/prompt`
สร้างข้อความสรุปดวงชะตาแบบ **Roman Ultra-Compact** เพื่อป้อนเข้า System Prompt ของ LLM/AI โดยตรง ประหยัด Token สูงสุด (~75–100 Tokens) พร้อมระบุ Orb, Deep Combust, Drishti และ Pratyantardasha (PD)
- ค่าเริ่มต้นคืนเป็น JSON: `{"status": "success", "prompt": "..."}`
- หากส่ง Query Parameter `?format=raw` หรือ Header `Accept: text/plain` จะส่งกลับเป็นข้อความ Plain Text ทันที

**ตัวอย่างข้อความผลลัพธ์ (Roman Ultra-Compact):**
```text
[Vedic Sidereal(Lahiri)/WholeSign]
ASC:Aq 03°39'23"|D9:Sc|Nak:Dhanishta(4,Ma) | MC:Sc 10°54'11"(H10) DSC:Le 03°39'23"(H7) IC:Ta 10°54'11"(H4)
Planets(D1:Sign,DMS,Dig|House|Nak(Pada,Lord)|D9:Sign,Dig):
Su:Sg 07°49'03"(Fri)|H11|Mula(3,Ke)|D9:Ge(Neu)
Mo:Ge 06°43'29"(Fri)|H5|Ardra(1,Ra)|D9:Sg(Neu)
Me[C]:Sg 19°57'00"(Neu)|H11|P.Ashadha(2,Ve)|D9:Vi(Exa)
Ve:Sc 13°01'15"(Neu)|H10|Anuradha(3,Sa)|D9:Li(Own)
Ma:Pi 22°07'37"(Fri)|H2|Revati(2,Me)|D9:Cp(Exa)
Ju[R]:Ta 03°43'04"(Ene)|H4|Krittika(3,Su)|D9:Aq(Neu)
Sa[C]:Sg 10°50'27"(Neu)|H11|Mula(4,Ke)|D9:Cn(Ene)
Ra[R,Varg]:Aq 14°33'09"(Own)|H1|Shatabhisha(3,Ra)|D9:Aq(Own)
Ke[Varg]:Le 14°33'09"(Neu)|H7|P.Phalguni(1,Ve)|D9:Le(Neu)
Lords(H:Sign,Lord->InH): 1:Aq,Sa->H11 2:Pi,Ju->H4 3:Ar,Ma->H2 4:Ta,Ve->H10 5:Ge,Me->H11 6:Cn,Mo->H5 7:Le,Su->H11 8:Vi,Me->H11 9:Li,Ve->H10 10:Sc,Ma->H2 11:Sg,Ju->H4 12:Cp,Sa->H11
Conj: H1(ASC+Ra[orb 10°54']) H11(Su+Me[orb 12°08'], Su+Sa[orb 3°01',DeepCombust])
Drishti(Special/Mutual): Sa->H1(Ra),H5(Mo),H8 | Ma->H5(Mo),H8,H9 | Ju->H8,H10(Ve),H12 | Mutual: H1(Ra)<->H7(Ke), H4(Ju)<->H10(Ve), H5(Mo)<->H11(Su,Me,Sa)
Dasha(Moon:Ardra): MD:Sa(2022-11-25..2041-11-24)/AD:Me(2025-11-27..2028-08-06)/PD:Ve
```

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
