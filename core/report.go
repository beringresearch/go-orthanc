package core

import (
	"bufio"
	"os"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func CreateStructuredReport(dicomPath string, outPath string) error {
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

	// Gather derived DICOM elements
	headerElements, err := derivedHeaderElements(enhancedSRIOD)
	if err != nil {
		return err
	}
	derivedElements, err := derivedMetadata(ds)
	if err != nil {
		return err
	}
	generatedElements, err := generateInstanceMetadata(srModality)
	if err != nil {
		return err
	}

	sr, err := generateSR()
	if err != nil {
		return err
	}

	// Assemble derived dataset from elements
	derivedDicom := dicom.Dataset{
		Elements: []*dicom.Element{},
	}
	derivedDicom.Elements = append(derivedDicom.Elements, headerElements...)
	derivedDicom.Elements = append(derivedDicom.Elements, derivedElements...)
	derivedDicom.Elements = append(derivedDicom.Elements, generatedElements...)
	derivedDicom.Elements = append(derivedDicom.Elements, sr)

	// ---------------------------------------
	// ----       Structured Report      -----
	// ---------------------------------------
	// Insert probabilities in a known format.
	// Potential use of SR templates here: https://dicom.nema.org/medical/dicom/current/output/chtml/part16/sect_ChestCADSRIODTemplates.html
	// For now keeping it as simple as possible - like a "JSON" of classification probabilities.
	// Front-end should expect this and be able to parse it

	// Write out dataset
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	bufOut := bufio.NewWriter(outFile)

	err = dicom.Write(bufOut, derivedDicom)
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

func generateSR() (*dicom.Element, error) {

	var measurementUnitCodeSequence []*dicom.Element

	codeValueEle, err := dicom.NewElement(tag.CodeValue, []string{"probability"})
	if err != nil {
		return nil, err
	}
	// TODO: define the schema we are using - private one?
	// codingSchemeDesignatorEle, err := dicom.NewElement(tag.CodingSchemeDesignator, []string{"UCUM"})
	// if err != nil {
	// 	return nil, err
	// }
	// codingSchemeVersionEle, err := dicom.NewElement(tag.CodingSchemeVersion, []string{"1.4"})
	// if err != nil {
	// 	return nil, err
	// }
	codeMeaningEle, err := dicom.NewElement(tag.CodeMeaning, []string{"abnormality classification"})
	if err != nil {
		return nil, err
	}

	measurementUnitCodeSequence = append(measurementUnitCodeSequence, codeValueEle)
	// measurementUnitCodeSequence = append(measurementUnitCodeSequence, codingSchemeDesignatorEle)
	// measurementUnitCodeSequence = append(measurementUnitCodeSequence, codingSchemeVersionEle)
	measurementUnitCodeSequence = append(measurementUnitCodeSequence, codeMeaningEle)

	// MeasuredValueSequence contains MeasurementUnitCodeSequence and NumericValue
	numericValueEle, err := dicom.NewElement(tag.NumericValue, []string{"0.75"})
	if err != nil {
		return nil, err
	}

	var measuredValueSequence []*dicom.Element
	measuredValueSequence = append(measuredValueSequence, measurementUnitCodeSequence...)
	measuredValueSequence = append(measuredValueSequence, numericValueEle)

	// Second probability class
	var measurementUnitCodeSequence2 []*dicom.Element
	codeValueEle2, err := dicom.NewElement(tag.CodeValue, []string{"probability"})
	if err != nil {
		return nil, err
	}
	codeMeaningEle2, err := dicom.NewElement(tag.CodeMeaning, []string{"normality classification"})
	if err != nil {
		return nil, err
	}
	measurementUnitCodeSequence2 = append(measurementUnitCodeSequence2, codeValueEle2)
	measurementUnitCodeSequence2 = append(measurementUnitCodeSequence2, codeMeaningEle2)

	numericValueEle2, err := dicom.NewElement(tag.NumericValue, []string{"0.25"})
	if err != nil {
		return nil, err
	}

	var measuredValueSequence2 []*dicom.Element
	measuredValueSequence2 = append(measuredValueSequence2, measurementUnitCodeSequence2...)
	measuredValueSequence2 = append(measuredValueSequence2, numericValueEle2)

	measuredValueSequenceEle, err := dicom.NewElement(
		tag.MeasuredValueSequence,
		[][]*dicom.Element{measuredValueSequence, measuredValueSequence2},
	)
	if err != nil {
		return nil, err
	}

	return measuredValueSequenceEle, nil
}
