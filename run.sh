#!/usr/bin/env bash

# ==============================================================================
# Script Name : run.sh
# Description : คอมไพล์ (go build) และรันไฟล์ Binary ของ lahiri-api
# Usage       : ./run.sh          (Build ใหม่และรัน)
#               ./run.sh --no-build (รันจากไฟล์ Binary เดิมทันที)
#               ./run.sh --build-only (Build อย่างเดียว ไม่รัน)
# ==============================================================================

set -e

APP_NAME="lahiri-api"
BUILD_DIR="bin"

# ตรวจสอบระบบปฏิบัติการเพื่อกำหนดนามสกุลไฟล์ (.exe สำหรับ Windows)
EXT=""
case "$(uname -s)" in
    CYGWIN*|MINGW*|MSYS*)
        EXT=".exe"
        ;;
esac

BINARY_PATH="${BUILD_DIR}/${APP_NAME}${EXT}"

# ค่าเริ่มต้นของ Environment Variables (สามารถกำหนดจากภายนอกก่อนรันสคริปต์ได้)
export PORT="${PORT:-8083}"
export APP_ENV="${APP_ENV:-development}"
export CGO_ENABLED=1

DO_BUILD=true
DO_RUN=true

for arg in "$@"; do
    case "$arg" in
        --no-build)
            DO_BUILD=false
            ;;
        --build-only)
            DO_RUN=false
            ;;
    esac
done

echo "=================================================================="
echo " 🚀 Lahiri Vedic Astrology API - Build & Run Script"
echo "=================================================================="

build_binary() {
    mkdir -p "${BUILD_DIR}"
    go build -ldflags="-s -w" -o "${BINARY_PATH}" .
}

# ------------------------------------------------------------------------------
# Step 1: Build Go Binary
# ------------------------------------------------------------------------------
if [[ "${DO_BUILD}" == "true" ]]; then
    echo ""
    echo "📦 Step 1: Building binary -> ${BINARY_PATH} (CGO_ENABLED=${CGO_ENABLED})..."
    build_binary
    echo "✅ Build completed successfully!"
elif [[ ! -f "${BINARY_PATH}" ]]; then
    echo "⚠️  Binary not found at ${BINARY_PATH}. Building now..."
    build_binary
else
    echo ""
    echo "⏭️  Skipping build step (using existing ${BINARY_PATH})"
fi

# ------------------------------------------------------------------------------
# Step 2: Run Binary
# ------------------------------------------------------------------------------
if [[ "${DO_RUN}" == "true" ]]; then
    echo ""
    echo "🟢 Step 2: Starting ${APP_NAME} on port ${PORT} (APP_ENV=${APP_ENV})..."
    echo "------------------------------------------------------------------"
    exec "./${BINARY_PATH}"
fi
