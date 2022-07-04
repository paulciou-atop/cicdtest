//Filename: cfgapi.c
//Update  :
//  2002/11/01: char GW_Model[16]="GW21I"; is changed to extern char GW_Model[16];
//  2003/01/24: in Read_com_cfg(), add request_frame.Server_Host_Name[63]= 2 ;
//  2003/01/27: Model name is changed from "GW21I" to "GW21R"
//  2003/06/06: update GW_Model
//  2003/06/05: response count of serial ports in Read_com_cfg()
//              update GetTotalPort(), return value is change to count of serial ports
//  2003/10/28: add Get_Cfg_Data(), Read_sys_cfg(). remove Read_com_cfg(), Read_ip_cfg()
//  2003/10/28: remark Set_RSMODE() in Configuate()
//  2004/04/16: use FOREPROM to remark t2f(), f2t()
//  2004/04/26: update Initial(), to add new model name "GW51C-MAXI", "GW51W-MAXI" from GW21CMAXI, GW21WMAXI
//  2004/05/05: in Read_sys_cfg(), add request_frame.Server_Host_Name[45]= country_code ;
//  2004/08/17: update DEF_SETPORTTYPE to DEF_SET485WIRES in Initial()
//  2004/09/16: fix a bug, must be FAR_strlen() not strlen(), in ZipPwd()
//  2004/11/22: update Initial() to support multiple model
//  2005/07/08: update AccountManager() to fix login bug
//  2005/10/22: add DEF_GW26A, it is almost tha same as DEF_GW21L
//  2005/11/09: update daprecv_post() to to prevent username and password overflow
//  2005/11/15: add DEF_GW21S001
//  2006/06/04: add DEF_SE2002
//  2006/08/17: add DEF_GW21S001_10M, update Initial() for setting 10M speed to GW21S001
//  2007/04/02: add report_config2(), for wireless LAN re-connect to report I am here
//  2007/05/18: add to service the delayed ack msg from invite command
//              delay time = ((ip>>24)&0x1F)<<1
//  2007/05/19: moved this define AP_Serial_No[] from dapdos.asm to here, and initial AP_Serial_No[0]=0
//              Note: AP_SERIAL_NO_LEN is also defined in dapdos.asm
//  2007/07/10: update DelayAckInviteAdd() to fix a bug if delay_ms==0
//  2007/07/25: update report_config2(), to report more times
//  2007/09/02: add to get model from EEPROM for DEF_GW21CMAXI + DEF_GW21SMAXI + DEF_GW21WMAXI + DEF_PHYSIO, in Initial()
//  2007/09/21: add SECTION_LOCK to replace _cli()
//  2007/10/23: update DelayAckInviteService()
//  2008/01/29: add SE5001 to support external WDT, append the model name with "-B"
//  2008/09/06: if shadow mode, append '-R' to the model name
//  2009/09/02: fix a bug that the monitor tool invits nothing if the DHCP cannot get a IP (IP=0.0.0.0)
//  2009/10/19: add download type to .Server_Host_Name[ 44], 
//              (0:unknown, 1:80186, 2:mega, 3:IDT, 4:PPC)
//  2010/03/31: add _IS_BACKUP for EEPROM backup
//  2011/03/14: for DEF_SE2002, add to call GetFeuIO() to decide 1/2 com port
//  2011/04/07: for shadow mode do not append 'R' to model name,
//              but by kernel version, odd:shadow mode, even:normal mode
//  2011/05/06: add _IS_BACKUP3 
//  2011/06/14: add _IS_BACKUP4 
//  2011/09/26: add report_gateway() for making gateway refresh its ARP table
//              send the 1-st after 3 sec, the 2-nd/3-rd after 10 sec, others every 5 min
//  2016/01/08: try to add command "@eeprom backup" to backup EEPROM to DOOO segment.
//              but it is failed. System will reboot.                
//  2016/01/18: add command "@ee 96" to show EEPROM word 96, 
//              add command "@mm 80000" to show Flash 0x8000:0000
//  2016/02/18: backup command "@eeprom backup" is workable
//  2016/02/18: mave backup commands from Set_Host_Name() to report_debug0()
//  2016/02/18: for _IS_BACKUP4, provide backup commands: "@backup eeprom", "@backup erase", "@backup check" 
//  2016/02/19: add commands: "@shadow on", "@shadow off", "@shadow check" 
//  2016/02/19: fix the bug that failed to write MAC address to flash when in shadow mode
//  2018/04/19: In LoadMacAddress(), if MAC is the same, then not overwrite the last 3 bytes of MAC address. 
//              (overwrite 2 words actually)

#include	"model.h"
#include	"tcpdef.h"
#include	"bootp.h"
#include	"mymodel.h"

                                    //---add 2006/09/21
#define SECTION_LOCK       _asm pushf _asm cli
#define SECTION_RELEASE    _asm popf

#if DEF_GW21SMAXI + DEF_GW21CMAXI + DEF_GW21WMAXI + DEF_GW21S001 + DEF_PHYSIO + DEF_SE2002 >= 1
#include        "tcpapi.h"
#include        "dapapi.h"
#include        "debug2.h"
#endif

//#define DEF_TOTALPORT          2

#if DEF_GETPORTTYPE >= 1
#define TYPE232  0x00
#define TYPE485  0x01
#endif

#if DEF_GET485WIRES >= 1
#define WIRENULL  0x00
#define WIRE2     0x00
#define WIRE4     0x01
#endif

//#define FOREPROM
#if DEF_GW21R + DEF_GW21SMAXI + DEF_GW21CMAXI +DEF_GW21S001 + DEF_PHYSIO + DEF_SE2002 >= 1
#define UTABLEKEY 1001

#elif DEF_GW21L + DEF_GW21W + DEF_GW21WMAXI + DEF_GW26A >= 1
#define UTABLEKEY  999
#endif

/*************************************************************************/
/*************************************************************************/
#if DEF_GW21S256 >= 1
#define USERMAX 1

#else
#define USERMAX 4
#endif

/*---remark at 2003/03/03
#if 0
#define MODELNAME1  "GW21I-VM"
#define MODELNAME2  "GW21I"
#else   //-----update at 2003/01/27
#define MODELNAME1  "GW21R-VM"
#define MODELNAME2  "GW21R"
#endif

#define DLLNAME1   "GW21LE"
#define DLLNAME2   "GW21LE" //2002.07.26 IVM->LE, hjh
---*/

#if 1                       //---add 2007/05/18
typedef struct 
{
//  unsigned short  delay_ms ;
    short           delay_ms ;      //---update 2007/10/23
    unsigned long   src_ip ;
} DELAY_ACK_INVITE ;

#define DELAY_ACK_INVITE_MAX        1
DELAY_ACK_INVITE AckInviteStru[ DELAY_ACK_INVITE_MAX] ;
#endif

typedef struct
{
   unsigned char name[6];
   unsigned short pwdzip;
} USERITEM;

typedef struct
{
   USERITEM user[USERMAX];
} USERTABLE;

USERTABLE utable;

extern short RETVALUE_IDPWD;
extern char ID_PSW_PAIR[60];
//extern char AP_Serial_No[30]; //update 2003/03/03                         
//extern char GW_Model[16];     //update 2003/03/03
extern char HostName[] ;        //add 2003/10/28
extern long    local_ipaddr;
extern unsigned long MY_IPADDR;
extern  long   default_router;      //---add 2011/09/23

unsigned eeprom_is_blank=1;

#if 0
extern char AP_Serial_No[30];

extern char GW_Model[16];        //update at 2002/11/01, defined in dapmain.asm
//char GW_Model[16]="GW21I";
char GW_Dll[16]="GW21LE";

char APMSGBUF[30]="tEST";

#else                           //---update 2004/03/01
#define GW_MODEL_LEN    16
#define GW_DLL_LEN      8
char GW_Model[ GW_MODEL_LEN+1];                          
char GW_Dll[ GW_DLL_LEN+1];

 #if 1                          //---add 2007/05/19, moved this define from dapdos.asm
#define AP_SERIAL_NO_LEN    120
char    AP_Serial_No[ AP_SERIAL_NO_LEN] ;
 #endif

#endif


