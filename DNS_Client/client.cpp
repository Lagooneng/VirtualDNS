#include<stdio.h>
#include<string.h>
#include<unistd.h>
#include<sys/socket.h>
#include<arpa/inet.h>
#include<iostream>
#define MAXBUFFER 100
using namespace std;

string getDNS(char *dns);
string regDNS(char *dns, char *ip);

int main(int argc, char *argv[]) {
    if( argc != 3 && argc != 4) {
        cout << "USAGE 1: ./client -get domainname" << endl;
        cout << "USAGE 2: ./client -reg domainname ip" << endl;
        return -1;
    }
    
    if( strcmp( argv[1], "-get" ) == 0 ) {     // get
        if( argc != 3 ) {
            cout << "USAGE 1: ./client -get domainname" << endl;
            return -1;
        }

        string str = getDNS( argv[2] );
        cout << str << endl;
    }
    else if( strcmp( argv[1], "-reg" ) == 0 ){ // register
        if( argc != 4 ) {
            cout << "USAGE 2: ./client -reg domainname ip" << endl;
            return -1;
        }

        string str = regDNS(argv[2], argv[3]);
        cout << str << endl;
    }
    
    return 0;
}

string getDNS(char *dns) {
    int cs;
    char recvBuf[100];
    struct sockaddr_in csa;

    memset(&csa, 0, sizeof(csa));
    csa.sin_family = AF_INET;
    csa.sin_addr.s_addr = inet_addr("127.0.0.1");
    csa.sin_port = htons(8000);

    cs = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
    connect(cs, (struct sockaddr *) &csa, sizeof(csa));

    string str = dns;
    str = "0" + str;
    int length = str.size();
    
    send(cs, str.c_str(), length + 1, 0);

    recv(cs, recvBuf, 100, 0);

    close(cs);

    return recvBuf;
}


string regDNS(char *dns, char *ip) {
    int cs;
    char recvBuf[100];
    string sendBuf = "";
    struct sockaddr_in csa;

    unsigned short dnLength, ipLength;  // 각각 2바이트 
    string dns_str = dns;
    string ip_str = ip;

    dnLength = (short)dns_str.size();
    ipLength = (short)ip_str.size();


    memset(&csa, 0, sizeof(csa));
    csa.sin_family = AF_INET;
    csa.sin_addr.s_addr = inet_addr("127.0.0.1");
    csa.sin_port = htons(8000);

    cs = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP);
    connect(cs, (struct sockaddr *) &csa, sizeof(csa));

    sendBuf += "1";
    sendBuf += ((dnLength >> 8) & 0xFF);
    sendBuf += (dnLength & 0xFF);
    sendBuf += ((ipLength >> 8) & 0xFF);
    sendBuf += (ipLength & 0xFF);

    sendBuf += dns;
    sendBuf += ip;

    int length = sendBuf.size();
    
    send(cs, sendBuf.c_str(), length + 1, 0);

    recv(cs, recvBuf, 100, 0);

    close(cs);

    return recvBuf;
}