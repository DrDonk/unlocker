$setup=$args[0]
Write-Host "Repairing VMware installation from $setup"
Start-Process -Wait msiexec.exe -ArgumentList "/x `"C:\Program Files (x86)\Common Files\VMware\InstallerCache\{9797B0AF-2570-4326-B9A8-C1B4E75DFFC0}.msi`" /qb"
Start-Process -Wait "$setup" -ArgumentList "/s /v `"EULAS_AGREED=1 REINSTALLMODE=vomus /qb`""
