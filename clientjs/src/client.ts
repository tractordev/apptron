// @ts-ignore
import * as qtalk from "../lib/qtalk.min.js";

(async ()=>{
  if (window) {
    // @ts-ignore
    window.$host = await connect(`ws://${window.location.host}/`)
  }
})()

export async function connect(url: string): Promise<Client> {
  return new Client(await qtalk.connect(url, new qtalk.JSONCodec()))
}

export class Client {
  rpc: any

  constructor(peer: qtalk.Peer) {
    this.rpc = peer.virtualize()
  }

  get screen(): screen {
    return new Screen(this.rpc)
  }

  get shell(): shell {
    return new Shell(this.rpc)
  }

  get window() {
    return {
      New: async (options: WindowOptions): Promise<Window> => {
        const w = await this.rpc.window.New(options)
        return new Window(this.rpc, w.ID)
      },
    }
  }

}

export interface screen {
  Displays(): Promise<Display[]>
}

class Screen {
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

class Shell {
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

export class Window {
  ID: number 
  rpc: any

  constructor(rpc: any, id: number) {
    this.rpc = rpc
    this.ID = id
  }

  // Destroy
  async destroy() {
    await this.rpc.window.Destroy(this.ID)
  }

  // Focus
  async focus() {
    await this.rpc.window.Focus(this.ID)
  }

  // GetOuterPosition
  async getOuterPosition(): Promise<Position> {
    return await this.rpc.window.GetOuterPosition(this.ID)
  }

  // GetOuterSize
  async getOuterSize(): Promise<Size> {
    return await this.rpc.window.GetOuterSize(this.ID)
  }

  // IsDestroyed
  async isDestroyed(): Promise<boolean> {
    return await this.rpc.window.IsDestroyed(this.ID)
  }

  // IsVisible
  async isVisible(): Promise<boolean> {
    return await this.rpc.window.IsVisible(this.ID)
  }

  // SetVisible
  async setVisible(visible: boolean) {
    return await this.rpc.window.SetVisible(this.ID, visible)
  }

  // SetMaximized
  async setMaximized(maximized: boolean) {
    return await this.rpc.window.SetMaximized(this.ID, maximized)
  }

  // SetMinimized
  async setMinimized(minimized: boolean) {
    return await this.rpc.window.SetMinimized(this.ID, minimized)
  }

  // SetFullscreen
  async setFullscreen(fullscreen: boolean) {
    return await this.rpc.window.SetFullscreen(this.ID, fullscreen)
  }

  // SetMinSize
  async setMinSize(size: Size) {
    return await this.rpc.window.SetMinSize(this.ID, size)
  }

  // SetMaxSize
  async setMaxSize(size: Size) {
    return await this.rpc.window.SetMaxSize(this.ID, size)
  }

  // SetResizable
  async setResizable(resizable: boolean) {
    return await this.rpc.window.SetResizable(this.ID, resizable)
  }

  // SetAlwaysOnTop
  async setAlwaysOnTop(always: boolean) {
    return await this.rpc.window.SetAlwaysOnTop(this.ID, always)
  }

  // SetSize
  async setSize(size: Size) {
    return await this.rpc.window.SetSize(this.ID, size)
  }

  // SetPosition
  async setPosition(position: Position) {
    return await this.rpc.window.SetPosition(this.ID, position)
  }

  // SetTitle
  async setTitle(title: string) {
    return await this.rpc.window.SetTitle(this.ID, title)
  }
}