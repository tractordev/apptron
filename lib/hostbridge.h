#include <stdint.h>

//
// Types
//

#define bool uint8_t

#if defined(_WIN32)
  #define OS_WINDOWS 1
#elif defined(__APPLE__)
  #define OS_MACOS 1
#else
  //#define OS_LINUX 1
  #error "[os.h] Unsuported operating system!"
#endif

typedef struct Position {
	double x;
	double y;
} Position;

typedef struct Size {
	double width;
	double height;
} Size;

typedef enum EventType {
	EventNone      = 0,
	EventClose     = 1,
	EventDestroyed = 2,
	EventFocused   = 3,
	EventBlurred   = 4,
	EventResized   = 5,
	EventMoved     = 6,
	EventMenuItem  = 7,
	EventShortcut  = 8,
} EventType;

typedef struct Event {
	int      event_type;
	int      window_id;
	Position position;
	Size     size;
	int      menu_id;
	char *   shortcut;
} Event;

// NOTE(nick): this has to be kept in sync with wry's EventLoop struct size
typedef struct EventLoop {
	unsigned char data[40];
} EventLoop;

// NOTE(nick): this has to be kept in sync with wry's Menu struct size
typedef struct Menu {
	#ifdef OS_MACOS
	unsigned char data[16];
	#elif OS_WINDOWS
	unsigned char data[64];
	#endif
} Menu;

// NOTE(nick): this has to be kept in sync with wry's ContextMenu struct size
typedef struct ContextMenu {
	#ifdef OS_MACOS
	unsigned char data[16];
	#elif OS_WINDOWS
	unsigned char data[64];
	#endif
} ContextMenu;

typedef struct Icon {
	unsigned char *data;
	int size;
} Icon;

typedef struct Window_Options {
	bool     always_on_top;
	bool     frameless;
	bool     fullscreen;
	Size     size;
	Size     min_size;
	Size     max_size;
	bool     maximized;
	Position position;
	bool     resizable;
	char *   title;
	bool     transparent;
	bool     visible;
	bool     center;
	Icon     icon;
	char *   url;
	char *   html;
	char *   script;
} Window_Options;

typedef struct Menu_Item {
	int  id;
	char *title;
	bool enabled;
	bool selected;
	char *accelerator;
} Menu_Item;

typedef struct Display {
	char *   name;
	Size     size;
	Position position;
	double   scale_factor;
} Display;

// @Cleanup: do we just want to make these be all be the generic Array?
typedef struct Array {
	void *data;
	int count;
} Array;

typedef struct StringArray {
	char **data;
	int count;
} StringArray;

typedef struct DisplayArray {
	Display *data;
	int      count;
} DisplayArray;

//
// Go Functions
//

typedef void (*closure)();

void go_app_main_loop();

//
// API Methods
//

EventLoop create_event_loop();
void run(EventLoop event_loop, void (*callback)(Event event));

void reset_temporary_storage();

int      window_create(EventLoop event_loop, Window_Options options, Menu menu);
bool     window_destroy(int window_id);
bool     window_set_title(int window_id, char *title);
bool     window_set_visible(int window_id, bool is_visible);
bool     window_set_focused(int window_id);
bool     window_set_fullscreen(int window_id, bool is_fullscreen);
bool     window_set_maximized(int window_id, bool is_maximized);
bool     window_set_minimized(int window_id, bool is_minimized);
bool     window_set_size(int window_id, Size size);
bool     window_set_min_size(int window_id, Size size);
bool     window_set_max_size(int window_id, Size size);
bool     window_set_resizable(int window_id, bool is_resizable);
bool     window_set_always_on_top(int window_id, bool is_on_top);
bool     window_set_position(int window_id, Position position);
Position window_get_outer_position(int window_id);
Size     window_get_outer_size(int window_id);
Position window_get_inner_position(int window_id);
Size     window_get_inner_size(int window_id);
double   window_get_dpi_scale(int window_id);
bool     window_is_visible(int window_id);

Menu menu_create();
bool menu_add_item(Menu menu, Menu_Item item);
bool menu_add_submenu(Menu menu, char *title, bool enabled, Menu submenu);

ContextMenu context_menu_create();
bool context_menu_add_item(ContextMenu menu, Menu_Item item);
bool context_menu_add_submenu(ContextMenu menu, char *title, bool enabled, ContextMenu submenu);

bool tray_set_system_tray(EventLoop event_loop, Icon icon, ContextMenu menu);

bool        shell_show_notification(char *title, char *subtitle, char *body);
bool        shell_show_dialog(char *title, char *body, char *level, char *buttons);
StringArray shell_show_file_picker(char *title, char *directory, char *filename, char *mode, char *filters);
char *      shell_read_clipboard();
bool        shell_write_clipboard(char *text);
bool        shell_register_shortcut(char *accelerator);
bool        shell_is_shortcut_registered(char *accelerator);
bool        shell_unregister_shortcut(char *accelerator);
bool        shell_unregister_all_shortcuts();

DisplayArray screen_get_available_displays();
