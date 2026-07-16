# README

```
C:\USERS\ADMIN\DESKTOP\FOLDER_CLEAN\CODES\CODINGENVHERE\WSML
├───cmd
│   └───wslbm
├───distros
├───internal
│   ├───distro
│   └───wslin
└───pkg
    └───downloader
```
The tree file window cli doesnt show file so do urself

This project
make os managing easier/install other os easier
wslbm
wsl(window subsystem linux) better more

internal/wslin/ :
`add.go` will add os NOT on the LIST wsl
`install.go` will add OS on the list wsl
`shut.go` will shutdown all os that are running currently / wslbm shut -s `<distros>`
`list.go` will list os On the list and on repo / wslbm -l -nv (wslbm -list -nv) will only show os that are on the list repo

internal/distro/ :
`registry.go` yk yk

internal/distros/:
`*.json` is MAIN repo of ROOTFS for examples
```json
{
    "verjson": "0.1", //Use for Speficing Version of os old
    "name": "Chimera",
    "url": "https://repo.chimera-linux.org/live/latest/chimera-linux-x86_64-ROOTFS-20251220-full.tar.gz",
    "installationtype": "tar"
}
```

cmd/wslbm:
`main.go` Main access to all of this above