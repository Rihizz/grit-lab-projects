#!/bin/bash

# Repository details
GITHUB_REPO_URL="https://github.com/Valkamo/filler"
BRANCH="main" # adjust if you need a different branch

# List of files or directories to fetch from the GitHub repo
FILES_TO_FETCH=(
    "docker_image/m1_robots"
    "docker_image/m1_game_engine"
)

# Create a temporary directory for cloning the GitHub repository
TEMP_DIR=$(mktemp -d)

# Clone only the specified branch and only the latest commit for speed
git clone --depth 1 --branch $BRANCH $GITHUB_REPO_URL $TEMP_DIR

# Fetch the files
for file in "${FILES_TO_FETCH[@]}"; do
    # Create the parent directory if it doesn't exist
    mkdir -p "$(dirname "$file")"
    
    if [[ -e "$TEMP_DIR/$file" ]]; then
        cp -r "$TEMP_DIR/$file" "$file"
        echo "Fetched $file from GitHub."
    else
        echo "Failed to fetch $file."
    fi
done

# Remove the temporary directory
rm -rf $TEMP_DIR

echo "Done fetching files."

cd docker_image
./build.sh