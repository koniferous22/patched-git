package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type FlagSectionParametrizedFlagTerm struct {
	ParametrizedFlag         string `json:"parametrizedFlag"`
	ParametrizedFlagArgument string `json:"parametrizedFlagArgument"`
	Tag                      string `json:"tag"`
}

type FlagSectionOptionalParametrizedFlagTerm struct {
	OptionalParametrizedFlag         string `json:"optionalParametrizedFlag"`
	OptionalParametrizedFlagArgument string `json:"optionalParametrizedFlagArgument"`
	Tag                              string `json:"tag"`
}

type FlagSectionSequenceTerm struct {
	FlagSectionFlagTerm                     string                           `json:"flagSectionFlagTerm,omitempty"`
	FlagSectionFixedTerm                    string                           `json:"flagSectionFixedTerm,omitempty"`
	FlagSectionSubstitutionTerm             string                           `json:"flagSectionSubstitutionTerm,omitempty"`
	FlagSectionSubstitutionKeyValueTerm     string                           `json:"flagSectionSubstitutionKeyValueTerm,omitempty"`
	FlagSectionParametrizedFlagTerm         *FlagSectionParametrizedFlagTerm `json:"flagSectionParametrizedFlagTerm,omitempty"`
	FlagSectionOptionalParametrizedFlagTerm *FlagSectionParametrizedFlagTerm `json:"flagSectionOptionalParametrizedFlagTerm,omitempty"`
	Tag                                     string                           `json:"tag"`
}

type FlagSectionSequence struct {
	FlagSectionSequenceTerms []FlagSectionSequenceTerm `json:"flagSectionSequenceTerms"`
}

type SynopsisFlagSection struct {
	FlagSectionSequences []FlagSectionSequence `json:"flagSectionSequences"`
	Tag                  string                `json:"tag"`
}

type SynopsisTerm struct {
	SynopsisFixedTerm        string               `json:"synopsisFixedTerm,omitempty"`
	SynopsisSubstitutionTerm string               `json:"synopsisSubstitutionTerm,omitempty"`
	SynopsisFlagSection      *SynopsisFlagSection `json:"synopsisFlagSection,omitempty"`
	Tag                      string               `json:"tag"`
}

type SynopsisJSON struct {
	SynopsisTerms []SynopsisTerm `json:"synopsisTerms"`
}

func ParseSynopsisJson(fp string) (*SynopsisJSON, error) {
	handleError := func(err error) error {
		return fmt.Errorf("error parsing synopsis\n%w", err)
	}
	contents, err := os.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("error reading synopsis JSON file:\n%w", err)
	}
	var data SynopsisJSON
	err = json.Unmarshal(contents, &data)
	if err != nil {
		return nil, handleError(fmt.Errorf("error parsing synopsis JSON data:\n%w", err))
	}
	return &data, nil
}

func (flagSectionSequence FlagSectionSequence) GetSpaceDelimitedFlagOptions() (*[]string, error) {
	if len(flagSectionSequence.FlagSectionSequenceTerms) == 0 {
		return nil, fmt.Errorf("empty \"FlagSectionSequence\"")
	}
	var result []string
	for i, sequenceTerm := range flagSectionSequence.FlagSectionSequenceTerms[1:] {
		previousTerm := flagSectionSequence.FlagSectionSequenceTerms[i]
		if sequenceTerm.Tag == "FlagSectionSubstitutionTerm" && previousTerm.Tag == "FlagSectionFlagTerm" {
			value := previousTerm.FlagSectionFlagTerm
			if value == "" {
				return nil, fmt.Errorf("encountered empty flag term")
			}
			result = append(result, value)
		}
	}
	return &result, nil
}

func (sj SynopsisJSON) GetSpaceDelimitedFlagOptions() (*[]string, error) {
	handleError := func(err error) error {
		return fmt.Errorf("invalid synopsis shape\n%w", err)
	}
	var result []string
	for _, synopsisTerm := range sj.SynopsisTerms {
		if synopsisTerm.Tag == "SynopsisFlagSection" {
			if synopsisTerm.SynopsisFlagSection == nil {
				return nil, handleError(fmt.Errorf("invalid \"SynopsisFlagSection\" tag"))
			}
			for _, flagSectionSequence := range synopsisTerm.SynopsisFlagSection.FlagSectionSequences {
				foundTerms, err := flagSectionSequence.GetSpaceDelimitedFlagOptions()
				if err != nil {
					return nil, handleError(err)
				}
				result = append(result, *foundTerms...)
			}
		}
	}
	return &result, nil
}
