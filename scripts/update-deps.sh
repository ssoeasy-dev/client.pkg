#!/bin/bash
# scripts/update-deps.sh
# Обновляет внутренние зависимости в go.mod и package.json всех модулей.
# Использование:
#   update-deps.sh --mode=(beta|stable) [pkg1 version1] [pkg2 version2] ...
#
# Режимы:
#   beta   - заменяет ВСЕ dev-версии на соответствующие beta-версии.
#   stable - заменяет beta-версии указанных пакетов на стабильные.

set -e

MODE=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        --mode=*)
            MODE="${1#*=}"
            shift
            ;;
        --mode)
            MODE="$2"
            shift 2
            ;;
        *)
            break
            ;;
    esac
done

if [[ "$MODE" != "beta" && "$MODE" != "stable" ]]; then
    echo "Ошибка: --mode должен быть 'beta' или 'stable'"
    exit 1
fi

SCOPE_GO="github.com/ssoeasy-dev/client.pkg"
SCOPE_NPM="@ssoeasy-dev"

cd "$(git rev-parse --show-toplevel)" || exit 1
mapfile -t ALL_PKGS < <(bash scripts/list-packages.sh)

# Функция обновления одного пакета
update_pkg() {
    local pkg_path="$1"
    local version="$2"
    local pattern

    if [[ "$MODE" == "beta" ]]; then
        pattern='v[0-9]+\.[0-9]+\.[0-9]+-dev-[^[:space:]]+'
        echo "Обновление $pkg_path до $version (dev → beta)"
    else
        pattern='v[0-9]+\.[0-9]+\.[0-9]+-beta\.[0-9]+'
        echo "Обновление $pkg_path до стабильной $version (beta → stable)"
    fi

    for mod in "${ALL_PKGS[@]}"; do
        # Go модуль
        if [[ -f "$mod/go.mod" ]]; then
            pkg_module="${SCOPE_GO}/${pkg_path}"
            if grep -q "${pkg_module} " "$mod/go.mod"; then
                echo "  Обновление в Go модуле $mod"
                (cd "$mod" && go mod edit -require "${pkg_module}@${version}")
            fi
        fi

        # npm пакет
        if [[ -f "$mod/package.json" ]]; then
            npm_name="${pkg_path##*/}"
            full_npm="${SCOPE_NPM}/${npm_name}"
            if jq -e ".dependencies[\"$full_npm\"]" "$mod/package.json" > /dev/null 2>&1; then
                echo "  Обновление в npm пакете $mod"
                jq ".dependencies[\"$full_npm\"] = \"$version\"" "$mod/package.json" > "$mod/package.json.tmp"
                mv "$mod/package.json.tmp" "$mod/package.json"
            fi
        fi
    done
}

# Применяем обновления для каждой пары пакет-версия
while [[ $# -ge 2 ]]; do
    update_pkg "$1" "$2"
    shift 2
done

if [[ $# -eq 1 ]]; then
    echo "Ошибка: пропущена версия для пакета $1" >&2
    exit 1
fi

# Обновляем lock-файлы
for mod in "${ALL_PKGS[@]}"; do
    if [[ -f "$mod/go.mod" ]]; then
        echo "  go mod tidy для $mod"
        (cd "$mod" && go mod tidy)
        git add "$mod/go.mod" "$mod/go.sum" 2>/dev/null || true
    fi
    if [[ -f "$mod/package.json" ]]; then
        echo "  pnpm install для $mod"
        (cd "$mod" && rm -rf node_modules pnpm-lock.yaml && pnpm install --lockfile-only)
        git add "$mod/package.json" "$mod/pnpm-lock.yaml" 2>/dev/null || true
    fi
done

# Коммитим изменения, если они есть
if ! git diff --cached --quiet; then
    if [[ "$MODE" == "beta" ]]; then
        git commit -m "chore(deps): update all dev dependencies to beta versions"
    else
        git commit -m "chore(deps): update beta dependencies to stable versions"
    fi
    echo "Изменения закоммичены."
else
    echo "Обновлений зависимостей не требуется."
fi
