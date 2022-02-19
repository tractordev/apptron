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
	EventResized   = 4,
	EventMoved     = 5,
	EventMenuItem  = 6,
} EventType;

typedef struct Event {
	int      event_type;
	int      window_id;
	Position position;
	Size     size;
	int      menu_id;
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

typedef struct Window_Options {
	bool transparent;
	bool decorations;
	char *html;
} Window_Options;

typedef struct Menu_Item {
	int  id;
	char *title;
	bool enabled;
	bool selected;
	char *accelerator;
} Menu_Item;

typedef struct Icon {
	unsigned char *data;
	int size;
} Icon;

typedef struct StringArray {
	char **data;
	int count;
} StringArray;

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
bool     window_set_fullscreen(int window_id, bool is_fullscreen);
Position window_get_outer_position(int window_id);
Size     window_get_outer_size(int window_id);
Position window_get_inner_position(int window_id);
Size     window_get_inner_size(int window_id);
double   window_get_dpi_scale(int window_id);

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

bool  shell_write_clipboard(char *text);
char *shell_read_clipboard();
