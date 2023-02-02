# dicomweb-cli
A CLI tool for interacting with DICOMweb services written in Golang.

Supports:
- Upload of DICOM studies using STOW-RS
- Download of DICOM studies using WADO-RS
- Querying of DICOM studies using QUDO-RS

## Install
Ensure `$GOPATH/bin` is on your PATH variable and then run `go install` from the root directory of this project.

The `dicomweb-cli` tool should be built and added to $GOPATH/bin

## Examples

```sh
dicomweb-cli configure
> Enter the DICOMweb server URL: http://localhost:8042/dicom-web

dicomweb-cli upload ./example.dcm
Uploaded  1.3.12.2.1107.5.4.3.123456789012345.19950922.121803.6.dcm

dicomweb-cli download 1.3.12.2.1107.5.4.3.123456789012345.19950922.121803.6 --output test_download.dcm
Downloaded  1.3.12.2.1107.5.4.3.123456789012345.19950922.121803.6

dicomweb-cli query 1.3.12.2.1107.5.4.3.123456789012345.19950922.121803.6
Metadata saved:  1.3.12.2.1107.5.4.3.123456789012345.19950922.121803.6
```
