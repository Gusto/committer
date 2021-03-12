#!/bin/bash -ex

VERSION="0.1.6"
GIT_PRE_COMMIT_HOOK=".git/hooks/pre-commit"
COMMITTER_YML="committer.yml"
COMMITTER_LOCATION="/usr/local/bin/committer"
DOWNLOAD_URL="https://github.com/Gusto/committer/releases/download/v$VERSION/committer-v$VERSION"
if [ ! -f $COMMITTER_LOCATION ]; then
  echo "Committer is not installed!"

  echo "Downloading $DOWNLOAD_URL to $COMMITTER_LOCATION..."
  curl -L -o $COMMITTER_LOCATION --fail $DOWNLOAD_URL

  echo "Making $COMMITTER_LOCATION executable"
  chmod +x $COMMITTER_LOCATION
fi
  # Write out committer.yml
if [ ! -f $COMMITTER_YML ]; then
  cat > $COMMITTER_YML <<EOL
tasks:
- name: Rubocop
  command: bundle exec rubocop  --color --force-exclusion --format simple
  files: '\.rb|\.rake'
  fix:
    command: bundle exec rubocop --color --force-exclusion --format simple --auto-correct
    output: '\[Corrected\]|=='
EOL
fi

  # Setup git hook
if [ ! -f $GIT_PRE_COMMIT_HOOK ]; then
  cat > $GIT_PRE_COMMIT_HOOK <<EOL
#!/bin/bash
/usr/local/bin/committer --fix --changed
EOL
  chmod +x $GIT_PRE_COMMIT_HOOK
fi
