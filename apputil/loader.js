window["$apptron"] = {
  ready: new Promise((resolve) => {
    import("/-/client.js").then(async (mod) => {
      window["$apptron"] = await mod.connect(`ws://${window.location.host}/-/ws`)
      window["$apptron"].ready = Promise.resolve()
      resolve()
    })
  })
}