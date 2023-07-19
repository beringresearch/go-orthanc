#!/bin/bash

source orthanc.env

git clone https://github.com/beringresearch/go-orthanc && \
cd go-orthanc && \
go install && \
cd && \
rm -r go-orthanc && \
mv go/bin/dicomweb-cli . && \
rm -r go

./dicomweb-cli configure http://localhost:8042/dicom-web

mkdir dcm

while read p; do
  OUT=$(basename $p)
  wget -r -N -c -np --user $MIMIC_USERNAME --password $MIMIC_PASSWORD https://physionet.org/files/mimic-cxr/2.0.0/files/$p -O dcm/$OUT
done <dcm_file_list.txt

for FILE in dcm/*; do
  ./dicomweb-cli upload $FILE
done
rm -r dcm