#ifdef MAC_ENCODE
#include	"macaddr.h"
static unsigned char mac[6]={0x00,0x00,0x00,0x00,0x00,0x00};
extern unsigned char ProgOneByteK( unsigned char id, unsigned char data, unsigned long addr) ;
#endif
/*************************************************************************/
/*************************************************************************/

unsigned char first_report=0x00;
unsigned char try_report=0x00;          //add 2003/10/28

int    handle;

long   bootp_address=0xffffffff,host_addr=0xffffffff,subnetmask=0xffffff00,netaddr;
long   debugflag=0,debugaddr;
short  localport=0xda92,bootp_port=0xda92;
struct Bootp_format request_frame,ack_frame,bootpframe;
unsigned char ghwadd[6],gipadd[4],ggwadd[4],gmask[4];

char far * receive_buf=(char far *)&bootpframe;
char far * request_buf=(char far *)&request_frame;
char far * ack_buf=(char far *)&ack_frame;

char far *arg=0;
int resetflag=0;

unsigned char   NEEDRUNAP=0x01 ;   //lbr

#if DEF_SE2002 >= 1
unsigned char TotalCom;         //---add 2011/03/14
#endif

#ifdef _IS_BACKUP4      //---add 2011/06/14
extern int     bk_eeprom_flag;
#endif

/////////////////////////////////////////////////////////////////////////////
void _loadds far Set_IP_Address();
void _loadds far Set_Default_Gateway();
void _loadds far Set_Netmask();
void _loadds far Set_Host_Name() ;
void _loadds far Configuate();
void _loadds far ConfigRoutine();
void _loadds far LoadMacAddress();
void _loadds far Set_HW_Address();
//void _loadds far Set_RSMODE();
//void _loadds far Read_com_cfg();
void _loadds far Set_RS232();
void _loadds far Read_sys_cfg();

unsigned char _loadds far getmycomtype(unsigned char port) ;
unsigned char _loadds far getmy485wire(unsigned char port) ;

int bdecode(unsigned char far *c_ary,unsigned char far *originalbyte) ;
int macdecode(unsigned char far *macaddr,unsigned char far *maccode) ;
#ifdef FOREPROM
void t2f(unsigned char *tptr,unsigned char *fptr) ;
int f2t(unsigned char *fptr,unsigned char *tptr) ;
#endif
/////////////////////////////////////////////////////////////////////////////

unsigned short _loadds far ZipPwd(unsigned char far *pwd)
{
        unsigned short val=0;
        int k;
        int i;
        
//      return (*pwd);

//      k=strlen(pwd);
        k=FAR_strlen(pwd);  //---fix a bug, must be FAR_strlen() not strlen(), 2004/09/16

        for(i=0;i<k;i++)
        {
           val+= (*(pwd+i))*(i+1);
        }
        return val;
}


void _loadds far update_utable()
{
    int i;
    unsigned short data;
    unsigned char *ptr;

        ptr=(unsigned char *)&utable.user[0].name[0];
        
        for(i=0;i<(sizeof(utable)/2);i++)
        {   
           data=(*(ptr+(i*2)))*256+(*(ptr+(i*2)+1));            
           WriteKerEEPROM(25+i,data);              
        }   
}

void _loadds far Set_utable_default()
{
     WriteKerEEPROM(24,0) ;

#if DEF_GW21W + DEF_GW21WMAXI >= 1
     WriteKerEEPROM(96,0) ;
#endif
}


void _loadds far load_utable()
{
    int i;
    unsigned short data;
    unsigned char *ptr;

    data=ReadKerEEPROM(24) ;
    ptr=(unsigned char *)&utable.user[0].name[0];

    if(data==UTABLEKEY)
    {
	   for(i=0;i<(sizeof(utable)/2);i++)
	   {
	      data=ReadKerEEPROM(25+i);
	      *(ptr+(i*2))=data/256;
	      *(ptr+(i*2)+1)=data%256;
	   }
	   eeprom_is_blank=0;
    }else
    {
       strcpy(utable.user[0].name,"admin");
       utable.user[0].pwdzip=ZipPwd("");
       for(i=1;i<USERMAX;i++)
       {
	      utable.user[i].name[0]=0x00;
       }
       update_utable();
       WriteKerEEPROM(24,UTABLEKEY) ;
       eeprom_is_blank=0x99;
    }
}

extern char far *FAR_strcpy( char far *dst, char far *src) ;
extern int FAR_strlen( char far *src) ;
extern int FAR_strcmp( char far *s1, char far *s2) ;
extern int FAR_strncmp( char far *s1, char far *s2, int maxcnt) ;  //---add 2005/07/08

void _loadds far AccountManager()
{  
    unsigned short retval=0;
    int i,j;
    unsigned char name[30];
    unsigned char pwd[30];
    
    if (ID_PSW_PAIR[0]==0x00) 
    {
        RETVALUE_IDPWD=retval;
        return;
    }
  
    i=0;
    while(i<60)
    {
        name[i]=ID_PSW_PAIR[i];
        //APMSGBUF[i]=name[i]; 
//      if(name[i]==0x20)
        if(name[i]==0x20 || name[i]==0x00 || i==6)      //---update 2005/07/08
        {
            name[i]=0x00;
            break;
        }
        i++;
    }

    if( FAR_strlen( name)==0) 
    {
        RETVALUE_IDPWD=retval;
        return;
    }     
    j=0;
    i++;
    
    while(i<60)
    {
        pwd[j]=ID_PSW_PAIR[i]; 
        //APMSGBUF[j]=pwd[j];
//      if(pwd[j]==0x00)
        if(pwd[j]==0x00 || j==20)           //---update 2005/07/08
        {
            pwd[j] = 0 ;
            break;
        }
        i++;        
        j++;
    }
    
    switch(RETVALUE_IDPWD)
    {
    case 0: //add account or change password
        j = -1 ;
        for(i=0;i<USERMAX;i++)
        {
//          if( FAR_strcmp( utable.user[i].name, name)==0)
            if( FAR_strncmp( utable.user[i].name, name, 6)==0)      //---update 2005/07/08
            {
                utable.user[i].pwdzip=ZipPwd(pwd);
                retval=1;
                update_utable();        
                break;                          
            }
            else
//          if(utable.user[i].name[0]=0x00)
            if(utable.user[i].name[0]==0x00)        //---fix a bug, 2005/07/08
            {
#if 0
                FAR_strcpy( utable.user[i].name, name);
                utable.user[i].pwdzip=ZipPwd(pwd);
                retval=2;
                update_utable();
                break;
#else //---update 2005/07/08
                if ( j<0) j = i ;
#endif
            }       
        }
        
        //---add 2005/07/08
        if ( i>=USERMAX && j>=0)
        {
                FAR_strcpy( utable.user[j].name, name);
                utable.user[j].pwdzip=ZipPwd(pwd);
                retval=2;
                update_utable();
        }

        break;
        
    case 1: //del account
        for(i=0;i<USERMAX;i++)
        {
//          if( FAR_strcmp( utable.user[i].name,name)==0)
            if( FAR_strncmp( utable.user[i].name, name, 6)==0)      //---update 2005/07/08
            {
                utable.user[i].name[0]=0x00;
                retval=1;
                update_utable();
                break;
            }       
        }   
        break;
        
    case 2: //verify account
        for(i=0;i<USERMAX;i++)
        {
//          if( FAR_strcmp( utable.user[i].name, name)==0)
            if( FAR_strncmp( utable.user[i].name, name, 6)==0)      //---update 2005/07/08
            {
                if(ZipPwd(pwd)==utable.user[i].pwdzip)
                {
                    retval=1;
                }
                break;
            }
        }               
        break;
   }
   RETVALUE_IDPWD=retval;
}
/////////////////////////////////////////////////////////////////////////////

void _loadds far Set_HW_Address()
{
    unsigned short wordtemp;
    unsigned char flag;
    unsigned char *ptr;
    
    ptr=&(bootpframe.Vendor_Specific_Area[60]);
    
    if((ptr[0]=='A')&&(ptr[1]=='T')&&(ptr[2]=='O')&&(ptr[3]=='P'))
    {
        wordtemp=bootpframe.Vendor_Specific_Area[7]*256+bootpframe.Vendor_Specific_Area[6];
        WriteHwEEPROM(0x01,wordtemp);
        wordtemp=bootpframe.Vendor_Specific_Area[9]*256+bootpframe.Vendor_Specific_Area[8];
        WriteHwEEPROM(0x02,wordtemp);
    }
}

