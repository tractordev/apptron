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
} EventType;

typedef struct Event {
	int      event_type;
	int      window_id;
	Position position;
	Size     size;
} Event;

// NOTE(nick): this has to be kept in sync with wry's EventLoop struct size
typedef struct EventLoop {
	unsigned char data[40];
} EventLoop;

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

EventLoop create_event_loop();

int      window_create(EventLoop event_loop, Window_Options options);
bool     window_destroy(int window_id);
bool     window_set_title(int window_id, char *title);
bool     window_set_visible(int window_id, bool is_visible);
bool     window_set_fullscreen(int window_id, bool is_fullscreen);
Position window_get_outer_position(int window_id);
Size     window_get_outer_size(int window_id);
Position window_get_inner_position(int window_id);
Size     window_get_inner_size(int window_id);
double   window_get_dpi_scale(int window_id);

void run(EventLoop event_loop, void (*callback)(Event event));
