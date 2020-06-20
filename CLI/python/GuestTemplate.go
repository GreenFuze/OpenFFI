package main

import (
	"fmt"
	"strings"
	"text/template"
)

const GuestTemplate = `
# Code generated by OpenFFI. DO NOT EDIT.
# Guest code for {{.ProtoIDLFilename}}

from {{.ProtobufFilename}} import *
from {{.ForeignFunctionFilename}} import {{.ForeignFunctionName}} as Foreign{{.ForeignFunctionName}}

# Call to foreign {{.ForeignFunctionName}}
def {{.ForeignFunctionName}}(paramsVal: bytes) -> ReturnVal:
	
	# TODO: try/catch?

	req = {{.ProtobufRequestStruct}}()
	req.ParseFromString(str(paramsVal))

	ret = {{.ProtobufResponseStruct}}()

	# python method to call a function without knowing its parameter names?

	{{range $index, $elem := .ExpandedReturn}}{{if $index}},{{end}} res.{{$elem}}{{end}} = {{.ForeignFunctionFilename}}.Foreign{{.ForeignFunctionName}}({{range $index, $elem := .ExpandedParameters}}{{if $index}},{{end}} req.{{$elem}} {{end}})

	return bytes(ret.SerializeToString(), 'utf-8')
`

//--------------------------------------------------------------------
type GuestTemplateGenerator struct {
	ProtoIDLFilename string
	ProtobufFilename string
	ForeignFunctionFilename string
	ForeignFunctionName string
	ProtobufRequestStruct string
	ProtobufResponseStruct string
	ExpandedParameters []string
	ExpandedReturn []string
}
//--------------------------------------------------------------------
func (this *GuestTemplateGenerator) Generate() (string, error){

	if this.ProtoIDLFilename == ""{
		return "", fmt.Errorf("ProtoIDLFilename is empty")
	}

	if this.ProtobufFilename == ""{
		return "", fmt.Errorf("ProtobufFilename is empty")
	}

	if this.ForeignFunctionFilename == ""{
		return "", fmt.Errorf("ForeignFunctionFilename is empty")
	}

	if this.ForeignFunctionName == ""{
		return "", fmt.Errorf("ForeignFunctionName is empty")
	}

	if this.ProtobufRequestStruct == ""{
		return "", fmt.Errorf("ProtobufRequestStruct is empty")
	}

	if this.ProtobufResponseStruct == ""{
		return "", fmt.Errorf("ProtobufResponseStruct is empty")
	}

	if this.ExpandedParameters == nil{
		return "", fmt.Errorf("ExpandedParameters is nil")
	}

	if this.ExpandedParameters == nil{
		return "", fmt.Errorf("ExpandedParameters is nil")
	}

	temp, err := template.New("Guest").Parse(GuestTemplate)
	if err != nil{
		return "", fmt.Errorf("Failed to parse GuestTemplate, err: %v", err)
	}

	strbuf := strings.Builder{}

	err = temp.Execute(&strbuf, this)
	if err != nil{
		return "", fmt.Errorf("Failed to execute guest template, err: %v", err)
	}

	return strbuf.String(), nil
}
//--------------------------------------------------------------------