void _loadds far Set_IP_Address()
{
    unsigned short wordtemp;
    wordtemp=bootpframe.Your_IP_Address[1]*256+bootpframe.Your_IP_Address[0];
    WriteMacEEPROM(0x10,wordtemp);
    wordtemp=bootpframe.Your_IP_Address[3]*256+bootpframe.Your_IP_Address[2];
    WriteMacEEPROM(0x11,wordtemp);
}

void _loadds far Set_Default_Gateway()
{
    unsigned short wordtemp;

    wordtemp=bootpframe.Gateway_IP_Address[1]*256+bootpframe.Gateway_IP_Address[0];
    WriteMacEEPROM(0x12,wordtemp);
    wordtemp=bootpframe.Gateway_IP_Address[3]*256+bootpframe.Gateway_IP_Address[2];
    WriteMacEEPROM(0x13,wordtemp);
}

void _loadds far Set_Netmask()
{
    unsigned short wordtemp;

    wordtemp=bootpframe.Vendor_Specific_Area[1]*256+bootpframe.Vendor_Specific_Area[0];
    WriteMacEEPROM(0x14,wordtemp);
    wordtemp=bootpframe.Vendor_Specific_Area[3]*256+bootpframe.Vendor_Specific_Area[2];
    WriteMacEEPROM(0x15,wordtemp);

}

#if 0  //---remark 2005/09/09
void _loadds far Set_RS232()
{
    unsigned short wordtemp;
    unsigned char flag;
    
    flag=bootpframe.Vendor_Specific_Area[19];
    if(flag==0x00)
    {
#if DEF_GW21R + DEF_GW21L + DEF_GW21W >= 1
        wordtemp=bootpframe.Server_Host_Name[24]*256+bootpframe.Server_Host_Name[25];
        WriteEEPROM(18,wordtemp);
#endif
        wordtemp=bootpframe.Vendor_Specific_Area[11]*256+bootpframe.Vendor_Specific_Area[10];
        WriteEEPROM(20,wordtemp);
        wordtemp=bootpframe.Vendor_Specific_Area[13]*256+bootpframe.Vendor_Specific_Area[12];
        WriteEEPROM(21,wordtemp);
        wordtemp=bootpframe.Vendor_Specific_Area[15]*256+bootpframe.Vendor_Specific_Area[14];
        WriteEEPROM(22,wordtemp);
        wordtemp=bootpframe.Vendor_Specific_Area[17]*256+bootpframe.Vendor_Specific_Area[16];
        WriteEEPROM(23,wordtemp);
    }
}
#endif

////////////////////////////////////////////////////////////////////////////
void _loadds far Configuate()
{
   Set_IP_Address();
   Set_Default_Gateway();
   Set_Netmask();
   Set_Host_Name() ;        //---add 2003/10/28

   Set_HW_Address();
// Set_RS232();
}

unsigned char _loadds far getmycomtype(unsigned char port)
{
#if DEF_GETPORTTYPE >= 1
    unsigned short mytype;
    
    mytype = Getporttype();
    if(mytype&(0x01<<(port)))
        return TYPE232;
    else
        return TYPE485;
#else
   return( 0) ;     //---not support
#endif
}

//////////////////////////////////////////////////////////////////////////////////////
int _loadds far setmy485wire(unsigned char port,unsigned char mywire)
//---ret=0:not support this function, 1:ok
{
#if DEF_SET485WIRES >= 1
    unsigned short mytype;
    unsigned char val;
    
    val=Get485type();
    
    if(mywire==WIRE4)
    {
        val=(val&(~(0x01<<(port))));
    }else
    {
        val=(val|(0x01<<(port)));
    }
    
    Set485type(val);
   return( 1) ;     //---ok
#else
   return( 0) ;     //---not support
#endif
}

unsigned char _loadds far getmy485wire(unsigned char port)
//---ret=0:not support this function or type unknown or 2 wires, 1:4 wires
{
#if DEF_GET485WIRES >= 1
    unsigned char mytype;
    unsigned char val;

    mytype= getmycomtype(port);
    if(mytype==TYPE485)
    {
        val=Get485type();
        if(val&(0x01<<port))
        {
            return WIRE2;
        }
        else
        {
            return WIRE4;
        }
    }
    return WIRENULL;
#else
   return 0;
#endif
}

void _loadds far daprecv_post(cb, arg)
CMDBLK far *cb;
void far *arg;
{
	long a,b;
	int delay=10000;
	int dd,d2;
	int i;
	unsigned short pwdret;
    int     shadow_flag;

#ifdef MAC_ENCODE
    unsigned short ret;
    unsigned char id,data;
    unsigned char far * addr;
    unsigned char *mptr;
    unsigned char is_virg;
#endif

    WD_FUN();  // Clear watch dog timer to prevent the system resetting

    switch( bootpframe.Server_Host_Name[26])
    {
    case 0x00:              //no pasword
        RETVALUE_IDPWD=0;
        pwdret=RETVALUE_IDPWD;
        break;

    case 0x01:              //only password, user is "admin"
//      _cli();
        SECTION_LOCK        //---update 2007/09/21
        
        RETVALUE_IDPWD=2;
        strcpy(ID_PSW_PAIR,"admin ");
//      strcpy(&ID_PSW_PAIR[6],&bootpframe.Server_Host_Name[27]);
        strncpy( &ID_PSW_PAIR[6], &bootpframe.Server_Host_Name[27], 32-6);  //---update 2005/11/09, to prevent overflow
        ID_PSW_PAIR[ 32] = '\0' ;

        AccountManager();
        pwdret=RETVALUE_IDPWD;

//      _sti();
        SECTION_RELEASE     //---update 2007/09/21

        break;
        
    case 0x02:              //username & password
//      _cli();
        SECTION_LOCK        //---update 2007/09/21

        RETVALUE_IDPWD=2;
//      strcpy(&ID_PSW_PAIR[0],&bootpframe.Server_Host_Name[27]);
        strncpy( &ID_PSW_PAIR[0], &bootpframe.Server_Host_Name[27], 32);  //---update 2005/11/09, to prevent overflow
        ID_PSW_PAIR[ 32] = '\0' ;

        AccountManager();
        pwdret=RETVALUE_IDPWD;

//      _sti();
        SECTION_RELEASE     //---update 2007/09/21

        break;
    }

    if((cb->recvlen==sizeof(struct Bootp_format))&&(bootpframe.Transaction_Id== MY_TRANSACTION_ID))  
    {
#if 0
        unsigned short ret;
        unsigned char id,data;
        unsigned char far * addr;
        unsigned char *mptr; 
        unsigned char is_virg;
#endif

//      if(fist_report==1)
        if(first_report>0)          //update 2003/07/25
        {
            switch (bootpframe.Op)
            {
            case INVITE:
                host_addr=cb->from_ipaddr;
#if 0
//              Read_com_cfg();
//              Read_ip_cfg() ;           //add 2003/10/22, for DHCP to update IP
                Read_sys_cfg();           //add 2003/10/28

                // is it the same subnet
                if((subnetmask&host_addr)!=netaddr) udp_send(handle,bootp_port,bootp_address,request_buf,sizeof(struct Bootp_format));
                udp_send(handle,bootp_port,host_addr,request_buf,sizeof(struct Bootp_format));

#else  //---update 2007/05/18
                Read_sys_cfg();           //add 2003/10/28
 #if 0
                if ( DelayAckInviteAdd( host_addr, (gipadd[3]&0x1F)<<1) == 0)       //---use ip as delay count, interval 2ms
                report_config3( host_addr, request_buf) ;
 #else //---update 2007/10/23
                DelayAckInviteAdd( host_addr, ((gipadd[3]&0x1F)+1)<<1) ;    //---use ip as delay count, interval 2ms
 #endif
#endif
                break;

            case DAPRESET:
                if( pwdret==0 ) break;
                host_addr=cb->from_ipaddr;
                if(!memcmp(bootpframe.C_H_A,ghwadd,6))
                {
#if 0
                    if((subnetmask&host_addr)!=netaddr) udp_send(handle,bootp_port,bootp_address,ack_buf,sizeof(struct Bootp_format));
                    udp_send(handle,bootp_port,host_addr,ack_buf,sizeof(struct Bootp_format));
#else  //---update 2007/05/18
                    report_config3( host_addr, ack_buf) ;
#endif
                    resetflag=1;
                }
                break;

            case DAPBEEP:
                host_addr=cb->from_ipaddr;
                if(!memcmp(bootpframe.C_H_A,ghwadd,6))
                {
                    BuzzerOn();
                }
                break;

            case REPLY:
                if( pwdret==0 ) break;
                host_addr=cb->from_ipaddr;
                if(!memcmp(bootpframe.C_H_A,ghwadd,6))
                {
 #if 1 //---add 2016/01/18, show EEPROM or Flash
                    if ( report_debug()) {
                        report_config3( host_addr, ack_buf) ;
                        break;
                    }
 #endif
                 
#if 0
                    if((subnetmask&host_addr)!=netaddr) udp_send(handle,bootp_port,bootp_address,ack_buf,sizeof(struct Bootp_format));
                    udp_send(handle,bootp_port,host_addr,ack_buf,sizeof(struct Bootp_format));
#else  //---update 2007/05/18
                    report_config3( host_addr, ack_buf) ;
#endif
                    Configuate();
                    resetflag=1;
                }
                break;

#ifdef MAC_ENCODE
            case MAC_ENCODE:        
                mptr=&bootpframe.Server_Host_Name[0];
                addr=(unsigned char far *)0xf000fef0;
                is_virg=1;
                
#ifdef FOREPROM
                {
                unsigned short wordtemp=0x59+0x02+0x12+0x00+0x60+0xe9;
                unsigned short data,sflag;
                int i;
                
                if(macdecode(&mac[1],mptr)!=1)
                {
                    if(macdecode(&mac[1],gmaccode)!=1) while(1);
                }
                t2f(&mac[1],mptr);
                wordtemp=mptr[1]*256+mptr[0];
                WriteKerEEPROM(25+(sizeof(utable)/2),wordtemp);
                wordtemp=mptr[3]*256+mptr[2];
                WriteKerEEPROM(25+(sizeof(utable)/2)+1,wordtemp);
                resetflag=1;
                }

#else
                for(i=0;i<256;i++)
                {
                    if((*(addr+i))!=0xff)
                    {
                        is_virg=0;
                        break;
                    }
                }
                if(is_virg==1)
                {
 #if 1 //---add 2016/02/19, for supporting shadow mode
                    if ( shadow_flag = ShadowMode_chk()) {
                        ShadowMode_off();
                    }
 #endif
                    for(i=0;i<256;i++)
                    {
                        data=*mptr;
                        ret=ProgOneByteK( 0x34, data,(unsigned long) addr ) ;
                        addr++;
                        mptr++;
                    }

 #if 1 //---add 2016/02/19, for supporting shadow mode
                    if ( shadow_flag) {
                        ShadowMode_on();
                    }
 #endif

                    resetflag=1;
                }
#endif
                break;
#endif
            } //switch
        }
    }
    
    debugaddr=subnetmask&host_addr;
    debugflag=1;
    if(SUCCESS==udp_receive(handle,receive_buf,sizeof(struct Bootp_format),daprecv_post,0))
    {
    }else
    {
    }
    if(resetflag==1)SYSTEM_RESET();
}

