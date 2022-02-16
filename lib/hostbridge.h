#include <stdint.h>

//
// Types
//

#define bool uint8_t

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
	unsigned char data[16];
} Menu;

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

//
// Go Functions
//

typedef void (*closure)();

void go_app_main_loop();

//
// API Methods
//

EventLoop create_event_loop();

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
bool menu_set_application_menu(Menu menu);

bool tray_set_system_tray(EventLoop event_loop, Icon icon, Menu_Item *item_data, int item_count);

void run(EventLoop event_loop, void (*callback)(Event event));

