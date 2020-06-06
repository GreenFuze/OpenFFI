#include "utils.h"
#include "../scope_guard.hpp"
#include <python3.7/Python.h>

std::string get_py_error(void)
{
	// error has occurred
	PyObject* ptype = nullptr;
	PyObject* pvalue = nullptr;
	PyObject* ptraceback = nullptr;
	PyErr_Fetch(&ptype, &pvalue, &ptraceback);
	
	scope_guard sg([&]()
	{
		if(ptype){Py_DecRef(ptype);}
		if(pvalue){Py_DecRef(pvalue);}
		if(ptraceback){Py_DecRef(ptraceback);}
	});

	if(!pvalue)
	{
		printf("******* NO ERROR FOUND?!\n");
		return std::string("No Error was found!");
	}

	// if pvalue is str, get text, if not, get msg from exception
	if(strcmp(pvalue->ob_type->tp_name, "str") == 0)
	{
		uint64_t size;
		const char* pmsg = PyUnicode_AsUTF8AndSize(pvalue, (Py_ssize_t*)&size);

		return std::string(pmsg, size);
	}
	

	PyObject* msg = PyObject_CallMethod(pvalue, "__str__", nullptr);
	if(!msg)
	{
		return std::string("Failed to get error description: call to method \"__str__\" failed");
	}

	scope_guard sgmsg([&]()
	{
		if(msg){Py_DecRef(msg);}
	});

	uint64_t size;
	const char* pmsg = PyUnicode_AsUTF8AndSize(msg, (Py_ssize_t*)&size);
	return std::string(pmsg, size);
}