void _loadds far report_config(unsigned long i)// 2001/2/8
{
//  first_report=1;
//  first_report+=1;         //update 2003/07/25
    try_report++;            //update 2003/10/28

#if 0
#if 0
    memcpy(request_frame.Server_Host_Name,GW_Model,1+strlen(GW_Model));
    memcpy(&request_frame.Server_Host_Name[16],GW_Dll,1+strlen(GW_Dll));
#else  //---update 2003/03/03
    strncpy( request_frame.Server_Host_Name, GW_Model, GW_MODEL_LEN);
    strncpy(&request_frame.Server_Host_Name[ GW_MODEL_LEN], GW_Dll, GW_DLL_LEN);
#endif
    GetAPSerialNo((char far *)&request_frame.FName[2]);

#else       //---update 2003/10/29
    Read_sys_cfg() ;
#endif

    if ( local_ipaddr==0)      //dhcp mode and not get ip in 3*2 sec, add 2003/10/28
    if ( try_report < 3)
    {
//      set_timer(18*2,(unsigned long) 0, report_config );
        set_timer(18*3,(unsigned long) 0, report_config );      //---update 2007/06/05
        return ;
    }

    //memcpy(&request_frame.FName[2],APMSGBUF,30);
    udp_send(handle,bootp_port,bootp_address,request_buf,sizeof(struct Bootp_format));//send to local    

    first_report++ ;
    if ( first_report<3)            //add 2003/07/25
    {
        set_timer(18*2,(unsigned long) 0, report_config );
    }
}

#if 1 //---add 2011/09/26, make gateway refresh its ARP table
void _loadds far report_gateway( unsigned long seq)
//---send the 1-st after 3 sec
//---send the 2-nd/3-rd after 10 sec
//---then send every 5 min
{
    if ( local_ipaddr && default_router)
    arp_output( default_router);

    if ((unsigned short)seq < 3) seq++;
    set_timer( ((unsigned short)seq<3)? 18*10: 18*60*5,(unsigned long)seq, report_gateway);
}
#endif

#if 0
void _loadds far Read_com_cfg()
{
   unsigned short wordtemp;

#if 0
    memcpy(&request_frame.Server_Host_Name[0],GW_Model,1+strlen(GW_Model));
    memcpy(&request_frame.Server_Host_Name[16],GW_Dll,1+strlen(GW_Dll));
#else  //---update 2003/03/03
    strncpy( request_frame.Server_Host_Name, GW_Model, GW_MODEL_LEN);
    strncpy(&request_frame.Server_Host_Name[ GW_MODEL_LEN], GW_Dll, GW_DLL_LEN);
#endif
    request_frame.Server_Host_Name[63] = GetTotalPort() ;   //count of serial ports, add 2003/06/06

#if DEF_GET485WIRES >= 1   
    request_frame.Server_Host_Name[24]=(getmycomtype(0)<<1)|getmy485wire(0);//for port 1
    request_frame.Server_Host_Name[25]=(getmycomtype(1)<<1)|getmy485wire(1);//for port 2
#endif

    wordtemp=GetVersion();
    request_frame.FName[0]=wordtemp/256;
    request_frame.FName[1]=wordtemp%256;

    GetAPSerialNo((unsigned long)&request_frame.FName[2]);

    if(eeprom_is_blank==0x99) return;  //do not report config, for test eeprom;

    wordtemp=ReadEEPROM(20);
    ack_frame.Vendor_Specific_Area[11]=wordtemp/256;
    request_frame.Vendor_Specific_Area[11]=wordtemp/256;
    ack_frame.Vendor_Specific_Area[10]=wordtemp%256;
    request_frame.Vendor_Specific_Area[10]=wordtemp%256;

    wordtemp=ReadEEPROM(21);
    ack_frame.Vendor_Specific_Area[13]=wordtemp/256;
    request_frame.Vendor_Specific_Area[13]=wordtemp/256;
    ack_frame.Vendor_Specific_Area[12]=wordtemp%256;
    request_frame.Vendor_Specific_Area[12]=wordtemp%256;

    wordtemp=ReadEEPROM(22);
    ack_frame.Vendor_Specific_Area[15]=wordtemp/256;
    request_frame.Vendor_Specific_Area[15]=wordtemp/256;
    ack_frame.Vendor_Specific_Area[14]=wordtemp%256;
    request_frame.Vendor_Specific_Area[14]=wordtemp%256;

    wordtemp=ReadEEPROM(23);
    ack_frame.Vendor_Specific_Area[17]=wordtemp/256;
    request_frame.Vendor_Specific_Area[17]=wordtemp/256;
    ack_frame.Vendor_Specific_Area[16]=wordtemp%256;
    request_frame.Vendor_Specific_Area[16]=wordtemp%256;
}
#endif

