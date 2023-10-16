package core

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/frame"
	"github.com/suyashkumar/dicom/pkg/tag"
	"github.com/suyashkumar/dicom/pkg/uid"
)

func CreateDerivedImage(dicomPath string, imagePath string, outPath string) error {
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

	const digitalXrIOD = "1.2.840.10008.5.1.4.1.1.1.1"

	SOPInstanceUID, err := generateUUID()
	if err != nil {
		return err
	}
	seriesInstanceUID, err := generateUUID()
	if err != nil {
		return err
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// ---------------------------------------
	// -------  DICOM header metadata --------
	// ---------------------------------------

	metadataVerEle, err := dicom.NewElement(tag.FileMetaInformationVersion, []byte{01})
	if err != nil {
		return err
	}
	metadataSOPClassUIDEle, err := dicom.NewElement(tag.MediaStorageSOPClassUID, []string{digitalXrIOD})
	if err != nil {
		return err
	}
	metadataSOPInstanceUIDEle, err := dicom.NewElement(tag.MediaStorageSOPInstanceUID, []string{SOPInstanceUID})
	if err != nil {
		return err
	}
	transferSyntaxEle, err := dicom.NewElement(tag.TransferSyntaxUID, []string{uid.ExplicitVRLittleEndian})
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
	sopClassUIDEle, err := dicom.NewElement(tag.SOPClassUID, []string{digitalXrIOD})
	if err != nil {
		return err
	}

	// // ---------------------------------------
	// // ----  Pulled from original DICOM ------
	// // ---------------------------------------

	// // TODO: handle missing tags

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
	// Assume the new generated generated DICOM is the first and only in its sequence (since seriesIntaceUID is newly generated)
	seriesNumberEle, err := dicom.NewElement(tag.SeriesNumber, []string{"1"})
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

	// Constant fields
	const XRModality = "DX"
	const manufacturer = "Bering Limited"
	const seriesDescription = "AI derived series"
	const manufacturerModelName = "BraveCX"
	const softwareVersions = "1.0.0"

	modalityEle, err := dicom.NewElement(tag.Modality, []string{XRModality})
	if err != nil {
		return err
	}
	manufacturerEle, err := dicom.NewElement(tag.Manufacturer, []string{manufacturer})
	if err != nil {
		return err
	}
	seriesDescriptionEle, err := dicom.NewElement(tag.SeriesDescription, []string{seriesDescription})
	if err != nil {
		return err
	}
	manufacturerModelNameEle, err := dicom.NewElement(tag.ManufacturerModelName, []string{manufacturerModelName})
	if err != nil {
		return err
	}
	softwareVersionsEle, err := dicom.NewElement(tag.SoftwareVersions, []string{softwareVersions})
	if err != nil {
		return err
	}

	// ---------------------------------------
	// ----        Derived image         -----
	// ---------------------------------------

	derivedImage := dicom.Dataset{
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
			modalityEle,
			manufacturerEle,
			seriesDescriptionEle,
			manufacturerModelNameEle,
			softwareVersionsEle,
		},
	}

	// To test including an image
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	bufRead := bufio.NewReader(imgFile)

	img, _, err := image.Decode(bufRead)
	if err != nil {
		return err
	}
	rows := img.Bounds().Max.Y
	cols := img.Bounds().Max.X
	bitsAllocated := 8
	bitsStored := 8
	highBit := bitsStored - 1

	pixelDataInfo := dicom.PixelDataInfo{
		Frames: []frame.Frame{
			{NativeData: frame.NativeFrame{
				Data:          make([][]int, rows*cols),
				Rows:          rows,
				Cols:          cols,
				BitsPerSample: bitsAllocated,
			}},
		},
		IsEncapsulated: false,
	}
	for i := 0; i < rows*cols; i++ {
		pixelDataInfo.Frames[0].NativeData.Data[i] = make([]int, 3)
	}

	// Fill pixel values
	for i := 0; i < rows*cols; i++ {
		row := i / cols
		col := i % cols
		pixel := img.At(col, row)
		r, g, b, _ := pixel.RGBA()
		r_8bit := r >> 8
		g_8bit := g >> 8
		b_8bit := b >> 8
		pixelDataInfo.Frames[0].NativeData.Data[i][0] = int(r_8bit)
		pixelDataInfo.Frames[0].NativeData.Data[i][1] = int(g_8bit)
		pixelDataInfo.Frames[0].NativeData.Data[i][2] = int(b_8bit)
	}

	pixelDataEle, err := dicom.NewElement(tag.PixelData, pixelDataInfo)
	if err != nil {
		return fmt.Errorf("failed to add pixeldata element: %s", err)
	}

	samplesPerPixelEle, err := dicom.NewElement(tag.SamplesPerPixel, []int{3})
	if err != nil {
		return fmt.Errorf("failed to add samplesPerPixel element: %s", err)
	}
	planarConfigEle, err := dicom.NewElement(tag.PlanarConfiguration, []int{0})
	if err != nil {
		return fmt.Errorf("failed to add planar config element: %s", err)
	}
	rowsEle, err := dicom.NewElement(tag.Rows, []int{rows})
	if err != nil {
		return fmt.Errorf("failed to add rows element: %s", err)
	}
	colsEle, err := dicom.NewElement(tag.Columns, []int{cols})
	if err != nil {
		return fmt.Errorf("failed to add columns element: %s", err)
	}
	bitsAllocEle, err := dicom.NewElement(tag.BitsAllocated, []int{bitsAllocated})
	if err != nil {
		return fmt.Errorf("failed to add bits allocated element: %s", err)
	}
	bitsStoredEle, err := dicom.NewElement(tag.BitsStored, []int{bitsStored})
	if err != nil {
		return fmt.Errorf("failed to add bits stored element: %s", err)
	}
	highBitEle, err := dicom.NewElement(tag.HighBit, []int{highBit})
	if err != nil {
		return fmt.Errorf("failed to add high bit element: %s", err)
	}
	pixelRepEle, err := dicom.NewElement(tag.PixelRepresentation, []int{0})
	if err != nil {
		return fmt.Errorf("failed to add pixel rep element: %s", err)
	}
	// numFramesEle, err := dicom.NewElement(tag.NumberOfFrames, []int{1})
	// if err != nil {
	// 	return fmt.Errorf("failed to add numFrames element: %s", err)
	// }
	photometricEle, err := dicom.NewElement(tag.PhotometricInterpretation, []string{"RGB"})
	if err != nil {
		return fmt.Errorf("failed to add photometric interpretation element: %s", err)
	}
	imageTypeEle, err := dicom.NewElement(tag.ImageType, []string{"DERIVED", "SECONDARY"})
	if err != nil {
		return fmt.Errorf("failed to add image type element: %s", err)
	}

	var imageElements []*dicom.Element
	imageElements = append(
		imageElements,
		imageTypeEle,
		pixelRepEle,
		samplesPerPixelEle,
		planarConfigEle,
		rowsEle,
		colsEle,
		bitsAllocEle,
		bitsStoredEle,
		highBitEle,
		// numFramesEle,
		photometricEle,
		pixelDataEle,
	)

	derivedImage.Elements = append(derivedImage.Elements, imageElements...)

	bufOut := bufio.NewWriter(outFile)

	err = dicom.Write(bufOut, derivedImage)
	if err != nil {
		return err
	}

	err = bufOut.Flush()
	if err != nil {
		return err
	}

	if err := outFile.Close(); err != nil {
		return err
	}

	return nil
}
