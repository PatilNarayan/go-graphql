@echo off
REM Set environment variables
set PERMIT_PROJECT=dev
set PERMIT_ENV=dev
set PERMIT_TOKEN=permit_key_mXEiRTTrdC78tEWyjy2dazPGFruDrhJswiNthVQ0YYqjT7VWH1dn9Wiw3vtVFUl25tweXwICha30RsPe0XRwwU
set PERMIT_PDP_ENDPOINT=https://api.permit.io
set DB_USERNAME=root
set DB_PASSWORD=root
set DB_HOST=localhost
set DB_PORT=3306
set DB_NAME=root

REM Optional: Display the set variables (for debugging)
echo Environment Variables Set:
echo PERMIT_PROJECT=%PERMIT_PROJECT%
echo PERMIT_ENV=%PERMIT_ENV%
echo PERMIT_PDP_ENDPOINT=%PERMIT_PDP_ENDPOINT%
echo DB_USERNAME=%DB_USERNAME%
echo DB_PASSWORD=%DB_PASSWORD%
echo DB_HOST=%DB_HOST%
echo DB_NAME=%DB_NAME%
echo DB_PORT=%DB_PORT%