void _loadds far Initial()
{
    int     i, j, ret, flag ;
    long a,b;
    char buf[32], *ptr;
    unsigned short wordtemp;
    unsigned char port0type,port1type;

    AP_Serial_No[0] = '\0' ;    //---add 2007/05/19
    
    DelayAckInviteIni() ;       //---add 2007/05/18

    memset( GW_Model, 0, GW_MODEL_LEN);     //---add 2008/09/06

#if DEF_GW21R >= 1      //333333333333333333333333333333333333333333
    switch(GetFeuIO()&0x01)
    {
	case 0://gw21i-vm
	    strcpy(GW_Model,MODELNAME1);
	    strcpy(GW_Dll,DLLNAME1);
	    break;
	case 1://gw21i
	    strcpy(GW_Model,MODELNAME2);
	    strcpy(GW_Dll,DLLNAME2);
	    break;
    }

#elif DEF_GW21CMAXI + DEF_GW21SMAXI + DEF_GW21WMAXI + DEF_PHYSIO >= 1    //333333333333333, add 2004/04/26
    ret = GetFeuIO() ;

#if 0                                       //---update 2006/06/07
    #if DEF_PHYSIO >= 1
        Setnewporttype( 0, ret & 0x03) ; // Com port 1 
        Setnewporttype( 1, ret & 0x03) ; // Com port 2
    #endif
#endif

 #if 0
    if (ret>=MODEL_POOL_MAX) ret = 0 ;      //overflow then set default
    strcpy( GW_Model, ModelPool[ ret]);
    strcpy( GW_Dll, DLLNAME1);

 #else  //---update 2004/11/22
    i = ret & 0xF ;
    j = (ret & 0x70)>>4 ;

  #if 0
    if ((ret & 0x80) || i>=MODEL_POOL_MAX)       //overflow then set default
        strcpy( GW_Model, "Unknown");
  #else						                     //---update 2007/09/02
    if (ret & 0x80)                              //overflow then set default
        strcpy( GW_Model, "Unknown");
    else
    if (ret & 0x40) {                            //---read model from EEPROM
        for ( i=0; i<GW_MODEL_LEN/2; i++) {
            *(ushort *)&GW_Model[ i*2] = ReadMacEEPROM( 41+i) ;
        }
        GW_Model[ GW_MODEL_LEN] = '\0' ; 
    }
    else
    if (i>=MODEL_POOL_MAX)                       //overflow then set default
        strcpy( GW_Model, "Unknown");
  #endif
  
    else
    {
        strcpy( GW_Model, ModelPool[ i]);
        if (j)
        {
            strcat( GW_Model, "-DC0") ;
            i = strlen( GW_Model) ;
            GW_Model[ i-1] += j ;
        }
    }
    strcpy( GW_Dll, DLLNAME1);
 #endif
    
#elif DEF_GW21S001 + DEF_SE2002 >= 1    //---add 2006/03/16
 #if 1
    ret = GetFeuIO() ;                      //---bit0-1: com port type, 0:RS232, 1:RS485, 2:RS422
                                            //---bit2:VCOM disable
    Setnewporttype( 0, (uchar)(ret & 0x03)) ;

  #if DEF_SE2002 >= 1                       //---add 2006/06/04, set COM2 as RS232/485/422
    Setnewporttype( 1, (unsigned char)((ret>>8) & 0x03)) ;
//  Getnewporttype( 1);		                //---add 2006/06/04, For debug
  #endif
    
  #if DEF_SE2002 >= 1
    TotalCom = (ret & 0x0800)? 1: 2;       //---add 2011/03/14, decide 1/2 com port
  #endif

  #if DEF_GW21S001_10M >= 1    //---add 2006/08/17, for setting 10M speed to GW21S001
    strncpy( GW_Model, MODELNAME8, GW_MODEL_LEN);
    GW_Model[ GW_MODEL_LEN] = '\0' ;
  #else
//  for ( i=0; i<8; i++) {
    for ( i=0; i<GW_MODEL_LEN/2; i++) {     //---update 2007/06/05
        *(ushort *)&GW_Model[ i*2] = ReadMacEEPROM( 41+i) ;
    }
    GW_Model[ GW_MODEL_LEN] = '\0' ;        //---add 2006/06/30

   #if DEF_GW21S001_WDT >= 1                //---add 2008/01/29, add SE5001 to support external WDT
    if (( i= strlen( GW_Model)) < GW_MODEL_LEN-2) {
        GW_Model[ i++] = '-' ;
        GW_Model[ i] = 'B' ;
    }
   #endif
   
   #if DEF_LAN_10M >= 1                     //---add 2019/02/13, force LAN 10MHz
    if (( i= strlen( GW_Model)) < GW_MODEL_LEN-2) {
        GW_Model[ i++] = '-' ;
        GW_Model[ i++] = '1' ;
        GW_Model[ i] = '0' ;
    }
   #endif

  #endif
    strcpy( GW_Dll, DLLNAME1);

 #else                      //---update 2006/03/21
 
 #if 0
    flag = ret ;
    _asm mov dx,0ff80h _asm mov al,'@' _asm out dx,al              //---add 2006/03/17
    _asm mov dx,0ff80h _asm mov al,flag _asm out dx,al              //---add 2006/03/17
 #endif

    i = ret & 0xF ;
//  j = (ret & 0x70)>>4 ;
    if ( i>=MODEL_POOL_MAX)                     //overflow then set default
        strcpy( GW_Model, "Unknown");
    else
    {
        strcpy( GW_Model, ModelPool[ i]);
        flag = 0 ;
 #if 0
        if (j & 0x01)                           //---show 232/485 or not
        {
            ret = Getporttype() ;
            if ( ret & 0x01)                    //---RS232
                 ptr = "-S2" ;
            else ptr = "-S5" ;
            strcat( GW_Model, ptr) ;
            flag = 1 ;
        }
        if ((j & 0x02)==0)                      //---show without VCOM
        {
            strcat( GW_Model, (flag)? "N": "-N") ;
        }
 #else
        if (ret & 0x10)                         //---is RS232/RS485 or not
        {
            ptr = "-S5" ;
            strcat( GW_Model, ptr) ;
            flag = 1 ;
        }
        if ((ret & 0x20)==0)                    //---is no VCOM or not
        {
            strcat( GW_Model, (flag)? "N": "-N") ;
        }
 #endif
    }
    strcpy( GW_Dll, DLLNAME1);
 #endif

#elif DEF_GW21L >= 1    //333333333333333333333333333333333333333333
    switch(GetTotalPort())
    {
    case 2:     //GW21L, 1 -> 2, update 2003/06/06
#if 0	
	     strcpy(GW_Model,MODELNAME1);
	     strcpy(GW_Dll,DLLNAME1);
#else   //---update 2003/03/03
	     strncpy( GW_Model, MODELNAME1, GW_MODEL_LEN);
	     GW_Model[ GW_MODEL_LEN] = '\0' ;
	     strncpy( GW_Dll, DLLNAME1, GW_DLL_LEN);
	     GW_Dll[ GW_DLL_LEN] = '\0' ;
#endif
         break;
    case 1:     //GW21E, 2 -> 1, update 2003/06/06
#if 0
	     strcpy(GW_Model,MODELNAME2);
	     strcpy(GW_Dll,DLLNAME2);
#else   //---update 2003/03/03

	     strncpy( GW_Model, MODELNAME2, GW_MODEL_LEN);
	     GW_Model[ GW_MODEL_LEN] = '\0' ;
	     strncpy( GW_Dll, DLLNAME2, GW_DLL_LEN);
	     GW_Dll[ GW_DLL_LEN] = '\0' ;
#endif
        break;
    default:    //GW21L/E, no used, update 2003/06/06
#if 0
	        strcpy(GW_Model,MODELNAME0);
	        strcpy(GW_Dll,DLLNAME0);        
#else   //---update 2003/03/03
         strncpy( GW_Model, MODELNAME0, GW_MODEL_LEN);
         GW_Model[ GW_MODEL_LEN] = '\0' ;
         strncpy( GW_Dll, DLLNAME0, GW_DLL_LEN);
         GW_Dll[ GW_DLL_LEN] = '\0' ;
#endif
        break;
    }

#elif DEF_GW26A >= 1    //333333333333333333333333333333333333333333, add 2005/10/22
    switch(Getporttype())
    {
    case 1://access
        strncpy( GW_Model, MODELNAME1, GW_MODEL_LEN);
        GW_Model[ GW_MODEL_LEN] = '\0' ;
        strncpy( GW_Dll, DLLNAME1, GW_DLL_LEN);
        GW_Dll[ GW_DLL_LEN] = '\0' ;
        break;
    case 3://io
    default:
        strncpy( GW_Model, MODELNAME2, GW_MODEL_LEN);
        GW_Model[ GW_MODEL_LEN] = '\0' ;
        strncpy( GW_Dll, DLLNAME2, GW_DLL_LEN);
        GW_Dll[ GW_DLL_LEN] = '\0' ;
        break;
    }

#else                       //333333333333333333333333333333333333333333
    strcpy(GW_Model,MODELNAME1);
    strcpy(GW_Dll,DLLNAME1);
#endif                      //333333333333333333333333333333333333333333

#if 0                                   //---remark 2011/04/07, by kernel version, odd:shadow mode, even:normal mode
#ifdef DEF_SHADOW_MODE
#if DEF_SHADOW_MODE >= 1                //---add 2008/09/06, if shadow mode, append 'R' to the model name
    if (( i= strlen( GW_Model)) < GW_MODEL_LEN-2) {
        for ( j=i-1; j>=0; j--)
        if ( GW_Model[ j] == '-') break;
        
        if ( j<0) GW_Model[ i++] = '-' ;
        GW_Model[ i] = 'R' ;
    }
#endif
#endif
#endif

// initial 485 type
//#if DEF_SETPORTTYPE >= 1
#if DEF_SET485WIRES >= 1            //---update 2004/08/17
    wordtemp=ReadEEPROM(18);

    if(getmycomtype(0)==TYPE485)
    {
       port0type=(unsigned char)(wordtemp/256);
       setmy485wire(0,(unsigned char)(port0type&0x01));
    }

    if(getmycomtype(1)==TYPE485)
    {
       port1type=(unsigned char)(wordtemp%256);
       setmy485wire(1,(unsigned char)(port1type&0x01));
    }
#endif

    request_frame.Op = REQUEST;
    ack_frame.Op = ACK;
    ack_frame.Htype = request_frame.Htype = ETHER_TYPE;
    ack_frame.Hlen  = request_frame.Hlen  = HW_ADDRESS_LEN;
    ack_frame.Transaction_Id = request_frame.Transaction_Id =  MY_TRANSACTION_ID;

#if 0
    Read_com_cfg();

    GetLocalEthAddr(ghwadd);
    memcpy(ack_frame.C_H_A,ghwadd,6);
    memcpy(request_frame.C_H_A,ghwadd,6);

    GetLocalIPAddr(gipadd);
    memcpy(ack_frame.Client_Address,gipadd,4);
    memcpy(request_frame.Client_Address,gipadd,4);

    GetDefaultGateway(ggwadd);
    memcpy(ack_frame.Gateway_IP_Address,ggwadd,4);
    memcpy(request_frame.Gateway_IP_Address,ggwadd,4);

    GetLocalNetmask(gmask);
    memcpy(ack_frame.Vendor_Specific_Area,gmask,4);
    memcpy(request_frame.Vendor_Specific_Area,gmask,4);
#else
//  Read_ip_cfg() ;                 //update 2003/10/22
    Read_sys_cfg();
#endif

    a=((gmask[0]&0x000000ff)<<0)|((gmask[1]&0x000000ff)<<8);
    b=((gmask[2]&0x000000ff)<<0)|((gmask[3]&0x000000ff)<<8);
    subnetmask=((b&0x0000ffff)<<16)|(a&0x0000ffff);
    a=((gipadd[0]&0x000000ff)<<0)|((gipadd[1]&0x000000ff)<<8);
    b=((gipadd[2]&0x000000ff)<<0)|((gipadd[3]&0x000000ff)<<8);
    netaddr=((b&0x0000ffff)<<16)|(a&0x0000ffff);
    netaddr=netaddr&subnetmask;
    debugaddr=netaddr;
    debugflag=1;

    if(udp_open((int *)&handle,(short *)&localport)== SUCCESS)
    {
       if(SUCCESS==udp_receive(handle,receive_buf,sizeof(struct Bootp_format),daprecv_post,0))
       {
       }else
       {
       }
       //report_config();
       //udp_send(handle,bootp_port,bootp_address,request_buf,sizeof(struct Bootp_format));//send to local
#if 0
#if DEF_GW21S256 >= 1
       set_timer( 20,(unsigned long) 0, report_config );
#else
       set_timer( 5,(unsigned long) 0, report_config );
#endif

#else
       set_timer( 18*2,(unsigned long) 0, report_config );  //more delay, 20 -> 18*1, update 2003/07/25
#endif
    }
    else
    {
    }

 #if 1 //---add 2011/09/26
    set_timer( 18*3,(unsigned long)0, report_gateway);
 #endif

//  DebugOpen() ;    // add 2003/03/25
//  DebugSend( "System Start", -1) ;
}

