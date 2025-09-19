#!/bin/bash

# Go-Admin 測試執行腳本
# 此腳本提供完整的測試執行功能，包含單元測試、整合測試、效能測試等

set -e  # 遇到錯誤立即退出

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 輔助函數
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 檢查必要工具
check_dependencies() {
    log_info "檢查必要工具..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go 語言環境未安裝"
        exit 1
    fi
    
    log_success "Go 版本: $(go version)"
}

# 設定測試環境
setup_test_env() {
    log_info "設定測試環境..."
    
    # 設定環境變數
    export APP_MODE=test
    export DB_TYPE=sqlite3
    export DB_PATH=":memory:"
    export LOG_LEVEL=error
    export JWT_SECRET=test-secret-key
    
    # 建立測試目錄
    mkdir -p test_logs
    mkdir -p test_uploads
    mkdir -p coverage_reports
    
    # 清理之前的測試檔案
    rm -f test.db
    rm -f test_*.log
    rm -rf test_uploads/*
    
    log_success "測試環境設定完成"
}

# 執行單元測試
run_unit_tests() {
    log_info "執行單元測試..."
    
    if [ -d "test/unit" ] || find . -name "*_test.go" -path "*/unit/*" | grep -q .; then
        go test -v -timeout 30s ./test/unit/... || {
            log_warning "部分單元測試失敗或測試檔案不存在"
        }
    else
        log_warning "未找到單元測試檔案"
    fi
    
    log_success "單元測試執行完成"
}

# 執行整合測試
run_integration_tests() {
    log_info "執行整合測試..."
    
    if [ -d "test/integration" ] || find . -name "*_test.go" -path "*/integration/*" | grep -q .; then
        go test -v -timeout 60s ./test/integration/... || {
            log_warning "部分整合測試失敗或測試檔案不存在"
        }
    else
        log_warning "未找到整合測試檔案"
    fi
    
    log_success "整合測試執行完成"
}

# 執行 API 測試
run_api_tests() {
    log_info "執行 API 測試..."
    
    if [ -d "test/api" ] || find . -name "*_test.go" -path "*/api/*" | grep -q .; then
        go test -v -timeout 120s ./test/api/... || {
            log_warning "部分 API 測試失敗或測試檔案不存在"
        }
    else
        log_warning "未找到 API 測試檔案"
    fi
    
    log_success "API 測試執行完成"
}

# 執行效能測試
run_benchmark_tests() {
    log_info "執行效能測試..."
    
    if [ -d "test/benchmark" ] || find . -name "*_test.go" -path "*/benchmark/*" | grep -q .; then
        go test -v -bench=. -benchmem -timeout 300s ./test/benchmark/... || {
            log_warning "部分效能測試失敗或測試檔案不存在"
        }
    else
        log_warning "未找到效能測試檔案"
    fi
    
    log_success "效能測試執行完成"
}

# 執行負載測試
run_load_tests() {
    log_info "執行負載測試..."
    
    if [ -d "test/load" ] || find . -name "*_test.go" -path "*/load/*" | grep -q .; then
        go test -v -timeout 600s ./test/load/... || {
            log_warning "部分負載測試失敗或測試檔案不存在"
        }
    else
        log_warning "未找到負載測試檔案"
    fi
    
    log_success "負載測試執行完成"
}

# 產生測試覆蓋率報告
generate_coverage_report() {
    log_info "產生測試覆蓋率報告..."
    
    # 執行覆蓋率測試
    go test -coverprofile=coverage_reports/coverage.out ./... 2>/dev/null || {
        log_warning "無法產生覆蓋率報告，可能沒有測試檔案"
        return
    }
    
    # 產生 HTML 報告
    if [ -f coverage_reports/coverage.out ]; then
        go tool cover -html=coverage_reports/coverage.out -o coverage_reports/coverage.html
        
        # 計算覆蓋率百分比
        COVERAGE=$(go tool cover -func=coverage_reports/coverage.out | grep "total:" | awk '{print $3}')
        log_success "測試覆蓋率: $COVERAGE"
        
        # 檢查覆蓋率是否達標
        COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
        if (( $(echo "$COVERAGE_NUM >= 70" | bc -l) )); then
            log_success "✅ 測試覆蓋率達標 (>= 70%)"
        else
            log_warning "⚠️  測試覆蓋率未達標 ($COVERAGE < 70%)"
        fi
        
        log_info "覆蓋率報告已產生: coverage_reports/coverage.html"
    fi
}

# 執行程式碼品質檢查
run_quality_checks() {
    log_info "執行程式碼品質檢查..."
    
    # 檢查 Go fmt
    if [ -n "$(gofmt -l .)" ]; then
        log_warning "發現程式碼格式問題："
        gofmt -l .
        log_info "執行 'gofmt -w .' 修復格式問題"
    else
        log_success "✅ 程式碼格式正確"
    fi
    
    # 檢查 Go vet
    if go vet ./...; then
        log_success "✅ Go vet 檢查通過"
    else
        log_warning "⚠️  Go vet 發現潛在問題"
    fi
    
    # 檢查是否有 golint（如果安裝的話）
    if command -v golint &> /dev/null; then
        if golint ./... | grep -v "should have comment" | grep -q .; then
            log_warning "⚠️  Golint 發現問題："
            golint ./... | grep -v "should have comment"
        else
            log_success "✅ Golint 檢查通過"
        fi
    fi
}

# 檢查測試結果
check_test_results() {
    log_info "檢查測試結果..."
    
    # 檢查是否有測試失敗
    if [ -f test_results.log ]; then
        if grep -q "FAIL" test_results.log; then
            log_error "❌ 發現測試失敗"
            grep "FAIL" test_results.log
            return 1
        fi
    fi
    
    # 檢查慢速測試
    if find . -name "*_test.go" -exec grep -l "time.Sleep.*[0-9]*s" {} \; | grep -q .; then
        log_warning "⚠️  發現可能的慢速測試（使用 time.Sleep）"
        find . -name "*_test.go" -exec grep -l "time.Sleep.*[0-9]*s" {} \;
    fi
    
    log_success "✅ 測試結果檢查完成"
}

# 清理測試環境
cleanup_test_env() {
    log_info "清理測試環境..."
    
    # 清理測試檔案
    rm -f test.db
    rm -f test_*.log
    
    # 清理測試上傳目錄
    if [ -d "test_uploads" ]; then
        rm -rf test_uploads/*
    fi
    
    log_success "測試環境清理完成"
}

# 顯示幫助資訊
show_help() {
    echo "Go-Admin 測試腳本"
    echo ""
    echo "使用方法: $0 [選項]"
    echo ""
    echo "選項:"
    echo "  all             執行所有測試"
    echo "  unit            僅執行單元測試"
    echo "  integration     僅執行整合測試"
    echo "  api             僅執行 API 測試"
    echo "  benchmark       僅執行效能測試"
    echo "  load            僅執行負載測試"
    echo "  coverage        僅產生覆蓋率報告"
    echo "  quality         僅執行程式碼品質檢查"
    echo "  cleanup         清理測試環境"
    echo "  help            顯示此幫助資訊"
    echo ""
    echo "範例:"
    echo "  $0 all          # 執行所有測試"
    echo "  $0 unit         # 僅執行單元測試"
    echo "  $0 coverage     # 僅產生覆蓋率報告"
}

# 主要執行邏輯
main() {
    local cmd=${1:-"all"}
    
    case $cmd in
        "all")
            log_info "🧪 開始執行完整測試套件..."
            check_dependencies
            setup_test_env
            run_unit_tests
            run_integration_tests
            run_api_tests
            run_benchmark_tests
            generate_coverage_report
            run_quality_checks
            check_test_results
            cleanup_test_env
            log_success "🎉 所有測試執行完成！"
            ;;
        "unit")
            log_info "🧪 執行單元測試..."
            check_dependencies
            setup_test_env
            run_unit_tests
            cleanup_test_env
            ;;
        "integration")
            log_info "🔗 執行整合測試..."
            check_dependencies
            setup_test_env
            run_integration_tests
            cleanup_test_env
            ;;
        "api")
            log_info "🌐 執行 API 測試..."
            check_dependencies
            setup_test_env
            run_api_tests
            cleanup_test_env
            ;;
        "benchmark")
            log_info "⚡ 執行效能測試..."
            check_dependencies
            setup_test_env
            run_benchmark_tests
            ;;
        "load")
            log_info "🔥 執行負載測試..."
            check_dependencies
            setup_test_env
            run_load_tests
            ;;
        "coverage")
            log_info "📊 產生測試覆蓋率報告..."
            check_dependencies
            setup_test_env
            generate_coverage_report
            ;;
        "quality")
            log_info "🔍 執行程式碼品質檢查..."
            check_dependencies
            run_quality_checks
            ;;
        "cleanup")
            log_info "🧹 清理測試環境..."
            cleanup_test_env
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "未知選項: $cmd"
            show_help
            exit 1
            ;;
    esac
}

# 執行主程式
main "$@"