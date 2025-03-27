for file in edges_reversed/*.txt; do
    sort "$file" -o "$file"
    echo "Sorted: $file"
done
