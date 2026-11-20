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
├── .github/
│   └── workflows/
│       └── ci.yml         # GitHub Actions CI (Vet, Race Test, Coverage & SonarCloud Scan)
├── internal/
│   ├── api/
│   │   ├── api.go         # HTTP Handlers (Controller / Presentation Layer), Validation, Routing
│   │   ├── api_test.go    # Integration & Unit Tests สำหรับ API Endpoints
│   │   ├── dto.go         # Request / Response Struct Type Aliases พร้อม Validate() method
│   │   ├── gexf.go        # Format GEXF 1.2 XML Exporter สำหรับนำเข้าซอฟต์แวร์ Gephi
│   │   ├── gexf_test.go   # Integration & Unit Tests สำหรับ Endpoint ส่งออก GEXF
│   │   ├── karma_csv.go   # Format Nodes/Edges CSV และ ZIP Archive สำหรับ Gephi Data Laboratory
│   │   ├── prompt.go      # Format AI Prompts (Presenter) Roman Ultra-Compact (Drishti, Orb, PD)
│   │   └── prompt_test.go # Unit Test สำหรับตรวจรูปแบบและฟังก์ชันสร้าง AI Prompt
│   └── calculator/
│       ├── config.go      # ข้อมูลราศี, นักษัตร, ดาวเจ้าเรือน, คำแปลไทย, รหัสย่อสากล
│       ├── dasha.go       # คำนวณวิมโสตรีทศา 120 ปี (Mahadasha, Antardasha, Pratyantardasha)
│       ├── dignity.go     # คำนวณมาตรฐานดาว D1/D9 (อุจจ์/เกษตร/มิตร/นิจ), ดาวดับ (Combust)
│       ├── gephi.go       # คำนวณโครงข่ายโหราศาสตร์พระเวท (Nodes/Bala, Yuti, Drishti, Parivartana, Graha Maitri)
│       ├── gephi_test.go  # Unit Tests สำหรับตรรกะกราฟ Gephi
│       ├── gephi_houses.go      # คำนวณโครงข่ายภพเรือน Bipartite (9 ดาว + 12 ภพ, House Lord & Placed In)
│       ├── gephi_houses_test.go # Unit Tests สำหรับตรรกะกราฟภพเรือน Gephi
│       ├── gephi_karma.go       # คำนวณผังกรรม Topology of Karma (True Nodes, D1+D9, Axis/Taskmaster/Dispositor)
│       ├── gephi_karma_test.go  # Unit Tests สำหรับตรรกะผังกรรม Gephi
│       ├── gephi_nakshatra.go       # คำนวณโครงข่ายนักษัตร Hub Bipartite (Active Nakshatra Hubs + Point -> Nakshatra)
│       ├── gephi_nakshatra_test.go  # Unit Tests สำหรับตรรกะกราฟนักษัตร Gephi
│       ├── helpers.go     # Concurrency Mutex (sweMu), Julian Day, Ayanamsa wrapper
│       ├── house.go       # คำนวณลัคนาและมุมภพ (ASC, MC, DSC, IC)
│       ├── nakshatra.go   # คำนวณ 27 นักษัตรฤกษ์ และ บาทฤกษ์ (Pada 1-4)
│       ├── navamsa.go     # คำนวณนวางศ์จักร (D9 Navamsa Chart)
│       ├── planet.go      # คำนวณตำแหน่งดาวเคราะห์, ราหู, เกตุ
│       ├── service.go     # Application / Domain Service Layer (Orchestration สำหรับทุก Use Case)
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
├── run.sh                 # Bash script สำหรับ Verify, Test + Coverage, Build และ Run binary
├── sonar-project.properties  # SonarCloud Configuration & Coverage Report Path
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

### 4.5 `POST /api/v1/hora/transitdasha`
คำนวณดาวจรในเดือนที่ระบุ (`transit_year` และ `transit_month`) พร้อมข้อมูลทศาของเจ้าชะตา ณ เวลาดาวจรนั้น (`dasha`: มหาทศา, อันตรทศา/อนุทศา, และปรัตยันตรทศา) รวมอยู่ใน `monthly_transit`

