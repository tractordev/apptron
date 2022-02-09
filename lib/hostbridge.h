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

typedef enum Event_Type {
	Event_Type__None      = 0,
	Event_Type__Close     = 1,
	Event_Type__Destroyed = 2,
	Event_Type__Focused   = 3,
	Event_Type__Resized   = 4,
	Event_Type__Moved     = 5,
} Event_Type;

typedef struct Event {
	int      event_type;
	int      window_id;
	Position position;
	Size     size;
} Event;

// NOTE(nick): this has to be kept in sync with wry's EventLoop struct size
typedef struct Event_Loop {
	unsigned char data[40];
} Event_Loop;

typedef struct Window_Options {
	bool transparent;
	bool decorations;
	char *html;
} Window_Options;

//
// Go Functions
//

typedef void (*closure)();

void go_app_main_loop();

//
// API Methods
//

Event_Loop create_event_loop();

int      window_create(Event_Loop event_loop, Window_Options options);
bool     window_destroy(int window_id);
bool     window_set_title(int window_id, char *title);
bool     window_set_visible(int window_id, bool is_visible);
bool     window_set_fullscreen(int window_id, bool is_fullscreen);
Position window_get_outer_position(int window_id);
Size     window_get_outer_size(int window_id);
Position window_get_inner_position(int window_id);
Size     window_get_inner_size(int window_id);
double   window_get_dpi_scale(int window_id);

void run(Event_Loop event_loop, void (*callback)(Event event));
