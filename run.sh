#!/usr/bin/env bash

# ==============================================================================
# Script Name : run.sh
# Description : ตรวจสอบโค้ด, รัน Unit Tests พร้อม Test Coverage, คอมไพล์ Binary
#               และเริ่มรันเซิร์ฟเวอร์ lahiri-api อย่างสมบูรณ์
# Usage       : ./run.sh              (Verify -> Test + Cover -> Build -> Run)
#               ./run.sh --build-only (Verify -> Test + Cover -> Build)
#               ./run.sh --no-build   (รันจากไฟล์ Binary เดิมทันที)
#               ./run.sh test|cover|build|ci|bench|dev|clean|help
# ==============================================================================

set -euo pipefail

# ย้ายมายังโฟลเดอร์ของโปรเจกต์เสมอไม่ว่าจะเรียกสคริปต์จากที่ไหน
cd "$(dirname "$0")"

APP_NAME="lahiri-api"
BUILD_DIR="bin"
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
MIN_COVERAGE="${MIN_COVERAGE:-0}"

# ค่าเริ่มต้นของ Environment Variables (สามารถกำหนดจากภายนอกก่อนรันสคริปต์ได้)
export PORT="${PORT:-8083}"
export APP_ENV="${APP_ENV:-development}"
export CGO_ENABLED=1

# ตรวจสอบระบบปฏิบัติการเพื่อกำหนดนามสกุลไฟล์และ Linker Flags (รองรับ libswe-dev บน Linux)
EXT=""
TEST_LDFLAGS=""
BUILD_LDFLAGS="-s -w"
case "$(uname -s)" in
    CYGWIN*|MINGW*|MSYS*)
        EXT=".exe"
        ;;
    Linux*)
        TEST_LDFLAGS="-extldflags=-Wl,--no-as-needed,-lm"
        BUILD_LDFLAGS="-s -w ${TEST_LDFLAGS}"
        ;;
esac

BINARY_PATH="${BUILD_DIR}/${APP_NAME}${EXT}"

print_step() {
    printf "\n==================================================================\n"
    printf " %s\n" "$1"
    printf "==================================================================\n"
}

require_toolchain() {
    if ! command -v go >/dev/null 2>&1; then
        echo "❌ ไม่พบคำสั่ง 'go' ใน PATH กรุณาติดตั้ง Go 1.26+ ก่อนใช้งาน" >&2
        exit 1
    fi
    if ! command -v gcc >/dev/null 2>&1 && ! command -v cc >/dev/null 2>&1 && ! command -v clang >/dev/null 2>&1; then
        echo "❌ ไม่พบ C Compiler (gcc/cc/clang) ซึ่งจำเป็นสำหรับ CGO ร่วมกับ Swiss Ephemeris (swephgo)" >&2
        echo "💡 กรุณารันผ่าน WSL (Ubuntu/Debian) พร้อมติดตั้ง build-essential และ libswe-dev" >&2
        exit 1
    fi
}

detect_race_flag() {
    local os_name
    os_name="$(uname -s 2>/dev/null || echo unknown)"
    if [[ "${CGO_ENABLED}" == "1" ]] && [[ "${os_name}" == Linux* || "${os_name}" == Darwin* ]]; then
        echo "-race"
    else
        echo ""
    fi
}

run_go_test() {
    if [[ -n "${TEST_LDFLAGS}" ]]; then
        go test -ldflags="${TEST_LDFLAGS}" "$@"
    else
        go test "$@"
    fi
}

verify_code() {
    print_step "🔍 [1/3] ตรวจสอบ Dependencies และ Static Analysis (go mod verify + go vet)"
    go mod verify
    go vet ./...
    echo "✅ ผ่านการตรวจสอบ go mod verify และ go vet เรียบร้อย"
}

