
          #!/bin/bash

          make coverage > coverage.txt
          cat coverage.txt

          COVERAGE_PERCENTAGE=$(grep "^total:" "coverage.txt" | awk "{print $NF}")

          if [ "$COVERAGE_PERCENTAGE" != "100.0%" ]; then
            echo "Coverage is not 100%. Actual coverage: $COVERAGE_PERCENTAGE"
            exit 1
          fi

          echo "Cheers! Code coverage is 100%!"
          
