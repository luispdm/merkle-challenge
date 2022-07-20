# Test Plan

## High-level objective
The objective of this test plan is to verify that the application (app) works as expected. That means validating that the product can act as a trustless web server. It will be validated that:
- it is possible to download a file per parts
- it is possible to reconstruct the root hash from any piece of the file
- files can be served in a fault-tolerant and tamper-safe manner

## Scope
The app also mimics a trusted source which users can get the root hash from. However, this part of the app will not be tested as it falls out of the scope of the product goals and specifications.
### In scope
The part of the app that returns content and hash of each piece of the file(s) (`/piece/:hash/:pieceId` endpoint) is the primary point of focus for this test plan.
### Out of scope
The part of the app that mimics the trusted source (`/hashes` endpoint) will not be tested: it will only be used to retrieve information that helps validate the rest of the app.

## Test item
The app is only made of a single Java web server, therefore, there is only one item to be tested.

## Approach
Functional tests, both positives, and negatives, following the black-box methodology as the web server source code is unknown.

## Criteria

### Exit criteria
95% of the tests pass **and** the test running the reconstruction of the root hash from every piece passes.

### Suspension criteria
20% of tests fail **or** the test running the reconstruction of the root hash from any piece fails.

### Considerations
The exit criteria and suspension criteria have very strict requirements because the app is a proof of concept: its functionality is limited.

## Test cases
### Pre-conditions
The web server will be spawned with the following command:
```bash
java -jar merkle-tree-java.jar icons_rgb_circle.png example-file-merkle-tree.png
```
By doing so, the `/hashes` endpoint will return [icons_rgb_circle.png](../problem/icons_rgb_circle.png) at index 0 and [example-file-merkle-tree.png](../problem/example-file-merkle-tree.png) at index 1.
### Legend
*IH* = `9b39e1edb4858f7a3424d5a3d0c4579332640e58e101c29f99314a12329fc60b` (root hash of icons_rgb_circle.png)

### Test cases document
<details>
<summary>Expand</summary>

