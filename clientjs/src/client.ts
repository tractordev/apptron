// @ts-ignore
import * as qtalk from "../lib/qtalk.min.js";


export async function connect(url: string): Promise<Client> {
  return new Client(await qtalk.connect(url, new qtalk.CBORCodec()))
}

export class Client {
  peer: qtalk.Peer
  rpc: any
  data: {[index: string]: ArrayBuffer}

  app: app
  menu: menu
  system: system
  shell: shell
  window: window

  onevent?: (e: Event) => void

  constructor(peer: qtalk.Peer) {
    this.data = {}
    this.peer = peer
    this.rpc = peer.virtualize()
    this.app = new AppModule(this)
    this.menu = new MenuModule(this)
    this.system = new SystemModule(this)
    this.shell = new ShellModule(this)
    this.window = new WindowModule(this)
    this.handleEvents(peer)
    this.peer.respond()
  }

  async serveData(d: ArrayBuffer): Promise<string> {
    const selector = toHexString(new Uint8Array(await crypto.subtle.digest("SHA-1", d)))
    if (selector in this.data) {
      return selector
    }
    this.data[selector] = d
    this.peer.handle(selector, {respondRPC: async (resp: any, call: any) => {
      await call.receive()
      const ch = await resp.continue()
      const data = new Uint8Array(this.data[selector])
      await ch.write(data)
      await ch.close()
    }})
    return selector
  }

  async handleEvents(peer: qtalk.Peer) {
    const resp = await peer.call("Listen")
    while (true) {
      const obj = await resp.receive()
      if (obj === null) {
        break;
      }
      const event = obj as Event
      if (this.onevent) {
        this.onevent(event)
      }
      switch (event.Name) {
        case "menu":
          if (this.menu.onclick)
            this.menu.onclick(event)
          break
        case "shortcut":
          if (this.shell.onshortcut)
            this.shell.onshortcut(event)
          break
        default:
          const w = this.window.windows[event.WindowID]
          if (w) {
            switch (event.Name) {
              case "close":
                if (w.onclose) w.onclose(event)
                break
              case "destroy":
                if (w.ondestroyed) w.ondestroyed(event)
                delete this.window.windows[event.WindowID]
                break
              case "focus":
                if (w.onfocused) w.onfocused(event)
                break
              case "blur":
                if (w.onblurred) w.onblurred(event)
                break
              case "resize":
                if (w.onresized) w.onresized(event)
                break
              case "move":
                if (w.onmoved) w.onmoved(event)
                break
            }
          }
      }
    }
  }
}

export interface app {
  Run(options: AppOptions): void
  Menu(): Promise<Menu>
  SetMenu(m: Menu): void 
  NewIndicator(icon: ArrayBuffer, items: MenuItem[]): void
}

class AppModule {
  rpc: any
  client: Client

  constructor(client: Client) {
    this.rpc = client.rpc
    this.client = client
  }

  Run(options: AppOptions): void {
    this.rpc.app.Run(options)
  }

  Menu(): Promise<Menu> {
    return this.rpc.app.Menu()
  }

  SetMenu(m: Menu): void {
    this.rpc.app.SetMenu(m)
  }

  async NewIndicator(icon: ArrayBuffer, items: MenuItem[]): Promise<void> {
    const selector = await this.client.serveData(icon)
    this.rpc.app.NewIndicator(selector, items)
  }
}

export interface menu {
  onclick?: (e: Event) => void

  New(items: MenuItem[]): Promise<Menu>
  Popup(items: MenuItem[]): Promise<number>
}

class MenuModule {
  rpc: any

  onclick?: (e: Event) => void

  constructor(client: Client) {
    this.rpc = client.rpc
  }

  New(items: MenuItem[]): Promise<Menu> {
    return this.rpc.menu.New(items)
  }

  Popup(items: MenuItem[]): Promise<number> {
    return this.rpc.menu.Popup(items)
  }
}

export interface system {
  Displays(): Promise<Display[]>
}

class SystemModule {
  rpc: any

  constructor(client: Client) {
    this.rpc = client.rpc
  }

  Displays(): Promise<Display[]> {
    return this.rpc.system.Displays()
  }
}

export interface shell {
  onshortcut?: (e: Event) => void

  ShowNotification(n: Notification): void
  ShowMessage(msg: MessageDialog): void
  ShowFilePicker(fd: FileDialog): Promise<string[]>
  ReadClipboard(): Promise<string>
  WriteClipboard(text: string): Promise<boolean>
  RegisterShortcut(accelerator: string): Promise<boolean>
  IsShortcutRegistered(accelerator: string): Promise<boolean>
  UnregisterShortcut(accelerator: string): Promise<boolean>
  UnregisterAllShortcuts(): Promise<boolean>
}

class ShellModule {
  rpc: any

  onshortcut?: (e: Event) => void

  constructor(client: Client) {
    this.rpc = client.rpc
  }

  ShowNotification(n: Notification): void {
    this.rpc.shell.ShowNotification(n)
  }

  ShowMessage(msg: MessageDialog): void {
    this.rpc.shell.ShowMessage(msg)
  }

  ShowFilePicker(fd: FileDialog): Promise<string[]> {
    return this.rpc.shell.ShowFilePicker(fd)
  }

  ReadClipboard(): Promise<string> {
    return this.rpc.shell.ReadClipboard()
  }

  WriteClipboard(text: string): Promise<boolean> {
    return this.rpc.shell.WriteClipboard(text)
  }

  RegisterShortcut(accelerator: string): Promise<boolean> {
    return this.rpc.shell.RegisterShortcut(accelerator)
  }

