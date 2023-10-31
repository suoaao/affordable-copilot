@echo off
setlocal

set folder=%userprofile%\AppData\Local\github-copilot
set jsonfile=%folder%\hosts.json

if not exist "%folder%" (
    mkdir "%folder%"
)

echo {"github.com":{"user":"my user name","oauth_token":"i am free","dev_override":{"copilot_token_url":"https://127.0.0.1:8080/copilot_internal/v2/token"}}} > "%jsonfile%"
echo done. please restart your ide.
pause