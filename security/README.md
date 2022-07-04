

# How to use



## Docker

#### build

```dockerfile
 docker build -t security .
```

#### run

```dockerfile
 docker run  --name=security -it -p8080:8080 security run 
```



## Run

First Please install  step 

https://smallstep.com/docs/step-cli/installation

```go
go run .\api\service\main\main.go  
```

or

```go
go run .\main.go run 
```

### Install Step  ca

#### Select

please select standalone

```shell
Use the arrow keys to navigate: ↓ ↑ → ← 
? What deployment type would you like to configure?:
  ▸ Standalone - step-ca instance you run yourself
    Linked - standalone, plus cloud configuration, reporting & alerting
    Hosted - fully-managed step-ca cloud instance run for you by smallstep
```

#### Input info

```shell
✔ Deployment Type: Standalone
What would you like to name your new PKI?
✔ (e.g. Smallstep): Atop Inc.
What DNS names or IP addresses would you like to add to your new CA?
✔ (e.g. ca.smallstep.com[,1.1.1.1,etc.]): localhost
What IP and port will your new CA bind to?
✔ (e.g. :443 or 127.0.0.1:443): 127.0.0.1:8443
What would you like to name the CA's first provisioner?
✔ (e.g. you@smallstep.com): atop@example.com
Choose a password for your CA keys and first provisioner.
✔ [leave empty and we'll generate one]: atop

Generating root certificate... done!
Generating intermediate certificate... done!

✔ Root certificate: C:\Users\asus\.step\certs\root_ca.crt
✔ Root private key: C:\Users\asus\.step\secrets\root_ca_key
✔ Root fingerprint: 393c0f8b6046b1c9bc53bf242864c8a0e7f1b877f41aaeccc8ad945c04a93d9d
✔ Intermediate certificate: C:\Users\asus\.step\certs\intermediate_ca.crt
✔ Intermediate private key: C:\Users\asus\.step\secrets\intermediate_ca_key
✔ Database folder: C:\Users\asus\.step\db
✔ Default configuration: C:\Users\asus\.step\config\defaults.json
✔ Certificate Authority configuration: C:\Users\asus\.step\config\ca.json

Your PKI is ready to go. To generate certificates for individual services see 'step help ca'.

FEEDBACK 😍 🍻
  The step utility is not instrumented for usage statistics. It does not phone
  home. But your feedback is extremely valuable. Any information you can provide
  regarding how you’re using `step` helps. Please send us a sentence or two,
  good or bad at feedback@smallstep.com or join GitHub Discussions
  https://github.com/smallstep/certificates/discussions and our Discord
  https://u.step.sm/discord.
```

#### Run Step Ca Service

```shell
badger 2022/04/23 11:30:14 INFO: All 0 tables opened in 0s
Please enter the password to decrypt C:\Users\Az\.step\secrets\intermediate_ca_key:atop
2022/04/23 11:30:35 Starting Smallstep CLI/0000000-dev (windows/amd64)
2022/04/23 11:30:35 Documentation: https://u.step.sm/docs/ca
2022/04/23 11:30:35 Community Discord: https://u.step.sm/discord
2022/04/23 11:30:35 Config file:
2022/04/23 11:30:35 The primary server URL is https://localhost:8443
2022/04/23 11:30:35 Root certificates are available at https://localhost:8443/roots.pem
2022/04/23 11:30:35 X.509 Root Fingerprint: fa089a8085e928f31d1e1065610c8e05bdf5979c5b6df0f604e286962590936d
2022/04/23 11:30:35 Serving HTTPS on 127.0.0.1:8443 ...
```



#### Run Bootstrap Client

```shell
✔ Provisioner: atop (JWK) [kid: EA2yGB6vUj-MFevP3CN_wB28AiLHKQ37RoxOvUeLDPM]
Please enter the password to decrypt the provisioner key:atop
✔ CA: https://localhost:8443
✔ Certificate: srv.crt
✔ Private Key: srv.key
The root certificate has been saved in root.crt.
badger 2022/04/23 11:42:29 INFO: Storing value log head: {Fid:0 Len:30 Offset:2896}
badger 2022/04/23 11:42:29 INFO: [Compactor: 173] Running compaction: {level:0 score:1.73 dropPrefixes:[]} for level: 0
badger 2022/04/23 11:42:29 INFO: LOG Compact 0->1, del 1 tables, add 1 tables, took 2.1374ms
badger 2022/04/23 11:42:29 INFO: [Compactor: 173] Compaction for level: 0 DONE
badger 2022/04/23 11:42:29 INFO: Force compaction on level 0 done

```

#### Run Security Service

```shell
2022/04/23 11:42:29 Start service.....8080
```

