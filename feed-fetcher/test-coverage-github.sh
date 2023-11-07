#!/bin/sh

# Calculate the total test coverage percentage from the 'cover.out' file
TOTAL=$(go tool cover -func=cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+')

# Compare the calculated coverage with the minimum coverage threshold
if (( $(echo "${TOTAL} > ${COVER_MIN}" | bc -l) )); then
  STATE=success
else
  STATE=failure
fi

# Update the GitHub repository status based on the coverage results
curl -X POST "https://${GITHUB_CREDS}@${REPO_API}" -d "{\"state\":\"${STATE}\",\"description\":\"Total coverage ${TOTAL}%, minimum ${COVER_MIN}%\",\"context\": \"quality gate\"}"