run_tests_with_coverage() {
    print_step "🧪 [2/3] รัน Unit Tests ทั้งหมด พร้อมวัด Test Coverage"

    local race_flag
    race_flag="$(detect_race_flag)"

    if [[ -n "${race_flag}" ]]; then
        echo "🛡️ เปิดใช้งาน Race Detector (-race)"
        run_go_test -count=1 -v -race -coverprofile="${COVERAGE_FILE}" -covermode=atomic ./...
    else
        echo "ℹ️ ข้าม -race (ไม่รองรับบนแพลตฟอร์มปัจจุบัน)"
        run_go_test -count=1 -v -coverprofile="${COVERAGE_FILE}" -covermode=atomic ./...
    fi

    print_step "📊 สรุปผล Test Coverage (${COVERAGE_FILE})"
    go tool cover -func="${COVERAGE_FILE}"

    local total_cov
    total_cov="$(go tool cover -func="${COVERAGE_FILE}" | awk '/^total:/ {print $NF}')"
    echo ""
    echo "✅ Total Statement Coverage: ${total_cov:-N/A}"

    if [[ -n "${MIN_COVERAGE}" ]] && [[ "${MIN_COVERAGE}" != "0" ]] && [[ -n "${total_cov:-}" ]]; then
        local numeric_cov="${total_cov%%%}"
        if awk -v cov="${numeric_cov}" -v min="${MIN_COVERAGE}" 'BEGIN { exit !(cov + 0 < min + 0) }'; then
            echo "❌ Test Coverage (${total_cov}) ต่ำกว่าเกณฑ์ขั้นต่ำที่กำหนด (${MIN_COVERAGE}%)" >&2
            exit 1
        fi
        echo "✅ ผ่านเกณฑ์ Test Coverage ขั้นต่ำ (${MIN_COVERAGE}%)"
    fi

    go tool cover -html="${COVERAGE_FILE}" -o "${COVERAGE_HTML}"
    echo "📄 สร้างรายงาน Coverage HTML สำเร็จ: ./${COVERAGE_HTML}"
}

run_benchmarks() {
    print_step "⚡ รัน Benchmarks"
    run_go_test -run='^$' -bench=. -benchmem ./...
}

build_binary() {
    print_step "📦 [3/3] Build Production Binary -> ${BINARY_PATH} (CGO_ENABLED=${CGO_ENABLED})"
    mkdir -p "${BUILD_DIR}"
    go build -v -ldflags="${BUILD_LDFLAGS}" -o "${BINARY_PATH}" .
    chmod +x "${BINARY_PATH}" 2>/dev/null || true
    echo "✅ Build สำเร็จ: ./${BINARY_PATH}"
}

clean_artifacts() {
    print_step "🧹 ล้างไฟล์ Build และ Test Artifacts"
    rm -f "${BINARY_PATH}" "${COVERAGE_FILE}" "${COVERAGE_HTML}"
    echo "✅ ล้างไฟล์เรียบร้อย"
}

show_help() {
    cat <<EOF
วิธีใช้งาน: ./run.sh [command|flags]

คำสั่งและตัวเลือกที่รองรับ:
  ./run.sh [run|start]      - Verify (mod verify + vet) -> Test + Coverage -> Build -> รันเซิร์ฟเวอร์ (ค่าเริ่มต้น)
  ./run.sh build | ci       - Verify -> Test + Coverage -> Build Binary (ไม่รันเซิร์ฟเวอร์)
  ./run.sh --build-only     - Verify -> Test + Coverage -> Build Binary (เทียบเท่า build/ci)
  ./run.sh --no-build       - รันจาก Binary เดิมทันที (หากยังไม่มี Binary จะทำการ Verify + Test + Build ให้อัตโนมัติ)
  ./run.sh --skip-test      - ข้ามขั้นตอน Unit Tests (ใช้ร่วมกับ run หรือ --build-only ได้)
  ./run.sh test | check     - ตรวจสอบโค้ด + รัน Unit Tests ทั้งหมด + สร้างรายงาน Coverage (.out & .html)
  ./run.sh cover | coverage - เทียบเท่าคำสั่ง test (สร้างทั้ง coverage.out และ coverage.html)
  ./run.sh bench            - ตรวจสอบโค้ด + รัน Unit Tests + Coverage + Benchmarks
  ./run.sh dev              - ตรวจสอบโค้ด แล้วรันเซิร์ฟเวอร์ด้วย go run .
  ./run.sh clean            - ลบไฟล์ Binary และรายงาน Coverage ชั่วคราว
  ./run.sh help | -h        - แสดงวิธีใช้งานนี้

ตัวแปรสภาพแวดล้อมเพิ่มเติม:
  PORT=<port>               - กำหนดพอร์ตเซิร์ฟเวอร์ (ค่าเริ่มต้น: 8083)
  APP_ENV=<env>             - กำหนดโหมด environment (ค่าเริ่มต้น: development)
  MIN_COVERAGE=<เปอร์เซ็นต์>  - กำหนดเกณฑ์ขั้นต่ำของ Test Coverage (เช่น MIN_COVERAGE=40 ./run.sh test)
EOF
}

