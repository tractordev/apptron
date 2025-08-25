(async () => {
    console.log("running alpine init.js");
    window.makeScreen();
    const w = window.wanix.instance;        
    const vm = (await w.readText("web/vm/new")).trim();
    await w.writeFile("task/1/ctl", `bind #console/data1 web/vm/${vm}/ttyS0`);
    await w.writeFile("task/1/ctl", `bind . web/vm/${vm}/fsys`);
    await w.writeFile("task/1/ctl", `bind #bundle/rootfs web/vm/${vm}/fsys`);
    await w.writeFile(`web/vm/${vm}/ctl`, "start");
})();