  IsShortcutRegistered(accelerator: string): Promise<boolean> {
    return this.rpc.shell.IsShortcutRegistered(accelerator)
  }

  UnregisterShortcut(accelerator: string): Promise<boolean> {
    return this.rpc.shell.UnregisterShortcut(accelerator)
  }
  
  UnregisterAllShortcuts(): Promise<boolean> {
    return this.rpc.shell.UnregisterAllShortcuts()
  }
}


export interface window {
  main: Window
  windows: {[id: number]: Window}

  New(options: WindowOptions): Promise<Window>
}

class WindowModule {
  rpc: any
  main: Window
  windows: {[id: number]: Window}

  constructor(client: Client) {
    this.rpc = client.rpc
    this.main = new Window(this.rpc, 0)
    this.windows = {0: this.main}
  }

  async New(options: WindowOptions): Promise<Window> {
    const w = await this.rpc.window.New(options)
    this.windows[w.ID] = new Window(this.rpc, w.ID)
    return this.windows[w.ID]
  }
}

export class Menu {
  ID: number
  rpc: any

  constructor(rpc: any, id: number) {
    this.rpc = rpc
    this.ID = id
  }
}


export class Window {
  ID: number 
  rpc: any

  onclose?: (e: Event) => void
  onmoved?: (e: Event) => void
  onresized?: (e: Event) => void
  ondestroyed?: (e: Event) => void
  onfocused?: (e: Event) => void
  onblurred?: (e: Event) => void

  constructor(rpc: any, id: number) {
    this.rpc = rpc
    this.ID = id
  }

  // Destroy
  destroy() {
    this.rpc.window.Destroy(this.ID)
  }

  // Focus
  focus() {
    this.rpc.window.Focus(this.ID)
  }

  // GetOuterPosition
  getOuterPosition(): Promise<Position> {
    return this.rpc.window.GetOuterPosition(this.ID)
  }

  // GetOuterSize
  getOuterSize(): Promise<Size> {
    return this.rpc.window.GetOuterSize(this.ID)
  }

  // IsDestroyed
  isDestroyed(): Promise<boolean> {
    return this.rpc.window.IsDestroyed(this.ID)
  }

  // IsVisible
  isVisible(): Promise<boolean> {
    return this.rpc.window.IsVisible(this.ID)
  }

  // SetVisible
  setVisible(visible: boolean) {
    return this.rpc.window.SetVisible(this.ID, visible)
  }

  // SetMaximized
  setMaximized(maximized: boolean) {
    return this.rpc.window.SetMaximized(this.ID, maximized)
  }

  // SetMinimized
  setMinimized(minimized: boolean) {
    return this.rpc.window.SetMinimized(this.ID, minimized)
  }

  // SetFullscreen
  setFullscreen(fullscreen: boolean) {
    return this.rpc.window.SetFullscreen(this.ID, fullscreen)
  }

  // SetMinSize
  setMinSize(size: Size) {
    return this.rpc.window.SetMinSize(this.ID, size)
  }

  // SetMaxSize
  setMaxSize(size: Size) {
    return this.rpc.window.SetMaxSize(this.ID, size)
  }

  // SetResizable
  setResizable(resizable: boolean) {
    return this.rpc.window.SetResizable(this.ID, resizable)
  }

  // SetAlwaysOnTop
  setAlwaysOnTop(always: boolean) {
    return this.rpc.window.SetAlwaysOnTop(this.ID, always)
  }

  // SetSize
  setSize(size: Size) {
    return this.rpc.window.SetSize(this.ID, size)
  }

  // SetPosition
  setPosition(position: Position) {
    return this.rpc.window.SetPosition(this.ID, position)
  }

  // SetTitle
  setTitle(title: string) {
    return this.rpc.window.SetTitle(this.ID, title)
  }
}

export interface MenuItem {
	ID:          number
	Title:       string
	Enabled:     boolean
	Selected:    boolean
	Accelerator: string
	SubMenu:     MenuItem[]
}

export interface Notification {
	Title:    string
	Subtitle: string // for MacOS only
	Body:     string
}

export interface MessageDialog {
	Title:   string
	Body:    string
	Level:   string // info, warning, error
	Buttons: string // ok, okcancel, yesno
}

export interface Size {
  Width:  number
  Height: number
}

export interface Position {
  X: number
  Y: number
}

export interface Display {
	Name:        string
	Size:        Size
	Position:    Position
	ScaleFactor: number
}

export interface AppOptions {
  Identifier:          string
	RunsAfterLastWindow: boolean
	AccessoryMode:       boolean
}

export interface WindowOptions {
	AlwaysOnTop: boolean
	Frameless:   boolean
	Fullscreen:  boolean
	Size:        Size
	MinSize:     Size
	MaxSize:     Size
	Maximized:   boolean
	Position:    Position
	Resizable:   boolean
	Title:       string
	Transparent: boolean
	Visible:     boolean
  Hidden:      boolean
	Center:      boolean
	IconSel:     string // TODO
	URL:         string
	HTML:        string
	Script:      string
}

export interface FileDialog {
	Title:     string
	Directory: string
	Filename:  string
	Mode:      string   // pickfile, pickfiles, pickfolder, savefile
	Filters:   string[] // each string is comma delimited (go,rs,toml) with optional label prefix (text:go,txt)
}


export interface Event {
	Type:     number
	Name:     string
	WindowID: number
	Position: Position
	Size:     Size
	MenuID:   number
	Shortcut: string
}

function toHexString(byteArray: Uint8Array): string {
  return Array.from(byteArray, function(byte: number) {
    return ('0' + (byte & 0xFF).toString(16)).slice(-2);
  }).join('')
}