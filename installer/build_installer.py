import glob
import os
import zipfile
import io
import base64
import shutil
import re
from typing import List
import platform
import subprocess


def zip_installer_files(files: List[str], root: str):
	# Create a file-like object in memory
	buffer = io.BytesIO()
	
	# Create a zip file object and write the files to it
	# Use the highest compression level
	with zipfile.ZipFile(buffer, "w", zipfile.ZIP_DEFLATED, compresslevel=9) as zf:
		for file in files:
			
			arcname = file
			is_sanity_test_file = 'lang-plugin-go/api/tests/sanity' in arcname
			if 'lang-plugin-go/api/tests/sanity' in arcname:
				arcname = arcname.replace('../../lang-plugin-go/api/tests/sanity', 'tests/go/sanity')
			
			# Read the file from the filename using the parameters and value "root+file"
			# Write the file into the zip with the filename written in "file" value
			if is_sanity_test_file:  # don't use root, as files are not in the METAFFI_HOME dir
				zf.write(file, arcname=arcname)
			else:
				zf.write(root + file, arcname=arcname)
	
	# Get the byte array from the buffer
	return buffer.getvalue()


def update_python_file(python_source_filename, windows_zip, ubuntu_zip):
	# Encode the binary data to base64 strings
	windows_zip_str = base64.b64encode(windows_zip)
	ubuntu_zip_str = base64.b64encode(ubuntu_zip)
	
	# Open the source file in read mode
	with open(python_source_filename, "r") as f:
		# Read the source code as a string
		source_code = f.read()
	
	# Find and replace the variables with the encoded strings
	source_code = re.sub(r"windows_x64_zip\s*=\s*.+", f"windows_x64_zip = {windows_zip_str}", source_code, count=1)
	source_code = re.sub(r"ubuntu_x64_zip\s*=\s*.+", f"ubuntu_x64_zip = {ubuntu_zip_str}", source_code, count=1)
	
	# Open the source file in write mode
	with open(python_source_filename, "w") as f:
		# Write the updated source code to the file
		f.write(source_code)


def get_windows_metaffi_files():
	files = []
	
	# metaffi
	files.extend(['xllr.dll', 'metaffi.exe', 'bin/vcruntime140_1d.dll', 'bin/vcruntime140d.dll', 'bin/boost_filesystem-mt-gd-x64.dll', 'bin/boost_program_options-mt-gd-x64.dll', 'bin/msvcp140d.dll', 'bin/ucrtbased.dll', 'include/cdt_capi_loader.c', 'include/cdt_capi_loader.h', 'include/cdt_structs.h', 'include/metaffi_primitives.h'])
	
	# python plugin
	files.extend(['xllr.python3.dll'])
	
	# go plugin
	files.extend(['xllr.go.dll', 'metaffi.compiler.go.dll'])
	
	# openjdk plugin
	files.extend(['xllr.openjdk.dll', 'xllr.openjdk.bridge.jar', 'xllr.openjdk.jni.bridge.dll'])
	
	# sanity tests
	go_sanity = glob.glob('../../lang-plugin-go/api/tests/sanity/**', recursive=True)
	if len(go_sanity) == 0:
		raise Exception('failed to find Go plugin sanity tests')
	
	go_sanity = [path for path in go_sanity if os.path.isfile(path) and not path.endswith('.pyc')]
	files.extend(go_sanity)
	
	return files


def get_ubuntu_metaffi_files():
	files = []
	
	# metaffi
	files.extend(['xllr.so', 'metaffi', 'lib/libstdc++.so.6.0.30', 'lib/libc.so.6', 'lib/libboost_thread-mt-d-x64.so.1.79.0', 'lib/libboost_program_options-mt-d-x64.so.1.79.0', 'lib/libboost_filesystem-mt-d-x64.so.1.79.0'])
	
	# python plugin
	files.extend(['xllr.python3.so'])
	
	# go plugin
	files.extend(['xllr.go.so', 'metaffi.compiler.go.so'])
	
	# openjdk plugin
	files.extend(['xllr.openjdk.so', 'xllr.openjdk.bridge.jar', 'xllr.openjdk.jni.bridge.so'])
	
	return files


# TODO: When running the executable installer, the test stage in the installer
# that checks if Python installed detects the temporary python within the installer
# def create_executables():
# 	import PyInstaller.__main__
#
# 	# Define the name of your script
# 	script_name = "install_metaffi.py"
#
# 	# make for Windows
# 	if platform.system() == 'Windows':
# 		PyInstaller.__main__.run([
# 			script_name,  # The name of your script
# 			"--uac-uiaccess",  # elevate process
# 			"--onefile",  # Create a single file executable
# 			"--name", "metaffi_installer.exe",  # The name of the output executable
# 		])
# 	else:
# 		print('Running in Ubuntu - skipping making executable installer for windows')
#
# 	# make for ubuntu
# 	if platform.system() == 'Windows':
# 		# NOTICE: assume wsl exists and its python has PyInstaller installed!
# 		command = f"wsl pyinstaller {script_name} --onefile --name metaffi_installer"
# 		# Run the command using subprocess.run
# 		output = subprocess.run(command, capture_output=True, text=True)
# 		if output.returncode != 0:
# 			raise Exception(f'pyinstaller via wsl failed. Error: {output.returncode}.\nstdout:{str(output.stdout)}\nstderr:{str(output.stderr)}')
# 	else:
# 		PyInstaller.__main__.run([
# 			script_name,  # The name of your script
# 			"--onefile",  # Create a single file executable
# 			"--name", "metaffi_installer",  # The name of the output executable
# 		])
#
# 	# cleanup
# 	shutil.rmtree('build')
# 	os.remove('metaffi_installer.exe.spec')
# 	os.remove('metaffi_installer.spec')
# 	shutil.move('dist/metaffi_installer.exe', 'metaffi_installer.exe')
# 	shutil.move('dist/metaffi_installer', 'metaffi_installer')
# 	shutil.rmtree('dist')


def main():
	windows_files = get_windows_metaffi_files()
	ubuntu_files = get_ubuntu_metaffi_files()
	
	windows_zip = zip_installer_files(windows_files, './../out/windows/x64/debug/')
	ubuntu_zip = zip_installer_files(ubuntu_files, './../out/ubuntu/x64/debug/')
	
	shutil.copy('install_metaffi_template.py', 'metaffi_installer.py')
	
	update_python_file('metaffi_installer.py', windows_zip, ubuntu_zip)
	
	
	print('Done')


if __name__ == '__main__':
	main()
