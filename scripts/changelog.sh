#!/bin/bash

# Create a new CHANGELOG.md file with the header "Changelog"
echo "# Changelog" > CHANGELOG.md

# Loop through each tag, sorted by version in descending order
git tag --sort=-version:refname | while read -r TAG; do
    # Add a new section in the changelog for the current tag
    echo -e "\n## $TAG\n" >> CHANGELOG.md

    # Check if the previous tag exists
    if [ "$PREVIOUS_TAG" != "" ]; then
        # If the previous tag exists, list the commits between the current and previous tag
        git log --no-merges --format="* %s (%h)" $TAG..$PREVIOUS_TAG >> CHANGELOG.md
    else
        # If there is no previous tag, list the commits for the current tag
        git log --no-merges --format="* %s (%h)" $TAG >> CHANGELOG.md
    fi

    # Set the current tag as the previous tag for the next iteration
    PREVIOUS_TAG=$TAG
done
