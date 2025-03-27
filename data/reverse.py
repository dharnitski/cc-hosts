import os

# Define paths
folder1 = 'edges'  # Source folder
folder2 = 'edges_reversed'  # Destination folder

# Create folder2 if it doesn't exist
os.makedirs(folder2, exist_ok=True)

# Process each file in folder1
for filename in os.listdir(folder1):
    input_path = os.path.join(folder1, filename)
    output_path = os.path.join(folder2, filename)
    
    with open(input_path, 'r') as infile, open(output_path, 'w') as outfile:
        for line in infile:
            # Split by tab and swap the two parts
            parts = line.strip().split('\t')
            if len(parts) == 2:
                swapped_line = f"{parts[1]}\t{parts[0]}\n"
                outfile.write(swapped_line)
            else:
                # Handle lines that don't match the expected format
                outfile.write(line)  # Or skip with: continue

print("Processing complete!")
