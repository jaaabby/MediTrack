#!/bin/bash

echo "Branch actual: $CF_PAGES_BRANCH"

if [ "$CF_PAGES_BRANCH" = "develop" ]; then
  npm run build:develop
elif [ "$CF_PAGES_BRANCH" = "test" ]; then
  npm run build:test
else
  npm run build:production
fi