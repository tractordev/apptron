// @ts-ignore
import * as qtalk from "../lib/qtalk.min.js";

(()=>{
  if (window) {
    window.requestAnimationFrame(async () => {
      // @ts-ignore
      window["$host"] = await connect(`ws://${window.location.host}/`)
    })
  }
})()

export async function connect(url: string): Promise<Client> {
  return new Client(await qtalk.connect(url, new qtalk.JSONCodec()))
}

export class Client {
  rpc: any
  app: app
  menu: menu
  screen: screen
  shell: shell
  window: window

  onevent?: (e: Event) => void

  constructor(peer: qtalk.Peer) {
    this.rpc = peer.virtualize()
    this.app = new AppModule(this.rpc)
    this.menu = new MenuModule(this.rpc)
    this.screen = new ScreenModule(this.rpc)
    this.shell = new ShellModule(this.rpc)
    this.window = new WindowModule(this.rpc)
    ;(async () => {
      const resp = await peer.call("Listen")
      while (true) {
        const e = await resp.receive()
        if (e === null) {
          break;
        }
        if (this.onevent) {
          this.onevent(e as Event)
        }
      }
    })()
    peer.respond()
  }
}

export interface app {
  Menu(): Promise<Menu>
  SetMenu(m: Menu): void 
  //NewIndicator(icon, items)
}

class AppModule {
  rpc: any

  constructor(rpc: any) {
    this.rpc = rpc
  }

  Menu(): Promise<Menu> {
    return this.rpc.app.Menu()
  }

  SetMenu(m: Menu): void {
    this.rpc.app.SetMenu(m)
  }
}

export interface menu {
  New(items: MenuItem[]): Promise<Menu>
}

class MenuModule {
  rpc: any

  constructor(rpc: any) {
    this.rpc = rpc
  }

  New(items: MenuItem[]): Promise<Menu> {
    return this.rpc.menu.New(items)
  }
}

export interface screen {
  Displays(): Promise<Display[]>
}

class ScreenModule {
  rpc: any

  constructor(rpc: any) {
    this.rpc = rpc
  }

  Displays(): Promise<Display[]> {
    return this.rpc.screen.Displays()
  }
}

export interface shell {
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

  constructor(rpc: any) {
    this.rpc = rpc
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
  New(options: WindowOptions): Promise<Window>
}

class WindowModule {
  rpc: any

  constructor(rpc: any) {
    this.rpc = rpc
  }

  async New(options: WindowOptions): Promise<Window> {
    const w = await this.rpc.window.New(options)
    return new Window(this.rpc, w.ID)
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

  onmoved?: (e: Event) => void
  onresized?: (e: Event) => void
  // TODO: more events

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