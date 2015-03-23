// list of json errors
package main

var (
	ERROR_INVALID_PATH     = M{"errors": M{"id": 1, "title": "Invalid Path", "detail": "Path was not able to validate to a valid location."}}
	ERROR_NO_OS_PERMISSION = M{"errors": M{"id": 2, "title": "No Permission", "detail": "Did not have enough OS permission to complete this action"}}
	ERROR_ZIP_COMMAND      = M{"errors": M{"id": 3, "title": "Zip Error", "detail": "Could not create the zip archive"}}
)