| Test case | Test description | Steps | Test data | Comment |
| --------- | ---------------- | ----- | --------- | ------- |
| **1** | Given the root hash *IH*, when a request is made to `/piece/IH/pieceId`, and the root hash from that piece is calculated, then the calculated root hash matches *IH*  | 1. Make a `GET` request to `/piece/IH/pieceId` <br> 2. Base64 decode the received content <br> 3. Calculate the SHA256 of the decoded content (leaf node) <br> 4. With the leaf node and the proofs, calculate the root hash <br> 5. Compare the calculated root hash to *IH* <br> 6. Repeat the steps above for every value of `pieceId` | `pieceId` = [0...16] | |
| **2** | Given the root hash *IH*, when a request is made to `/piece/IH/pieceId`, and the piece hash is calculated, and a request is made to `/piece/IH/pieceSiblingId`, then the piece hash matches the first element of the `proofs` array of its sibling | 1. Make a `GET` request to `/piece/IH/pieceId` <br> 2. Base64 decode the received content <br> 3. Calculate the SHA256 of the decoded content <br> 4. Make a `GET` request to `/piece/IH/pieceSiblingId` <br> 5. Compare the piece hash to its sibling's first element of the `proofs` array <br> 6. Repeat the steps above for every value of `pieceId` | `pieceId` = [0...15] <br> `pieceSiblingId` = piece id of the current piece's sibling | The test must skip the verification for the last piece because its sibling is off the end of the file |
| **3** | Given the root hash *IH*, when a request is made to `/piece/IH/pieceId`, and a request is made to `/piece/IH/pieceSiblingId`, then the uncles of the `proofs` array of both pieces match | 1. Make a `GET` request to `/piece/IH/pieceId` <br> 2. Make a `GET` request to `/piece/IH/pieceSiblingId` <br> 3. Compare the elements of the `proofs` array of the first response to the ones of the second response, except for the first element <br> 4. Repeat the steps above for every value of `pieceId` | `pieceId` = [0...14] (with increments of +2 on each iteration) <br> `pieceSiblingId` = `pieceId` +1 | The test must skip the verification for the last piece because its sibling is off the end of the file |
| **4** | Given the correct piece hash `CPH`, and the correct piece index `CPI`, and the correct piece proofs `CPP`, when the root hash from that piece is calculated, then the calculated root hash matches *IH* | 1. With `CPH` and `CPP`, calculate the root hash <br> 2. Compare the calculated root hash to *IH*  | `CPI` = 0 <br> `CPH` = hash of `CPI` <br> `CPP` = proofs returned by `/piece/IH/CPI` | |
| **5** | Given the wrong piece hash `WPH`, and the wrong piece index `WPI`, and the wrong piece proofs `WPP`, when the root hash from that piece is calculated, then the calculated root hash does not match *IH* | 1. With `WPH` and `WPP`, calculate the root hash <br> 2. Compare the calculated root hash to *IH*  | `WPI` = 2 <br> `WPH` = a hash that does not belong to the current file <br> `WPP` = an array made of hashes that do not belong to the current file | |
| **6** | Given the correct parent hash `CPTH`, and the correct child's index `CCI`, and the correct child's proofs `CCP`, when the root hash from that parent hash is calculated, then the calculated root hash matches *IH* | 1. With `CPTH` and `CCP`, calculate the root hash <br> 2. Compare the calculated root hash to *IH*  | `CCI` = 4 <br> `CPTH` = hash (hash(`CCI`)+hash(`CCI` sibling)) <br> `CCP` = proofs returned by `/piece/IH/CCI` | |
| **7** | Given the wrong parent hash `WPTH`, and the wrong child's index `WCI`, and the wrong child's proofs `WCP`, when the root hash from that parent hash is calculated, then the calculated root hash does not match *IH* | 1. With `WPTH` and `WCP`, calculate the root hash <br> 2. Compare the calculated root hash to *IH*  | `WCI` = 5 <br> `WPTH` = a hash that does not belong to the current file <br> `WCP` = an array made of hashes that do not belong to the current file | |
| **8** | Given the root hash *IH*, when a request is made to `/piece/IH/pieceId`, then the size of `content` is exactly 1 KB | 1. Make a `GET` request to `/piece/IH/pieceId` <br> 2. Get the size of `content` <br> 3. Compare the size of `content` to 1 KB <br> 4. Repeat the steps above for every value of `pieceId` | `pieceId` = [0...16] | |
| **9** | Given the root hash *IH*, and the wrong piece index `WPI`, when a request is made to `/piece/IH/WPI`, then the server returns the status code `SC` | 1. Make a `GET` request to `/piece/IH/WPI` <br> 2. Compare the status code to `SC` | `WPI` = -1 <br> `SC` = 404 | |
| **10** | Given the root hash *IH*, and the wrong piece index `WPI`, when a request is made to `/piece/IH/WPI`, then the server returns the status code `SC` | 1. Make a `GET` request to `/piece/IH/WPI` <br> 2. Compare the status code to `SC` | `WPI` = 17 <br> `SC` = 404 | |
| **11** | Given the root hash *IH*, and the wrong piece index `WPI`, when a request is made to `/piece/IH/WPI`, then the server returns the status code `SC` | 1. Make a `GET` request to `/piece/IH/WPI` <br> 2. Compare the status code to `SC` | `WPI` = abc <br> `SC` = 400 | |
| **12** | Given the root hash *IH*, and the wrong piece index `WPI`, when a request is made to `/piece/IH/WPI`, then the server returns the status code `SC` | 1. Make a `GET` request to `/piece/IH/WPI` <br> 2. Compare the status code to `SC` | `WPI` = 🫠 <br> `SC` = 400 | |
| **13** | Given the root hash *IH*, and the piece index `PI`, when a request is made to `/piece/IH/PI`, then the first element of the `proofs` array matches `S0` | 1. Make a `GET` request to `/piece/IH/PI` <br> 2. Compare the first element of the `proofs` array to `S0` | `PI` = 16 <br> `S0` = 0...0 (32 0s) | |
| **14** | Given the index `I`, when a request is made to `/hashes`, and the values of `hash` and `pieces` at index `I` are retrieved, and a request is made to `/piece/hash/pieceId`, then the server returns the status code `SC` | 1. Make a `GET` request to `/hashes` <br> 2. Retrieve the value of `hash` and `pieces` at index `I` of the array <br> 3. Make a `GET` request to `/piece/hash/pieceId` <br> 4. Compare the status code to `SC` <br> 5. Repeat the steps above for `pieceId` = [0...`pieces`-1] | `I` = 1 <br> `SC` = 200 | |

</details>

## Test environments and CI
The CI system to be used will be [Circle CI](https://circleci.com).

The CI pipeline will be configured to run the tests on every commit of a pull request targeting the `main` branch of the repo. Every test execution will generate a junit `.xml` test report that will be stored in Circle CI.

Two custom Docker images will be generated each time the CI pipeline is triggered: one for the tests and one for the web server. The custom containers will be part of the same docker-compose to enable the test container to make requests to the web server container.

The Docker images can be stored in an image registry and pulled by Circle CI, instead of being built on every execution, provided that this approach is cost-effective and speeds things up.

## Test execution process
1. The [test cases](#test-cases-document) will be translated into the Go programming language
2. In parallel, the CI pipeline will be set up

## Deliverables
At the end of the test execution cycle, a PDF file with test results and defects reports will be delivered.
