#!/bin/bash
# scripts/build-release-plan.sh
# Принимает список изменённых пакетов, строит граф зависимостей,
# выполняет топологическую сортировку и возвращает уровни в JSON.

set -e

cd "$(git rev-parse --show-toplevel)" || exit 1

CHANGED_PKGS=("$@")
SCOPE_GO="github.com/ssoeasy-dev/client.pkg"
SCOPE_NPM="@ssoeasy-dev"

mapfile -t ALL_PKGS < <(bash scripts/list-packages.sh)

declare -A is_changed
for pkg in "${CHANGED_PKGS[@]}"; do
    is_changed["$pkg"]=1
done

declare -A module_to_path deps indegree adj

# Собираем информацию о модулях и зависимостях
for pkg in "${ALL_PKGS[@]}"; do
    internal_deps=()

    if [[ -f "$pkg/go.mod" ]]; then
        module=$(grep -E "^module\s+" "$pkg/go.mod" | awk '{print $2}')
        module_to_path["$module"]="$pkg"
        # Извлекаем зависимости с префиксом SCOPE_GO
        mapfile -t go_deps < <(sed -n '/^require/,/^)/p' "$pkg/go.mod" | \
            grep -oE "${SCOPE_GO}[^[:space:]]*" | \
            sed "s|${SCOPE_GO}/||" | sort -u)
        internal_deps+=("${go_deps[@]}")
    fi

    if [[ -f "$pkg/package.json" ]]; then
        name=$(jq -r '.name' "$pkg/package.json")
        module_to_path["$name"]="$pkg"
        # Извлекаем зависимости с префиксом SCOPE_NPM (полное имя пакета)
        mapfile -t npm_deps < <(jq -r --arg scope "$SCOPE_NPM" '
            .dependencies // {} | to_entries[] | select(.key | startswith($scope + "/")) | .key
        ' "$pkg/package.json" 2>/dev/null)
        internal_deps+=("${npm_deps[@]}")
    fi

    # Убираем дубликаты и пустые строки
    deps["$pkg"]=$(printf '%s\n' "${internal_deps[@]}" | sort -u | sed '/^$/d')
done

# Строим граф только между изменёнными пакетами
for pkg in "${!is_changed[@]}"; do
    for dep_name in ${deps["$pkg"]}; do
        [[ -z "$dep_name" ]] && continue
        dep_path="${module_to_path["$dep_name"]}"
        if [[ -n "$dep_path" && -n "${is_changed["$dep_path"]}" ]]; then
            # pkg зависит от dep_path → dep_path должен быть выпущен раньше
            adj["$dep_path"]="${adj["$dep_path"]} $pkg"
            indegree["$pkg"]=$((indegree["$pkg"] + 1))
        fi
    done
done

# Топологическая сортировка (алгоритм Кана)
queue=()
for pkg in "${!is_changed[@]}"; do
    if [[ ${indegree["$pkg"]:-0} -eq 0 ]]; then
        queue+=("$pkg")
    fi
done

levels=()
while [[ ${#queue[@]} -gt 0 ]]; do
    level_json=$(printf '%s\n' "${queue[@]}" | jq -R . | jq -s -c .)
    levels+=("$level_json")

    new_queue=()
    for pkg in "${queue[@]}"; do
        for neighbor in ${adj["$pkg"]}; do
            indegree["$neighbor"]=$((indegree["$neighbor"] - 1))
            if [[ ${indegree["$neighbor"]} -eq 0 ]]; then
                new_queue+=("$neighbor")
            fi
        done
    done
    queue=("${new_queue[@]}")
done

if [[ ${#levels[@]} -eq 0 ]]; then
    echo "[]"
else
    printf '%s\n' "${levels[@]}" | jq -s -c .
fi
