package app

import (
	"fmt"
)

import (
	"tractor.dev/apptron/bridge/api/menu"
	"tractor.dev/apptron/bridge/win32"
	"tractor.dev/apptron/bridge/platform"
)

var (
	mainMenu *menu.Menu
	//app      cocoa.NSApplication
)

func init() {
	//app = cocoa.NSApp()

	//syscall.Syscall(procSetProcessDpiAwarenessContext.Addr(), DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2, 0, 0, 0)


	/*
  SetProcessDpiAwarenessContext(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2);

  //Win32RegisterWindowClass(TrayWindowClassName, GetModuleHandle(0), TrayWindowCallback);

  WNDCLASSEX WindowClass = {sizeof(WindowClass)};
  WindowClass.style = Style;
  WindowClass.lpfnWndProc = Callback;
  WindowClass.hInstance = HInstance;
  WindowClass.hIcon = LoadIcon(WindowClass.hInstance, MAKEINTRESOURCE(101));
  WindowClass.hCursor = LoadCursor(0, IDC_ARROW);
  WindowClass.lpszClassName = Name;

  return RegisterClassEx(&WindowClass) != 0;

  TrayWindow = CreateWindowEx(0, TrayWindowClassName, 0, 0,
                              0, 0, 1, 1,
                              0, 0, GetModuleHandle(0), 0);
  */
}

func Menu() *menu.Menu {
	return mainMenu
}

func SetMenu(menu *menu.Menu) error {
	//app.SetMainMenu(menu.NSMenu)
	mainMenu = menu
	return nil
}

func NewIndicator(icon []byte, items []menu.Item) {
	fmt.Println("NewIndicator", icon)

	/*
	obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
	obj.Retain()
	//obj.Button().SetTitle(i.Text)
	data := mac.NSData_WithBytes(icon, uint64(len(icon)))
	image := cocoa.NSImage_InitWithData(data)
	image.SetSize(mac.Size(16.0, 16.0))
	image.SetTemplate(true)
	obj.Button().SetImage(image)
	obj.Button().SetImagePosition(cocoa.NSImageOnly)

	menu := menu.New(items)
	obj.SetMenu(menu.NSMenu)
	*/
}

func Run(options Options) error {
  win32.SetProcessDpiAwarenessContext(win32.DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)

	win32.SetupTray()
  /*
	win32.Main()

	for {
		win32.PollEvents()
	}
	*/

	platform.Start()
	return nil
}
