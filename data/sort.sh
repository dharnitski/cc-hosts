for file in edges_reversed/*.txt; do
    sort -n "$file" -o "$file"
    echo "Sorted: $file"
done