### 4.6 `POST /api/v1/hora/transitdasha_prompt`
สร้างข้อความสรุปดาวจรในเดือนและทศา/อนุทศา ณ เวลาดาวจร แบบ **Roman Ultra-Compact** พร้อมระบุความเชื่อมโยงเรือนภพกำเนิดสู่เรือนภพจร (`[nH..->tH..]`), `TransitConj`, และ `TransitDrishti` (รองรับ `?format=raw` เพื่อรับ Plain Text ทันที)

### 4.7 `POST /api/v1/hora/prompt`
สร้างข้อความสรุปดวงชะตาแบบ **Roman Ultra-Compact** เพื่อป้อนเข้า System Prompt ของ LLM/AI โดยตรง ประหยัด Token สูงสุด (~75–100 Tokens) พร้อมระบุวันเวลาเกิด (Birth Date/Time, TimeZone), Orb, Deep Combust, Drishti และ Pratyantardasha (PD)
- ค่าเริ่มต้นคืนเป็น JSON: `{"status": "success", "birth_date": "1988-12-23", "birth_time": "10:21", "time_zone": 7, "prompt": "..."}`
- หากส่ง Query Parameter `?format=raw` หรือ Header `Accept: text/plain` จะส่งกลับเป็นข้อความ Plain Text ทันที

**ตัวอย่างข้อความผลลัพธ์ (Roman Ultra-Compact):**
```text
[Vedic Sidereal(Lahiri)/WholeSign | Birth:1988-12-23 10:21 (TZ+7)]
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

### 4.8 `POST /api/v1/hora/gephi`
ส่งออกข้อมูลโครงสร้างเครือข่ายโหราศาสตร์พระเวทสำหรับ 2 คนขึ้นไป (เช่น Triad Synastry 3 คน = 30 Nodes) เป็นไฟล์ **GEXF 1.3 XML** (`application/gexf+xml`, `attachment; filename="vedic-synastry.gexf"`), รูปแบบ **JSON** ผ่าน `?format=json`, หรือไฟล์ **ZIP สำหรับ Cosmograph (run.cosmograph.app)** และ Gephi Data Laboratory ผ่าน `?format=zip` หรือ `?format=cosmograph` (พร้อม `?format=nodes_csv` และ `?format=links_csv`)
- **Cosmograph (run.cosmograph.app) Ready:** ดาวน์โหลดไฟล์ `vedic-synastry.zip` ซึ่งบรรจุ `nodes.csv` (มี `id`, `label`, `person`, `body`, `rasi`, `dignity`, `bala`, `size`, `color`, `house`) และ `links.csv` (มี `source`, `target`, `weight`, `relation`, `relation_thai`, `edge_type`, `color`) ลากวางลงใน Cosmograph เพื่อรัน GPU Force-directed 2D/3D Graph ได้ทันที
- **P0: มาตรฐาน Edge Direction & GEXF 1.3 Mixed Graph:** กำหนดกราฟเป็น `defaultedgetype="mixed"` เพื่อรองรับทั้งเส้นแบบสองทาง (`type="undirected"` สำหรับ `Yuti`, `Parivartana`, `Nakshatra Compatibility`) และเส้นทิศทางเดียว (`type="directed"` สำหรับ `Graha Drishti`) แก้ไขปัญหา Gephi SEVERE Error `Edge type doesn't fit with default set (DIRECTED), edge is ignored` ทำให้ Gephi โหลดเส้นกุมทั้งหมดได้ครบ 100%
- **P0: การกุมกันในเรือนราศี D1 (Yuti in Rasi):** ระบุ `relation_thai="กุม (Yuti)"`, `rasi`, `rasi_name`, `rasi_thai` และ Label เช่น `Yuti in Cancer / กุมในกรกฎ (1.11°)` บนเส้นเชื่อม ให้สามารถตรวจสอบและฟิลเตอร์ตามเรือนราศีได้ทันทีใน Gephi Data Laboratory (เน้น D1 ราศีจักรล้วน)
- **P0: Birth Metadata ในโหนด:** ตรวจสอบตำแหน่งดาวย้อนหลังด้วย `birth_date`, `birth_time`, `time_zone`, `latitude`, `birth_longitude`, `julian_day`, `ayanamsa`
- **P0: Scoring Formula รายเส้น:** ระบุ `score`, `base_weight`, `scoring_formula` อย่างชัดเจนและตรวจสอบซ้ำได้
- **P1: Nakshatra Compatibility:** โหนดมีคุณสมบัติ `nakshatra_gana` (Deva/Manushya/Rakshasa), `nakshatra_yoni` (14 สัตว์), `nakshatra_nadi` (Adi/Madhya/Antya) และมีเส้นเชื่อมข้ามดวงระหว่างดวงจันทร์ (`Nakshatra Compatibility`) ตามเกณฑ์อัษฐกูฏ 36 คะแนน พร้อม `tara_name`, `tara_score`, `gana_compat`, `yoni_compat`, `nadi_compat`
- **P1: Graha Drishti Details:** เพิ่ม `drishti_rule`, `is_mutual`, และ `aspect_angle` (90°, 120°, 180°, 210°, 240°)
- **P2: Vimshottari Dasha Active Timing:** โหนดระบุดาวที่กำลังเสวยอายุ MD/AD/PD ด้วย `is_dasha_active`, `dasha_role`, `dasha_period`
- **P2: Pairwise Synastry Metrics:** สรุปสถิติความสัมพันธ์รายคู่ (`PairwiseSynastryMetric`: Harmony Ratio, Ashtakoota 36 คะแนน และคำวินิจฉัย) ฝังในคำอธิบาย GEXF และเข้าถึงได้โดยตรงผ่าน `?format=json`

