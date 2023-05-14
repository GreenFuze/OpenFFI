#pragma once
#include <compiler/idl_plugin_interface.h>
#include <memory>
#include <type_traits>
#include <boost/dll.hpp>

//--------------------------------------------------------------------
class idl_plugin_interface_wrapper : public idl_plugin_interface
{
public:
	static std::shared_ptr<idl_plugin_interface_wrapper> load(const std::string& idl_name);
	
	
private:
	std::shared_ptr<boost::dll::shared_library> mod;
	std::shared_ptr<boost::dll::detail::import_type<void(const char*, uint32_t, const char*, uint32_t, char**, uint32_t*, char**, uint32_t*)>::type> pparse_idl;
	std::shared_ptr<boost::dll::detail::import_type<void(void)>::type> pinit_plugin;
	
protected:
	idl_plugin_interface_wrapper() = default;
	
public:
	explicit idl_plugin_interface_wrapper(const std::string& idl_extension);

	void init() override;

	virtual void parse_idl(const std::string& idl_name,
	               const std::string& idl,
	               std::string& out_idl_def_json,
	               std::string& out_err);
	
	/**
	 * Compile IDL to code called from XLLR to the foreign function
	 */
	void parse_idl(const char* idl_name, uint32_t idl_name_length,
	               const char* idl, uint32_t idl_length,
	               char** out_idl_def_json, uint32_t* out_idl_def_json_length,
	               char** out_err, uint32_t* out_err_len) override;
};
//--------------------------------------------------------------------
class json_idl_plugin : public idl_plugin_interface_wrapper
{
public:
	explicit json_idl_plugin() = default;
	void init() override{}
	
	void parse_idl(const std::string& idl_name,
	               const std::string& idl,
	               std::string& out_idl_def_json,
	               std::string& out_err) override
	{
		out_idl_def_json = idl;
		out_err.clear();
	}
	
	void parse_idl(const char* idl_name, uint32_t idl_name_length,
	               const char* idl, uint32_t idl_length,
	               char** out_idl_def_json, uint32_t* out_idl_def_json_length,
	               char** out_err, uint32_t* out_err_len) override
	{
		*out_idl_def_json = const_cast<char*>(idl);
		*out_idl_def_json_length = idl_length;
		*out_err = nullptr;
		*out_err_len = 0;
	}
};
//--------------------------------------------------------------------