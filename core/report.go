package core

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func CreateStructuredReport(dicomPath string) error {
	f, err := os.Open(dicomPath)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		return err
	}

	ds, err := dicom.Parse(f, info.Size(), nil)
	if err != nil {
		return err
	}

	const enhancedSRIOD = "1.2.840.10008.5.1.4.1.1.88.22"
	const explicitVRLittleEndian = "1.2.840.10008.1.2.1"

	const SRModality = "SR"
	const manufacturer = "Bering Limited"
	const seriesDescription = "AI derived series"
	const manufacturerModelName = "BraveCX"
	const softwareVersions = "1.0.0"

	// // StudyDate from original DICOM
	// studyInstanceUID := ""
	// studyID := ""
	// studyDate := ""
	// studyTime := ""

	// accessionNumber := ""
	// patientName := ""
	// patientID := ""
	// patientBirthDate := ""
	// patientSex := ""
	// patientAge := ""

	// // Series date is new - current date
	// seriesInstanceUID := ""
	// seriesNumber := "99"
	// instanceNumber := "1"
	// seriesDate := ""
	// contentDate := ""
	// seriesTime := ""
	// contentTime := ""

	// SOPInstanceUID := "2.25.116240234176243277889131258530491654266"
	SOPInstanceUID, err := generateUUID()
	if err != nil {
		return err
	}
	seriesInstanceUID, err := generateUUID()
	if err != nil {
		return err
	}

	outFile, err := os.Create("test_sr.dcm")
	if err != nil {
		return err
	}
	defer outFile.Close()

	metadataVerEle, err := dicom.NewElement(tag.FileMetaInformationVersion, []byte{01})
	if err != nil {
		return err
	}
	metadataSOPClassUIDEle, err := dicom.NewElement(tag.MediaStorageSOPClassUID, []string{enhancedSRIOD})
	if err != nil {
		return err
	}
	metadataSOPInstanceUIDEle, err := dicom.NewElement(tag.MediaStorageSOPInstanceUID, []string{SOPInstanceUID})
	if err != nil {
		return err
	}
	transferSyntaxEle, err := dicom.NewElement(tag.TransferSyntaxUID, []string{explicitVRLittleEndian})
	if err != nil {
		return err
	}
	// TODO
	// implementationClassUIDEle, err := dicom.NewElement(tag.ImplementationClassUID, []string{"1.2.276.0.7230010.3.0.3.6.6"})
	// if err != nil {
	// 	return err
	// }
	// implementationVersionNameEle, err := dicom.NewElement(tag.ImplementationVersionName, []string{"OFFIS_DCMTK_366"})
	// if err != nil {
	// 	return err
	// }
	sopInstanceUIDEle, err := dicom.NewElement(tag.SOPInstanceUID, []string{SOPInstanceUID})
	if err != nil {
		return err
	}
	sopClassUIDEle, err := dicom.NewElement(tag.SOPClassUID, []string{enhancedSRIOD})
	if err != nil {
		return err
	}

	// ---------------------------------------
	// ----  Pulled from original DICOM ------
	// ---------------------------------------

	// TODO: handle missing tags

	// Study fields
	studyInstanceUIDEle, err := ds.FindElementByTag(tag.StudyInstanceUID)
	if err != nil {
		return err
	}
	studyInstanceIDEle, err := ds.FindElementByTag(tag.StudyID)
	if err != nil {
		return err
	}
	studyDateEle, err := ds.FindElementByTag(tag.StudyDate)
	if err != nil {
		return err
	}
	studyTimeEle, err := ds.FindElementByTag(tag.StudyTime)
	if err != nil {
		return err
	}

	// Patient fields
	accessionNumberEle, err := ds.FindElementByTag(tag.AccessionNumber)
	if err != nil {
		return err
	}
	patientNameEle, err := ds.FindElementByTag(tag.PatientName)
	if err != nil {
		return err
	}
	patientIDEle, err := ds.FindElementByTag(tag.PatientID)
	if err != nil {
		return err
	}
	patientBirthDateEle, err := ds.FindElementByTag(tag.PatientBirthDate)
	if err != nil {
		return err
	}
	patientSexEle, err := ds.FindElementByTag(tag.PatientSex)
	if err != nil {
		return err
	}
	// patientAgeEle, err := ds.FindElementByTag(tag.PatientAge)
	// if err != nil {
	// 	return err
	// }

	// ---------------------------------------
	// ----       Generated fields      ------
	// ---------------------------------------

	// Generate series datetimes
	currentTime := time.Now()
	currentDate := currentTime.Format("20060602")           // yyyyMMdd
	currentTimestamp := currentTime.Format("150405.000000") // HHmmss.SSSSSS

	seriesInstanceUIDEle, err := dicom.NewElement(tag.SeriesInstanceUID, []string{seriesInstanceUID})
	if err != nil {
		return err
	}
	// TODO - way to find this out?
	seriesNumberEle, err := dicom.NewElement(tag.SeriesNumber, []string{"99"})
	if err != nil {
		return err
	}
	instanceNumberEle, err := dicom.NewElement(tag.InstanceNumber, []string{"1"})
	if err != nil {
		return err
	}
	seriesDateEle, err := dicom.NewElement(tag.SeriesDate, []string{currentDate})
	if err != nil {
		return err
	}
	contentDateEle, err := dicom.NewElement(tag.ContentDate, []string{currentDate})
	if err != nil {
		return err
	}
	seriesTimeEle, err := dicom.NewElement(tag.SeriesTime, []string{currentTimestamp})
	if err != nil {
		return err
	}
	contentTimeEle, err := dicom.NewElement(tag.ContentTime, []string{currentTimestamp})
	if err != nil {
		return err
	}

	// seriesNumber := "99"
	// instanceNumber := "1"
	// seriesDate := ""
	// contentDate := ""
	// seriesTime := ""
	// contentTime := ""

	// Generate series datetimes

	structuredReport := dicom.Dataset{
		Elements: []*dicom.Element{
			metadataVerEle,
			metadataSOPClassUIDEle,
			metadataSOPInstanceUIDEle,
			transferSyntaxEle,
			// implementationClassUIDEle,
			// implementationVersionNameEle,
			sopInstanceUIDEle,
			sopClassUIDEle,
			studyInstanceUIDEle,
			studyInstanceIDEle,
			studyDateEle,
			studyTimeEle,
			accessionNumberEle,
			patientNameEle,
			patientIDEle,
			patientBirthDateEle,
			patientSexEle,
			// patientAgeEle,
			seriesInstanceUIDEle,
			seriesNumberEle,
			seriesDateEle,
			instanceNumberEle,
			contentDateEle,
			seriesTimeEle,
			contentTimeEle,
		},
	}

	err = dicom.Write(outFile, structuredReport)
	if err != nil {
		return err
	}

	if err := outFile.Close(); err != nil {
		return err
	}

	return nil
}

// UUID-based UID generation
// https://stackoverflow.com/questions/10295792/how-to-generate-sopinstance-uid-for-dicom-file
// https://stackoverflow.com/questions/46304306/how-to-generate-unique-dicom-uid
func generateUUID() (string, error) {
	id := uuid.New()
	idBinary, err := id.MarshalBinary()
	if err != nil {
		return "", err
	}

	idInt := new(big.Int)
	idInt.SetBytes(idBinary)

	return fmt.Sprintf("2.25.%d", idInt), nil
}