### 4.9 `POST /api/v1/hora/gephi_houses`
ส่งออกข้อมูลโครงสร้างเครือข่ายภพเรือนโหราศาสตร์พระเวท (Bipartite House-Planet Network) สำหรับ 2 คนขึ้นไป เป็นไฟล์ **GEXF 1.2 XML** (`application/gexf+xml`, `attachment; filename="vedic-houses-synastry.gexf"`) เปิดในซอฟต์แวร์ Gephi ได้ทันที
- **Nodes (คนละ 21 โหนด):** 9 ดาวเคราะห์ (`<PersonID>_Sun` .. `<PersonID>_Ketu`) และ 12 ภพเรือน (`<PersonID>_H1` .. `<PersonID>_H12`) แยกด้วย `node_kind` (`Planet`/`House`), `house`, `house_name_thai`, `house_group`, `sign_lord`, `rasi`, `dignity`, `bala`
- **Edges (Directed House Lord Flow — Intra-chart & Inter-person):**
  - `House Lord`: `House -> LordPlanet` ภายในดวงเจ้าชะตา
  - `Placed In`: `Planet -> TargetHouse` ทั้งในดวงตัวเองและข้ามบุคคลตามราศีที่ดาวสถิต (น้ำหนัก `4.0–10.0` แปรผกผันกับ Orb)
- **Edge Type (3 สถานะ):** คำนวณจากกลุ่มภพเรือน (`Trikona/Labha +1`, `Dusthana -1`, `Kendra/Dhana 0`) รวมกับกำลัง/ความเป็นมิตรของดาว (`Harmonious` สีน้ำเงิน / `Neutral` สีเทา / `Challenging` สีแดง)

### 4.10 `POST /api/v1/hora/gephi_karma`
ส่งออก**ผังกรรม (Topology of Karma)** สำหรับ 2 คนขึ้นไป โดยใช้ **True Nodes** ของราหู/เกตุ และซ้อนชั้น **D-1 + D-9 Navamsa**:
- **Format รองรับ:** ค่าเริ่มต้นเป็น GEXF 1.2 XML (`vedic-karma-topology.gexf`) และรองรับ `?format=zip` (`karma_nodes.csv` + `karma_edges.csv` 2 ไฟล์), `?format=nodes_csv`, `?format=edges_csv`, `?format=json`
- **Nodes (คนละ 10 โหนด):** มี `Graha_Type` (`Benefic`/`Malefic`/`Karmic Axis`), `Dignity_Score` (`-5..+5`), `Force_Role` (`Attract`/`Repel` สำหรับ ForceAtlas2), พิกัด D1 และ D9 (`D9_Rasi`, `D9_Degree`, `D9_Dignity`)
- **Edges (3 ระบบกรรมข้ามบุคคล):** `Karmic Axis Alignment` (กุม/เล็งราหู-เกตุ W=10 ใน D1, W=8 ใน D9), `Karmic Taskmaster` (Saturn เกณฑ์ 3/7/10 W=7..9), และ `Nakshatra Dispositor` (Star Lord มาตรฐานดี -> เจ้าของนักษัตร W=6..9)

