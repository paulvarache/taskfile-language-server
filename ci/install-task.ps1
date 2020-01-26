$clnt = new-object System.Net.WebClient
$url = "https://github.com/go-task/task/releases/download/v2.8.0/task_windows_amd64.zip"
$file = (Get-Location).Path + "\ci\task_windows_amd64.zip"
$clnt.DownloadFile($url,$file)

# Unzip the file to current location
$shell_app=new-object -com shell.application 
$zip_file = $shell_app.namespace((Get-Location).Path + "\ci\task_windows_amd64.zip")
$out = (Get-Location).Path + "\ci"
$destination = $shell_app.namespace($out) 
$destination.Copyhere($zip_file.items())

Write-Host "##vso[task.prependpath]$out"