/////////////////////////////////////////////////////////////////////////////
void _loadds far ConfigRoutine()
{
    Initial();
    load_utable();
}

//---add 2003/10/28
void _loadds far Read_sys_cfg()
{
    unsigned short  wordtemp ;
    unsigned char   dhcp_flag, country ;
    int     total_port ;

    strncpy( request_frame.Server_Host_Name, GW_Model, GW_MODEL_LEN);
    total_port = GetTotalPort() ;               //count of serial ports, add 2003/01/24
    request_frame.Server_Host_Name[63] = total_port ;

    wordtemp = GetVersion();
    request_frame.FName[0] =  (unsigned char)(wordtemp/256);
    request_frame.FName[1] =  (unsigned char)(wordtemp%256);
    GetAPSerialNo((char far *)&request_frame.FName[2]);

    GetLocalEthAddr(ghwadd);
    memcpy(ack_frame.C_H_A,ghwadd,6);
    memcpy(request_frame.C_H_A,ghwadd,6);

    GetLocalIPAddr(gipadd);
    memcpy(ack_frame.Client_Address,gipadd,4);
    memcpy(request_frame.Client_Address,gipadd,4);

    GetDefaultGateway(ggwadd);
    memcpy(ack_frame.Gateway_IP_Address,ggwadd,4);
    memcpy(request_frame.Gateway_IP_Address,ggwadd,4);

    GetLocalNetmask(gmask);
    memcpy(ack_frame.Vendor_Specific_Area,gmask,4);
    memcpy(request_frame.Vendor_Specific_Area,gmask,4);

#if 1   
    //---add 2003/10/28
    GetHostName( (char far *)&HostName[0]) ;
    memcpy( &ack_frame.Server_Host_Name[46], HostName, HOST_NAME_LEN);
    memcpy( &request_frame.Server_Host_Name[46], HostName, HOST_NAME_LEN);
    
    //---add 2003/10/28
    dhcp_flag = ( MY_IPADDR)? 0: 1 ; 
    ack_frame.Server_Host_Name[ 62] = dhcp_flag ;
    request_frame.Server_Host_Name[ 62] = dhcp_flag ;

    //---add 2004/05/05
    country = GetCountry() & 0xFF ;
    ack_frame.Server_Host_Name[ 45] = country ;
    request_frame.Server_Host_Name[ 45] = country ;
#endif

    //---add download type, (0:unknown, 1:80186, 2:mega, 3:IDT, 4:PPC), 2009/10/19, 
    ack_frame.Server_Host_Name[ 44] = 1 ;
    request_frame.Server_Host_Name[ 44] = 1 ;
    
    //---for manufactured testing, gwcfg.exe
    wordtemp=ReadEEPROM(20);        //---for port 1
    ack_frame.Vendor_Specific_Area[11] = (unsigned char)(wordtemp/256);
    request_frame.Vendor_Specific_Area[11] = (unsigned char)(wordtemp/256);
    ack_frame.Vendor_Specific_Area[10] = (unsigned char)(wordtemp%256);
    request_frame.Vendor_Specific_Area[10] = (unsigned char)(wordtemp%256);

    wordtemp=ReadEEPROM(21);
    ack_frame.Vendor_Specific_Area[13] = (unsigned char)(wordtemp/256);
    request_frame.Vendor_Specific_Area[13] = (unsigned char)(wordtemp/256);
    ack_frame.Vendor_Specific_Area[12] = (unsigned char)(wordtemp%256);
    request_frame.Vendor_Specific_Area[12] = (unsigned char)(wordtemp%256);

    if ( total_port > 1)
    {
    wordtemp=ReadEEPROM(22);        //---for port 2
    ack_frame.Vendor_Specific_Area[15] = (unsigned char)(wordtemp/256);
    request_frame.Vendor_Specific_Area[15] = (unsigned char)(wordtemp/256);
    ack_frame.Vendor_Specific_Area[14] = (unsigned char)(wordtemp%256);
    request_frame.Vendor_Specific_Area[14] = (unsigned char)(wordtemp%256);

    wordtemp=ReadEEPROM(23);
    ack_frame.Vendor_Specific_Area[17] = (unsigned char)(wordtemp/256);
    request_frame.Vendor_Specific_Area[17] = (unsigned char)(wordtemp/256);
    ack_frame.Vendor_Specific_Area[16] = (unsigned char)(wordtemp%256);
    request_frame.Vendor_Specific_Area[16] = (unsigned char)(wordtemp%256);
    }
}

