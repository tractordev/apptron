//
// Types
//

typedef enum bool {
    false = 0,
    true = 1,
} bool;

typedef struct Vector2 {
    double x;
    double y;
} Vector2;

static void invoke(void (*f)()) {
    f();
}

typedef void (*closure)();

// NOTE(nick): this has to be kept in sync with wry's EventLoop struct size
typedef struct Event_Loop {
    unsigned char data[40];
} Event_Loop;

//
// Go Functions
//

void go_main_loop();

//
// API Methods
//

Event_Loop create_event_loop();

int     create_window(Event_Loop event_loop);
bool    destroy_window(int window_id);
bool    window_set_title(int window_id, char *title);
bool    window_set_foucs(int window_id, bool is_focused);
Vector2 window_get_outer_position(int window_id);
Vector2 window_get_outer_size(int window_id);
Vector2 window_get_inner_position(int window_id);
Vector2 window_get_inner_size(int window_id);
double  window_get_dpi_scale(int window_id);

void run(Event_Loop event_loop, void (*callback)());
