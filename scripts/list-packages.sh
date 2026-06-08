#!/bin/bash
# Поиск всех пакетов (Go и npm) в репозитории.
# Исключает директории, содержащие "example" или "examples".

cd "$(git rev-parse --show-toplevel)" || exit 1

packages=()
while IFS= read -r -d '' dir; do
    # Проверяем, есть ли go.mod или package.json
    if [[ -f "$dir/go.mod" || -f "$dir/package.json" ]]; then
        # Убираем начальный "./"
        pkg_path="${dir#./}"
        packages+=("$pkg_path")
    fi
done < <(find . -type d \
    -not -path "*/.*" \
    -not -path "*/node_modules*" \
    -not -path "*/example*" \
    -not -path "*/examples*" \
    -print0)

if [[ "$1" == "--json" ]]; then
    printf '%s\n' "${packages[@]}" | jq -R . | jq -s -c .
else
    printf '%s\n' "${packages[@]}"
fi
