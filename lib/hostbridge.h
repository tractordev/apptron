
//void hello(char *name);
//void gomain();

static void invoke(void (*f)()) {
    f();
}

typedef void (*closure)();

typedef void (*Main_Loop_Callback)(int event_type);

void go_main_loop();

int run(Main_Loop_Callback callback);

int window_create(int width, int height, char *title);

typedef enum bool {
    false = 0,
    true = 1,
} bool;

bool window_set_title(int window_id, char *title);
