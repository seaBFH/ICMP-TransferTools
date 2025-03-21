FROM ubuntu:22.04
# add network tools, tcpdump,curl,bcc-tools,net-tools
RUN apt-get update && apt-get install -y iputils-ping tcpdump curl bpfcc-tools net-tools python3 python3-pip
RUN pip3 install impacket
# add conda
ADD ./ICMP-ReceiveFile.py /tools/ICMP-ReceiveFile.py
ADD ./ICMP-SendFile.py /tools/ICMP-SendFile.py
ADD ./Invoke-IcmpDownload.ps1 /tools/Invoke-IcmpDownload.ps1
ADD ./Invoke-IcmpUpload.ps1 /tools/Invoke-IcmpUpload.ps1
ADD ./invoke/dist /tools/invoke
ENTRYPOINT ["sleep", "infinity"]


