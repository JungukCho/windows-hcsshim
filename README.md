## Explore hcsshim for connectivity between windows containers on windows 
1. This repo walks through hnsshim (hcn) to play with connectivity between windows containers on windows 2022 preview. This page are compiled based on resources from NPM and hcsshim teams.

## Prerequisite for setup win 2022
1. [Git](https://git-scm.com/download/win)
2. [Golang](https://golang.org/dl/)
3. vscode - go extension


## Set up docker on win 2022
1. [Get started: Prep Windows for containers](https://docs.microsoft.com/en-us/virtualization/windowscontainers/quick-start/set-up-environment?tabs=Windows-Server)


## Play with HNS to control networks
1. Create two dockers
* One for a server
  - [Many windows container examples](https://github.com/MicrosoftDocs/Virtualization-Documentation.git)
```Dockerfile
# This dockerfile utilizes components licensed by their respective owners/authors.
# Prior to utilizing this file or resulting images please review the respective licenses at: http://nginx.org/LICENSE
FROM mcr.microsoft.com/windows/server/insider:10.0.20344.1
LABEL Description="Nginx" Vendor="Nginx" Version="1.0.13"
RUN powershell -Command \
        $ErrorActionPreference = 'Stop'; \
        Invoke-WebRequest -Method Get -Uri http://nginx.org/download/nginx-1.9.13.zip -OutFile c:\nginx-1.9.13.zip ; \
        Expand-Archive -Path c:\nginx-1.9.13.zip -DestinationPath c:\ ; \
        Remove-Item c:\nginx-1.9.13.zip -Force

WORKDIR /nginx-1.9.13
CMD ["/nginx-1.9.13/nginx.exe"]
```
* Build docker and run it
```shell
docker build -t nginx -f Dockerfile .
docker run --name nginx -d nginx 
```

* One for a client - to keep running docker container
```shell
docker run -d --name client  mcr.microsoft.com/windows/server/insider:10.0.20344.1 ping -t localhost
```

* Test connectivity between docker
```shell
docker exec -it client cmd.exe
# you can access to console of nginx docker and type ipconfig
> powershell
> curl.exe <nginx ip address>
```
* Print current endpoint status
  - nginx endpoint id : `edffc876-3051-46ec-80b3-820dda20aa14`
  - client endpoint id : `c047102f-2bb1-4f64-8f30-9ede02f464f1`
```shell
# on powershell
PS C:\Users\azure> Get-HnsEndpoint


ID                 : edffc876-3051-46ec-80b3-820dda20aa14
Name               : Ethernet
Version            : 64424509440
AdditionalParams   :
Resources          : @{AdditionalParams=; AllocationOrder=3; Allocators=System.Object[]; CompartmentOperationTime=0; Flags=0;
                     Health=; ID=78E66F4E-E634-414C-B26B-3FCC84A4A7B2; PortOperationTime=0; State=1; SwitchOperationTime=0;
                     VfpOperationTime=0; parentId=5CA27F17-3705-4BE7-8F3B-963E16533E12}
State              : 2
VirtualNetwork     : 625d8985-3f12-469e-a98d-2baac83dee45
VirtualNetworkName : nat
MacAddress         : 00-15-5D-90-AD-7C
EnableInternalDNS  : True
IPAddress          : 172.19.6.74
PrefixLength       : 20
GatewayAddress     : 172.19.0.1
IPSubnetId         : 35f30e75-364d-4e6d-a459-4ccdce1077e7
DNSServerList      : 172.19.0.1,168.63.129.16
DNSSuffix          : uuz1fvwpe0zejibongrs2s5cjg.bx.internal.cloudapp.net
SharedContainers   : {49c836b9725e69736157468f33b2a41b0753609c5d49565b46fb37505a5a3b9f}

ID                 : c047102f-2bb1-4f64-8f30-9ede02f464f1
Name               : Ethernet
Version            : 64424509440
AdditionalParams   :
Resources          : @{AdditionalParams=; AllocationOrder=3; Allocators=System.Object[]; CompartmentOperationTime=0; Flags=0;
                     Health=; ID=4C355094-4C39-4B90-9A7E-D3D71BAD6CA5; PortOperationTime=0; State=1; SwitchOperationTime=0;
                     VfpOperationTime=0; parentId=5CA27F17-3705-4BE7-8F3B-963E16533E12}
State              : 2
VirtualNetwork     : 625d8985-3f12-469e-a98d-2baac83dee45
VirtualNetworkName : nat
MacAddress         : 00-15-5D-90-AD-B2
EnableInternalDNS  : True
IPAddress          : 172.19.3.166
PrefixLength       : 20
GatewayAddress     : 172.19.0.1
IPSubnetId         : 35f30e75-364d-4e6d-a459-4ccdce1077e7
DNSServerList      : 172.19.0.1,168.63.129.16
DNSSuffix          : uuz1fvwpe0zejibongrs2s5cjg.bx.internal.cloudapp.net
SharedContainers   : {49e8cbcc30485056cc3e95a6d39f432a8ec2dda522efaaabb2f8dd513b7a5f9c}
```



3. Run simple go example to use hns to control network policy. 
It needs some manual option to change the code. Please check `main.go`
```shell
go build -o hns main.go
```

* Output comparison after blocking and allowing traffic
```
* Output from `Get-HnsEndpoint` after blocking and allowing traffic
  - nginx endpoint id : `edffc876-3051-46ec-80b3-820dda20aa14`
  - client endpoint id : `c047102f-2bb1-4f64-8f30-9ede02f464f1`
```shell
# After block traffic
PS C:\k\debug> Get-HnsEndpoint
ID                 : edffc876-3051-46ec-80b3-820dda20aa14
Name               : Ethernet
Version            : 64424509440
AdditionalParams   :
Resources          : @{AdditionalParams=; AllocationOrder=6; Allocators=System.Object[]; CompartmentOperationTime=0; Flags=0; Health=; ID=78E66F4E-E634-414C-B26B-3FCC84A4A7B2;
                     PortOperationTime=0; State=1; SwitchOperationTime=0; VfpOperationTime=0; parentId=5CA27F17-3705-4BE7-8F3B-963E16533E12}
State              : 2
VirtualNetwork     : 625d8985-3f12-469e-a98d-2baac83dee45
VirtualNetworkName : nat
Policies           : {@{Action=Block; Direction=In; LocalAddresses=172.19.6.74; LocalPorts=80; Priority=100; Protocols=6; Scope=0; Type=ACL}, @{Action=Block; Direction=Out;
                     Priority=100; Protocols=6; Scope=0; Type=ACL}, @{Action=Block; Direction=In; LocalAddresses=172.19.6.74; LocalPorts=80; Priority=100; Protocols=6; Scope=0;
                     Type=ACL}, @{Action=Block; Direction=Out; Priority=100; Protocols=6; Scope=0; Type=ACL}}
MacAddress         : 00-15-5D-90-AD-7C
EnableInternalDNS  : True
IPAddress          : 172.19.6.74
PrefixLength       : 20
GatewayAddress     : 172.19.0.1
IPSubnetId         : 35f30e75-364d-4e6d-a459-4ccdce1077e7
DNSServerList      : 172.19.0.1,168.63.129.16
DNSSuffix          : uuz1fvwpe0zejibongrs2s5cjg.bx.internal.cloudapp.net
SharedContainers   : {49c836b9725e69736157468f33b2a41b0753609c5d49565b46fb37505a5a3b9f}

ID                 : c047102f-2bb1-4f64-8f30-9ede02f464f1
Name               : Ethernet
Version            : 64424509440
AdditionalParams   :
Resources          : @{AdditionalParams=; AllocationOrder=3; Allocators=System.Object[]; CompartmentOperationTime=0; Flags=0; Health=; ID=4C355094-4C39-4B90-9A7E-D3D71BAD6CA5;
                     PortOperationTime=0; State=1; SwitchOperationTime=0; VfpOperationTime=0; parentId=5CA27F17-3705-4BE7-8F3B-963E16533E12}
State              : 2
VirtualNetwork     : 625d8985-3f12-469e-a98d-2baac83dee45
VirtualNetworkName : nat
MacAddress         : 00-15-5D-90-AD-B2
EnableInternalDNS  : True
IPAddress          : 172.19.3.166
PrefixLength       : 20
GatewayAddress     : 172.19.0.1
IPSubnetId         : 35f30e75-364d-4e6d-a459-4ccdce1077e7
DNSServerList      : 172.19.0.1,168.63.129.16
DNSSuffix          : uuz1fvwpe0zejibongrs2s5cjg.bx.internal.cloudapp.net
SharedContainers   : {49e8cbcc30485056cc3e95a6d39f432a8ec2dda522efaaabb2f8dd513b7a5f9c}


# After allow traffic
PS C:\k\debug> Get-HnsEndpoint
ID                 : edffc876-3051-46ec-80b3-820dda20aa14
Name               : Ethernet
Version            : 64424509440
AdditionalParams   :
Resources          : @{AdditionalParams=; AllocationOrder=7; Allocators=System.Object[]; CompartmentOperationTime=0; Flags=0; Health=; ID=78E66F4E-E634-414C-B26B-3FCC84A4A7B2;
                     PortOperationTime=0; State=1; SwitchOperationTime=0; VfpOperationTime=0; parentId=5CA27F17-3705-4BE7-8F3B-963E16533E12}
State              : 2
VirtualNetwork     : 625d8985-3f12-469e-a98d-2baac83dee45
VirtualNetworkName : nat
Policies           : {@{Action=Allow; Direction=In; LocalAddresses=172.19.6.74; LocalPorts=80; Priority=100; Protocols=6; Scope=0; Type=ACL}, @{Action=Allow; Direction=Out;
                     Priority=100; Protocols=6; Scope=0; Type=ACL}}
MacAddress         : 00-15-5D-90-AD-7C
EnableInternalDNS  : True
IPAddress          : 172.19.6.74
PrefixLength       : 20
GatewayAddress     : 172.19.0.1
IPSubnetId         : 35f30e75-364d-4e6d-a459-4ccdce1077e7
DNSServerList      : 172.19.0.1,168.63.129.16
DNSSuffix          : uuz1fvwpe0zejibongrs2s5cjg.bx.internal.cloudapp.net
SharedContainers   : {49c836b9725e69736157468f33b2a41b0753609c5d49565b46fb37505a5a3b9f}

ID                 : c047102f-2bb1-4f64-8f30-9ede02f464f1
Name               : Ethernet
Version            : 64424509440
AdditionalParams   :
Resources          : @{AdditionalParams=; AllocationOrder=3; Allocators=System.Object[]; CompartmentOperationTime=0; Flags=0; Health=; ID=4C355094-4C39-4B90-9A7E-D3D71BAD6CA5;
                     PortOperationTime=0; State=1; SwitchOperationTime=0; VfpOperationTime=0; parentId=5CA27F17-3705-4BE7-8F3B-963E16533E12}
State              : 2
VirtualNetwork     : 625d8985-3f12-469e-a98d-2baac83dee45
VirtualNetworkName : nat
MacAddress         : 00-15-5D-90-AD-B2
EnableInternalDNS  : True
IPAddress          : 172.19.3.166
PrefixLength       : 20
GatewayAddress     : 172.19.0.1
IPSubnetId         : 35f30e75-364d-4e6d-a459-4ccdce1077e7
DNSServerList      : 172.19.0.1,168.63.129.16
DNSSuffix          : uuz1fvwpe0zejibongrs2s5cjg.bx.internal.cloudapp.net
SharedContainers   : {49e8cbcc30485056cc3e95a6d39f432a8ec2dda522efaaabb2f8dd513b7a5f9c}

```

## Tested windows docker images
```shell
1.. Compatible containers
* [Announcing a New Windows Server Container Image Preview
](https://techcommunity.microsoft.com/t5/containers/announcing-a-new-windows-server-container-image-preview/ba-p/2304897?ranMID=24542&ranEAID=je6NUbpObpQ&ranSiteID=je6NUbpObpQ-7G741ImEBgESfTcCBeFawQ&epi=je6NUbpObpQ-7G741ImEBgESfTcCBeFawQ&irgwc=1&OCID=AID2200057_aff_7593_1243925&tduid=(ir__ms1mv9qee9kfqib1cydiwswhr22xr3pgbcmzyur300)(7593)(1243925)(je6NUbpObpQ-7G741ImEBgESfTcCBeFawQ)()&irclickid=_ms1mv9qee9kfqib1cydiwswhr22xr3pgbcmzyur300)
* [Nano Server Insider](https://hub.docker.com/_/microsoft-windows-nanoserver-insider)
* [Windows Server Insider](https://hub.docker.com/_/microsoft-windows-server-insider)
* [Windows Server Core Insider](https://hub.docker.com/_/microsoft-windows-servercore-insider)
```