(async () => {
    console.log("running apptron init.js");
    window.makeScreen();
    const w = window.wanix.instance;        
    const vm = (await w.readText("vm/new/default")).trim();
    await w.writeFile("task/1/ctl", `bind #console/data1 vm/${vm}/ttyS0`);
    await w.writeFile("task/1/ctl", `bind . vm/${vm}/fsys`);
    await w.writeFile("task/1/ctl", `bind #ramfs vm/${vm}/fsys/#ramfs`);
    await w.writeFile("task/1/ctl", `bind #pipe vm/${vm}/fsys/#pipe`);
    await w.writeFile("task/1/ctl", `bind #| vm/${vm}/fsys/#|`);
    await w.writeFile("task/1/ctl", `bind #bundle/rootfs vm/${vm}/fsys`);
    const cmdline = [
        "init=/bin/init",
        "rw",
        "root=host9p",
        "rootfstype=9p",
        `rootflags=trans=virtio,version=9p2000.L,aname=vm/${vm}/fsys,cache=none`,
    ];
    const ctlcmd = ["start", "-append", `'${cmdline.join(" ")}'`];
    if (window.wanix.config.network) {
        ctlcmd.push("-netdev");
        ctlcmd.push(`user,type=virtio,relay_url=${window.wanix.config.network}`);
    }
    await w.writeFile(`vm/${vm}/ctl`, ctlcmd.join(" "));
})();
