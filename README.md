# go-acvp

Go FIPS-140-3 Automated Cryptographic Validation Protocol (ACVP) test data.

## Usage

```bash
# Make sure you've created config.json in this directory, pointing at your
# ACVTS creds.

ACVPTOOL=/tmp/boringssl/util/fipstool/acvp/acvptool/acvptool
WRAPPER=/tmp/go/src/fips.test

# Fetch full vectors for all algorithms
go run ./cmd/fetch -tool $ACVPTOOL -wrapper $WRAPPER

# Process full vectors for all algorithms, verifying solutions with the ACVTS
go run ./cmd/process -tool $ACVPTOOL -wrapper $WRAPPER -upload

# But: full vectors/solutions are too big for CI. Trim down to 1 exemplar per
# test type per algorithm. This will generate trimmed vectors/*.bz2 for each alg.
go run ./cmd/trim

# Create expected answers by re-processing the trimmed vectors.
# This will generate trimmed expected/*.bz2 for each alg.
go run ./cmd/process -tool $ACVPTOOL -wrapper $WRAPPER

# Commit the .bz2 files. You're done
```

# License

All data obtained from the NIST Automated Cryptographic Validation Testing
System (ACVTS) demo server is reproduced here under the same license as the
[ACVP-Server].

> NIST-developed software is provided by NIST as a public service. You may use, copy, and distribute copies of the software in any medium, provided that you keep intact this entire notice. You may improve, modify, and create derivative works of the software or any portion of the software, and you may copy and distribute such modifications or works. Modified works should carry a notice stating that you changed the software and should note the date and nature of any such change. Please explicitly acknowledge the National Institute of Standards and Technology as the source of the software. 
>
> NIST-developed software is expressly provided "AS IS." NIST MAKES NO WARRANTY OF ANY KIND, EXPRESS, IMPLIED, IN FACT, OR ARISING BY OPERATION OF LAW, INCLUDING, WITHOUT LIMITATION, THE IMPLIED WARRANTY OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND DATA ACCURACY. NIST NEITHER REPRESENTS NOR WARRANTS THAT THE OPERATION OF THE SOFTWARE WILL BE UNINTERRUPTED OR ERROR-FREE, OR THAT ANY DEFECTS WILL BE CORRECTED. NIST DOES NOT WARRANT OR MAKE ANY REPRESENTATIONS REGARDING THE USE OF THE SOFTWARE OR THE RESULTS THEREOF, INCLUDING BUT NOT LIMITED TO THE CORRECTNESS, ACCURACY, RELIABILITY, OR USEFULNESS OF THE SOFTWARE.
>
>You are solely responsible for determining the appropriateness of using and distributing the software and you assume all risks associated with its use, including but not limited to the risks and costs of program errors, compliance with applicable laws, damage to or loss of data, programs or equipment, and the unavailability or interruption of operation. This software is not intended to be used in any situation where a failure could cause risk of injury or damage to property. The software developed by NIST employees is not subject to copyright protection within the United States.

[ACVP-Server]: https://github.com/usnistgov/ACVP-Server/tree/master?tab=readme-ov-file#license
