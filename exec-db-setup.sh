#!/bin/bash
set -euo pipefail

bash exec-migration.sh
bash exec-ddl.sh