//---add 2003/10/28
void _loadds far Set_Host_Name()
{
    int     i;
    
    memcpy( HostName, &bootpframe.Server_Host_Name[46], HOST_NAME_LEN);
    HostName[ HOST_NAME_LEN] = '\0' ;

#if 0 //---move to report_debug0(), 2016/02/18

#ifdef _IS_BACKUP       //---add 2010/03/31
#define BK_MAGIC_CODES  "ATOP EE Backup"
    if ( FAR_memcmp( (char far *)0xD0000000L, (char far *)BK_MAGIC_CODES, 14) == 0)
    if ( FAR_strcmp( (char far *)HostName, (char far *)"@erase backup")==0) {
        EraseSector( 5) ;     //---erase D000 segment
        FAR_strcpy( (char far *)&HostName[0], (char far *)"@erase ok");
    }
#endif

 #ifdef _IS_BACKUP3       //---add 2011/05/06, for testing
    if ( FAR_strcmp( (char far *)HostName, (char far *)"@erase eeprom")==0) {
        for ( i=0; i<256; i++) {
            WriteMacEEPROM( i, 0) ;
        }
    }
 #endif

 #ifdef _IS_BACKUP4      //---add 2011/06/14
    if ( bk_eeprom_flag)
    if ( FAR_strcmp( (char far *)HostName, (char far *)"@erase backup")==0) {
        EraseSector( 5) ;     //---erase D000 segment
        FAR_strcpy( (char far *)&HostName[0], (char far *)"@erase ok");
    }

  #if 1 //---add 2016/01/08, This is failed. System will reboot. 
    if ( FAR_strcmp( (char far *)HostName, (char far *)"@eeprom backup")==0) {
        EE_Backup( 1) ;     //---backup EEPROM to D000 segment
        FAR_strcpy( (char far *)&HostName[0], (char far *)"@backup ok");
        SetHostName( (char far *)&HostName[0]) ;
        return ;
    }
  #endif 
 #endif
#endif //---move to report_debug0(), 2016/02/18

    SetHostName( (char far *)&HostName[0]) ;
}

//---add 2003/10/28
void _loadds far Get_Cfg_Data( struct Bootp_format far *pcfg)
{
    Read_sys_cfg() ;
    FAR_memcpy( (char far *)pcfg, (char far *)&request_frame, sizeof( struct Bootp_format)) ;
}

#ifdef MAC_ENCODE
void _loadds far LoadMacAddress()
{
    int i;
    unsigned char far * cptr=(unsigned char far *)0xf000fef0;
    unsigned short wordtemp=0x59+0x02+0x12;
    unsigned short data;

#ifdef FOREPROM
    {
	unsigned short data;
	unsigned char fbuf[4];

	    data=ReadKerEEPROM(25+(sizeof(utable)/2)) ;
	    fbuf[0]=data%256;
	    fbuf[1]=data/256;
	    data=ReadKerEEPROM(25+(sizeof(utable)/2)+1) ;
	    fbuf[2]=data%256;
	    fbuf[3]=data/256;
	    if(f2t(fbuf,&mac[1])==0)
	    {
		if(macdecode(&mac[1],gmaccode)!=1) while(1);
	    }
    }
#else
    if(macdecode(&mac[1],cptr)!=1)
    {
	if(macdecode(&mac[1],gmaccode)!=1) while(1);
    }
#endif

#if 1 //---add 2018/04/19, In LoadMacAddress(), if MAC is the same, then no overwrite
    data = ReadMacEEPROM( 1);
    if ( (data >> 8) == mac[1]) {
        data = ReadMacEEPROM( 2);
        if ( data == (mac[3]<<8) + mac[2]) {
            return;
        }
    }
#endif

    wordtemp=mac[1]*256+mac[0];
    WriteHwEEPROM(0x01,wordtemp);
    wordtemp=mac[3]*256+mac[2];
    WriteHwEEPROM(0x02,wordtemp);
}

int bdecode(unsigned char far *c_ary,unsigned char far *originalbyte)
{
	int i;
	unsigned char ptrptr; //1~5
	unsigned char ptr;
	unsigned char bitidx;
	unsigned char shiftval;
	unsigned char mask=0x01;
	unsigned char chksum=0x00;
	unsigned char tmp=0x00;

	ptrptr=(c_ary[0]%5)+1;
	ptr=c_ary[ptrptr]%16+8;
	bitidx=c_ary[ptrptr+1]%8;
	shiftval=c_ary[ptrptr+2];
	mask=mask<<bitidx;

	for(i=0;i<8;i++)
	{
		tmp=(tmp<<1);
	    if(c_ary[ptr]&mask)
		{
		   tmp+=1;
		}
		chksum+=c_ary[ptr];
		ptr++;
	}
	if(chksum==c_ary[ptr])
	{
		*originalbyte=tmp-shiftval;
		return 1;
	}
	else return 0;
}


int macdecode(unsigned char far *macaddr,unsigned char far *maccode)
{
	int k;
	unsigned char mac1[3];
	unsigned char mac2[3];
	unsigned char tmp;
	static unsigned char sum;

	tmp=0x00;
	for(k=0;k<3;k++)
	{
	if(bdecode(&maccode[k*32],&mac1[k])!=1) return 0;
		tmp+=mac1[k];
	}

    if(bdecode(&maccode[3*32],&sum)!=1) return 0;
	if(tmp!=sum) return 0;
    tmp+=tmp;

	for(k=0;k<3;k++)
	{
	if(bdecode(&maccode[(k+4)*32],&mac2[k])!=1) return 0;
		tmp+=mac2[k];
	}
    if(bdecode(&maccode[7*32],&sum)!=1) return 0;
	if(tmp!=sum)return 0;

	macaddr[0]=mac2[0];
	macaddr[1]=mac1[1];
	macaddr[2]=mac2[2];

    return 1;
}

#ifdef FOREPROM
void t2f(unsigned char *tptr,unsigned char *fptr)
{
	unsigned char rd;
	unsigned char sum=0;
	unsigned char key[4]={0x59,0xc3,0x68,0x17};
	int i;

	rd=((*fptr)>>6)&0x03;

	for(i=0;i<3;i++)
	{
	   *(fptr+1+i)=*(tptr+i)+key[(rd+i)%4];
	   sum+=((*(fptr+1+i)>>rd)+key[i]);
	}
	*fptr=((*fptr&0xc0)|sum&0x3f);
}

int f2t(unsigned char *fptr,unsigned char *tptr)
{
	unsigned char rd;
	unsigned char sum=0;
	unsigned char key[4]={0x59,0xc3,0x68,0x17};
	int i;

	rd=((*fptr)>>6)&0x03;

	for(i=0;i<3;i++)
	{
	   *(tptr+i)=*(fptr+1+i)-key[(rd+i)%4];
	   sum+=((*(fptr+1+i)>>rd)+key[i]);
	}
	if(*fptr==((*fptr&0xc0)|sum&0x3f)) return 1;
	else return 0;
}
#endif
#endif

#if DEF_GW21W + DEF_GW21WMAXI >= 1
//---add 2007/04/02, for LAN re-connect to report I am here
#if 0
void report_config2()
{
    Read_sys_cfg() ;

    if ( local_ipaddr)
    udp_send(handle,bootp_port,bootp_address,request_buf,sizeof(struct Bootp_format));
}