run_pipeline() {
    require_toolchain

    echo "=================================================================="
    echo " 🚀 Lahiri Vedic Astrology API - Verify, Test, Build & Run"
    echo "=================================================================="

    if [[ "${DO_BUILD}" == "true" ]]; then
        verify_code
        if [[ "${DO_TEST}" == "true" ]]; then
            run_tests_with_coverage
        else
            echo ""
            echo "⏭️  ข้ามขั้นตอน Unit Tests (--skip-test)"
        fi
        build_binary
    elif [[ ! -f "${BINARY_PATH}" ]]; then
        echo ""
        echo "⚠️  ไม่พบไฟล์ Binary ที่ ${BINARY_PATH} — กำลังตรวจสอบ ทดสอบ และ Build ใหม่..."
        verify_code
        run_tests_with_coverage
        build_binary
    else
        echo ""
        echo "⏭️  ข้ามขั้นตอน Verify / Test / Build (ใช้ Binary เดิม: ${BINARY_PATH})"
    fi

    if [[ "${DO_RUN}" == "true" ]]; then
        print_step "🟢 กำลังเริ่มทำงาน ${APP_NAME} บนพอร์ต ${PORT} (APP_ENV=${APP_ENV})"
        exec "./${BINARY_PATH}"
    fi
}

ACTION="run"
DO_TEST=true
DO_BUILD=true
DO_RUN=true

if [[ $# -gt 0 ]]; then
    case "$1" in
        run|start|build|ci|test|check|cover|coverage|bench|dev|clean|help|-h|--help)
            ACTION="$1"
            shift
            ;;
    esac
fi

for arg in "$@"; do
    case "$arg" in
        --no-build)
            DO_BUILD=false
            ;;
        --build-only)
            DO_RUN=false
            ;;
        --skip-test)
            DO_TEST=false
            ;;
        -h|--help)
            ACTION="help"
            ;;
        *)
            echo "❌ ไม่รู้จักตัวเลือก '${arg}'" >&2
            echo ""
            show_help
            exit 1
            ;;
    esac
done

case "${ACTION}" in
    help|-h|--help)
        show_help
        ;;
    clean)
        clean_artifacts
        ;;
    test|check|cover|coverage)
        require_toolchain
        verify_code
        run_tests_with_coverage
        ;;
    bench)
        require_toolchain
        verify_code
        run_tests_with_coverage
        run_benchmarks
        ;;
    dev)
        require_toolchain
        print_step "🔧 โหมดพัฒนา (Development)"
        verify_code
        echo "🚀 กำลังรันเซิร์ฟเวอร์ด้วย go run . (PORT=${PORT}, APP_ENV=${APP_ENV})..."
        if [[ -n "${TEST_LDFLAGS}" ]]; then
            exec go run -ldflags="${TEST_LDFLAGS}" .
        else
            exec go run .
        fi
        ;;
    build|ci)
        DO_RUN=false
        run_pipeline
        ;;
    run|start)
        run_pipeline
        ;;
esac
