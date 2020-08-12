package openffi_cli

import (
	"fmt"
	"strings"
	"text/template"
)

//--------------------------------------------------------------------
type TemplateParameters struct {
	ProtoIDLFilename       		string
	ProtoIDLFilenameNoExtension string
	ProtobufFilename       		string
	TargetLanguage				string
	protoTypeToTargetType 		func(protoType string, isarray bool)(typeStr string, isComplexType bool)
	modifyParameterName			func(string)string

	Modules []*TemplateModuleParameters
}
//--------------------------------------------------------------------
type TemplateModuleParameters struct {
	Name string
	Functions []*TemplateFunctionParameters
}
//--------------------------------------------------------------------
type TemplateFunctionParameters struct {
	ForeignFunctionName    		string
	OriginalForeignFunctionName string
	ProtobufRequestStruct  		string
	ProtobufResponseStruct 		string
	ExpandedParameters     []*TemplateFunctionParameterData
	ExpandedReturn         []*TemplateFunctionParameterData
}
//--------------------------------------------------------------------
type TemplateFunctionParameterData struct {
	Name string
	Type string
	IsComplex bool
	ParamPass *PassMethod
	IsArray bool
}
func (this *TemplateFunctionParameterData) NameDereferenceIfNeeded(prefix string) string{
	if this.IsComplex && !this.IsArray{
		return "*"+prefix+this.Name
	} else {
		return prefix+this.Name
	}
}
func (this *TemplateFunctionParameterData) NamePointerIfNeeded(prefix string) string{
	if this.IsComplex && !this.IsArray{
		return "&"+prefix+this.Name
	} else {
		return prefix+this.Name
	}
}
func (this *TemplateFunctionParameterData) TypePointerIfNeeded(prefix string) string{
	if this.IsComplex && !this.IsArray{
		return "*"+prefix+this.Type
	} else {
		return prefix+this.Type
	}
}
//--------------------------------------------------------------------
func NewTemplateParameters(protoIDLFilename string,
							protobufFilenameSuffix string,
							targetLanguage string,
							protoTypeToTargetType func(string, bool)(string, bool),
							modifyParameterName func(string)string) (*TemplateParameters, error){

	extensionIndex := strings.LastIndex(protoIDLFilename, ".")
	if extensionIndex == -1{
		return nil, fmt.Errorf("Cannot find extension in proto filename: %v", protoIDLFilename)
	}

	protoFilenameWithoutExtension := protoIDLFilename[:extensionIndex]

	if modifyParameterName == nil{ // set default modifyParameterName() to not change anything
		modifyParameterName = func(p string)string{ return p }
	}

	gtp := &TemplateParameters{
		ProtoIDLFilename: protoIDLFilename,
		ProtoIDLFilenameNoExtension: protoFilenameWithoutExtension,
		ProtobufFilename: protoFilenameWithoutExtension + protobufFilenameSuffix,
		TargetLanguage: targetLanguage,
		protoTypeToTargetType: protoTypeToTargetType,
		modifyParameterName: modifyParameterName,
	}

	gtp.Modules = make([]*TemplateModuleParameters, 0)

	return gtp, nil
}
//--------------------------------------------------------------------
func NewTemplateFunctionParameterData(p *ParameterData,
									protoTypeToTargetType func(string, bool)(string, bool),
									modifyParameterName func(string)string) *TemplateFunctionParameterData{

	htfp := &TemplateFunctionParameterData{
		Name: modifyParameterName(p.Name),
	}

	htfp.Type, htfp.IsComplex = protoTypeToTargetType(p.Type, p.IsArray)
	htfp.IsArray = p.IsArray
	htfp.ParamPass = p.PassParam

	return htfp
}
//--------------------------------------------------------------------
func (this *TemplateParameters) AddModule(m *Module){

	// add modules
	modParams := &TemplateModuleParameters{
		Name:      m.Name,
		Functions: make([]*TemplateFunctionParameters, 0),
	}

	// for each module, add the function

	for _, f := range m.Functions{

		funcParams := &TemplateFunctionParameters{
			ForeignFunctionName: this.modifyParameterName(f.Name),
			OriginalForeignFunctionName: f.Name,
			ProtobufRequestStruct: this.modifyParameterName(f.RequestName),
			ProtobufResponseStruct: this.modifyParameterName(f.ResponseName),
			ExpandedParameters: make([]*TemplateFunctionParameterData, 0),
			ExpandedReturn: make([]*TemplateFunctionParameterData, 0),
		}

		// generate parameters
		for _, p := range f.Parameters{
			funcParams.ExpandedParameters = append(funcParams.ExpandedParameters,
														NewTemplateFunctionParameterData(p, this.protoTypeToTargetType, this.modifyParameterName))
		}

		for _, r := range f.Return{
			funcParams.ExpandedReturn = append(funcParams.ExpandedReturn,
														NewTemplateFunctionParameterData(r, this.protoTypeToTargetType, this.modifyParameterName))
		}

		modParams.Functions = append(modParams.Functions, funcParams)
	}

	this.Modules = append(this.Modules, modParams)
}
//--------------------------------------------------------------------
func (this *TemplateParameters) Generate(templateName string, templateText string) (string, error){

	if this.ProtoIDLFilename == ""{
		return "", fmt.Errorf("ProtoIDLFilename is empty")
	}

	if this.ProtobufFilename == ""{
		return "", fmt.Errorf("ProtobufFilename is empty")
	}

	if this.Modules == nil || len(this.Modules) == 0{
		return "", fmt.Errorf("No functions defined")
	}

	for _, m := range this.Modules {

		for _, f := range m.Functions {
			if f.ForeignFunctionName == "" {
				return "", fmt.Errorf("ForeignFunctionName is empty")
			}

			if f.ProtobufRequestStruct == "" {
				return "", fmt.Errorf("ProtobufRequestStruct is empty")
			}

			if f.ProtobufResponseStruct == "" {
				return "", fmt.Errorf("ProtobufResponseStruct is empty")
			}

			if f.ExpandedParameters == nil {
				return "", fmt.Errorf("ExpandedParameters is nil")
			}

			if f.ExpandedParameters == nil {
				return "", fmt.Errorf("ExpandedParameters is nil")
			}
		}
	}


	temp, err := template.New(templateName).Parse(templateText)
	if err != nil{
		return "", fmt.Errorf("Failed to parse template \"%v\", err: %v", templateName, err)
	}

	strbuf := strings.Builder{}

	err = temp.Execute(&strbuf, this)
	if err != nil{
		return "", fmt.Errorf("Failed to execute guest template, err: %v", err)
	}

	return strbuf.String(), nil
}
//--------------------------------------------------------------------


