package main

const GuestTemplate = `//+build guest

// Code generated by OpenFFI. DO NOT EDIT.
// Guest code for {{.ProtoIDLFilename}}

package main

import "C"
import "fmt"
import "github.com/golang/protobuf/proto"

{{range $mindex, $m := .Modules}}
{{if $m.Name}}	import . "{{$m.Name}}"{{end}}
{{end}}

func main(){} // main function must be declared to create dynamic library

func errToOutError(out_err **C.char, out_err_len *C.ulonglong, is_error *C.char, customText string, err error){
	*is_error = 1
	txt := customText+err.Error()
	*out_err = C.CString(txt)
	*out_err_len = C.ulonglong(len(txt))
}

func panicHandler(out_err **C.char, out_err_len *C.ulonglong, is_error *C.char){
	
	if rec := recover(); rec != nil{
		fmt.Println("Caught Panic")

		msg := "Panic in Go function. Panic Data: "
		switch recType := rec.(type){
			case error: msg += (rec.(error)).Error()
			case string: msg += rec.(string)
			default: msg += fmt.Sprintf("Panic with type: %v - %v", recType, rec)
		}

		*is_error = 1
		*out_err = C.CString(msg)
		*out_err_len = C.ulonglong(len(msg))
	}
}

// add functions
{{range $mindex, $m := .Modules}}

{{range $findex, $f := $m.Functions}}

// Call to foreign {{.ForeignFunctionName}}
//export Foreign{{$f.ForeignFunctionName}}
func Foreign{{$f.ForeignFunctionName}}(in_params *C.char, in_params_len C.ulonglong, out_params **C.char, out_params_len *C.ulonglong, out_ret **C.char, out_ret_len *C.ulonglong, is_error *C.char){

	// catch panics and return them as errors
	defer panicHandler(out_ret, out_ret_len, is_error)
	
	*is_error = 0

	// deserialize parameters
	inParams := C.GoStringN(in_params, C.int(in_params_len))
	req := {{$f.ProtobufRequestStruct}}{}
	err := proto.Unmarshal([]byte(inParams), &req)
	if err != nil{
		errToOutError(out_ret, out_ret_len, is_error, "Failed to unmarshal parameters", err)
		return
	}
	
	// call original function
	{{range $index, $elem := $f.ExpandedReturn}}{{if $index}},{{end}}{{$elem.Name}}{{end}}{{if $f.ExpandedReturn}} := {{end}}{{$f.ForeignFunctionName}}({{range $index, $elem := $f.ExpandedParameters}}{{if $index}},{{end}}{{$elem.NameDereferenceIfNeeded "req."}}{{end}})
	
	ret := {{$f.ProtobufResponseStruct}}{}

	// === fill out_ret
	// if one of the returned parameters is of interface type Error, check if error, and if so, return error
	{{range $index, $elem := $f.ExpandedReturn}}
	if err, isError := interface{}({{$elem.Name}}).(error); isError{
		errToOutError(out_ret, out_ret_len, is_error, "Error returned", err)
		return
	} else {
		ret.{{$elem.Name}} = {{$elem.NamePointerIfNeeded ""}}
	}	
	{{end}}

	// serialize results
	serializedRet, err := proto.Marshal(&ret)
	if err != nil{
		errToOutError(out_ret, out_ret_len, is_error, "Failed to marshal return values into protobuf", err)
		return
	}

	// write serialized results to out_ret
	serializedRetStr := string(serializedRet)
	*out_ret = C.CString(serializedRetStr)
	*out_ret_len = C.ulonglong(len(serializedRetStr))

	// === fill out_params
	serializedParams, err := proto.Marshal(&req)
	if err != nil{
		errToOutError(out_ret, out_ret_len, is_error, "Failed to marshal parameter values into protobuf", err)
		return
	}
	
	if out_params != nil && out_params_len != nil{
		// write serialized parameters to out_params
		serializedParamsStr := string(serializedParams)
		*out_params = C.CString(serializedParamsStr)
		*out_params_len = C.ulonglong(len(serializedParamsStr))
	}
	
}

{{end}}

{{end}}

`