### 4.11 `POST /api/v1/hora/gephi_nakshatra`
ส่งออกข้อมูลโครงสร้างเครือข่าย**นักษัตรเป็น Hub กลาง** (`Shared Bipartite Nakshatra Hub Graph`) สำหรับ 2 คนขึ้นไป เป็นไฟล์ **GEXF 1.2 XML** (`application/gexf+xml`, `attachment; filename="vedic-nakshatra-synastry.gexf"`) เปิดใน Gephi ได้ทันที
- **Nodes (`node_kind`):** โหนดจุดของทุกคนคนละ 10 จุด (`Point`: `Lagna` + 9 ดาว) + โหนด Hub ของนักษัตรที่ถูกใช้งานจริง (`Nakshatra`: `Nak_01_Ashwini` .. `Nak_27_Revati`, ไม่มีโหนดลอยว่างเปล่า) พร้อม `occupied_count`, `nakshatra_lord`, `nakshatra_pada`
- **Edges (`Point -> Nakshatra`):** เส้นเชื่อมแบบมีทิศทาง (`Occupied Nakshatra`) ระบุ `pada` 1–4, `nakshatra_lord`, น้ำหนัก `3.5–10.0` (เน้น `Moon`, `Lagna` และกำลังดาว/Vargottama) และสีความเข้ากันได้กับเจ้าเรือนนักษัตร (`Harmonious` Blue / `Neutral` Gray / `Challenging` Red)

---

## 5. Build & Execution Rules

1. **WSL / Linux Environment:**
   - เนื่องจาก `swephgo` มี CGO Dependency (`libswe-dev`), การ Build และ Run จะต้องทำผ่านสภาพแวดล้อมที่มี GCC/Clang เช่น **WSL (Ubuntu/Debian)**
2. **การรันสคริปต์:**
   - Verify -> Test + Coverage (`-race`) -> Build -> Run: `./run.sh` (หรือ `./run.sh run` / `start`)
   - Verify -> Test + Coverage -> Build Binary (ไม่รันเซิร์ฟเวอร์): `./run.sh --build-only` (หรือ `./run.sh build` / `ci`)
   - รันจาก Binary เดิมทันที (หากยังไม่มี Binary จะ Verify + Test + Build ให้อัตโนมัติ): `./run.sh --no-build`
   - ข้ามขั้นตอน Unit Tests (ใช้ร่วมกับ `run` หรือ `--build-only`): `./run.sh --skip-test`
   - ตรวจสอบโค้ด + รัน Unit Tests พร้อม Coverage (`coverage.out` และ `coverage.html`): `./run.sh test` (หรือ `cover` / `coverage` / `check` หรือ `./test.sh`)
   - กำหนดเกณฑ์ Coverage ขั้นต่ำ / ตัวเลือกอื่น: `MIN_COVERAGE=40 ./run.sh test`, `./run.sh bench`, `./run.sh dev`, `./run.sh clean`, `./run.sh help`
   - บน Linux/WSL สคริปต์จะแนบ `-ldflags="-extldflags=-Wl,--no-as-needed,-lm"` ให้อัตโนมัติเพื่อรองรับ `libswe-dev`
3. **Concurrency Safety:**
   - การเรียกใช้ฟังก์ชัน Swiss Ephemeris (`swephgo`) ต้องถูกล็อกผ่าน `sweMu` (Mutex ใน `helpers.go`) เสมอ เนื่องจาก C Library มี Shared Global State
