#!/bin/bash

# i18n 多语言测试脚本
# 用于快速测试所有支持的语言环境

set -e

echo "=========================================="
echo "  Podman-TUI i18n 多语言测试"
echo "=========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 支持的语言列表
declare -A LANGUAGES=(
    ["en_US.UTF-8"]="English"
    ["zh_CN.UTF-8"]="简体中文"
    ["zh_TW.UTF-8"]="繁體中文"
    ["ja_JP.UTF-8"]="日本語"
    ["ko_KR.UTF-8"]="한국어"
    ["ru_RU.UTF-8"]="Русский"
)

# 测试翻译示例
TEST_KEYS=("Cancel" "HELP" "add connection" "delete" "info")

echo "1️⃣  运行单元测试..."
echo "----------------------------------------"
go test ./i18n/ -v
echo ""

echo "2️⃣  检查翻译完整性..."
echo "----------------------------------------"
go run ./i18n/tools/check_translations.go ./i18n/
echo ""

echo "3️⃣  测试各语言环境..."
echo "----------------------------------------"

for lang in "${!LANGUAGES[@]}"; do
    name="${LANGUAGES[$lang]}"
    echo -e "${BLUE}测试语言: $name ($lang)${NC}"
    
    # 设置语言环境并运行示例
    export LANG=$lang
    
    # 运行简单的翻译测试
    go run -C ./i18n/example main.go 2>/dev/null | head -10
    
    echo -e "${GREEN}✓ $name 测试完成${NC}"
    echo ""
done

echo "=========================================="
echo -e "${GREEN}✅ 所有测试完成！${NC}"
echo "=========================================="
echo ""
echo "下一步："
echo "  1. 在不同语言环境下运行 podman-tui"
echo "     LANG=zh_CN.UTF-8 ./bin/darwin/podman-tui"
echo ""
echo "  2. 开始集成 i18n 到主程序"
echo "     参考: i18n/INTEGRATION.md"
echo ""
echo "  3. 查看快速入门指南"
echo "     参考: i18n/QUICKSTART.md"
echo ""
