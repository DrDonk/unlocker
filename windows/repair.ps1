$setup=$args[0]
Write-Host "Repairing VMware installation from $setup"
$installpath=Get-ItemPropertyValue -Path 'HKLM:\SOFTWARE\WOW6432Node\VMware, Inc.\VMware Player' -Name InstallPath
Write-Host "$installpath"
$installpath64=Get-ItemPropertyValue -Path 'HKLM:\SOFTWARE\WOW6432Node\VMware, Inc.\VMware Player' -Name InstallPath64
Write-Host "$installpath64"
Stop-Process -Name vmware-tray -Force -ErrorAction SilentlyContinue
Stop-Service -Name VMAuthdService -Force
(Get-Service VMAuthdService).WaitForStatus('Stopped')
Stop-Service -Name VMUSBArbService -Force
(Get-Service VMUSBArbService).WaitForStatus('Stopped')
$service = Get-Service -Name VMwareHostd -ErrorAction SilentlyContinue
if($service -ne $null)
{
    Stop-Service -Name VMwareHostd -Force
    (Get-Service VMwareHostd).WaitForStatus('Stopped')
}

$vmwarevmx=Join-Path -Path $installpath64 -ChildPath "vmware-vmx.exe"
Write-Host $vmwarevmx
Remove-Item -Path $vmwarevmx -Force

$vmwarevmxdebug=Join-Path -Path $installpath64 -ChildPath "vmware-vmx-debug.exe"
Write-Host $vmwarevmxdebug
Remove-Item -Path $vmwarevmxdebug -Force

$vmwarevmxstats=Join-Path -Path $installpath64 -ChildPath "vmware-vmx-stats.exe"
$isstats=Get-Item $vmwarevmxstats
if($isstats -ne $null) {
    Write-Host $vmwarevmxstats
    Remove-Item -Path $vmwarevmxstats -Force
}
$vmwarebase=Join-Path -Path $installpath -ChildPath "vmwarebase.dll"
Write-Host $vmwarebase
Remove-Item -Path $vmwarebase -Force
Start-Process -Wait $setup -ArgumentList "/s /v `"EULAS_AGREED=1 REINSTALLMODE=vomus /qn`""