#else  //---update 2007/07/25, to report more times
void report_config2( unsigned short report_seq)
{
    static unsigned short old_seq = 1000 ;
    
    if ( local_ipaddr) {
        if ( report_seq) {      //---for not calling in the first time
            Read_sys_cfg() ;
            udp_send(handle,bootp_port,bootp_address,request_buf,sizeof(struct Bootp_format));
        }

        if ( old_seq==0)            //---need to renew a request
            report_seq = 0 ;
        else
        if ( old_seq<3)            //---already in timer, will be called next time
        if ( report_seq==0) {       //---a new request
            old_seq = 0 ;
            return ;
        }
        
        if ( report_seq<3)         //---add 2003/07/25
        {
//          BuzzerOn();
            set_timer( ((report_seq==0)? 18: 18*10), (unsigned long)(report_seq+1), report_config2);

            old_seq = report_seq+1 ;
            return ;
        }
    }
    old_seq = 1000 ;
}
#endif
#endif

#if 1                       //---add 2007/05/18
int report_config3( unsigned long host_addr, char far *snd_buf)
{
    //--- is it the same subnet
//  if ((subnetmask & host_addr) != netaddr) 
    if ( (*(long *)&gipadd[0] == 0) || (subnetmask & host_addr) != netaddr)     //---fix a bug that the monitor tool invits nothing if the DHCP cannot get a IP (IP=0.0.0.0), 2009/09/02
    udp_send( handle, bootp_port, bootp_address, snd_buf, sizeof(struct Bootp_format));
    udp_send( handle, bootp_port, host_addr, snd_buf, sizeof(struct Bootp_format));
}

int DelayAckInviteIni() {
    memset( (char *)&AckInviteStru[0], 0, sizeof(DELAY_ACK_INVITE)*DELAY_ACK_INVITE_MAX) ;
}

int DelayAckInviteAdd( unsigned long src_ip, int delay_ms) {
    int     i ;
    
    if ( delay_ms)      //---add 2007/07/10, fix a bug if delay_ms==0
    
    for ( i=0; i<DELAY_ACK_INVITE_MAX; i++)
    if ( AckInviteStru[i].delay_ms==0) {
        AckInviteStru[i].delay_ms = delay_ms ;
        AckInviteStru[i].src_ip = src_ip ;
        return( 1) ;
    }
    return( 0) ;
}

int DelayAckInviteService() {
    int     i, ret ;
    
    for ( i=0; i<DELAY_ACK_INVITE_MAX; i++)
    if ( AckInviteStru[i].delay_ms) {
        AckInviteStru[i].delay_ms-- ;

//      if ( AckInviteStru[i].delay_ms==0) {
        if ( AckInviteStru[i].delay_ms<=0) {            //---update 2007/10/23
        
            Read_sys_cfg();
            ret = report_config3( AckInviteStru[i].src_ip, request_buf);

#if 1 //---add 2007/10/23
            if ( ret==0 || AckInviteStru[i].delay_ms<=-3) {     //---send ok or too many failed
                AckInviteStru[i].delay_ms = 0 ;
            }
            else
            if ( AckInviteStru[i].delay_ms==0) AckInviteStru[i].delay_ms = -1 ;
#endif
        }
    }
}
#endif

#if 1 //---add 2016/01/18, show EEPROM or Flash
int report_debug() {
    int     i;
    unsigned long laddr;
    unsigned short udata;
    char far *tmpstr ;
    unsigned char far *tmpbuf;

 #ifdef _IS_BACKUP4      //---add 2016/02/18
    unsigned char far *addr_ptr, *id_ptr ;

    extern unsigned char EE_Table[];
    extern int     bk_eeprom_flag;

#define BK_MAGIC_CODES  "ATOP EEBackup4"
 #endif    

    unsigned long FAR_axtol( char far *str, int len);
    int FAR_dumpstr( char far *src, int cnt, char far *dst);

    tmpstr = (char far *)&bootpframe.Server_Host_Name[46];    

    if ( FAR_memcmp( tmpstr, (char far *)"@ee ", 4)==0) {
        //---different gateway makes SerialManager auto updated
        *(unsigned long *)&ack_frame.Gateway_IP_Address[0] = *(unsigned long *)&bootpframe.Gateway_IP_Address[0];
    
        laddr = FAR_axtol( &tmpstr[4], HOST_NAME_LEN-4);
        tmpstr = (char far *)&ack_frame.Server_Host_Name[46];
        
        for ( i=0; i<HOST_NAME_LEN/4; i++) {
            if ( laddr + i >= 512) {            //---EEPROM address out of range
                *(unsigned long *)tmpstr = 0x20202020;
                continue;                
            }
            
            udata = ReadMacEEPROM( (unsigned int)laddr + i) ;
            FAR_dumpstr( (char far *)&udata, 2, tmpstr + i*4);
        }
        return 1;
    }
    else
    if ( FAR_memcmp( tmpstr, (char far *)"@mm ", 4)==0) {
        laddr = FAR_axtol( &tmpstr[4], HOST_NAME_LEN-4);
        tmpbuf = (unsigned char far *)(((laddr & 0xFFFF0000) << 12) + (laddr & 0xFFFF));
        tmpstr = (char far *)&ack_frame.Server_Host_Name[46];
            
        FAR_dumpstr( (char far *)tmpbuf, HOST_NAME_LEN/2, tmpstr);
        return 1;
    }

#if 1 //---add 2016/02/18

#ifdef _IS_BACKUP       //---add 2010/03/31
#define BK_MAGIC_CODES  "ATOP EE Backup"
    if ( FAR_memcmp( (char far *)0xD0000000L, (char far *)BK_MAGIC_CODES, 14) == 0)
    if ( FAR_strcmp( tmpstr, (char far *)"@erase backup")==0) {
        EraseSector( 5) ;     //---erase D000 segment

        SetHostName( (char far *)"@erase ok") ;
        return 1;
    }
#endif

 #ifdef _IS_BACKUP3       //---add 2011/05/06, for testing
    if ( FAR_strcmp( tmpstr, (char far *)"@erase eeprom")==0) {
        for ( i=0; i<256; i++) {
            WriteMacEEPROM( i, 0) ;
        }
        return 1;
    }
 #endif

 #ifdef _IS_BACKUP4      //---add 2011/06/14
    if ( FAR_strcmp( tmpstr, (char far *)"@backup erase")==0) {
        if ( bk_eeprom_flag) {
            EE_Backup( 0) ;       //---erase D000 segment
    
            SetHostName( (char far *)"#Erase ok") ;
        }
        else {
            SetHostName( (char far *)"#No backup") ;
        }
        return 1;
    }
    else
    if ( FAR_strcmp( tmpstr, (char far *)"@backup eeprom")==0) {
        EE_Backup( 1) ;     //---backup EEPROM to D000 segment
        SetHostName( (char far *)"#Backup ok") ;    //---Setting this hostname is useless because running after EEPROM backup
        return 1;
    }
    else
    if ( FAR_strcmp( tmpstr, (char far *)"@backup check")==0) {
        if ( bk_eeprom_flag) {
            SetHostName( (char far *)"#Have backup") ;
        }
        else {
            SetHostName( (char far *)"#No backup") ;
        }
        return 1;
    }
 #endif
#endif //---add 2016/02/18

 #if 1 //---add for shadow mode, 2016/02/19
    if ( FAR_strcmp( tmpstr, (char far *)"@shadow on")==0) {
        ShadowMode_on();
        return 1;
    }
    else
    if ( FAR_strcmp( tmpstr, (char far *)"@shadow off")==0) {
        ShadowMode_off();
        return 1;
    }
    else
    if ( FAR_strcmp( tmpstr, (char far *)"@shadow check")==0) {
 #if 0
        tmpbuf = (unsigned char far *)0xD0000000L; 
        udata = *tmpbuf;   
        (*tmpbuf)++;
        SetHostName( (udata == *tmpbuf) ? (char far *)"#Shadow off" : (char far *)"#Shadow on") ;
        *tmpbuf = udata & 0xFF;
 #else
        SetHostName( (ShadowMode_chk()==0) ? (char far *)"#Shadow off" : (char far *)"#Shadow on") ;
 #endif
        return 1;
    }
 #endif

    return 0;    
}
#endif
