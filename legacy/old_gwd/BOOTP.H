// begin lbr    only for host    
#define INVITE 0x02
#define ACK    0x03
#define REQUEST 0x01
#define REPLY   0x00
#define DAPRESET   0x05
#define DAPBEEP   0x07
#define NAK    0x08

#if DEF_MAC_ENCODE >= 1
#define MAC_ENCODE  0xff
#endif

#define BOOTP_FRAME_LEN 300
#define HW_ADDRESS_LEN 6    
#define ETHER_TYPE 0x01          
#define MY_TRANSACTION_ID 0xda92
#define SERVER_HOST_NAME "my dap boot server"   
#define BROADCAST {0xff,0xff,0xff,0xff};

struct	Bootp_format		// 300 bytes in length
{
 unsigned char Op;       // specifies whether the message is a request
 unsigned char Htype;     // specifies the network hardware type
 unsigned char Hlen;     // specifies the length of a hardware address   
 unsigned char Resv;    
 long Transaction_Id; // contains an integer that diskless machines use to match responses whith request
 unsigned char Second[2];
 unsigned char Unused[2];
 unsigned char Client_Address[4];  
 unsigned char Your_IP_Address[4];
 long Server_IP_Address;
 unsigned char Gateway_IP_Address[4];
 unsigned char C_H_A[16];
 unsigned char Server_Host_Name[64];
 unsigned char FName[128];
 unsigned char Vendor_Specific_Area[64];// used to pass Net mask 
// int BUFID;       // this is not a field of bootp frame, only for identifing
};
               
               
// end lbr   
               
