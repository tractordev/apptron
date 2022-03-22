window["$host"] = {
  ready: new Promise((resolve) => {
    import("/-/client.js").then(async (mod) => {
      window["$host"] = await mod.connect(`ws://${window.location.host}/-/ws`)
      window["$host"].ready = Promise.resolve()
      resolve()
    })
  })
}